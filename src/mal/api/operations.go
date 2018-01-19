/**
 * MIT License
 *
 * Copyright (c) 2017 - 2018 CNES
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
package api

import (
	"errors"
	. "mal"
	"sync/atomic"
)

type OperationHandler interface {
	onMessage(msg *Message) error
	onClose() error
}

type OperationContext struct {
	Ctx       *Context
	Uri       *URI
	handlers  map[ULong]OperationHandler
	txcounter uint64
}

func NewOperationContext(ctx *Context, service string) (*OperationContext, error) {
	// TODO (AF): Verify the uri
	uri := ctx.NewURI(service)
	handlers := make(map[ULong]OperationHandler)
	ictx := &OperationContext{ctx, uri, handlers, 0}
	err := ctx.RegisterEndPoint(uri, ictx)
	if err != nil {
		return nil, err
	}
	return ictx, nil
}

func (ictx *OperationContext) register(tid ULong, handler OperationHandler) error {
	// TODO (AF): Synchronization
	old := ictx.handlers[tid]
	if old != nil {
		// TODO (AF): Log an error
		return errors.New("Handler already registered for this transaction")
	}
	ictx.handlers[tid] = handler
	return nil
}

func (ictx *OperationContext) deregister(tid ULong) error {
	// TODO (AF): Synchronization
	if ictx.handlers[tid] == nil {
		// TODO (AF): Log an error
		return errors.New("No handler registered for this transaction")
	}
	delete(ictx.handlers, tid)
	return nil
}

func (ictx *OperationContext) TransactionId() ULong {
	return ULong(atomic.AddUint64(&ictx.txcounter, 1))
}

func (ictx *OperationContext) Close() error {
	return ictx.Ctx.UnregisterEndPoint(ictx.Uri)
}

const (
	_CREATED byte = iota
	_INITIATED
	_ACKNOWLEDGED
	_PROGRESSING
	_REGISTERED
	_FINAL
	_CLOSED
)

type Operation interface {
	GetTid() ULong
	finalize() error
	Close() error
	Reset() error
}

type OperationX struct {
	ictx        *OperationContext
	tid         ULong
	ch          chan *Message
	urito       *URI
	area        UShort
	areaVersion UOctet
	service     UShort
	operation   UShort
	status      byte
}

// Finalize the operation
func (op *OperationX) finalize() error {
	op.status = _FINAL
	if op.ch != nil {
		return op.ictx.deregister(op.tid)
	}
	return nil
}

func (op *OperationX) GetTid() ULong {
	return op.tid
}

// Closes the operation.
// Be careful a closed operation cannot be used anymore.
func (op *OperationX) Close() error {
	if op.status == _CLOSED {
		return nil
	}
	op.status = _CLOSED
	if op.ch != nil {
		var err error = nil
		if (op.status != _CREATED) && (op.status != _FINAL) {
			err = op.ictx.deregister(op.tid)
		}
		close(op.ch)
		op.ch = nil
		return err
	}
	return nil
}

// Resets the operation for a new use, a new TransactionId is allocated
func (op *OperationX) Reset() error {
	if op.status != _FINAL {
		return errors.New("Bad operation status")
	}
	// Gets a new TransactionId for operation
	op.tid = op.ictx.TransactionId()
	op.status = _CREATED
	return nil
}

// ================================================================================
// SendOperation

type SendOperation interface {
	Operation
	Send(body []byte) error
}

type SendOperationX struct {
	OperationX
}

func (ictx *OperationContext) NewSendOperation(urito *URI, area UShort, areaVersion UOctet, service UShort, operation UShort) (SendOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	op := &SendOperationX{OperationX: OperationX{ictx, tid, nil, urito, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *SendOperationX) Send(body []byte) error {
	if op.status != _CREATED {
		return errors.New("Bad operation status")
	}

	op.status = _INITIATED
	msg := &Message{
		UriFrom:          op.ictx.Uri,
		UriTo:            op.urito,
		InteractionType:  MAL_INTERACTIONTYPE_SEND,
		InteractionStage: MAL_IP_STAGE_SEND,
		ServiceArea:      op.area,
		AreaVersion:      op.areaVersion,
		Service:          op.service,
		Operation:        op.operation,
		TransactionId:    op.tid,
		Body:             body,
	}
	// This operation doesn't wait any reply, so we don't need to register it.
	// Send the SEND MAL message
	err := op.ictx.Ctx.Send(msg)
	op.status = _FINAL
	return err
}

func (op *SendOperationX) onMessage(msg *Message) error {
	// TODO (AF): Should never reveive messages, log an error
	return nil
}

func (op *SendOperationX) onClose() error {
	// TODO (AF): Should never be called, log an error
	return nil
}

// ================================================================================
// SubmitOperation

type SubmitOperation interface {
	Operation
	Submit(body []byte) (*Message, error)
}

type SubmitOperationX struct {
	OperationX
}

func (ictx *OperationContext) NewSubmitOperation(urito *URI, area UShort, areaVersion UOctet, service UShort, operation UShort) (SubmitOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &SubmitOperationX{OperationX: OperationX{ictx, tid, ch, urito, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *SubmitOperationX) Submit(body []byte) (*Message, error) {
	if op.status != _CREATED {
		return nil, errors.New("Bad operation status")
	}

	op.status = _INITIATED
	msg := &Message{
		UriFrom:          op.ictx.Uri,
		UriTo:            op.urito,
		InteractionType:  MAL_INTERACTIONTYPE_SUBMIT,
		InteractionStage: MAL_IP_STAGE_SUBMIT,
		ServiceArea:      op.area,
		AreaVersion:      op.areaVersion,
		Service:          op.service,
		Operation:        op.operation,
		TransactionId:    op.tid,
		Body:             body,
	}
	// Registers this Submit Operation in OperationContext
	err := op.ictx.register(op.tid, op)
	if err != nil {
		op.status = _FINAL
		return nil, err
	}
	// Send the SUBMIT MAL message
	err = op.ictx.Ctx.Send(msg)
	if err != nil {
		op.status = _FINAL
		return nil, err
	}

	// Waits for the SUBMIT_ACK MAL message
	msg, more := <-op.ch
	if !more {
		logger.Debugf("Operation ends: %s, %s", op.ictx.Uri, op.tid)
		op.status = _FINAL
		return nil, errors.New("Operation ends")
	}
	// Verify that the message is ok (ack or error)
	if msg.InteractionStage != MAL_IP_STAGE_INVOKE_ACK {
		logger.Errorf("Bad return message, stage %d", msg.InteractionStage)
		op.status = _FINAL
		return nil, errors.New("Bad return message")
	}
	// This operation should not received anymore messages,
	// unregisters this Submit Operation in OperationContext
	op.status = _FINAL
	op.ictx.deregister(op.tid)

	return msg, nil
}

func (op *SubmitOperationX) onMessage(msg *Message) error {
	// TODO (AF): Verify the message: service area, version, service, operation
	if msg.InteractionType != MAL_INTERACTIONTYPE_SUBMIT {
		// TODO (AF): log an error
		return errors.New("Bad message")
	}
	op.ch <- msg
	return nil
}

func (op *SubmitOperationX) onClose() error {
	// TODO (AF):
	return nil
}

// ================================================================================
// RequestOperation

type RequestOperation interface {
	Operation
	Request(body []byte) (*Message, error)
}

type RequestOperationX struct {
	OperationX
}

func (ictx *OperationContext) NewRequestOperation(urito *URI, area UShort, areaVersion UOctet, service UShort, operation UShort) (RequestOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &RequestOperationX{OperationX: OperationX{ictx, tid, ch, urito, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *RequestOperationX) Request(body []byte) (*Message, error) {
	if op.status != _CREATED {
		return nil, errors.New("Bad operation status")
	}

	msg := &Message{
		UriFrom:          op.ictx.Uri,
		UriTo:            op.urito,
		InteractionType:  MAL_INTERACTIONTYPE_REQUEST,
		InteractionStage: MAL_IP_STAGE_REQUEST,
		ServiceArea:      op.area,
		AreaVersion:      op.areaVersion,
		Service:          op.service,
		Operation:        op.operation,
		TransactionId:    op.tid,
		Body:             body,
	}

	// Registers this Request Operation in OperationContext
	err := op.ictx.register(op.tid, op)
	if err != nil {
		return nil, err
	}
	op.status = _INITIATED

	// Send the REQUEST MAL message
	err = op.ictx.Ctx.Send(msg)
	if err != nil {
		op.status = _FINAL
		return nil, err
	}

	// Waits for the RESPONSE MAL message
	msg, more := <-op.ch
	if !more {
		logger.Debugf("Operation ends: %s, %s", op.ictx.Uri, op.tid)
		// TODO (AF): Returns an error
		op.status = _FINAL
		return nil, errors.New("Operation ends")
	}
	op.status = _FINAL

	// TODO (AF): Verify that the message is ok (response or error)
	if msg.InteractionStage == MAL_IP_STAGE_REQUEST_RESPONSE {
		// This operation should not received anymore messages,
		// unregisters this Request Operation in OperationContext
		op.ictx.deregister(op.tid)
		return msg, nil
	} else {
		// TODO (AF): Close the operation
		return nil, errors.New("Bad operation stage")
	}
}

func (op *RequestOperationX) onMessage(msg *Message) error {
	// TODO (AF): Verify the message: service area, version, service, operation
	if msg.InteractionType != MAL_INTERACTIONTYPE_REQUEST {
		// TODO (AF): log an error
		return errors.New("Bad message")
	}
	op.ch <- msg
	return nil
}

func (op *RequestOperationX) onClose() error {
	// TODO (AF):
	return nil
}

// ================================================================================
// Invoke Operation

type InvokeOperation interface {
	Operation
	Invoke(body []byte) (*Message, error)
	GetResponse() (*Message, error)
}

type InvokeOperationX struct {
	OperationX
	// TODO (AF): Handling of response (see api1)
	response *Message
}

func (ictx *OperationContext) NewInvokeOperation(urito *URI, area UShort, areaVersion UOctet, service UShort, operation UShort) (InvokeOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &InvokeOperationX{OperationX: OperationX{ictx, tid, ch, urito, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *InvokeOperationX) Invoke(body []byte) (*Message, error) {
	if op.status != _CREATED {
		return nil, errors.New("Bad operation status")
	}

	msg := &Message{
		UriFrom:          op.ictx.Uri,
		UriTo:            op.urito,
		InteractionType:  MAL_INTERACTIONTYPE_INVOKE,
		InteractionStage: MAL_IP_STAGE_INVOKE,
		ServiceArea:      op.area,
		AreaVersion:      op.areaVersion,
		Service:          op.service,
		Operation:        op.operation,
		TransactionId:    op.tid,
		Body:             body,
	}

	// Registers this Invoke Operation in OperationContext
	err := op.ictx.register(op.tid, op)
	if err != nil {
		return nil, err
	}
	op.status = _INITIATED

	// Send the INVOKE MAL message
	err = op.ictx.Ctx.Send(msg)
	if err != nil {
		op.status = _FINAL
		return nil, err
	}

	// Waits for the ACKNOWLEDGE MAL message
	msg, more := <-op.ch
	if !more {
		logger.Debugf("Operation ends: %s, %s", op.ictx.Uri, op.tid)
		// TODO (AF): Returns an error
		op.status = _FINAL
		return nil, errors.New("Operation ends")
	}
	op.status = _ACKNOWLEDGED

	// TODO (AF): Verify that the message is ok (ack or error)
	if msg.InteractionStage != MAL_IP_STAGE_INVOKE_ACK {
		// TODO (AF): Close the operation
		return nil, errors.New("Bad operation stage")
	}

	return msg, nil
}

// Returns the response.
func (op *InvokeOperationX) GetResponse() (*Message, error) {
	if (op.status == _FINAL) && (op.response != nil) {
		return op.response, nil
	}

	if op.status != _ACKNOWLEDGED {
		return nil, errors.New("Bad operation status")
	}

	// Waits for next MAL message
	msg, more := <-op.ch
	if !more {
		logger.Debugf("Operation ends: %s, %s", op.ictx.Uri, op.tid)
		// TODO (AF): Returns an error
		op.status = _FINAL
		return nil, errors.New("Operation ends")
	}

	if msg.InteractionStage == MAL_IP_STAGE_INVOKE_RESPONSE {
		op.response = msg
		op.status = _FINAL
		op.ictx.deregister(op.tid)
		return msg, nil
	} else {
		// TODO (AF): Close the operation
		op.status = _FINAL
		return nil, errors.New("Bad operation stage")
	}
}

func (op *InvokeOperationX) onMessage(msg *Message) error {
	// TODO (AF): Verify the message: service area, version, service, operation
	if msg.InteractionType != MAL_INTERACTIONTYPE_INVOKE {
		// TODO (AF): log an error
		return errors.New("Bad message")
	}
	op.ch <- msg
	return nil
}

func (op *InvokeOperationX) onClose() error {
	// TODO (AF):
	return nil
}

// ================================================================================
// Progress Operation

type ProgressOperation interface {
	Operation
	Progress(body []byte) (*Message, error)
	GetUpdate() (*Message, error)
	GetResponse() (*Message, error)
}

type ProgressOperationX struct {
	OperationX
	response *Message
}

func (ictx *OperationContext) NewProgressOperation(urito *URI, area UShort, areaVersion UOctet, service UShort, operation UShort) (ProgressOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &ProgressOperationX{OperationX: OperationX{ictx, tid, ch, urito, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *ProgressOperationX) Progress(body []byte) (*Message, error) {
	if op.status != _CREATED {
		return nil, errors.New("Bad operation status")
	}

	msg := &Message{
		UriFrom:          op.ictx.Uri,
		UriTo:            op.urito,
		InteractionType:  MAL_INTERACTIONTYPE_PROGRESS,
		InteractionStage: MAL_IP_STAGE_PROGRESS,
		ServiceArea:      op.area,
		AreaVersion:      op.areaVersion,
		Service:          op.service,
		Operation:        op.operation,
		TransactionId:    op.tid,
		Body:             body,
	}

	// Registers this Submit Operation in OperationContext
	err := op.ictx.register(op.tid, op)
	if err != nil {
		return nil, err
	}
	op.status = _INITIATED

	// Send the SUBMIT MAL message
	err = op.ictx.Ctx.Send(msg)
	if err != nil {
		op.status = _FINAL
		return nil, err
	}
	// Waits for the PROGRESS_ACK MAL message
	msg, more := <-op.ch
	if !more {
		logger.Debugf("Operation ends: %s, %s", op.ictx.Uri, op.tid)
		// TODO (AF): Returns an error
		op.status = _FINAL
		return nil, errors.New("Operation ends")
	}
	op.status = _ACKNOWLEDGED

	// TODO (AF): Verify that the message is ok (ack or error)
	if msg.InteractionStage != MAL_IP_STAGE_PROGRESS_ACK {
		// TODO (AF): Close the operation
		return nil, errors.New("Bad operation stage")
	}

	return msg, nil
}

// Returns next update or nil if there is no more update.
func (op *ProgressOperationX) GetUpdate() (*Message, error) {
	if (op.status != _ACKNOWLEDGED) && (op.status != _PROGRESSING) {
		return nil, errors.New("Bad operation status")
	}

	// Waits for next MAL message
	msg, more := <-op.ch
	if !more {
		logger.Debugf("Operation ends: %s, %s", op.ictx.Uri, op.tid)
		// TODO (AF): Returns an error
		op.status = _FINAL
		return nil, errors.New("Operation ends")
	}

	if msg.InteractionStage == MAL_IP_STAGE_PROGRESS_UPDATE {
		return msg, nil
	} else if msg.InteractionStage == MAL_IP_STAGE_PROGRESS_RESPONSE {
		op.response = msg
		op.status = _FINAL
		op.ictx.deregister(op.tid)
		return nil, nil
	} else {
		// TODO (AF): Close the operation
		return nil, errors.New("Bad operation stage")
	}
}

// Returns the response.
func (op *ProgressOperationX) GetResponse() (*Message, error) {
	if (op.status == _FINAL) && (op.response != nil) {
		return op.response, nil
	}

	if (op.status != _ACKNOWLEDGED) && (op.status != _PROGRESSING) {
		return nil, errors.New("Bad operation status")
	}

	// Waits for next MAL message
	msg, more := <-op.ch
	if !more {
		logger.Debugf("Operation ends: %s, %s", op.ictx.Uri, op.tid)
		// TODO (AF): Returns an error
		op.status = _FINAL
		return nil, errors.New("Operation ends")
	}

	if msg.InteractionStage == MAL_IP_STAGE_PROGRESS_RESPONSE {
		op.response = msg
		op.status = _FINAL
		op.ictx.deregister(op.tid)
		return msg, nil
	} else {
		// TODO (AF): Close the operation
		return nil, errors.New("Bad operation stage")
	}
}

func (op *ProgressOperationX) onMessage(msg *Message) error {
	// TODO (AF): Verify the message: service area, version, service, operation
	if msg.InteractionType != MAL_INTERACTIONTYPE_PROGRESS {
		// TODO (AF): log an error
		return errors.New("Bad message")
	}
	op.ch <- msg
	return nil
}

func (op *ProgressOperationX) onClose() error {
	// TODO (AF):
	return nil
}

// ================================================================================
// Subscriber Operation

type SubscriberOperation interface {
	Operation
	Register(body []byte) error
	GetNotify() (*Message, error)
	Deregister(body []byte) error
}

type SubscriberOperationX struct {
	OperationX
}

func (ictx *OperationContext) NewSubscriberOperation(urito *URI, area UShort, areaVersion UOctet, service UShort, operation UShort) (SubscriberOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &SubscriberOperationX{OperationX: OperationX{ictx, tid, ch, urito, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *SubscriberOperationX) Register(body []byte) error {
	// TODO (AF): Can we register anew a Subscriber?
	if op.status != _CREATED {
		return errors.New("Bad operation status")
	}

	msg := &Message{
		UriFrom:          op.ictx.Uri,
		UriTo:            op.urito,
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_REGISTER,
		ServiceArea:      op.area,
		AreaVersion:      op.areaVersion,
		Service:          op.service,
		Operation:        op.operation,
		TransactionId:    op.tid,
		Body:             body,
	}

	// Registers this Operation in OperationContext
	err := op.ictx.register(op.tid, op)
	if err != nil {
		return err
	}
	op.status = _INITIATED

	// Send the REGISTER MAL message
	err = op.ictx.Ctx.Send(msg)
	if err != nil {
		op.status = _FINAL
		return err
	}

	// Waits for the REGISTER_ACK MAL message
	msg, more := <-op.ch
	if !more {
		logger.Debugf("Operation ends: %s, %s", op.ictx.Uri, op.tid)
		// TODO (AF): Returns an error
		op.status = _FINAL
		return errors.New("Operation ends")
	}
	op.status = _REGISTERED

	// TODO (AF): Verify that the message is ok (ack or error)
	if msg.InteractionStage != MAL_IP_STAGE_PUBSUB_REGISTER_ACK {
		// TODO (AF): Close the operation
		return errors.New("Bad operation stage")
	}

	return nil
}

// Returns next update or nil if there is no more update.
func (op *SubscriberOperationX) GetNotify() (*Message, error) {
	if op.status != _REGISTERED {
		return nil, errors.New("Bad operation status")
	}

	// Waits for next MAL message
	msg, more := <-op.ch
	if !more {
		logger.Debugf("Operation ends: %s, %s", op.ictx.Uri, op.tid)
		// TODO (AF): Returns an error
		op.status = _FINAL
		return nil, errors.New("Operation ends")
	}

	if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_NOTIFY {
		return msg, nil
	} else {
		// TODO (AF): Close the operation
		return nil, errors.New("Bad operation stage")
	}
}

func (op *SubscriberOperationX) Deregister(body []byte) error {
	if op.status != _REGISTERED {
		return errors.New("Bad operation status")
	}

	msg := &Message{
		UriFrom:          op.ictx.Uri,
		UriTo:            op.urito,
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_DEREGISTER,
		ServiceArea:      op.area,
		AreaVersion:      op.areaVersion,
		Service:          op.service,
		Operation:        op.operation,
		TransactionId:    op.tid,
		Body:             body,
	}

	// Send the DEREGISTER MAL message
	err := op.ictx.Ctx.Send(msg)
	if err != nil {
		op.status = _FINAL
		return err
	}

	// Removes useless notify waiting messages
	for {
		// Waits for the PROGRESS_ACK MAL message
		msg, more := <-op.ch
		if !more {
			logger.Debugf("Operation ends: %s, %s", op.ictx.Uri, op.tid)
			// TODO (AF): Returns an error
			op.status = _FINAL
			return errors.New("Operation ends")
		}

		if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_NOTIFY {
			continue
		}

		op.ictx.deregister(op.tid)
		op.status = _FINAL

		// TODO (AF): Verify that the message is ok (ack or error)
		if msg.InteractionStage != MAL_IP_STAGE_PUBSUB_DEREGISTER_ACK {
			// TODO (AF): Close the operation
			return errors.New("Bad operation stage")
		}

		return nil
	}
}

func (op *SubscriberOperationX) onMessage(msg *Message) error {
	// TODO (AF): Verify the message: service area, version, service, operation
	if msg.InteractionType != MAL_INTERACTIONTYPE_PUBSUB {
		// TODO (AF): log an error
		return errors.New("Bad message")
	}
	op.ch <- msg
	return nil
}

func (op *SubscriberOperationX) onClose() error {
	// TODO (AF):
	return nil
}

// ================================================================================
// Publisher Operation

type PublisherOperation interface {
	Operation
	Register(body []byte) error
	Publish(body []byte) error
	Deregister(body []byte) error
}

type PublisherOperationX struct {
	OperationX
}

func (ictx *OperationContext) NewPublisherOperation(urito *URI, area UShort, areaVersion UOctet, service UShort, operation UShort) (PublisherOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &PublisherOperationX{OperationX: OperationX{ictx, tid, ch, urito, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *PublisherOperationX) Register(body []byte) error {
	// TODO (AF): We can register anew a Subscriber
	if op.status != _CREATED {
		return errors.New("Bad operation status")
	}

	msg := &Message{
		UriFrom:          op.ictx.Uri,
		UriTo:            op.urito,
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER,
		ServiceArea:      op.area,
		AreaVersion:      op.areaVersion,
		Service:          op.service,
		Operation:        op.operation,
		TransactionId:    op.tid,
		Body:             body,
	}

	// Registers this Operation in OperationContext
	err := op.ictx.register(op.tid, op)
	if err != nil {
		return err
	}
	op.status = _INITIATED

	// Send the PUBLISH_REGISTER MAL message
	err = op.ictx.Ctx.Send(msg)
	if err != nil {
		op.status = _FINAL
		return err
	}

	// Waits for the PUBLISH_REGISTER_ACK MAL message
	msg, more := <-op.ch
	if !more {
		logger.Debugf("Operation ends: %s, %s", op.ictx.Uri, op.tid)
		// TODO (AF): Returns an error
		op.status = _FINAL
		return errors.New("Operation ends")
	}
	op.status = _REGISTERED

	// TODO (AF): Verify that the message is ok (ack or error)
	if msg.InteractionStage != MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER_ACK {
		// TODO (AF): Close the operation
		return errors.New("Bad operation stage")
	}

	return nil
}

func (op *PublisherOperationX) Publish(body []byte) error {
	if op.status != _REGISTERED {
		return errors.New("Bad operation status")
	}

	msg := &Message{
		UriFrom:          op.ictx.Uri,
		UriTo:            op.urito,
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_PUBLISH,
		ServiceArea:      op.area,
		AreaVersion:      op.areaVersion,
		Service:          op.service,
		Operation:        op.operation,
		TransactionId:    op.tid,
		Body:             body,
	}

	// Send the MAL message
	err := op.ictx.Ctx.Send(msg)
	if err != nil {
		op.status = _FINAL
		return err
	}

	return nil
}

func (op *PublisherOperationX) Deregister(body []byte) error {
	if op.status != _REGISTERED {
		return errors.New("Bad operation status")
	}

	msg := &Message{
		UriFrom:          op.ictx.Uri,
		UriTo:            op.urito,
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER,
		ServiceArea:      op.area,
		AreaVersion:      op.areaVersion,
		Service:          op.service,
		Operation:        op.operation,
		TransactionId:    op.tid,
		Body:             body,
	}

	// Send the DEREGISTER MAL message
	err := op.ictx.Ctx.Send(msg)
	if err != nil {
		op.status = _FINAL
		return err
	}

	// Waits for the PUBLISH_DEREGISTER_ACK MAL message
	msg, more := <-op.ch
	if !more {
		logger.Debugf("Operation ends: %s, %s", op.ictx.Uri, op.tid)
		// TODO (AF): Returns an error
		op.status = _FINAL
		return errors.New("Operation ends")
	}
	op.status = _FINAL
	op.ictx.deregister(op.tid)

	// TODO (AF): Verify that the message is ok (ack or error)
	if msg.InteractionStage != MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER_ACK {
		// TODO (AF): Close the operation
		return errors.New("Bad operation stage")
	}

	return nil
}

func (op *PublisherOperationX) onMessage(msg *Message) error {
	// TODO (AF): Verify the message: service area, version, service, operation
	if msg.InteractionType != MAL_INTERACTIONTYPE_PUBSUB {
		// TODO (AF): log an error
		return errors.New("Bad message")
	}
	op.ch <- msg
	return nil
}

func (op *PublisherOperationX) onClose() error {
	// TODO (AF):
	return nil
}

// ================================================================================
// Defines Listener interface used by context to route MAL messages

func (ictx *OperationContext) OnMessage(msg *Message) error {
	to, ok := ictx.handlers[msg.TransactionId]
	if ok {
		logger.Debugf("onMessage %t", to)
		to.onMessage(msg)
		logger.Debugf("OnMessageMessage transmitted: %s", msg)
	} else {
		logger.Debugf("Cannot route message to: %s?TransactionId=", msg.UriTo, msg.TransactionId)
	}
	return nil
}

func (ictx *OperationContext) OnClose() error {
	logger.Infof("close EndPoint: %s", ictx.Uri)
	for tid, handler := range ictx.handlers {
		logger.Debugf("close operation: %d", tid)
		err := handler.onClose()
		if err != nil {
			// TODO (AF): print an error message
		}
	}
	return nil
}
