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
	OnMessage(msg *Message) error
	OnClose() error
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
	// TODO (AF): Should define PubSub status
)

type Operation interface {
	GetTid() ULong
}

type OperationX struct {
	ictx        *OperationContext
	tid         ULong
	ch          chan *Message
	area        UShort // TODO (AF): Should be in OperationContext or Context?
	areaVersion UOctet // TODO (AF): Should be in OperationContext or Context?
	service     UShort // TODO (AF): Should be in OperationContext or Context?
	operation   UShort
	status      byte
}

func (op *OperationX) GetTid() ULong {
	return op.tid
}

// ================================================================================
// SendOperation

type SendOperation interface {
	Operation
	Send(msg *Message) error
}

type SendOperationX struct {
	OperationX
}

func (ictx *OperationContext) NewSendOperation(area UShort, areaVersion UOctet, service UShort, operation UShort) (SendOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	op := &SendOperationX{OperationX: OperationX{ictx, tid, nil, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *SendOperationX) Send(msg *Message) error {
	if op.status != _CREATED {
		return errors.New("Bad operation status")
	}

	msg.ServiceArea = op.area
	msg.AreaVersion = op.areaVersion
	msg.Service = op.service
	msg.Operation = op.operation
	msg.InteractionType = MAL_INTERACTIONTYPE_SEND
	msg.InteractionStage = MAL_IP_STAGE_SEND
	msg.TransactionId = op.tid
	msg.UriFrom = op.ictx.Uri

	// This operation doesn't wait any reply, so we don't need to register it.
	op.status = _INITIATED
	op.status = _FINAL

	return op.ictx.Ctx.Send(msg)
}

func (op *SendOperationX) OnMessage(msg *Message) error {
	// TODO (AF): Should never reveive messages, log an error
	return nil
}

func (op *SendOperationX) OnClose() error {
	// TODO (AF): Should never be called, log an error
	return nil
}

// ================================================================================
// SubmitOperation

type SubmitOperation interface {
	Operation
	Submit(msg *Message) (*Message, error)
}

type SubmitOperationX struct {
	OperationX
}

func (ictx *OperationContext) NewSubmitOperation(area UShort, areaVersion UOctet, service UShort, operation UShort) (SubmitOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &SubmitOperationX{OperationX: OperationX{ictx, tid, ch, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *SubmitOperationX) Submit(msg *Message) (*Message, error) {
	if op.status != _CREATED {
		return nil, errors.New("Bad operation status")
	}

	msg.ServiceArea = op.area
	msg.AreaVersion = op.areaVersion
	msg.Service = op.service
	msg.Operation = op.operation
	msg.InteractionType = MAL_INTERACTIONTYPE_SUBMIT
	msg.InteractionStage = MAL_IP_STAGE_SUBMIT
	msg.TransactionId = op.tid
	msg.UriFrom = op.ictx.Uri

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

	// Waits for the SUBMIT_ACK MAL message
	msg, more := <-op.ch
	if !more {
		logger.Debugf("Operation ends: %s, %s", op.ictx.Uri, op.tid)
		// TODO (AF): Returns an error
		op.status = _FINAL
		return nil, errors.New("Operation ends")
	}
	op.status = _ACKNOWLEDGED

	// TODO (AF): Verify that the message is ok (ack or error)

	// TODO (AF): This operation should not received anymore messages,
	// unregisters this Submit Operation in OperationContext
	op.status = _FINAL
	op.ictx.deregister(op.tid)

	return msg, nil
}

func (op *SubmitOperationX) OnMessage(msg *Message) error {
	// TODO (AF): Verify the message: service area, version, service, operation
	if msg.InteractionType != MAL_INTERACTIONTYPE_SUBMIT {
		// TODO (AF): log an error
		return errors.New("Bad message")
	}
	op.ch <- msg
	return nil
}

func (op *SubmitOperationX) OnClose() error {
	// TODO (AF):
	return nil
}

// ================================================================================
// RequestOperation

type RequestOperation interface {
	Operation
	Request(msg *Message) (*Message, error)
}

type RequestOperationX struct {
	OperationX
}

func (ictx *OperationContext) NewRequestOperation(area UShort, areaVersion UOctet, service UShort, operation UShort) (RequestOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &RequestOperationX{OperationX: OperationX{ictx, tid, ch, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *RequestOperationX) Request(msg *Message) (*Message, error) {
	if op.status != _CREATED {
		return nil, errors.New("Bad operation status")
	}

	msg.ServiceArea = op.area
	msg.AreaVersion = op.areaVersion
	msg.Service = op.service
	msg.Operation = op.operation
	msg.InteractionType = MAL_INTERACTIONTYPE_REQUEST
	msg.InteractionStage = MAL_IP_STAGE_REQUEST
	msg.TransactionId = op.tid
	msg.UriFrom = op.ictx.Uri

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

func (op *RequestOperationX) OnMessage(msg *Message) error {
	// TODO (AF): Verify the message: service area, version, service, operation
	if msg.InteractionType != MAL_INTERACTIONTYPE_REQUEST {
		// TODO (AF): log an error
		return errors.New("Bad message")
	}
	op.ch <- msg
	return nil
}

func (op *RequestOperationX) OnClose() error {
	// TODO (AF):
	return nil
}

// ================================================================================
// Invoke Operation

type InvokeOperation interface {
	Operation
	Invoke(msg *Message) (*Message, error)
	GetResponse() (*Message, error)
}

type InvokeOperationX struct {
	OperationX
	// TODO (AF): Handling of response (see api1)
	response *Message
}

func (ictx *OperationContext) NewInvokeOperation(area UShort, areaVersion UOctet, service UShort, operation UShort) (InvokeOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &InvokeOperationX{OperationX: OperationX{ictx, tid, ch, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *InvokeOperationX) Invoke(msg *Message) (*Message, error) {
	if op.status != _CREATED {
		return nil, errors.New("Bad operation status")
	}

	msg.ServiceArea = op.area
	msg.AreaVersion = op.areaVersion
	msg.Service = op.service
	msg.Operation = op.operation
	msg.InteractionType = MAL_INTERACTIONTYPE_INVOKE
	msg.InteractionStage = MAL_IP_STAGE_INVOKE
	msg.TransactionId = op.tid
	msg.UriFrom = op.ictx.Uri

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

func (op *InvokeOperationX) OnMessage(msg *Message) error {
	// TODO (AF): Verify the message: service area, version, service, operation
	if msg.InteractionType != MAL_INTERACTIONTYPE_INVOKE {
		// TODO (AF): log an error
		return errors.New("Bad message")
	}
	op.ch <- msg
	return nil
}

func (op *InvokeOperationX) OnClose() error {
	// TODO (AF):
	return nil
}

// ================================================================================
// Progress Operation

type ProgressOperation interface {
	Operation
	Progress(msg *Message) (*Message, error)
	GetUpdate() (*Message, error)
	GetResponse() (*Message, error)
}

type ProgressOperationX struct {
	OperationX
	response *Message
}

func (ictx *OperationContext) NewProgressOperation(area UShort, areaVersion UOctet, service UShort, operation UShort) (ProgressOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &ProgressOperationX{OperationX: OperationX{ictx, tid, ch, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *ProgressOperationX) Progress(msg *Message) (*Message, error) {
	if op.status != _CREATED {
		return nil, errors.New("Bad operation status")
	}

	msg.ServiceArea = op.area
	msg.AreaVersion = op.areaVersion
	msg.Service = op.service
	msg.Operation = op.operation
	msg.InteractionType = MAL_INTERACTIONTYPE_PROGRESS
	msg.InteractionStage = MAL_IP_STAGE_PROGRESS
	msg.TransactionId = op.tid
	msg.UriFrom = op.ictx.Uri

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

func (op *ProgressOperationX) OnMessage(msg *Message) error {
	// TODO (AF): Verify the message: service area, version, service, operation
	if msg.InteractionType != MAL_INTERACTIONTYPE_PROGRESS {
		// TODO (AF): log an error
		return errors.New("Bad message")
	}
	op.ch <- msg
	return nil
}

func (op *ProgressOperationX) OnClose() error {
	// TODO (AF):
	return nil
}

// ================================================================================
// Subscriber Operation

type SubscriberOperation interface {
	Operation
	Register(msg *Message) error
	GetNotify() (*Message, error)
	Deregister(msg *Message) error
}

type SubscriberOperationX struct {
	OperationX
}

func (ictx *OperationContext) NewSubscriberOperation(area UShort, areaVersion UOctet, service UShort, operation UShort) (SubscriberOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &SubscriberOperationX{OperationX: OperationX{ictx, tid, ch, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *SubscriberOperationX) Register(msg *Message) error {
	// TODO (AF): Can we register anew a Subscriber?
	if op.status != _CREATED {
		return errors.New("Bad operation status")
	}

	msg.ServiceArea = op.area
	msg.AreaVersion = op.areaVersion
	msg.Service = op.service
	msg.Operation = op.operation
	msg.InteractionType = MAL_INTERACTIONTYPE_PUBSUB
	msg.InteractionStage = MAL_IP_STAGE_PUBSUB_REGISTER
	msg.TransactionId = op.tid
	msg.UriFrom = op.ictx.Uri

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

func (op *SubscriberOperationX) Deregister(msg *Message) error {
	if op.status != _REGISTERED {
		return errors.New("Bad operation status")
	}

	msg.ServiceArea = op.area
	msg.AreaVersion = op.areaVersion
	msg.Service = op.service
	msg.Operation = op.operation
	msg.InteractionType = MAL_INTERACTIONTYPE_PUBSUB
	msg.InteractionStage = MAL_IP_STAGE_PUBSUB_DEREGISTER
	msg.TransactionId = op.tid
	msg.UriFrom = op.ictx.Uri

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

func (op *SubscriberOperationX) OnMessage(msg *Message) error {
	// TODO (AF): Verify the message: service area, version, service, operation
	if msg.InteractionType != MAL_INTERACTIONTYPE_PUBSUB {
		// TODO (AF): log an error
		return errors.New("Bad message")
	}
	op.ch <- msg
	return nil
}

func (op *SubscriberOperationX) OnClose() error {
	// TODO (AF):
	return nil
}

// ================================================================================
// Publisher Operation

type PublisherOperation interface {
	Operation
	Register(msg *Message) error
	Publish(msg *Message) error
	Deregister(msg *Message) error
}

type PublisherOperationX struct {
	OperationX
}

func (ictx *OperationContext) NewPublisherOperation(area UShort, areaVersion UOctet, service UShort, operation UShort) (PublisherOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &PublisherOperationX{OperationX: OperationX{ictx, tid, ch, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *PublisherOperationX) Register(msg *Message) error {
	// TODO (AF): We can register anew a Subscriber
	if op.status != _CREATED {
		return errors.New("Bad operation status")
	}

	msg.ServiceArea = op.area
	msg.AreaVersion = op.areaVersion
	msg.Service = op.service
	msg.Operation = op.operation
	msg.InteractionType = MAL_INTERACTIONTYPE_PUBSUB
	msg.InteractionStage = MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER
	msg.TransactionId = op.tid
	msg.UriFrom = op.ictx.Uri

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

func (op *PublisherOperationX) Publish(msg *Message) error {
	if op.status != _REGISTERED {
		return errors.New("Bad operation status")
	}

	msg.ServiceArea = op.area
	msg.AreaVersion = op.areaVersion
	msg.Service = op.service
	msg.Operation = op.operation
	msg.InteractionType = MAL_INTERACTIONTYPE_PUBSUB
	msg.InteractionStage = MAL_IP_STAGE_PUBSUB_PUBLISH
	msg.TransactionId = op.tid
	msg.UriFrom = op.ictx.Uri

	// Send the MAL message
	err := op.ictx.Ctx.Send(msg)
	if err != nil {
		op.status = _FINAL
		return err
	}

	return nil
}

func (op *PublisherOperationX) Deregister(msg *Message) error {
	if op.status != _REGISTERED {
		return errors.New("Bad operation status")
	}

	msg.ServiceArea = op.area
	msg.AreaVersion = op.areaVersion
	msg.Service = op.service
	msg.Operation = op.operation
	msg.InteractionType = MAL_INTERACTIONTYPE_PUBSUB
	msg.InteractionStage = MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER
	msg.TransactionId = op.tid
	msg.UriFrom = op.ictx.Uri

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

func (op *PublisherOperationX) OnMessage(msg *Message) error {
	// TODO (AF): Verify the message: service area, version, service, operation
	if msg.InteractionType != MAL_INTERACTIONTYPE_PUBSUB {
		// TODO (AF): log an error
		return errors.New("Bad message")
	}
	op.ch <- msg
	return nil
}

func (op *PublisherOperationX) OnClose() error {
	// TODO (AF):
	return nil
}

// ================================================================================
// Defines Listener interface used by context to route MAL messages

func (ictx *OperationContext) OnMessage(msg *Message) error {
	to, ok := ictx.handlers[msg.TransactionId]
	if ok {
		logger.Debugf("OnMessage %t", to)
		to.OnMessage(msg)
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
		err := handler.OnClose()
		if err != nil {
			// TODO (AF): print an error message
		}
	}
	return nil
}
