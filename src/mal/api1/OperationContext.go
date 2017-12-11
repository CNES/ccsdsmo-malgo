/**
 * MIT License
 *
 * Copyright (c) 2017 CNES
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
package api1

import (
	"errors"
	"fmt"
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
	_FINAL
	// TODO (AF): Should define PubSub status
)

type Operation struct {
	ictx        *OperationContext
	tid         ULong
	ch          chan *Message
	area        UShort // TODO (AF): Should be in OperationContext or Context?
	areaVersion UOctet // TODO (AF): Should be in OperationContext or Context?
	service     UShort // TODO (AF): Should be in OperationContext or Context?
	operation   UShort
	status      byte
}

// ================================================================================
// SendOperation

type SendOperation struct {
	Operation
}

func (ictx *OperationContext) NewSendOperation(area UShort, areaVersion UOctet, service UShort, operation UShort) (*SendOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	op := &SendOperation{Operation: Operation{ictx, tid, nil, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *SendOperation) Send(msg *Message) error {
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

func (op *SendOperation) OnMessage(msg *Message) error {
	// TODO (AF): Should never reveive messages, log an error
	return nil
}

func (op *SendOperation) OnClose() error {
	// TODO (AF): Should never be called, log an error
	return nil
}

// ================================================================================
// SubmitOperation

type SubmitOperation struct {
	Operation
}

func (ictx *OperationContext) NewSubmitOperation(area UShort, areaVersion UOctet, service UShort, operation UShort) (*SubmitOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &SubmitOperation{Operation: Operation{ictx, tid, ch, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *SubmitOperation) Submit(msg *Message) (*Message, error) {
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
		fmt.Println("Operation ends: ", op.ictx.Uri, op.tid)
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

func (op *SubmitOperation) OnMessage(msg *Message) error {
	// TODO (AF): Verify the message: service area, version, service, operation
	if msg.InteractionType != MAL_INTERACTIONTYPE_SUBMIT {
		// TODO (AF): log an error
		return errors.New("Bad message")
	}
	op.ch <- msg
	return nil
}

func (op *SubmitOperation) OnClose() error {
	// TODO (AF):
	return nil
}

// ================================================================================
// Progress Operation

type ProgressOperation struct {
	Operation
	response *Message
}

func (ictx *OperationContext) NewProgressOperation(area UShort, areaVersion UOctet, service UShort, operation UShort) (*ProgressOperation, error) {
	// Gets a new TransactionId for operation
	tid := ictx.TransactionId()
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	op := &ProgressOperation{Operation: Operation{ictx, tid, ch, area, areaVersion, service, operation, _CREATED}}
	return op, nil
}

func (op *ProgressOperation) Progress(msg *Message) (*Message, error) {
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
		fmt.Println("Operation ends: ", op.ictx.Uri, op.tid)
		// TODO (AF): Returns an error
		op.status = _FINAL
		return nil, errors.New("Operation ends")
	}
	op.status = _ACKNOWLEDGED

	// TODO (AF): Verify that the message is ok (ack or error)
	if msg.InteractionStage != MAL_IP_STAGE_SUBMIT_ACK {
		// TODO (AF): Close the operation
		return nil, errors.New("Bad operation stage")
	}

	return msg, nil
}

// Returns next update or nil if there is no more update.
func (op *ProgressOperation) GetUpdate() (*Message, error) {
	if (op.status != _ACKNOWLEDGED) && (op.status != _PROGRESSING) {
		return nil, errors.New("Bad operation status")
	}

	// Waits for next MAL message
	msg, more := <-op.ch
	if !more {
		fmt.Println("Operation ends: ", op.ictx.Uri, op.tid)
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
func (op *ProgressOperation) GetResponse() (*Message, error) {
	if (op.status == _FINAL) && (op.response != nil) {
		return op.response, nil
	}

	if (op.status != _ACKNOWLEDGED) && (op.status != _PROGRESSING) {
		return nil, errors.New("Bad operation status")
	}

	// Waits for next MAL message
	msg, more := <-op.ch
	if !more {
		fmt.Println("Operation ends: ", op.ictx.Uri, op.tid)
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

func (op *ProgressOperation) OnMessage(msg *Message) error {
	// TODO (AF): Verify the message: service area, version, service, operation
	if msg.InteractionType != MAL_INTERACTIONTYPE_PROGRESS {
		// TODO (AF): log an error
		return errors.New("Bad message")
	}
	op.ch <- msg
	return nil
}

func (op *ProgressOperation) OnClose() error {
	// TODO (AF):
	return nil
}

// ================================================================================
// Defines Listener interface used by context to route MAL messages

func (ictx *OperationContext) OnMessage(msg *Message) error {
	to, ok := ictx.handlers[msg.TransactionId]
	if ok {
		fmt.Printf("%t\n", to)
		to.OnMessage(msg)
		fmt.Println("Message transmitted: ", msg)
	} else {
		fmt.Println("Cannot route message to: ", msg.UriTo, "?TransactionId=", msg.TransactionId)
	}
	return nil
}

func (ictx *OperationContext) OnClose() error {
	fmt.Println("close EndPoint: ", ictx.Uri)
	for tid, handler := range ictx.handlers {
		fmt.Println("close operation: ", tid)
		err := handler.OnClose()
		if err != nil {
			// TODO (AF): print an error message
		}
	}
	return nil
}
