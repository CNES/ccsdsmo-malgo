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
package api2

import (
	. "mal"
)

type handler interface {
}

type provider interface {
	handler
}

type consumer interface {
	handler
}

// Defines a generic root Transaction interface (context of an interaction)
type Transaction interface {
	getURI() *URI
	getTid() ULong
}

// ================================================================================
// MAL Send interaction handler

type SendTransaction interface {
	Transaction
}

type SendHandler interface {
	provider
	OnSend(msg *Message, transaction SendTransaction) error
}

// ================================================================================
// MAL Submit interaction handlers

type SubmitTransaction interface {
	Transaction
	Ack(err error) error
}

type SubmitHandler interface {
	provider
	OnSubmit(msg *Message, transaction SubmitTransaction) error
}

//type ConsumerSubmit interface {
//	consumer
//	OnAck(endpoint *EndPoint, msg *Message)
//}

// ================================================================================
// MAL Request interaction handlers

type RequestTransaction interface {
	Reply([]byte, error) error
}

type RequestHandler interface {
	provider
	OnRequest(msg *Message, transaction RequestTransaction) error
}

//type ConsumerRequest interface {
//	consumer
//	OnResponse(endpoint *EndPoint, msg *Message)
//}

// ================================================================================
// MAL Invoke interaction handlers

type InvokeTransaction interface {
	Ack(error) error
	Reply([]byte, error) error
}

type InvokeHandler interface {
	provider
	OnInvoke(msg *Message, transaction InvokeTransaction) error
}

//type ConsumerInvoke interface {
//	consumer
//	OnAck(endpoint *EndPoint, msg *Message)
//	OnResponse(endpoint *EndPoint, msg *Message)
//}

// ================================================================================
// MAL Progress interaction handlers

type ProgressTransaction interface {
	Ack(error) error
	Update([]byte, error) error
	Reply([]byte, error) error
}

type ProgressHandler interface {
	provider
	OnProgress(msg *Message, transaction ProgressTransaction) error
}

//type ConsumerProgress interface {
//	consumer
//	OnAck(endpoint *EndPoint, msg *Message)
//	OnUpdate(endpoint *EndPoint, msg *Message)
//	OnResponse(endpoint *EndPoint, msg *Message)
//}

/// ================================================================================
// MAL PubSub interaction handlers

type BrokerSubscriberTransaction interface {
	RegisterAck(error) error
	Notify([]byte, error) error
	DeregisterAck(error) error
}

type BrokerPublisherTransaction interface {
	PublishRegisterAck(error) error
	PublishError(error) error
	PublishDeregisterAck(error) error
}

type BrokerHandler interface {
	handler
	OnRegister(msg *Message, transaction SubscriberTransaction) error
	OnDeregister(msg *Message, transaction SubscriberTransaction) error
	OnPublishRegister(msg *Message, transaction PublisherTransaction) error
	OnPublishDeregister(msg *Message, transaction PublisherTransaction) error
	OnPublish(msg *Message, transaction PublisherTransaction) error
}

//type ProviderPubSub interface {
//	provider
//	OnPublishRegisterAck(msg *Message)
//	OnPublishDeregisterAck(msg *Message)
//	OnPublishError(msg *Message)
//}
//
//type ConsumerPubSub interface {
//	consumer
//	OnRegisterAck(msg *Message)
//	OnDeregister(msg *Message)
//	OnNotify(msg *Message)
//}
