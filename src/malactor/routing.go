package malactor

import (
	"errors"
	. "mal"
)

const (
	MAL_ACTOR_PROVIDER_SEND_HANDLER UOctet = iota
	MAL_ACTOR_PROVIDER_SUBMIT_HANDLER
	MAL_ACTOR_PROVIDER_REQUEST_HANDLER
	MAL_ACTOR_PROVIDER_INVOKE_HANDLER
	MAL_ACTOR_PROVIDER_PROGRESS_HANDLER
	MAL_ACTOR_PROVIDER_PUBSUB_HANDLER
	MAL_ACTOR_CONSUMER_SEND_HANDLER
	MAL_ACTOR_CONSUMER_SUBMIT_HANDLER
	MAL_ACTOR_CONSUMER_REQUEST_HANDLER
	MAL_ACTOR_CONSUMER_INVOKE_HANDLER
	MAL_ACTOR_CONSUMER_PROGRESS_HANDLER
	MAL_ACTOR_CONSUMER_PUBSUB_HANDLER
	MAL_ACTOR_BROKER_PUBSUB_HANDLER
)

// Note (AF): this interface may be closer to the "service activator" design pattern than to the "router" one.
type Routing struct {
	handlers map[uint64]*handlerDesc
	endpoint *EndPoint
}

type handlerDesc struct {
	handlerType UOctet
	area        UShort
	areaVersion UOctet
	service     UShort
	operation   UShort
	handlerPtr  handler
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

func (routing *Routing) RegisterHandler(hdltype UOctet, area UShort, areaVersion UOctet, service UShort, operation UShort, hdlptr handler) error {
	key := key(area, areaVersion, service, operation)

	old := routing.handlers[key]
	if old != nil {
		return errors.New("MAL handler already registered")
	}

	var desc = &handlerDesc{
		handlerType: hdltype,
		area:        area,
		areaVersion: areaVersion,
		service:     service,
		operation:   operation,
		handlerPtr:  hdlptr,
	}

	routing.handlers[key] = desc
	return nil
}

func (routing *Routing) GetHandler(hdltype UOctet, area UShort, areaVersion UOctet, service UShort, operation UShort) (handler, error) {
	key := key(area, areaVersion, service, operation)

	desc := routing.handlers[key]
	if desc == nil {
		return nil, errors.New("MAL handler not registered")
	}

	if desc.handlerType != hdltype {
		return nil, errors.New("Bad Type for registered MAL handler")
	}

	return desc.handlerPtr, nil
}

func (routing *Routing) DeregisterHandler(area UShort, areaVersion UOctet, service UShort, operation UShort) error {
	key := key(area, areaVersion, service, operation)

	desc := routing.handlers[key]
	if desc == nil {
		return errors.New("MAL handler not registered")
	}

	delete(routing.handlers, key)
	return nil
}

// ================================================================================
// Public interface allowing to register and deregister handlers (provider, consumer and broker).

func (routing *Routing) RegisterProviderSend(area UShort, areaVersion UOctet, service UShort, operation UShort, provider *ProviderSend) error {
	return routing.RegisterHandler(MAL_ACTOR_PROVIDER_SEND_HANDLER, area, areaVersion, service, operation, provider)
}

func (routing *Routing) RegisterProviderSubmit(area UShort, areaVersion UOctet, service UShort, operation UShort, provider *ProviderSubmit) error {
	return routing.RegisterHandler(MAL_ACTOR_PROVIDER_SUBMIT_HANDLER, area, areaVersion, service, operation, provider)
}

func (routing *Routing) RegisterConsumerSubmit(area UShort, areaVersion UOctet, service UShort, operation UShort, consumer *ConsumerSubmit) error {
	return routing.RegisterHandler(MAL_ACTOR_CONSUMER_SUBMIT_HANDLER, area, areaVersion, service, operation, consumer)
}

func (routing *Routing) RegisterProviderRequest(area UShort, areaVersion UOctet, service UShort, operation UShort, provider *ProviderRequest) error {
	return routing.RegisterHandler(MAL_ACTOR_PROVIDER_REQUEST_HANDLER, area, areaVersion, service, operation, provider)
}

func (routing *Routing) RegisterConsumerRequest(area UShort, areaVersion UOctet, service UShort, operation UShort, consumer *ConsumerRequest) error {
	return routing.RegisterHandler(MAL_ACTOR_CONSUMER_REQUEST_HANDLER, area, areaVersion, service, operation, consumer)
}

func (routing *Routing) RegisterProviderInvoke(area UShort, areaVersion UOctet, service UShort, operation UShort, provider *ProviderInvoke) error {
	return routing.RegisterHandler(MAL_ACTOR_PROVIDER_INVOKE_HANDLER, area, areaVersion, service, operation, provider)
}

func (routing *Routing) RegisterConsumerInvoke(area UShort, areaVersion UOctet, service UShort, operation UShort, consumer *ConsumerInvoke) error {
	return routing.RegisterHandler(MAL_ACTOR_CONSUMER_INVOKE_HANDLER, area, areaVersion, service, operation, consumer)
}

func (routing *Routing) RegisterProviderProgress(area UShort, areaVersion UOctet, service UShort, operation UShort, provider *ProviderProgress) error {
	return routing.RegisterHandler(MAL_ACTOR_PROVIDER_PROGRESS_HANDLER, area, areaVersion, service, operation, provider)
}

func (routing *Routing) RegisterConsumerprogress(area UShort, areaVersion UOctet, service UShort, operation UShort, consumer *ConsumerProgress) error {
	return routing.RegisterHandler(MAL_ACTOR_CONSUMER_PROGRESS_HANDLER, area, areaVersion, service, operation, consumer)
}

func (routing *Routing) RegisterProviderPubSub(area UShort, areaVersion UOctet, service UShort, operation UShort, provider *ProviderPubSub) error {
	return routing.RegisterHandler(MAL_ACTOR_PROVIDER_PUBSUB_HANDLER, area, areaVersion, service, operation, provider)
}

func (routing *Routing) RegisterConsumerPubSub(area UShort, areaVersion UOctet, service UShort, operation UShort, consumer *ConsumerPubSub) error {
	return routing.RegisterHandler(MAL_ACTOR_CONSUMER_PUBSUB_HANDLER, area, areaVersion, service, operation, consumer)
}

func (routing *Routing) RegisterBrokerPubSub(area UShort, areaVersion UOctet, service UShort, operation UShort, broker *BrokerPubSub) error {
	return routing.RegisterHandler(MAL_ACTOR_BROKER_PUBSUB_HANDLER, area, areaVersion, service, operation, broker)
}

// ================================================================================
// Method allowing to handle MAL messages by appropriate handlers

// TODO (AF): use i18n
const (
	MAL_ROUTING_NO_HANDLER_MSG              = " *** mal_routing_handle: Error NO HANDLER CORRESPONDING TO THIS MESSAGE"
	MAL_ROUTING_BAD_IP_STAGE_MSG            = " *** mal_routing_handle: Error BAD IP STAGE"
	MAL_ROUTING_UNKNOW_INTERACTION_TYPE_MSG = " *** mal_routing_handle: Error UNKNOW INTERACTION TYPE"
)

func (routing *Routing) Handle(msg *Message) error {
	switch msg.InteractionType {
	case MAL_INTERACTIONTYPE_SEND:
		if msg.InteractionStage == MAL_IP_STAGE_SEND {
			handler, err := routing.GetHandler(
				MAL_ACTOR_PROVIDER_SEND_HANDLER,
				msg.ServiceArea,
				msg.AreaVersion,
				msg.Service,
				msg.Operation)

			if (handler != nil) && (err == nil) {
				provider := handler.(ProviderSend)
				provider.OnSend(routing.endpoint, msg)
			} else {
				//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
				//        rc = MAL_ROUTING_NO_HANDLER;
				return errors.New(MAL_ROUTING_NO_HANDLER_MSG)
			}
		} else {
			//      clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_BAD_IP_STAGE_MSG);
			//      rc = MAL_ROUTING_BAD_IP_STAGE;
			return errors.New(MAL_ROUTING_BAD_IP_STAGE_MSG)
		}
		break
	}
	//  case MAL_INTERACTIONTYPE_SUBMIT: {
	//    switch (mal_message_get_interaction_stage(message)) {
	//    case MAL_IP_STAGE_SUBMIT: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_SUBMIT/MAL_IP_STAGE_SUBMIT\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_PROVIDER_SUBMIT_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.provider_submit_handler.on_submit(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_SUBMIT_ACK: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_SUBMIT/MAL_IP_STAGE_SUBMIT_ACK\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_CONSUMER_SUBMIT_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc= handler->spec.consumer_submit_handler.on_ack(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    default:
	//      clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_BAD_IP_STAGE_MSG);
	//      rc = MAL_ROUTING_BAD_IP_STAGE;
	//    }
	//    break;
	//  }
	//  case MAL_INTERACTIONTYPE_REQUEST: {
	//    switch (mal_message_get_interaction_stage(message)) {
	//    case MAL_IP_STAGE_REQUEST: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_REQUEST/MAL_IP_STAGE_REQUEST\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_PROVIDER_REQUEST_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.provider_request_handler.on_request(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_REQUEST_RESPONSE: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_REQUEST/MAL_IP_STAGE_REQUEST_RESPONSE\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_CONSUMER_REQUEST_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc =handler->spec.consumer_request_handler.on_response(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    default:
	//      clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_BAD_IP_STAGE_MSG);
	//      rc = MAL_ROUTING_BAD_IP_STAGE;
	//    }
	//    break;
	//  }
	//  case MAL_INTERACTIONTYPE_INVOKE: {
	//    switch (mal_message_get_interaction_stage(message)) {
	//    case MAL_IP_STAGE_INVOKE: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_INVOKE/MAL_IP_STAGE_INVOKE\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_PROVIDER_INVOKE_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.provider_invoke_handler.on_invoke(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_INVOKE_ACK: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_INVOKE/MAL_IP_STAGE_INVOKE_ACK\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_CONSUMER_INVOKE_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.consumer_invoke_handler.on_ack(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_INVOKE_RESPONSE: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_INVOKE/MAL_IP_STAGE_INVOKE_RESPONSE\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_CONSUMER_INVOKE_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.consumer_invoke_handler.on_response(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    default:
	//      clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_BAD_IP_STAGE_MSG);
	//      rc = MAL_ROUTING_BAD_IP_STAGE;
	//    }
	//    break;
	//  }
	//  case MAL_INTERACTIONTYPE_PROGRESS: {
	//    switch (mal_message_get_interaction_stage(message)) {
	//    case MAL_IP_STAGE_PROGRESS: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_PROGRESS/MAL_IP_STAGE_PROGRESS\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_PROVIDER_PROGRESS_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.provider_progress_handler.on_progress(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_PROGRESS_ACK: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_PROGRESS/MAL_IP_STAGE_PROGRESS_ACK\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_CONSUMER_PROGRESS_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.consumer_progress_handler.on_ack(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_PROGRESS_UPDATE: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_PROGRESS/MAL_IP_STAGE_PROGRESS_UPDATE\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_CONSUMER_PROGRESS_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.consumer_progress_handler.on_update(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_PROGRESS_RESPONSE: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_PROGRESS/MAL_IP_STAGE_PROGRESS_RESPONSE\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_CONSUMER_PROGRESS_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.consumer_progress_handler.on_response(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    default:
	//      clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_BAD_IP_STAGE_MSG);
	//      rc = MAL_ROUTING_BAD_IP_STAGE;
	//    }
	//    break;
	//  }
	//  case MAL_INTERACTIONTYPE_PUBSUB: {
	//    switch (mal_message_get_interaction_stage(message)) {
	//    case MAL_IP_STAGE_PUBSUB_REGISTER: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_PUBSUB/MAL_IP_STAGE_PUBSUB_REGISTER\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_BROKER_PUBSUB_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.broker_pubsub_handler.on_register(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_PUBSUB_REGISTER_ACK: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_PUBSUB/MAL_IP_STAGE_PUBSUB_REGISTER_ACK\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_CONSUMER_PUBSUB_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.consumer_pubsub_handler.on_register_ack(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER: {
	//      clog_error(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_PUBSUB/MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER\n");
	//      // Should never happen with ZMQ Transport
	//
	//        mal_routing_handler_t* handler = mal_routing_get_handler(
	//            self,
	//            MAL_ACTOR_BROKER_PUBSUB_HANDLER,
	//            mal_message_get_service_area(message),
	//            mal_message_get_area_version(message),
	//            mal_message_get_service(message),
	//            mal_message_get_operation(message));
	//
	//        if (handler != NULL) {
	//          rc = handler->spec.broker_pubsub_handler.on_publish_register(self->state, mal_ctx, self->mal_endpoint, message);
	//        } else {
	//          clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//          rc = MAL_ROUTING_NO_HANDLER;
	//        }
	//      break;
	//    }
	//    case MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER_ACK: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_PUBSUB/MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER_ACK\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_PROVIDER_PUBSUB_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.provider_pubsub_handler.on_publish_register_ack(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_PUBSUB_PUBLISH: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_PUBSUB/MAL_IP_STAGE_PUBSUB_PUBLISH\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_BROKER_PUBSUB_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.broker_pubsub_handler.on_publish(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_PUBSUB_NOTIFY: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_PUBSUB/MAL_IP_STAGE_PUBSUB_NOTIFY\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_CONSUMER_PUBSUB_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.consumer_pubsub_handler.on_notify(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_PUBSUB_DEREGISTER: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_PUBSUB/MAL_IP_STAGE_PUBSUB_DEREGISTER\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_BROKER_PUBSUB_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.broker_pubsub_handler.on_deregister(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_PUBSUB_DEREGISTER_ACK: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_PUBSUB/MAL_IP_STAGE_PUBSUB_DEREGISTER_ACK\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_CONSUMER_PUBSUB_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.consumer_pubsub_handler.on_deregister_ack(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER: {
	//      clog_error(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_PUBSUB/MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER\n");
	//      // Should never happen  with ZMQ Transport
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_BROKER_PUBSUB_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.broker_pubsub_handler.on_publish_deregister(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    case MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER_ACK: {
	//      clog_debug(mal_logger, " *** mal_routing_handle: MAL_INTERACTIONTYPE_PUBSUB/MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER_ACK\n");
	//
	//      mal_routing_handler_t* handler = mal_routing_get_handler(
	//          self,
	//          MAL_ACTOR_PROVIDER_PUBSUB_HANDLER,
	//          mal_message_get_service_area(message),
	//          mal_message_get_area_version(message),
	//          mal_message_get_service(message),
	//          mal_message_get_operation(message));
	//
	//      if (handler != NULL) {
	//        rc = handler->spec.provider_pubsub_handler.on_publish_deregister_ack(self->state, mal_ctx, self->mal_endpoint, message);
	//      } else {
	//        clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_NO_HANDLER_MSG);
	//        rc = MAL_ROUTING_NO_HANDLER;
	//      }
	//      break;
	//    }
	//    default:
	//      clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_BAD_IP_STAGE_MSG);
	//      rc = MAL_ROUTING_BAD_IP_STAGE;
	//    }
	//    break;
	//  }
	//  default:
	//    clog_error(mal_logger, "ERROR: %s\n", MAL_ROUTING_UNKNOW_INTERACTION_TYPE_MSG);
	//    rc = MAL_ROUTING_UNKNOW_INTERACTION_TYPE;
	//  }
	//
	//  return rc;
	return nil
}
