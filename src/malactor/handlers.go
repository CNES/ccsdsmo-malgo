package malactor

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

// ================================================================================
// MAL Send interaction handler

type ProviderSend interface {
	provider
	OnSend(endpoint *EndPoint, msg *Message)
}

// ================================================================================
// MAL Submit interaction handlers

type ProviderSubmit interface {
	provider
	OnSubmit(endpoint *EndPoint, msg *Message)
}

type ConsumerSubmit interface {
	consumer
	OnAck(endpoint *EndPoint, msg *Message)
}

// ================================================================================
// MAL Request interaction handlers

type ProviderRequest interface {
	provider
	OnRequest(endpoint *EndPoint, msg *Message)
}

type ConsumerRequest interface {
	consumer
	OnResponse(endpoint *EndPoint, msg *Message)
}

// ================================================================================
// MAL Invoke interaction handlers

type ProviderInvoke interface {
	provider
	OnInvoke(endpoint *EndPoint, msg *Message)
}

type ConsumerInvoke interface {
	consumer
	OnAck(endpoint *EndPoint, msg *Message)
	OnResponse(endpoint *EndPoint, msg *Message)
}

// ================================================================================
// MAL Progress interaction handlers

type ProviderProgress interface {
	provider
	OnProgress(endpoint *EndPoint, msg *Message)
}

type ConsumerProgress interface {
	consumer
	OnAck(endpoint *EndPoint, msg *Message)
	OnUpdate(endpoint *EndPoint, msg *Message)
	OnResponse(endpoint *EndPoint, msg *Message)
}

// ================================================================================
// MAL PubSub interaction handlers

type ProviderPubSub interface {
	provider
	OnPublishRegisterAck(endpoint *EndPoint, msg *Message)
	OnPublishDeregisterAck(endpoint *EndPoint, msg *Message)
	OnPublishError(endpoint *EndPoint, msg *Message)
}

type ConsumerPubSub interface {
	consumer
	OnRegisterAck(endpoint *EndPoint, msg *Message)
	OnDeregister(endpoint *EndPoint, msg *Message)
	OnNotify(endpoint *EndPoint, msg *Message)
}

type BrokerPubSub interface {
	handler
	OnRegister(endpoint *EndPoint, msg *Message)
	OnDeregister(endpoint *EndPoint, msg *Message)
	OnPublishRegister(endpoint *EndPoint, msg *Message)
	OnPublishDeregister(endpoint *EndPoint, msg *Message)
	OnPublish(endpoint *EndPoint, msg *Message)
	OnNotifyError(endpoint *EndPoint, msg *Message)
}
