package event

import (
  "errors"
  "github.com/CNES/ccsdsmo-malgo/mal"
  malapi "github.com/CNES/ccsdsmo-malgo/mal/api"
  "github.com/CNES/ccsdsmo-malgo/com"
)

var Cctx *malapi.ClientContext
func Init(cctxin *malapi.ClientContext) error {
  if cctxin == nil {
    return errors.New("Illegal nil client context in Init")
  }
  Cctx = cctxin
  return nil
}

// generated code for operation monitorEvent
type MonitorEventSubscriberOperation struct {
  op malapi.SubscriberOperation
}
func NewMonitorEventSubscriberOperation(providerURI *mal.URI) (*MonitorEventSubscriberOperation, error) {
  op := Cctx.NewSubscriberOperation(providerURI, com.AREA_NUMBER, com.AREA_VERSION, SERVICE_NUMBER, MONITOREVENT_OPERATION_NUMBER)
  consumer := &MonitorEventSubscriberOperation{op}
  return consumer, nil
}
func (receiver *MonitorEventSubscriberOperation) Register(subscription *mal.Subscription) error {
  // create a body for the operation call
  body := receiver.op.NewBody()
  // encode in parameters
  err := body.EncodeLastParameter(subscription, false)
  if err != nil {
    return err
  }

  // operation call
  resp, err := receiver.op.Register(body)
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return err
    }
    // TODO Error handling un PUBSUB operations
    return err
  }

  // decode out parameters
  return nil
}
func (receiver *MonitorEventSubscriberOperation) GetNotify() (*mal.Identifier, *mal.UpdateHeaderList, *com.ObjectDetailsList, mal.ElementList, error) {

  // operation call
  resp, err := receiver.op.GetNotify()
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return nil, nil, nil, nil, err
    }
    // TODO Error handling un PUBSUB operations
    return nil, nil, nil, nil, err
  }

  // decode out parameters
  outElem_subscriptionid, err := resp.DecodeParameter(mal.NullIdentifier)
  if err != nil {
    return nil, nil, nil, nil, err
  }
  outParam_subscriptionid, ok := outElem_subscriptionid.(*mal.Identifier)
  if !ok {
    err = errors.New("unexpected type for parameter subscriptionid")
    return nil, nil, nil, nil, err
  }

  outElem_updateHeaders, err := resp.DecodeParameter(mal.NullUpdateHeaderList)
  if err != nil {
    return nil, nil, nil, nil, err
  }
  outParam_updateHeaders, ok := outElem_updateHeaders.(*mal.UpdateHeaderList)
  if !ok {
    err = errors.New("unexpected type for parameter updateHeaders")
    return nil, nil, nil, nil, err
  }

  outElem_eventLinks, err := resp.DecodeParameter(com.NullObjectDetailsList)
  if err != nil {
    return nil, nil, nil, nil, err
  }
  outParam_eventLinks, ok := outElem_eventLinks.(*com.ObjectDetailsList)
  if !ok {
    err = errors.New("unexpected type for parameter eventLinks")
    return nil, nil, nil, nil, err
  }

  outElem_eventBody, err := resp.DecodeLastParameter(nil, true)
  if err != nil {
    return nil, nil, nil, nil, err
  }
  outParam_eventBody, ok := outElem_eventBody.(mal.ElementList)
  if !ok {
    if outElem_eventBody == mal.NullElement {
      outParam_eventBody = mal.NullElementList
    } else {
      err = errors.New("unexpected type for parameter eventBody")
      return nil, nil, nil, nil, err
    }
  }

  return outParam_subscriptionid, outParam_updateHeaders, outParam_eventLinks, outParam_eventBody, nil
}
func (receiver *MonitorEventSubscriberOperation) Deregister(subscriptionids *mal.IdentifierList) error {
  // create a body for the operation call
  body := receiver.op.NewBody()
  // encode in parameters
  err := body.EncodeLastParameter(subscriptionids, false)
  if err != nil {
    return err
  }

  // operation call
  resp, err := receiver.op.Deregister(body)
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return err
    }
    // TODO Error handling un PUBSUB operations
    return err
  }

  // decode out parameters
  return nil
}

// generated code for operation monitorEvent
type MonitorEventPublisherOperation struct {
  op malapi.PublisherOperation
}
func NewMonitorEventPublisherOperation(providerURI *mal.URI) (*MonitorEventPublisherOperation, error) {
  op := Cctx.NewPublisherOperation(providerURI, com.AREA_NUMBER, com.AREA_VERSION, SERVICE_NUMBER, MONITOREVENT_OPERATION_NUMBER)
  consumer := &MonitorEventPublisherOperation{op}
  return consumer, nil
}
func (receiver *MonitorEventPublisherOperation) Register(entitykeys *mal.EntityKeyList) error {
  // create a body for the operation call
  body := receiver.op.NewBody()
  // encode in parameters
  err := body.EncodeLastParameter(entitykeys, false)
  if err != nil {
    return err
  }

  // operation call
  resp, err := receiver.op.Register(body)
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return err
    }
    // TODO Error handling un PUBSUB operations
    return err
  }

  // decode out parameters
  return nil
}
func (receiver *MonitorEventPublisherOperation) Publish(updateHeaders *mal.UpdateHeaderList, eventLinks *com.ObjectDetailsList, eventBody mal.ElementList) error {
  // create a body for the operation call
  body := receiver.op.NewBody()
  // encode in parameters
  err := body.EncodeParameter(updateHeaders)
  if err != nil {
    return err
  }
  err = body.EncodeParameter(eventLinks)
  if err != nil {
    return err
  }
  err = body.EncodeLastParameter(eventBody, true)
  if err != nil {
    return err
  }

  // operation call
  err = receiver.op.Publish(body)
  if err != nil {
    // TODO Error handling un PUBSUB operations
    return err
  }

  return nil
}
func (receiver *MonitorEventPublisherOperation) Deregister() error {

  // operation call
  resp, err := receiver.op.Deregister(nil)
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return err
    }
    // TODO Error handling un PUBSUB operations
    return err
  }

  // decode out parameters
  return nil
}
