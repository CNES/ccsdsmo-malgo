/**
 * MIT License
 *
 * Copyright (c) 2017 - 2020 CNES
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
	. "github.com/CNES/ccsdsmo-malgo/mal"
)

const (
	_CREATED byte = iota
	_INITIATED
	_ACKNOWLEDGED
	_PROGRESSING
	_REGISTER_INITIATED
	_REGISTERED
	_REREGISTER_INITIATED
	_DEREGISTER_INITIATED
	_FINAL
	_CLOSED
)

type Operation interface {
	// Get current TransactionId
	GetTid() ULong
	// Returns a new Body ready to encode
	NewBody() Body
	// Interrupt the operation during a blocking processing.
	Interrupt()
	// Reset the operation in order to reuse it
	Reset() error
	// Close the operation
	Close() error
}

type OperationX struct {
	cctx        *ClientContext
	tid         ULong
	ch          chan *Message
	urito       *URI
	area        UShort
	areaVersion UOctet
	service     UShort
	operation   UShort
	status      byte
}

// Verifies that the incoming message corresponds to the initiated operation
func (op *OperationX) verify(msg *Message) bool {
	if (msg.ServiceArea == op.area) && (msg.AreaVersion == op.areaVersion) &&
		(msg.Service == op.service) && (msg.Operation == op.operation) {
		return true
	}
	return false
}

// Finalize the operation
func (op *OperationX) finalize() {
	op.status = _FINAL
	if op.ch != nil {
		// This operation should not received anymore messages, unregisters it
		// in OperationContext
		op.cctx.deregisterOp(op.tid)
	}
}

func (op *OperationX) GetTid() ULong {
	return op.tid
}

func (op *OperationX) NewBody() Body {
	return op.cctx.Ctx.NewBody()
}

// Interrupts the operation.
func (op *OperationX) Interrupt() {
	if op.status == _CLOSED {
		return
	}
	if op.ch != nil {
		op.ch <- nil
	}
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
			err = op.cctx.deregisterOp(op.tid)
		}
		close(op.ch)
		op.ch = nil
		return err
	}
	return nil
}

// Resets the operation for a new use, a new TransactionId is allocated.
// Be careful, the operation must be in a FINAL state
func (op *OperationX) Reset() error {
	if op.status != _FINAL {
		return errors.New("Bad operation status")
	}
	// Gets a new TransactionId for operation
	op.tid = op.cctx.TransactionId()
	op.status = _CREATED
	return nil
}

// ================================================================================
// SendOperation

type SendOperation interface {
	Operation
	Send(body Body) error
}

type SendOperationX struct {
	OperationX
}

func (cctx *ClientContext) NewSendOperation(urito *URI, area UShort, areaVersion UOctet, service UShort, operation UShort) SendOperation {
	// Gets a new TransactionId for operation
	tid := cctx.TransactionId()
	op := &SendOperationX{OperationX: OperationX{cctx, tid, nil, urito, area, areaVersion, service, operation, _CREATED}}
	return op
}

func (op *SendOperationX) Send(body Body) error {
	if op.status != _CREATED {
		return errors.New("Bad operation status")
	}
	op.status = _INITIATED

	msg := &Message{
		UriFrom:          op.cctx.Uri,
		UriTo:            op.urito,
		AuthenticationId: op.cctx.AuthenticationId,
		EncodingId:       op.cctx.EncodingId,
		QoSLevel:         op.cctx.QoSLevel,
		Priority:         op.cctx.Priority,
		Domain:           op.cctx.Domain,
		NetworkZone:      op.cctx.NetworkZone,
		Session:          op.cctx.Session,
		SessionName:      op.cctx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_SEND,
		InteractionStage: MAL_IP_STAGE_SEND,
		TransactionId:    op.tid,
		ServiceArea:      op.area,
		Service:          op.service,
		Operation:        op.operation,
		AreaVersion:      op.areaVersion,
		Body:             body,
	}
	// This operation doesn't wait any reply, so we don't need to register it.
	// Send the SEND MAL message
	err := op.cctx.Ctx.Send(msg)
	op.status = _FINAL
	return err
}

func (op *SendOperationX) onMessage(msg *Message) {
	logger.Errorf("SendOperation.onMessage: should never happened!")
}

// This function is called when underlying ClientContext is closed
func (op *SendOperationX) onClose() {
	// Nothing to do
}

// ================================================================================
// SubmitOperation

type SubmitOperation interface {
	Operation
	Submit(body Body) (*Message, error)
}

type SubmitOperationX struct {
	OperationX
}

func (cctx *ClientContext) NewSubmitOperation(urito *URI, area UShort, areaVersion UOctet, service UShort, operation UShort) SubmitOperation {
	// Gets a new TransactionId for operation
	tid := cctx.TransactionId()
	// Normally should not receive more than one message
	ch := make(chan *Message)
	op := &SubmitOperationX{OperationX: OperationX{cctx, tid, ch, urito, area, areaVersion, service, operation, _CREATED}}
	return op
}

func (op *SubmitOperationX) Submit(body Body) (*Message, error) {
	if op.status != _CREATED {
		return nil, errors.New("Bad operation status")
	}
	op.status = _INITIATED

	msg := &Message{
		UriFrom:          op.cctx.Uri,
		UriTo:            op.urito,
		AuthenticationId: op.cctx.AuthenticationId,
		EncodingId:       op.cctx.EncodingId,
		QoSLevel:         op.cctx.QoSLevel,
		Priority:         op.cctx.Priority,
		Domain:           op.cctx.Domain,
		NetworkZone:      op.cctx.NetworkZone,
		Session:          op.cctx.Session,
		SessionName:      op.cctx.SessionName,
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
	err := op.cctx.registerOp(op.tid, op)
	if err != nil {
		op.finalize()
		return nil, err
	}
	// Send the SUBMIT MAL message
	err = op.cctx.Ctx.Send(msg)
	if err != nil {
		op.finalize()
		return nil, err
	}

	// Waits for the SUBMIT_ACK MAL message
	msg, more := <-op.ch
	if !more {
		op.finalize()
		logger.Errorf("SubmitOperation.Sumit: Operation ends: %s, %s", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation ends")
	}
	if msg == nil {
		op.finalize()
		logger.Warnf("SubmitOperation.Sumit: Operation interupted (%s, %s)", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation interupted")
	}
	if msg.InteractionStage != MAL_IP_STAGE_SUBMIT_ACK {
		op.finalize()
		logger.Errorf("SubmitOperation.Sumit: Bad return message, operation (%s, %s), stage %d", op.cctx.Uri, op.tid, msg.InteractionStage)
		return nil, errors.New("Bad return message")
	}
	op.finalize()
	// Verify that the message is ok (ack or error)
	if msg.IsErrorMessage {
		return msg, errors.New("Error message")
	} else {
		return msg, nil
	}
}

func (op *SubmitOperationX) onMessage(msg *Message) {
	// Verify the message: service area, version, service, operation
	if op.verify(msg) && (msg.InteractionType == MAL_INTERACTIONTYPE_SUBMIT) {
		op.ch <- msg
	} else {
		logger.Errorf("SUBMIT Operation (%s,%d) receives Bad message: %+v", *op.urito, op.tid, msg)
	}
}

// This function is called when underlying ClientContext is closed
func (op *SubmitOperationX) onClose() {
	op.Close()
}

// ================================================================================
// RequestOperation

type RequestOperation interface {
	Operation
	Request(body Body) (*Message, error)
}

type RequestOperationX struct {
	OperationX
}

func (cctx *ClientContext) NewRequestOperation(urito *URI, area UShort, areaVersion UOctet, service UShort, operation UShort) RequestOperation {
	// Gets a new TransactionId for operation
	tid := cctx.TransactionId()
	// Normally should not receive more than one message
	ch := make(chan *Message)
	op := &RequestOperationX{OperationX: OperationX{cctx, tid, ch, urito, area, areaVersion, service, operation, _CREATED}}
	return op
}

func (op *RequestOperationX) Request(body Body) (*Message, error) {
	if op.status != _CREATED {
		return nil, errors.New("Bad operation status")
	}
	op.status = _INITIATED

	msg := &Message{
		UriFrom:          op.cctx.Uri,
		UriTo:            op.urito,
		AuthenticationId: op.cctx.AuthenticationId,
		EncodingId:       op.cctx.EncodingId,
		QoSLevel:         op.cctx.QoSLevel,
		Priority:         op.cctx.Priority,
		Domain:           op.cctx.Domain,
		NetworkZone:      op.cctx.NetworkZone,
		Session:          op.cctx.Session,
		SessionName:      op.cctx.SessionName,
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
	err := op.cctx.registerOp(op.tid, op)
	if err != nil {
		op.finalize()
		return nil, err
	}
	// Send the REQUEST MAL message
	err = op.cctx.Ctx.Send(msg)
	if err != nil {
		op.finalize()
		return nil, err
	}

	// Waits for the RESPONSE MAL message
	msg, more := <-op.ch
	if !more {
		op.finalize()
		logger.Debugf("RequetOperation.Request: Operation ends: %s, %s", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation ends")
	}
	if msg == nil {
		op.finalize()
		logger.Warnf("RequetOperation.Request: Operation interupted (%s, %s)", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation interupted")
	}
	// Verify the message stage
	if msg.InteractionStage != MAL_IP_STAGE_REQUEST_RESPONSE {
		op.finalize()
		logger.Errorf("RequetOperation.Request: Bad return message, operation (%s, %s), stage %d", op.cctx.Uri, op.tid, msg.InteractionStage)
		return nil, errors.New("Bad return message")
	}
	op.finalize()
	// Verify that the message is ok (ack or error)
	if msg.IsErrorMessage {
		return msg, errors.New("Error message")
	} else {
		return msg, nil
	}
}

func (op *RequestOperationX) onMessage(msg *Message) {
	// Verify the message: service area, version, service, operation
	if op.verify(msg) && (msg.InteractionType == MAL_INTERACTIONTYPE_REQUEST) {
		op.ch <- msg
	} else {
		logger.Errorf("REQUEST Operation (%s,%d) receives Bad message: %+v", *op.urito, op.tid, msg)
	}
}

// This function is called when underlying ClientContext is closed
func (op *RequestOperationX) onClose() {
	op.Close()
}

// ================================================================================
// Invoke Operation

type InvokeOperation interface {
	Operation
	Invoke(body Body) (*Message, error)
	GetResponse() (*Message, error)
}

type InvokeOperationX struct {
	OperationX
	// TODO (AF): Handling of response (see api1)
	response *Message
}

func (cctx *ClientContext) NewInvokeOperation(urito *URI, area UShort, areaVersion UOctet, service UShort, operation UShort) InvokeOperation {
	// Gets a new TransactionId for operation
	tid := cctx.TransactionId()
	// Normally should not receive more than 2 messages, and it waits first message (ack) synchronously.
	ch := make(chan *Message)
	op := &InvokeOperationX{OperationX: OperationX{cctx, tid, ch, urito, area, areaVersion, service, operation, _CREATED}}
	return op
}

func (op *InvokeOperationX) Invoke(body Body) (*Message, error) {
	if op.status != _CREATED {
		return nil, errors.New("Bad operation status")
	}
	op.status = _INITIATED

	msg := &Message{
		UriFrom:          op.cctx.Uri,
		UriTo:            op.urito,
		AuthenticationId: op.cctx.AuthenticationId,
		EncodingId:       op.cctx.EncodingId,
		QoSLevel:         op.cctx.QoSLevel,
		Priority:         op.cctx.Priority,
		Domain:           op.cctx.Domain,
		NetworkZone:      op.cctx.NetworkZone,
		Session:          op.cctx.Session,
		SessionName:      op.cctx.SessionName,
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
	err := op.cctx.registerOp(op.tid, op)
	if err != nil {
		op.finalize()
		return nil, err
	}
	// Send the INVOKE MAL message
	err = op.cctx.Ctx.Send(msg)
	if err != nil {
		op.finalize()
		return nil, err
	}

	// Waits for the INVOKE_ACK MAL message
	msg, more := <-op.ch
	if !more {
		op.finalize()
		logger.Debugf("InvokeOperation.Invoke: Operation ends: %s, %s", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation ends")
	}
	if msg == nil {
		op.finalize()
		logger.Warnf("InvokeOperation.Invoke: Operation interupted (%s, %s)", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation interupted")
	}
	// Verify the message stage
	if msg.InteractionStage != MAL_IP_STAGE_INVOKE_ACK {
		op.finalize()
		logger.Errorf("InvokeOperation.Invoke: Bad return message, operation (%s, %s), stage %d", op.cctx.Uri, op.tid, msg.InteractionStage)
		return nil, errors.New("Bad return message")
	}
	op.status = _ACKNOWLEDGED
	// Verify that the message is ok (ack or error)
	if msg.IsErrorMessage {
		op.finalize()
		return msg, errors.New("Error message")
	} else {
		return msg, nil
	}
}

// Returns the response.
func (op *InvokeOperationX) GetResponse() (*Message, error) {
	if (op.status == _FINAL) && (op.response != nil) {
		if op.response.IsErrorMessage {
			return op.response, errors.New("Error message")
		} else {
			return op.response, nil
		}
	}
	if op.status != _ACKNOWLEDGED {
		return nil, errors.New("Bad operation status")
	}

	// Waits for next MAL message
	msg, more := <-op.ch
	if !more {
		op.finalize()
		logger.Debugf("InvokeOperation.GetResponse: Operation ends: %s, %s", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation ends")
	}
	if msg == nil {
		op.finalize()
		logger.Warnf("InvokeOperation.GetResponse: Operation interupted (%s, %s)", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation interupted")
	}
	// Verify the message stage
	if msg.InteractionStage != MAL_IP_STAGE_INVOKE_RESPONSE {
		op.finalize()
		logger.Errorf("InvokeOperation.GetResponse: Bad return message, operation (%s, %s), stage %d", op.cctx.Uri, op.tid, msg.InteractionStage)
		return nil, errors.New("Bad return message")
	}
	op.finalize()
	op.response = msg
	// Verify that the message is ok (ack or error)
	if msg.IsErrorMessage {
		return msg, errors.New("Error message")
	} else {
		return msg, nil
	}
}

func (op *InvokeOperationX) onMessage(msg *Message) {
	// Verify the message: service area, version, service, operation
	if op.verify(msg) && (msg.InteractionType == MAL_INTERACTIONTYPE_INVOKE) {
		op.ch <- msg
	} else {
		logger.Errorf("INVOKE Operation (%s,%d) receives Bad message: %+v", *op.urito, op.tid, msg)
	}
}

// This function is called when underlying ClientContext is closed
func (op *InvokeOperationX) onClose() {
	op.Close()
}

// ================================================================================
// Progress Operation

type ProgressOperation interface {
	Operation
	Progress(body Body) (*Message, error)
	GetUpdate() (*Message, error)
	GetResponse() (*Message, error)
}

type ProgressOperationX struct {
	OperationX
	response *Message
}

func (cctx *ClientContext) NewProgressOperation(urito *URI, area UShort, areaVersion UOctet, service UShort, operation UShort) ProgressOperation {
	// Gets a new TransactionId for operation
	tid := cctx.TransactionId()
	// Be careful: Depending of the application logic this channel can receive an arbitrary number of Update messages.
	// The size of this channel must be large enough to buffer these messages in order to avoid blocking of the underlying
	// MAL context thread.
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &ProgressOperationX{OperationX: OperationX{cctx, tid, ch, urito, area, areaVersion, service, operation, _CREATED}}
	return op
}

func (op *ProgressOperationX) Progress(body Body) (*Message, error) {
	if op.status != _CREATED {
		return nil, errors.New("Bad operation status")
	}
	op.status = _INITIATED

	msg := &Message{
		UriFrom:          op.cctx.Uri,
		UriTo:            op.urito,
		AuthenticationId: op.cctx.AuthenticationId,
		EncodingId:       op.cctx.EncodingId,
		QoSLevel:         op.cctx.QoSLevel,
		Priority:         op.cctx.Priority,
		Domain:           op.cctx.Domain,
		NetworkZone:      op.cctx.NetworkZone,
		Session:          op.cctx.Session,
		SessionName:      op.cctx.SessionName,
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
	err := op.cctx.registerOp(op.tid, op)
	if err != nil {
		op.finalize()
		return nil, err
	}
	// Send the SUBMIT MAL message
	err = op.cctx.Ctx.Send(msg)
	if err != nil {
		op.finalize()
		return nil, err
	}

	// Waits for the PROGRESS_ACK MAL message
	msg, more := <-op.ch
	if !more {
		op.finalize()
		logger.Debugf("ProgressOperation.Progress: Operation ends: %s, %s", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation ends")
	}
	if msg == nil {
		op.finalize()
		logger.Warnf("ProgressOperation.Progress: Operation interupted (%s, %s)", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation interupted")
	}
	// Verify the message stage
	if msg.InteractionStage != MAL_IP_STAGE_PROGRESS_ACK {
		op.finalize()
		logger.Errorf("ProgressOperation.Progress: Bad return message, operation (%s, %s), stage %d", op.cctx.Uri, op.tid, msg.InteractionStage)
		return nil, errors.New("Bad return message")
	}
	op.status = _ACKNOWLEDGED
	// Verify that the message is ok (ack or error)
	if msg.IsErrorMessage {
		op.finalize()
		return msg, errors.New("Error message")
	} else {
		return msg, nil
	}
}

// Returns next update or nil if there is no more update.
func (op *ProgressOperationX) GetUpdate() (*Message, error) {
	if (op.status != _ACKNOWLEDGED) && (op.status != _PROGRESSING) {
		return nil, errors.New("Bad operation status")
	}

	// Waits for next MAL message
	msg, more := <-op.ch
	if !more {
		op.finalize()
		logger.Debugf("ProgressOperation.GetUpdate: Operation ends: %s, %s", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation ends")
	}
	if msg == nil {
		logger.Warnf("ProgressOperation.GetUpdate: Operation interupted (%s, %s)", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation interupted")
	}

	if (msg.InteractionStage != MAL_IP_STAGE_PROGRESS_UPDATE) &&
		(msg.InteractionStage != MAL_IP_STAGE_PROGRESS_RESPONSE) {
		op.finalize()
		logger.Errorf("ProgressOperation.GetUpdate: Bad return message, operation (%s, %s), stage %d", op.cctx.Uri, op.tid, msg.InteractionStage)
		return nil, errors.New("Bad return message")
	}

	if msg.InteractionStage == MAL_IP_STAGE_PROGRESS_UPDATE {
		op.status = _PROGRESSING
		// Verify that the message is ok (ack or error)
		if msg.IsErrorMessage {
			op.finalize()
			return msg, errors.New("Error message")
		} else {
			return msg, nil
		}
	}
	// msg.InteractionStage == MAL_IP_STAGE_PROGRESS_RESPONSE {
	op.response = msg
	op.finalize()
	return nil, nil
}

// Returns the response.
func (op *ProgressOperationX) GetResponse() (*Message, error) {
	if (op.status == _FINAL) && (op.response != nil) {
		if op.response.IsErrorMessage {
			return op.response, errors.New("Error message")
		} else {
			return op.response, nil
		}
	}
	if (op.status != _ACKNOWLEDGED) && (op.status != _PROGRESSING) {
		return nil, errors.New("Bad operation status")
	}

	// Waits for next MAL message
	msg, more := <-op.ch
	if !more {
		op.finalize()
		logger.Debugf("ProgressOperation.GetResponse: Operation ends: %s, %s", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation ends")
	}
	if msg == nil {
		logger.Warnf("ProgressOperation.GetResponse: Operation interupted (%s, %s)", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation interupted")
	}

	if msg.InteractionStage != MAL_IP_STAGE_PROGRESS_RESPONSE {
		op.finalize()
		logger.Errorf("ProgressOperation.GetResponse: Bad return message, operation (%s, %s), stage %d", op.cctx.Uri, op.tid, msg.InteractionStage)
		return nil, errors.New("Bad return message")
	}
	op.finalize()
	op.response = msg
	// Verify that the message is ok (ack or error)
	if msg.IsErrorMessage {
		return msg, errors.New("Error message")
	} else {
		return msg, nil
	}
}

func (op *ProgressOperationX) onMessage(msg *Message) {
	// Verify the message: service area, version, service, operation
	if op.verify(msg) && (msg.InteractionType == MAL_INTERACTIONTYPE_PROGRESS) {
		op.ch <- msg
	} else {
		logger.Errorf("PROGRESS Operation (%s,%d) receives Bad message: %+v", *op.urito, op.tid, msg)
	}
}

// This function is called when underlying ClientContext is closed
func (op *ProgressOperationX) onClose() {
	op.Close()
}

// ================================================================================
// Subscriber Operation

type SubscriberOperation interface {
	Operation
	Register(body Body) (*Message, error)
	GetNotify() (*Message, error)
	Deregister(body Body) (*Message, error)
}

type SubscriberOperationX struct {
	OperationX
}

func (cctx *ClientContext) NewSubscriberOperation(urito *URI, area UShort, areaVersion UOctet, service UShort, operation UShort) SubscriberOperation {
	// Gets a new TransactionId for operation
	tid := cctx.TransactionId()
	// Be careful: Depending of the application logic this channel can receive an arbitrary number of Notify messages.
	// The size of this channel must be large enough to buffer these messages in order to avoid blocking of the underlying
	// MAL context thread.
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &SubscriberOperationX{OperationX: OperationX{cctx, tid, ch, urito, area, areaVersion, service, operation, _CREATED}}
	return op
}

func (op *SubscriberOperationX) Register(body Body) (*Message, error) {
	// TODO (AF): Be careful we can register anew a Subscriber
	if op.status != _CREATED {
		return nil, errors.New("Bad operation status")
	}
	op.status = _REGISTER_INITIATED

	msg := &Message{
		UriFrom:          op.cctx.Uri,
		UriTo:            op.urito,
		AuthenticationId: op.cctx.AuthenticationId,
		EncodingId:       op.cctx.EncodingId,
		QoSLevel:         op.cctx.QoSLevel,
		Priority:         op.cctx.Priority,
		Domain:           op.cctx.Domain,
		NetworkZone:      op.cctx.NetworkZone,
		Session:          op.cctx.Session,
		SessionName:      op.cctx.SessionName,
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
	err := op.cctx.registerOp(op.tid, op)
	if err != nil {
		op.finalize()
		return nil, err
	}
	// Send the REGISTER MAL message
	err = op.cctx.Ctx.Send(msg)
	if err != nil {
		op.finalize()
		return nil, err
	}

	// Waits for the REGISTER_ACK MAL message
	msg, more := <-op.ch
	if !more {
		op.finalize()
		logger.Debugf("SubscriberOperation.Register: Operation ends: %s, %s", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation ends")
	}
	if msg == nil {
		op.finalize()
		logger.Warnf("SubscriberOperation.Register: Operation interupted (%s, %s)", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation interupted")
	}
	// Verify the message stage
	if msg.InteractionStage != MAL_IP_STAGE_PUBSUB_REGISTER_ACK {
		op.finalize()
		logger.Errorf("SubscriberOperation.Register: Bad return message, operation (%s, %s), stage %d", op.cctx.Uri, op.tid, msg.InteractionStage)
		return nil, errors.New("Bad return message")
	}
	// Verify that the message is ok (ack or error)
	if msg.IsErrorMessage {
		op.finalize()
		return msg, errors.New("Error message")
	} else {
		op.status = _REGISTERED
		return msg, nil
	}
}

// Returns next notify.
func (op *SubscriberOperationX) GetNotify() (*Message, error) {
	// TODO (AF): May be we have to allow a timeout to this operation.
	if (op.status != _REGISTERED) && (op.status != _REREGISTER_INITIATED) && (op.status != _DEREGISTER_INITIATED) {
		return nil, errors.New("Bad operation status")
	}
	// TODO (AF): Handle _REREGISTER_INITIATED and _DEREGISTER_INITIATED status

	// Waits for next MAL message
	msg, more := <-op.ch
	if !more {
		op.finalize()
		logger.Debugf("SubscriberOperation.GetNotify: Operation ends: %s, %s", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation ends")
	}
	if msg == nil {
		logger.Warnf("SubscriberOperation.GetNotify: Operation interupted (%s, %s)", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation interupted")
	}
	// Verify the message stage
	if msg.InteractionStage != MAL_IP_STAGE_PUBSUB_NOTIFY {
		op.finalize()
		logger.Errorf("SubscriberOperation.GetNotify: Bad return message, operation (%s, %s), stage %d", op.cctx.Uri, op.tid, msg.InteractionStage)
		return nil, errors.New("Bad return message")
	}
	// Verify that the message is ok (ack or error)
	if msg.IsErrorMessage {
		op.finalize()
		return msg, errors.New("SubscriberOperation.GetNotify: Error message")
	} else {
		return msg, nil
	}
}

func (op *SubscriberOperationX) Deregister(body Body) (*Message, error) {
	if op.status != _REGISTERED {
		return nil, errors.New("Bad operation status")
	}
	op.status = _DEREGISTER_INITIATED

	msg := &Message{
		UriFrom:          op.cctx.Uri,
		UriTo:            op.urito,
		AuthenticationId: op.cctx.AuthenticationId,
		EncodingId:       op.cctx.EncodingId,
		QoSLevel:         op.cctx.QoSLevel,
		Priority:         op.cctx.Priority,
		Domain:           op.cctx.Domain,
		NetworkZone:      op.cctx.NetworkZone,
		Session:          op.cctx.Session,
		SessionName:      op.cctx.SessionName,
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
	err := op.cctx.Ctx.Send(msg)
	if err != nil {
		op.finalize()
		return nil, err
	}

	// Waits for the DEREGISTER_ACK MAL message, removing useless notify waiting messages
	for {
		msg, more := <-op.ch
		if !more {
			op.finalize()
			logger.Debugf("SubscriberOperation.Deregister: Operation ends: %s, %s", op.cctx.Uri, op.tid)
			return nil, errors.New("Operation ends")
		}
		if msg == nil {
			op.finalize()
			logger.Warnf("SubscriberOperation.Deregister: Operation interupted (%s, %s)", op.cctx.Uri, op.tid)
			return nil, errors.New("Operation interupted")
		}
		if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_NOTIFY {
			continue
		}
		// Verify the message stage
		if msg.InteractionStage != MAL_IP_STAGE_PUBSUB_DEREGISTER_ACK {
			op.finalize()
			logger.Errorf("SubscriberOperation.Deregister: Bad return message, operation (%s, %s), stage %d", op.cctx.Uri, op.tid, msg.InteractionStage)
			return nil, errors.New("Bad return message")
		}
		op.finalize()
		if msg.IsErrorMessage {
			return msg, errors.New("Error message")
		} else {
			return msg, nil
		}
	}
}

func (op *SubscriberOperationX) onMessage(msg *Message) {
	// Verify the message: service area, version, service, operation
	if op.verify(msg) && (msg.InteractionType == MAL_INTERACTIONTYPE_PUBSUB) {
		op.ch <- msg
	} else {
		logger.Errorf("PUBSUB Operation (%s,%d) receives Bad message: %+v", *op.urito, op.tid, msg)
	}
}

// This function is called when underlying ClientContext is closed
func (op *SubscriberOperationX) onClose() {
	op.Close()
}

// ================================================================================
// Publisher Operation

type PublisherOperation interface {
	Operation
	Register(body Body) (*Message, error)
	Publish(body Body) error
	GetPublishError() (*Message, error)
	Deregister(body Body) (*Message, error)
}

type PublisherOperationX struct {
	OperationX
	pe_ch chan *Message
}

func (cctx *ClientContext) NewPublisherOperation(urito *URI, area UShort, areaVersion UOctet, service UShort, operation UShort) PublisherOperation {
	// Gets a new TransactionId for operation
	tid := cctx.TransactionId()
	// In normal operation should not receive more than one message (Acknowledge of Register and
	// Deregister messages).
	ch := make(chan *Message)
	// Make channel for PublishError
	pe_ch := make(chan *Message, 5)
	op := &PublisherOperationX{OperationX: OperationX{cctx, tid, ch, urito, area, areaVersion, service, operation, _CREATED}, pe_ch: pe_ch}
	return op
}

func (op *PublisherOperationX) Register(body Body) (*Message, error) {
	// TODO (AF): Be careful currently we cannot register anew a publisher
	if op.status != _CREATED {
		return nil, errors.New("Bad operation status")
	}
	op.status = _REGISTER_INITIATED

	msg := &Message{
		UriFrom:          op.cctx.Uri,
		UriTo:            op.urito,
		AuthenticationId: op.cctx.AuthenticationId,
		EncodingId:       op.cctx.EncodingId,
		QoSLevel:         op.cctx.QoSLevel,
		Priority:         op.cctx.Priority,
		Domain:           op.cctx.Domain,
		NetworkZone:      op.cctx.NetworkZone,
		Session:          op.cctx.Session,
		SessionName:      op.cctx.SessionName,
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
	err := op.cctx.registerOp(op.tid, op)
	if err != nil {
		op.finalize()
		return nil, err
	}
	// Send the PUBLISH_REGISTER MAL message
	err = op.cctx.Ctx.Send(msg)
	if err != nil {
		op.finalize()
		return nil, err
	}
	// Waits for the PUBLISH_REGISTER_ACK MAL message
	msg, more := <-op.ch
	if !more {
		op.finalize()
		logger.Debugf("PublisherOperation.Register: Operation ends: %s, %s", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation ends")
	}
	if msg == nil {
		op.finalize()
		logger.Warnf("PublisherOperation.Register: Operation interupted (%s, %s)", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation interupted")
	}
	// Verify the message stage
	if msg.InteractionStage != MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER_ACK {
		op.finalize()
		logger.Errorf("PublisherOperation.Register: Bad return message, operation (%s, %s), stage %d", op.cctx.Uri, op.tid, msg.InteractionStage)
		return nil, errors.New("Bad return message")
	}
	// Verify that the message is ok (ack or error)
	if msg.IsErrorMessage {
		op.finalize()
		return msg, errors.New("Error message")
	} else {
		op.status = _REGISTERED
		return msg, nil
	}
}

func (op *PublisherOperationX) Publish(body Body) error {
	if op.status != _REGISTERED {
		return errors.New("Bad operation status")
	}

	msg := &Message{
		UriFrom:          op.cctx.Uri,
		UriTo:            op.urito,
		AuthenticationId: op.cctx.AuthenticationId,
		EncodingId:       op.cctx.EncodingId,
		QoSLevel:         op.cctx.QoSLevel,
		Priority:         op.cctx.Priority,
		Domain:           op.cctx.Domain,
		NetworkZone:      op.cctx.NetworkZone,
		Session:          op.cctx.Session,
		SessionName:      op.cctx.SessionName,
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
	err := op.cctx.Ctx.Send(msg)
	if err != nil {
		op.finalize()
		return err
	}

	return nil
}

func (op *PublisherOperationX) GetPublishError() (*Message, error) {
	select {
	case msg, ok := <-op.pe_ch:
		if ok {
			return msg, nil
		} else {
			return nil, errors.New("Operation ends")
		}
	default:
		return nil, nil
	}
}

func (op *PublisherOperationX) Deregister(body Body) (*Message, error) {
	if op.status != _REGISTERED {
		return nil, errors.New("Bad operation status")
	}
	op.status = _DEREGISTER_INITIATED

	msg := &Message{
		UriFrom:          op.cctx.Uri,
		UriTo:            op.urito,
		AuthenticationId: op.cctx.AuthenticationId,
		EncodingId:       op.cctx.EncodingId,
		QoSLevel:         op.cctx.QoSLevel,
		Priority:         op.cctx.Priority,
		Domain:           op.cctx.Domain,
		NetworkZone:      op.cctx.NetworkZone,
		Session:          op.cctx.Session,
		SessionName:      op.cctx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER,
		ServiceArea:      op.area,
		AreaVersion:      op.areaVersion,
		Service:          op.service,
		Operation:        op.operation,
		TransactionId:    op.tid,
		Body:             body,
	}
	// Send the PUBLISH_DEREGISTER MAL message
	err := op.cctx.Ctx.Send(msg)
	if err != nil {
		op.finalize()
		return nil, err
	}

	// Waits for the PUBLISH_DEREGISTER_ACK MAL message
	msg, more := <-op.ch
	if !more {
		op.finalize()
		logger.Debugf("PublisherOperation.Deregister: Operation ends: %s, %s", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation ends")
	}
	if msg == nil {
		op.finalize()
		logger.Warnf("PublisherOperation.Deregister: Operation interupted (%s, %s)", op.cctx.Uri, op.tid)
		return nil, errors.New("Operation interupted")
	}
	// Verify the message stage
	if msg.InteractionStage != MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER_ACK {
		op.finalize()
		logger.Errorf("PublisherOperation.Deregister: Bad return message, operation (%s, %s), stage %d", op.cctx.Uri, op.tid, msg.InteractionStage)
		return nil, errors.New("Bad return message")
	}
	op.finalize()
	if msg.IsErrorMessage {
		return msg, errors.New("Error message")
	} else {
		return msg, nil
	}
}

func (op *PublisherOperationX) onMessage(msg *Message) {
	// Verify the message: service area, version, service, operation
	if op.verify(msg) && (msg.InteractionType == MAL_INTERACTIONTYPE_PUBSUB) {
		if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH && msg.IsErrorMessage {
			// It is a PUBLISH_ERROR, keep it ans continue
			op.pe_ch <- msg
		} else {
			op.ch <- msg
		}
	} else {
		logger.Errorf("PUBSUB Operation (%s,%d) receives Bad message: %+v", *op.urito, op.tid, msg)
	}
}

// This function is called when underlying ClientContext is closed
func (op *PublisherOperationX) onClose() {
	op.Close()
}
