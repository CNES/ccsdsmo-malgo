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

// TODO (AF): Is this interface useful?
type service interface {
}

// TODO (AF): Is this interface useful?
type provider interface {
	service
}

// TODO (AF): Is this interface useful?
type consumer interface {
	service
}

type ProviderContext struct {
	HandlerContext
}

// TODO (AF): Merge with NewHandlerContext
func NewProviderContext(ctx *Context, service string) (*ProviderContext, error) {
	// TODO (AF): Verify the uri
	uri := ctx.NewURI(service)
	// TODO (AF): Fix length of channel?
	ch := make(chan *Message, 10)
	handlers := make(map[uint64](*handlerDesc))
	pctx := &ProviderContext{HandlerContext{ctx, uri, ch, handlers}}
	err := ctx.RegisterEndPoint(uri, pctx)
	if err != nil {
		return nil, err
	}
	return pctx, nil
}

func (pctx *ProviderContext) register(stype InteractionType, area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	key := key(area, areaVersion, service, operation)
	old := pctx.handlers[key]

	if old != nil {
		logger.Errorf("MAL service already registered: %d", key)
		return errors.New("MAL service already registered")
	} else {
		logger.Debugf("MAL service registered: %d", key)
	}

	var desc = &handlerDesc{
		handlerType: stype,
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

// ================================================================================
// MAL Send interaction provider

type SendProvider interface {
	provider
	OnSend(msg *Message, transaction SendTransaction) error
}

// Registers a SendProvider
func (pctx *ProviderContext) RegisterSendProvider(area UShort, areaVersion UOctet, service UShort, operation UShort, provider SendProvider) error {

	handler := func(msg *Message, tx Transaction) error {
		return provider.OnSend(msg, tx.(SendTransaction))
	}

	return pctx.register(MAL_INTERACTIONTYPE_SEND, area, areaVersion, service, operation, handler)
}

// ================================================================================
// MAL Submit interaction providers

type SubmitProvider interface {
	provider
	OnSubmit(msg *Message, transaction SubmitTransaction) error
}

// Registers a SubmitProvider
func (pctx *ProviderContext) RegisterSubmitProvider(area UShort, areaVersion UOctet, service UShort, operation UShort, provider SubmitProvider) error {

	handler := func(msg *Message, tx Transaction) error {
		return provider.OnSubmit(msg, tx.(SubmitTransaction))
	}

	return pctx.register(MAL_INTERACTIONTYPE_SUBMIT, area, areaVersion, service, operation, handler)
}

//type ConsumerSubmit interface {
//	consumer
//	OnAck(endpoint *EndPoint, msg *Message)
//}

// ================================================================================
// MAL Request interaction providers

type RequestProvider interface {
	provider
	OnRequest(msg *Message, transaction RequestTransaction) error
}

// Registers a RequestProvider
func (pctx *ProviderContext) RegisterRequestProvider(area UShort, areaVersion UOctet, service UShort, operation UShort, provider RequestProvider) error {

	handler := func(msg *Message, tx Transaction) error {
		return provider.OnRequest(msg, tx.(RequestTransaction))
	}

	return pctx.register(MAL_INTERACTIONTYPE_REQUEST, area, areaVersion, service, operation, handler)
}

//type ConsumerRequest interface {
//	consumer
//	OnResponse(endpoint *EndPoint, msg *Message)
//}

// ================================================================================
// MAL Invoke interaction providers

type InvokeProvider interface {
	provider
	OnInvoke(msg *Message, transaction InvokeTransaction) error
}

// Registers an InvokeProvider
func (pctx *ProviderContext) RegisterInvokeProvider(area UShort, areaVersion UOctet, service UShort, operation UShort, provider InvokeProvider) error {

	handler := func(msg *Message, tx Transaction) error {
		return provider.OnInvoke(msg, tx.(InvokeTransaction))
	}

	return pctx.register(MAL_INTERACTIONTYPE_INVOKE, area, areaVersion, service, operation, handler)
}

//type ConsumerInvoke interface {
//	consumer
//	OnAck(endpoint *EndPoint, msg *Message)
//	OnResponse(endpoint *EndPoint, msg *Message)
//}

// ================================================================================
// MAL Progress interaction providers

type ProgressProvider interface {
	provider
	OnProgress(msg *Message, transaction ProgressTransaction) error
}

// Registers a ProgressProvider
func (pctx *ProviderContext) RegisterProgressProvider(area UShort, areaVersion UOctet, service UShort, operation UShort, provider ProgressProvider) error {

	handler := func(msg *Message, tx Transaction) error {
		return provider.OnProgress(msg, tx.(ProgressTransaction))
	}

	return pctx.register(MAL_INTERACTIONTYPE_PROGRESS, area, areaVersion, service, operation, handler)
}

// TODO (AF): May it makes sense to implements such an interface for Progress interaction

//type ConsumerProgress interface {
//	consumer
//	OnAck(endpoint *EndPoint, msg *Message)
//	OnUpdate(endpoint *EndPoint, msg *Message)
//	OnResponse(endpoint *EndPoint, msg *Message)
//}

/// ================================================================================
// MAL PubSub interaction providers

type Broker interface {
	service
	OnRegister(msg *Message, transaction SubscriberTransaction) error
	OnDeregister(msg *Message, transaction SubscriberTransaction) error
	OnPublishRegister(msg *Message, transaction PublisherTransaction) error
	OnPublishDeregister(msg *Message, transaction PublisherTransaction) error
	OnPublish(msg *Message, transaction PublisherTransaction) error
}

// Registers a broker
func (pctx *ProviderContext) RegisterBroker(area UShort, areaVersion UOctet, service UShort, operation UShort, broker Broker) error {

	handler := func(msg *Message, tx Transaction) error {
		switch msg.InteractionStage {
		case MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER:
			// TODO (AF): use a goroutine
			return broker.OnPublishRegister(msg, tx.(PublisherTransaction))
		case MAL_IP_STAGE_PUBSUB_PUBLISH:
			// TODO (AF): use a goroutine
			return broker.OnPublish(msg, tx.(PublisherTransaction))
		case MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER:
			// TODO (AF): use a goroutine
			return broker.OnPublishDeregister(msg, tx.(PublisherTransaction))
		case MAL_IP_STAGE_PUBSUB_REGISTER:
			// TODO (AF): use a goroutine
			return broker.OnRegister(msg, tx.(SubscriberTransaction))
		case MAL_IP_STAGE_PUBSUB_DEREGISTER:
			// TODO (AF): use a goroutine
			return broker.OnDeregister(msg, tx.(SubscriberTransaction))
		default:
			// TODO (AF): Log an error, May be wa should not return this error
			return errors.New("Bad interaction stage for PubSub")
		}
	}

	return pctx.register(MAL_INTERACTIONTYPE_PUBSUB, area, areaVersion, service, operation, handler)
}
