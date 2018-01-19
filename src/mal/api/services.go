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

const (
	_SEND_PROVIDER UOctet = iota
	_SUBMIT_PROVIDER
	_REQUEST_PROVIDER
	_INVOKE_PROVIDER
	_PROGRESS_PROVIDER
	_PUBSUB_PROVIDER
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

type sdesc struct {
	stype       UOctet
	area        UShort
	areaVersion UOctet
	service     UShort
	operation   UShort
	shdl        service
}

type ProviderContext struct {
	Ctx      *Context
	Uri      *URI
	ch       chan *Message
	services map[uint64](*sdesc)
}

func NewProviderContext(ctx *Context, service string) (*ProviderContext, error) {
	// TODO (AF): Verify the uri
	uri := ctx.NewURI(service)
	// TODO (AF): Fix length of channel?
	ch := make(chan *Message, 10)
	services := make(map[uint64](*sdesc))
	pctx := &ProviderContext{ctx, uri, ch, services}
	err := ctx.RegisterEndPoint(uri, pctx)
	if err != nil {
		return nil, err
	}
	return pctx, nil
}

func (pctx *ProviderContext) register(stype UOctet, area UShort, areaVersion UOctet, service UShort, operation UShort, shdl service) error {
	key := key(area, areaVersion, service, operation)
	old := pctx.services[key]

	if old != nil {
		logger.Errorf("MAL service already registered: %d", key)
		return errors.New("MAL service already registered")
	} else {
		logger.Debugf("MAL service registered: %d", key)
	}

	var desc = &sdesc{
		stype:       stype,
		area:        area,
		areaVersion: areaVersion,
		service:     service,
		operation:   operation,
		shdl:        shdl,
	}

	pctx.services[key] = desc
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

// Re<gisters a SendProvider
func (pctx *ProviderContext) RegisterSendProvider(area UShort, areaVersion UOctet, service UShort, operation UShort, provider SendProvider) error {
	return pctx.register(_SEND_PROVIDER, area, areaVersion, service, operation, provider)
}

// ================================================================================
// MAL Submit interaction providers

type SubmitProvider interface {
	provider
	OnSubmit(msg *Message, transaction SubmitTransaction) error
}

// Registers a SubmitProvider
func (pctx *ProviderContext) RegisterSubmitProvider(area UShort, areaVersion UOctet, service UShort, operation UShort, provider SubmitProvider) error {
	return pctx.register(_SUBMIT_PROVIDER, area, areaVersion, service, operation, provider)
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
	return pctx.register(_REQUEST_PROVIDER, area, areaVersion, service, operation, provider)
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
	return pctx.register(_INVOKE_PROVIDER, area, areaVersion, service, operation, provider)
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
	return pctx.register(_PROGRESS_PROVIDER, area, areaVersion, service, operation, provider)
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
	return pctx.register(_PUBSUB_PROVIDER, area, areaVersion, service, operation, broker)
}

//type Publisher interface {
//	provider
//	OnPublishRegisterAck(msg *Message)
//	OnPublishDeregisterAck(msg *Message)
//	OnPublishError(msg *Message)
//}

// TODO (AF): May it makes sense to implements such an interface for Progress interaction

//type Subscriber interface {
//	consumer
//	OnRegisterAck(msg *Message)
//	OnDeregister(msg *Message)
//	OnNotify(msg *Message)
//}

// Registers a subscriber
//func (pctx *ProviderContext) RegisterSubscriber(area UShort, areaVersion UOctet, service UShort, operation UShort, subscriber Subscriber) error {
//	return pctx.register(_PUBSUB_PROVIDER, area, areaVersion, service, operation, subscriber)
//}

// ================================================================================
// Defines Listener interface used by context to route MAL messages

func (pctx *ProviderContext) getProvider(stype UOctet, area UShort, areaVersion UOctet, service UShort, operation UShort) (service, error) {
	key := key(area, areaVersion, service, operation)

	to, ok := pctx.services[key]
	if ok {
		if to.stype == stype {
			return to.shdl, nil
		} else {
			logger.Debugf("Bad service type: %d should be %d", to.stype, stype)
			return nil, errors.New("Bad handler type")
		}
	} else {
		logger.Debugf("MAL service not registered: %d", key)
		return nil, errors.New("MAL service not registered")
	}
}

func (pctx *ProviderContext) OnMessage(msg *Message) error {
	switch msg.InteractionType {
	// TODO (AF): May be we have to test for each interaction if the stage is valid
	case MAL_INTERACTIONTYPE_SEND:
		provider, err := pctx.getProvider(_SEND_PROVIDER, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			return err
		}
		sendProvider := provider.(SendProvider)
		transaction := &SendTransactionX{TransactionX{pctx.Ctx, pctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		// TODO (AF): use a goroutine
		return sendProvider.OnSend(msg, transaction)
	case MAL_INTERACTIONTYPE_SUBMIT:
		provider, err := pctx.getProvider(_SUBMIT_PROVIDER, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			return err
		}
		submitProvider := provider.(SubmitProvider)
		transaction := &SubmitTransactionX{TransactionX{pctx.Ctx, pctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		// TODO (AF): use a goroutine
		return submitProvider.OnSubmit(msg, transaction)
	case MAL_INTERACTIONTYPE_REQUEST:
		provider, err := pctx.getProvider(_REQUEST_PROVIDER, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			return err
		}
		requestProvider := provider.(RequestProvider)
		transaction := &RequestTransactionX{TransactionX{pctx.Ctx, pctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		// TODO (AF): use a goroutine
		return requestProvider.OnRequest(msg, transaction)
	case MAL_INTERACTIONTYPE_INVOKE:
		provider, err := pctx.getProvider(_INVOKE_PROVIDER, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			return err
		}
		invokeProvider := provider.(InvokeProvider)
		transaction := &InvokeTransactionX{TransactionX{pctx.Ctx, pctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		// TODO (AF): use a goroutine
		return invokeProvider.OnInvoke(msg, transaction)
	case MAL_INTERACTIONTYPE_PROGRESS:
		provider, err := pctx.getProvider(_PROGRESS_PROVIDER, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			return err
		}
		progressProvider := provider.(ProgressProvider)
		transaction := &ProgressTransactionX{TransactionX{pctx.Ctx, pctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		// TODO (AF): use a goroutine
		return progressProvider.OnProgress(msg, transaction)
	case MAL_INTERACTIONTYPE_PUBSUB:
		provider, err := pctx.getProvider(_PUBSUB_PROVIDER, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			return err
		}
		broker := provider.(Broker)
		if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER {
			transaction := &PublisherTransactionX{TransactionX{pctx.Ctx, pctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			// TODO (AF): use a goroutine
			return broker.OnPublishRegister(msg, transaction)
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH {
			transaction := &PublisherTransactionX{TransactionX{pctx.Ctx, pctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			// TODO (AF): use a goroutine
			return broker.OnPublish(msg, transaction)
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER {
			transaction := &PublisherTransactionX{TransactionX{pctx.Ctx, pctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			// TODO (AF): use a goroutine
			return broker.OnPublishDeregister(msg, transaction)
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_REGISTER {
			transaction := &SubscriberTransactionX{TransactionX{pctx.Ctx, pctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			// TODO (AF): use a goroutine
			return broker.OnRegister(msg, transaction)
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_DEREGISTER {
			transaction := &SubscriberTransactionX{TransactionX{pctx.Ctx, pctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			// TODO (AF): use a goroutine
			return broker.OnDeregister(msg, transaction)
		} else {
			// TODO (AF): Log an error, May be wa should not return this error
			return errors.New("Bad interaction stage for PubSub")
		}
	default:
		logger.Warnf("Cannot route message to: %s", *msg.UriTo)
	}

	return nil
}

func (pctx *ProviderContext) OnClose() error {
	logger.Infof("close EndPoint: %s", pctx.Uri)
	// TODO (AF): Close services ?
	//	for key, shdl := range pctx.services {
	//		fmt.Println("close service: ", key)
	//		err := shdl.OnClose()
	//		if err != nil {
	//			// TODO (AF): print an error message
	//		}
	//	}
	return nil
}
