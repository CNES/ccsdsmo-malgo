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
package api1

// Note (AF): This API is now deprecated, a merged API allows to use either handler
// based and interface based interfaces in a same way.

import (
	"errors"
	"fmt"
	. "mal"
)

const (
	_SEND_HANDLER UOctet = iota
	_SUBMIT_HANDLER
	_REQUEST_HANDLER
	_INVOKE_HANDLER
	_PROGRESS_HANDLER
	_BROKER_HANDLER
)

type handlerDesc struct {
	handlerType UOctet
	area        UShort
	areaVersion UOctet
	service     UShort
	operation   UShort
	handler     Handler
}

type ProviderContext struct {
	Ctx      *Context
	Uri      *URI
	ch       chan *Message
	handlers map[uint64](*handlerDesc)
}

func NewProviderContext(ctx *Context, service string) (*ProviderContext, error) {
	// TODO (AF): Verify the uri
	uri := ctx.NewURI(service)
	// TODO (AF): Fix length of channel?
	ch := make(chan *Message, 10)
	handlers := make(map[uint64](*handlerDesc))
	pctx := &ProviderContext{ctx, uri, ch, handlers}
	err := ctx.RegisterEndPoint(uri, pctx)
	if err != nil {
		return nil, err
	}
	return pctx, nil
}

func key(area UShort, areaVersion UOctet, service UShort, operation UShort) uint64 {
	key := uint64(area) << 8
	key |= uint64(areaVersion)
	key <<= 16
	key |= uint64(service)
	key <<= 16
	key |= uint64(operation)

	return key
}

func (pctx *ProviderContext) register(hdltype UOctet, area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	key := key(area, areaVersion, service, operation)
	old := pctx.handlers[key]

	if old != nil {
		fmt.Println("MAL handler already registered:", key)
		return errors.New("MAL handler already registered")
	} else {
		fmt.Println("MAL handler registered:", key)
	}

	var desc = &handlerDesc{
		handlerType: hdltype,
		area:        area,
		areaVersion: areaVersion,
		service:     service,
		operation:   operation,
		handler:     handler,
	}

	pctx.handlers[key] = desc
	return nil
}

func (pctx *ProviderContext) Close() error {
	return pctx.Ctx.UnregisterEndPoint(pctx.Uri)
}

// Defines a generic root Transaction interface
type Transaction interface {
	getURI() *URI
	getTid() ULong
}

// Defines a generic root Transaction structure
type TransactionX struct {
	pctx        *ProviderContext
	urifrom     *URI
	tid         ULong
	area        UShort
	areaVersion UOctet
	service     UShort
	operation   UShort
}

func (tx *TransactionX) getURI() *URI {
	return tx.urifrom
}

func (tx *TransactionX) getTid() ULong {
	return tx.tid
}

// Defines a generic root handler interface
type Handler func(*Message, Transaction) error

// ================================================================================
// SendHandler

type SendTransaction interface {
	Transaction
}

type SendTransactionX struct {
	TransactionX
}

type SendHandler func(*Message, SendTransaction) error

// TODO (AF):
//func (pctx *ProviderContext) RegisterSendHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler SendHandler) error {
func (pctx *ProviderContext) RegisterSendHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return pctx.register(_SEND_HANDLER, area, areaVersion, service, operation, handler)
}

// ================================================================================
// SubmitHandler

type SubmitTransaction interface {
	Transaction
	Ack(err error) error
}

type SubmitTransactionX struct {
	TransactionX
}

func (tx *SubmitTransactionX) Ack(err error) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_SUBMIT,
		InteractionStage: MAL_IP_STAGE_SUBMIT_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.pctx.Uri,
		UriTo:            tx.urifrom,
	}
	if err != nil {
		msg.IsErrorMessage = true
		msg.Body = []byte(err.Error())
	}
	return tx.pctx.Ctx.Send(msg)
}

type SubmitHandler func(*Message, SubmitTransaction) error

// TODO (AF):
//func (pctx *ProviderContext) RegisterSendHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler SendHandler) error {
func (pctx *ProviderContext) RegisterSubmitHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return pctx.register(_SUBMIT_HANDLER, area, areaVersion, service, operation, handler)
}

// ================================================================================
// RequestHandler

type RequestTransaction interface {
	Reply([]byte, error) error
}

type RequestTransactionX struct {
	TransactionX
}

func (tx *RequestTransactionX) Reply(body []byte, err error) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_REQUEST,
		InteractionStage: MAL_IP_STAGE_REQUEST_RESPONSE,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.pctx.Uri,
		UriTo:            tx.urifrom,
	}
	if err != nil {
		msg.IsErrorMessage = true
		msg.Body = []byte(err.Error())
	} else {
		msg.Body = body
	}
	return tx.pctx.Ctx.Send(msg)
}

type RequestHandler func(*Message, RequestTransaction) error

// TODO (AF):
//func (pctx *ProviderContext) RegisterRequestHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler RequestHandler) error {
func (pctx *ProviderContext) RegisterRequestHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return pctx.register(_REQUEST_HANDLER, area, areaVersion, service, operation, handler)
}

// ================================================================================
// InvokeHandler

type InvokeTransaction interface {
	Ack(error) error
	Reply([]byte, error) error
}

type InvokeTransactionX struct {
	TransactionX
}

func (tx *InvokeTransactionX) Ack(err error) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_INVOKE,
		InteractionStage: MAL_IP_STAGE_INVOKE_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.pctx.Uri,
		UriTo:            tx.urifrom,
	}
	if err != nil {
		msg.IsErrorMessage = true
		msg.Body = []byte(err.Error())
	}
	return tx.pctx.Ctx.Send(msg)
}

func (tx *InvokeTransactionX) Reply(body []byte, err error) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_INVOKE,
		InteractionStage: MAL_IP_STAGE_INVOKE_RESPONSE,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.pctx.Uri,
		UriTo:            tx.urifrom,
	}
	if err != nil {
		msg.IsErrorMessage = true
		msg.Body = []byte(err.Error())
	} else {
		msg.Body = body
	}
	return tx.pctx.Ctx.Send(msg)
}

type InvokeHandler func(*Message, InvokeTransaction) error

// TODO (AF):
//func (pctx *ProviderContext) RegisterInvokeHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler InvokeHandler) error {
func (pctx *ProviderContext) RegisterInvokeHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return pctx.register(_INVOKE_HANDLER, area, areaVersion, service, operation, handler)
}

// ================================================================================
// ProgressHandler

type ProgressTransaction interface {
	Ack(error) error
	Update([]byte, error) error
	Reply([]byte, error) error
}

type ProgressTransactionX struct {
	TransactionX
}

func (tx *ProgressTransactionX) Ack(err error) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PROGRESS,
		InteractionStage: MAL_IP_STAGE_PROGRESS_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.pctx.Uri,
		UriTo:            tx.urifrom,
	}
	if err != nil {
		msg.IsErrorMessage = true
		msg.Body = []byte(err.Error())
	}
	return tx.pctx.Ctx.Send(msg)
}

func (tx *ProgressTransactionX) Update(body []byte, err error) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PROGRESS,
		InteractionStage: MAL_IP_STAGE_PROGRESS_UPDATE,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.pctx.Uri,
		UriTo:            tx.urifrom,
	}
	if err != nil {
		msg.IsErrorMessage = true
		msg.Body = []byte(err.Error())
	} else {
		msg.Body = body
	}
	return tx.pctx.Ctx.Send(msg)
}

func (tx *ProgressTransactionX) Reply(body []byte, err error) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PROGRESS,
		InteractionStage: MAL_IP_STAGE_PROGRESS_RESPONSE,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.pctx.Uri,
		UriTo:            tx.urifrom,
	}
	if err != nil {
		msg.IsErrorMessage = true
		msg.Body = []byte(err.Error())
	} else {
		msg.Body = body
	}
	return tx.pctx.Ctx.Send(msg)
}

type ProgressHandler func(*Message, ProgressTransaction) error

// TODO (AF):
//func (pctx *ProviderContext) RegisterSendHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler SendHandler) error {
func (pctx *ProviderContext) RegisterProgressHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return pctx.register(_PROGRESS_HANDLER, area, areaVersion, service, operation, handler)
}

// ================================================================================
// BrokerHandler: There is only one handler but 2 transactions type depending of the
// incoming interaction.

type BrokerTransaction interface {
	Transaction
	AckRegister(error) error
	AckDeregister(error) error
}

type BrokerHandler func(*Message, BrokerTransaction) error

// TODO (AF):
//func (pctx *ProviderContext) RegisterBrokerHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler BrokerHandler) error {
func (pctx *ProviderContext) RegisterBrokerHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return pctx.register(_BROKER_HANDLER, area, areaVersion, service, operation, handler)
}

// SubscriberTransaction

type SubscriberTransaction interface {
	BrokerTransaction
	Notify([]byte, error) error
}

type SubscriberTransactionX struct {
	TransactionX
}

func (tx *SubscriberTransactionX) AckRegister(err error) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_REGISTER_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.pctx.Uri,
		UriTo:            tx.urifrom,
	}
	if err != nil {
		msg.IsErrorMessage = true
		msg.Body = []byte(err.Error())
	}
	return tx.pctx.Ctx.Send(msg)
}

func (tx *SubscriberTransactionX) Notify(body []byte, err error) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_NOTIFY,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.pctx.Uri,
		UriTo:            tx.urifrom,
	}
	if err != nil {
		msg.IsErrorMessage = true
		msg.Body = []byte(err.Error())
	} else {
		msg.Body = body
	}
	return tx.pctx.Ctx.Send(msg)
}

func (tx *SubscriberTransactionX) AckDeregister(err error) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_DEREGISTER_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.pctx.Uri,
		UriTo:            tx.urifrom,
	}
	if err != nil {
		msg.IsErrorMessage = true
		msg.Body = []byte(err.Error())
	}
	return tx.pctx.Ctx.Send(msg)
}

// PublisherTransaction

type PublisherTransaction interface {
	BrokerTransaction
}

type PublisherTransactionX struct {
	TransactionX
}

func (tx *PublisherTransactionX) AckRegister(err error) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.pctx.Uri,
		UriTo:            tx.urifrom,
	}
	if err != nil {
		msg.IsErrorMessage = true
		msg.Body = []byte(err.Error())
	}
	return tx.pctx.Ctx.Send(msg)
}

func (tx *PublisherTransactionX) AckDeregister(err error) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.pctx.Uri,
		UriTo:            tx.urifrom,
	}
	if err != nil {
		msg.IsErrorMessage = true
		msg.Body = []byte(err.Error())
	}
	return tx.pctx.Ctx.Send(msg)
}

// ================================================================================
// Defines Listener interface used by context to route MAL messages

func (pctx *ProviderContext) GetHandler(hdltype UOctet, area UShort, areaVersion UOctet, service UShort, operation UShort) (Handler, error) {
	key := key(area, areaVersion, service, operation)

	to, ok := pctx.handlers[key]
	if ok {
		if to.handlerType == hdltype {
			return to.handler, nil
		} else {
			fmt.Println("Bad handler type:", to.handlerType, " should be ", hdltype)
			return nil, errors.New("Bad handler type")
		}
	} else {
		fmt.Println("MAL handler not registered:", key)
		return nil, errors.New("MAL handler not registered")
	}

	// TODO (AF):
	//	if desc.handlerType != hdltype {
	//		return nil, errors.New("Bad Type for registered MAL handler")
	//	}
	//
	//	return desc.handlerPtr, nil
}

func (pctx *ProviderContext) OnMessage(msg *Message) error {
	switch msg.InteractionType {
	case MAL_INTERACTIONTYPE_SEND:
		handler, err := pctx.GetHandler(_SEND_HANDLER, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			return err
		}
		transaction := &SendTransactionX{TransactionX{pctx, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		// TODO (AF): use a goroutine
		return handler(msg, transaction)
	case MAL_INTERACTIONTYPE_SUBMIT:
		handler, err := pctx.GetHandler(_SUBMIT_HANDLER, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			return err
		}
		transaction := &SubmitTransactionX{TransactionX{pctx, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		// TODO (AF): use a goroutine
		return handler(msg, transaction)
	case MAL_INTERACTIONTYPE_REQUEST:
		handler, err := pctx.GetHandler(_REQUEST_HANDLER, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			return err
		}
		transaction := &RequestTransactionX{TransactionX{pctx, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		// TODO (AF): use a goroutine
		return handler(msg, transaction)
	case MAL_INTERACTIONTYPE_INVOKE:
		handler, err := pctx.GetHandler(_INVOKE_HANDLER, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			return err
		}
		transaction := &InvokeTransactionX{TransactionX{pctx, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		// TODO (AF): use a goroutine
		return handler(msg, transaction)
	case MAL_INTERACTIONTYPE_PROGRESS:
		handler, err := pctx.GetHandler(_PROGRESS_HANDLER, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			return err
		}
		transaction := &ProgressTransactionX{TransactionX{pctx, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		// TODO (AF): use a goroutine
		return handler(msg, transaction)
	case MAL_INTERACTIONTYPE_PUBSUB:
		handler, err := pctx.GetHandler(_BROKER_HANDLER, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			return err
		}
		var transaction Transaction
		if (msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER) ||
			(msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH) ||
			(msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER) {
			transaction = &PublisherTransactionX{TransactionX{pctx, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		} else if (msg.InteractionStage == MAL_IP_STAGE_PUBSUB_REGISTER) ||
			(msg.InteractionStage == MAL_IP_STAGE_PUBSUB_DEREGISTER) {
			transaction = &SubscriberTransactionX{TransactionX{pctx, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		} else {
			// TODO (AF): Log an error, May be wa should not return this error
			return errors.New("Bad interaction stage for PubSub")
		}
		// TODO (AF): use a goroutine
		return handler(msg, transaction)
	default:
		fmt.Println("Cannot route message to: ", *msg.UriTo)
	}
	//	// TODO (AF): calculate the key from the message content (see Routing).
	//	key := key(msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
	//	to, ok := pctx.handlers[key]
	//	if ok {
	//		fmt.Printf("%t\n", to)
	//		// TODO (AF): Route message
	//		to.OnMessage(msg)
	//		fmt.Println("Message transmitted: ", msg)
	//	} else {
	//		fmt.Println("Cannot route message to: ", *msg.UriTo, "?TransactionId=", msg.TransactionId, "?key=", key)
	//	}
	return nil
}

func (pctx *ProviderContext) OnClose() error {
	fmt.Println("close EndPoint: ", pctx.Uri)
	// TODO (AF): Close handlers ?
	//	for key, handler := range pctx.handlers {
	//		fmt.Println("close handler: ", key)
	//		err := handler.OnClose()
	//		if err != nil {
	//			// TODO (AF): print an error message
	//		}
	//	}
	return nil
}
