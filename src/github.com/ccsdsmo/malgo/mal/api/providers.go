package api

import (
	. "github.com/ccsdsmo/malgo/mal"
)

// Register an handler for Send interaction
func (cctx *ClientContext) RegisterSendHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler ProviderHandler) error {
	return cctx.registerProviderHandler(MAL_INTERACTIONTYPE_SEND, area, areaVersion, service, operation, handler)
}

// Register an handler for Submit interaction
func (cctx *ClientContext) RegisterSubmitHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler ProviderHandler) error {
	return cctx.registerProviderHandler(MAL_INTERACTIONTYPE_SUBMIT, area, areaVersion, service, operation, handler)
}

// Register an handler for Request interaction
func (cctx *ClientContext) RegisterRequestHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler ProviderHandler) error {
	return cctx.registerProviderHandler(MAL_INTERACTIONTYPE_REQUEST, area, areaVersion, service, operation, handler)
}

// Register an handler for Invoke interaction
func (cctx *ClientContext) RegisterInvokeHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler ProviderHandler) error {
	return cctx.registerProviderHandler(MAL_INTERACTIONTYPE_INVOKE, area, areaVersion, service, operation, handler)
}

// Register an handler for progress interaction
func (cctx *ClientContext) RegisterProgressHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler ProviderHandler) error {
	return cctx.registerProviderHandler(MAL_INTERACTIONTYPE_PROGRESS, area, areaVersion, service, operation, handler)
}

// ================================================================================
// Broker: There is only one handler but 2 transactions type depending of the
// incoming interaction (Publisher / Subscriber).

//func (cctx *ClientContext) RegisterBrokerHandler(area UShort, areaVersion UOctet, service UShort, operation UShort,
//	handler registerHandler, handler degisterHandler, handler registerPublishHandler, handler publishHandler, handler deregisterPublishHandler) error {
//	err := cctx.registerProviderHandler(MAL_INTERACTIONTYPE_PUBSUB, area, areaVersion, service, operation, registerHandler)
//	err = cctx.registerProviderHandler(MAL_INTERACTIONTYPE_PUBSUB, area, areaVersion, service, operation, degisterHandler)
//	err = cctx.registerProviderHandler(MAL_INTERACTIONTYPE_PUBSUB, area, areaVersion, service, operation, registerPublishHandler)
//	err = cctx.registerProviderHandler(MAL_INTERACTIONTYPE_PUBSUB, area, areaVersion, service, operation, publishHandler)
//	err = cctx.registerProviderHandler(MAL_INTERACTIONTYPE_PUBSUB, area, areaVersion, service, operation, deregisterPublishHandler)
//
//	return err
//}

func (cctx *ClientContext) RegisterBrokerHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler ProviderHandler) error {
	return cctx.registerProviderHandler(MAL_INTERACTIONTYPE_PUBSUB, area, areaVersion, service, operation, handler)
}
