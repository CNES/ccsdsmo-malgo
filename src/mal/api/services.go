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
	. "mal"
)

// TODO (AF): Is this interface useful?
type provider interface {
	service
}

// TODO (AF): Is this interface useful?
type consumer interface {
	service
}

// ================================================================================
// MAL Send interaction provider

type SendProvider interface {
	provider
	OnSend(msg *Message, transaction SendTransaction) error
}

// Registers a SendProvider
func (cctx *ClientContext) RegisterSendProvider(area UShort, areaVersion UOctet, service UShort, operation UShort, provider SendProvider) error {
	return cctx.registerService(MAL_INTERACTIONTYPE_SEND, area, areaVersion, service, operation, provider)
}

// ================================================================================
// MAL Submit interaction providers

type SubmitProvider interface {
	provider
	OnSubmit(msg *Message, transaction SubmitTransaction) error
}

// Registers a SubmitProvider
func (cctx *ClientContext) RegisterSubmitProvider(area UShort, areaVersion UOctet, service UShort, operation UShort, provider SubmitProvider) error {
	return cctx.registerService(MAL_INTERACTIONTYPE_SUBMIT, area, areaVersion, service, operation, provider)
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
func (cctx *ClientContext) RegisterRequestProvider(area UShort, areaVersion UOctet, service UShort, operation UShort, provider RequestProvider) error {
	return cctx.registerService(MAL_INTERACTIONTYPE_REQUEST, area, areaVersion, service, operation, provider)
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
func (cctx *ClientContext) RegisterInvokeProvider(area UShort, areaVersion UOctet, service UShort, operation UShort, provider InvokeProvider) error {
	return cctx.registerService(MAL_INTERACTIONTYPE_INVOKE, area, areaVersion, service, operation, provider)
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
func (cctx *ClientContext) RegisterProgressProvider(area UShort, areaVersion UOctet, service UShort, operation UShort, provider ProgressProvider) error {
	return cctx.registerService(MAL_INTERACTIONTYPE_PROGRESS, area, areaVersion, service, operation, provider)
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
func (cctx *ClientContext) RegisterBroker(area UShort, areaVersion UOctet, service UShort, operation UShort, broker Broker) error {
	return cctx.registerService(MAL_INTERACTIONTYPE_PUBSUB, area, areaVersion, service, operation, broker)
}
