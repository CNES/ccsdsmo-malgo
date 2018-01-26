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
)

type handlerDesc struct {
	handlerType InteractionType
	area        UShort
	areaVersion UOctet
	service     UShort
	operation   UShort
	handler     Handler
}

type HandlerContext struct {
	Ctx *Context
	Uri *URI
	handlers map[uint64](*handlerDesc)
}

func NewHandlerContext(ctx *Context, service string) (*HandlerContext, error) {
	// TODO (AF): Verify the uri
	uri := ctx.NewURI(service)
	handlers := make(map[uint64](*handlerDesc))
	hctx := &HandlerContext{ctx, uri, handlers}
	err := ctx.RegisterEndPoint(uri, hctx)
	if err != nil {
		return nil, err
	}
	return hctx, nil
}

func (hctx *HandlerContext) register(hdltype InteractionType, area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	key := key(area, areaVersion, service, operation)
	old := hctx.handlers[key]

	if old != nil {
		logger.Errorf("MAL handler already registered: %d", key)
		return errors.New("MAL handler already registered")
	} else {
		logger.Debugf("MAL handler registered: %d", key)
	}

	var desc = &handlerDesc{
		handlerType: hdltype,
		area:        area,
		areaVersion: areaVersion,
		service:     service,
		operation:   operation,
		handler:     handler,
	}

	hctx.handlers[key] = desc
	return nil
}

func (hctx *HandlerContext) Close() error {
	return hctx.Ctx.UnregisterEndPoint(hctx.Uri)
}

// Defines a generic root handler interface
type Handler func(*Message, Transaction) error

// ================================================================================
// SendHandler

// TODO (AF):
//type SendHandler func(*Message, SendTransaction) error

//func (hctx *ProviderContext) RegisterSendHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler SendHandler) error {
func (hctx *HandlerContext) RegisterSendHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return hctx.register(MAL_INTERACTIONTYPE_SEND, area, areaVersion, service, operation, handler)
}

// ================================================================================
// SubmitHandler

// TODO (AF):
//type SubmitHandler func(*Message, SubmitTransaction) error

//func (hctx *ProviderContext) RegisterSendHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler SendHandler) error {
func (hctx *HandlerContext) RegisterSubmitHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return hctx.register(MAL_INTERACTIONTYPE_SUBMIT, area, areaVersion, service, operation, handler)
}

// ================================================================================
// RequestHandler

// TODO (AF):
//type RequestHandler func(*Message, RequestTransaction) error

//func (hctx *ProviderContext) RegisterRequestHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler RequestHandler) error {
func (hctx *HandlerContext) RegisterRequestHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return hctx.register(MAL_INTERACTIONTYPE_REQUEST, area, areaVersion, service, operation, handler)
}

// ================================================================================
// InvokeHandler

// TODO (AF):
//type InvokeHandler func(*Message, InvokeTransaction) error

//func (hctx *ProviderContext) RegisterInvokeHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler InvokeHandler) error {
func (hctx *HandlerContext) RegisterInvokeHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return hctx.register(MAL_INTERACTIONTYPE_INVOKE, area, areaVersion, service, operation, handler)
}

// ================================================================================
// ProgressHandler

// TODO (AF):
//type ProgressHandler func(*Message, ProgressTransaction) error

//func (hctx *ProviderContext) RegisterSendHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler SendHandler) error {
func (hctx *HandlerContext) RegisterProgressHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return hctx.register(MAL_INTERACTIONTYPE_PROGRESS, area, areaVersion, service, operation, handler)
}

// ================================================================================
// BrokerHandler: There is only one handler but 2 transactions type depending of the
// incoming interaction.

// TODO (AF):
//type BrokerHandler func(*Message, BrokerTransaction) error

//func (hctx *ProviderContext) RegisterBrokerHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler BrokerHandler) error {
func (hctx *HandlerContext) RegisterBrokerHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return hctx.register(MAL_INTERACTIONTYPE_PUBSUB, area, areaVersion, service, operation, handler)
}

// ================================================================================
// Defines Listener interface used by context to route MAL messages

func (hctx *HandlerContext) getHandler(hdltype InteractionType, area UShort, areaVersion UOctet, service UShort, operation UShort) (Handler, error) {
	key := key(area, areaVersion, service, operation)

	to, ok := hctx.handlers[key]
	if ok {
		if to.handlerType == hdltype {
			return to.handler, nil
		} else {
			logger.Errorf("Bad handler type: %d should be %d", to.handlerType, hdltype)
			return nil, errors.New("Bad handler type")
		}
	} else {
		logger.Errorf("MAL handler not registered: %d", key)
		return nil, errors.New("MAL handler not registered")
	}
}

func (hctx *HandlerContext) OnMessage(msg *Message) error {
	handler, err := hctx.getHandler(msg.InteractionType, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
	if err != nil {
		return err
	}
	var transaction Transaction
	switch msg.InteractionType {
	case MAL_INTERACTIONTYPE_SEND:
		transaction = &SendTransactionX{TransactionX{hctx.Ctx, hctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
	case MAL_INTERACTIONTYPE_SUBMIT:
		transaction = &SubmitTransactionX{TransactionX{hctx.Ctx, hctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
	case MAL_INTERACTIONTYPE_REQUEST:
		transaction = &RequestTransactionX{TransactionX{hctx.Ctx, hctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
	case MAL_INTERACTIONTYPE_INVOKE:
		transaction = &InvokeTransactionX{TransactionX{hctx.Ctx, hctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
	case MAL_INTERACTIONTYPE_PROGRESS:
		transaction = &ProgressTransactionX{TransactionX{hctx.Ctx, hctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
	case MAL_INTERACTIONTYPE_PUBSUB:
		if (msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER) ||
			(msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH) ||
			(msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER) {
			transaction = &PublisherTransactionX{TransactionX{hctx.Ctx, hctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		} else if (msg.InteractionStage == MAL_IP_STAGE_PUBSUB_REGISTER) ||
			(msg.InteractionStage == MAL_IP_STAGE_PUBSUB_DEREGISTER) {
			transaction = &SubscriberTransactionX{TransactionX{hctx.Ctx, hctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		} else {
			// TODO (AF): Log an error, May be we should not return this error
			return errors.New("Bad interaction stage for PubSub")
		}
	default:
		logger.Debugf("Cannot route message to: %s", *msg.UriTo)
		return nil
	}

	// TODO (AF): use a goroutine
	return handler(msg, transaction)
}

func (hctx *HandlerContext) OnClose() error {
	logger.Infof("close EndPoint: %s", hctx.Uri)
	// TODO (AF): Close handlers ?
	//	for key, handler := range hctx.handlers {
	//		fmt.Println("close handler: ", key)
	//		err := handler.OnClose()
	//		if err != nil {
	//			// TODO (AF): print an error message
	//		}
	//	}
	return nil
}
