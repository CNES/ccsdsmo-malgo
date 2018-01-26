/**
 * MIT License
 *
 * Copyright (c) 2018 CNES
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
	"mal/debug"
	"sync/atomic"
)

var (
	logger debug.Logger = debug.GetLogger("mal.api")
)

// TODO (AF): Is this interface useful? provider/consumer?
type service interface {
}

type OperationHandler interface {
	onMessage(msg *Message)
	onClose()
}

type sDesc struct {
	// TODO (AF): Not really needed, implicit with handler interface.
	stype InteractionType
	// TODO (AF): Not really needed, these fields are included in the correponding key.
	area        UShort
	areaVersion UOctet
	service     UShort
	operation   UShort
	// TODO (AF): Could be directly referenced in map.
	shdl service
}

type ClientContext struct {
	Ctx        *Context
	Uri        *URI
	operations map[ULong]OperationHandler
	services   map[uint64](*sDesc)
	txcounter  uint64
}

func NewClientContext(ctx *Context, service string) (*ClientContext, error) {
	// TODO (AF): Verify the uri
	uri := ctx.NewURI(service)
	operations := make(map[ULong]OperationHandler)
	services := make(map[uint64](*sDesc))
	cctx := &ClientContext{ctx, uri, operations, services, 0}
	err := ctx.RegisterEndPoint(uri, cctx)
	if err != nil {
		return nil, err
	}
	return cctx, nil
}

func (cctx *ClientContext) TransactionId() ULong {
	return ULong(atomic.AddUint64(&cctx.txcounter, 1))
}

func (cctx *ClientContext) registerOp(tid ULong, handler OperationHandler) error {
	// TODO (AF): Synchronization
	old := cctx.operations[tid]
	if old != nil {
		logger.Warnf("Handler already registered for this transaction: %d", tid)
		return errors.New("Handler already registered for this transaction")
	}
	cctx.operations[tid] = handler
	return nil
}

func (cctx *ClientContext) deregisterOp(tid ULong) error {
	// TODO (AF): Synchronization
	if cctx.operations[tid] == nil {
		logger.Warnf("No handler registered for this transaction: %d", tid)
		return errors.New("No handler registered for this transaction")
	}
	delete(cctx.operations, tid)
	return nil
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

func (cctx *ClientContext) registerService(hdltype InteractionType, area UShort, areaVersion UOctet, service UShort, operation UShort, shdl service) error {
	// TODO (AF): Synchronization
	key := key(area, areaVersion, service, operation)
	old := cctx.services[key]

	if old != nil {
		logger.Errorf("MAL handler already registered: %d", key)
		return errors.New("MAL handler already registered")
	} else {
		logger.Debugf("MAL handler registered: %d", key)
	}

	var desc = &sDesc{
		stype:       hdltype,
		area:        area,
		areaVersion: areaVersion,
		service:     service,
		operation:   operation,
		shdl:        shdl,
	}

	cctx.services[key] = desc
	return nil
}

func (cctx *ClientContext) deregisterService(hdltype InteractionType, area UShort, areaVersion UOctet, service UShort, operation UShort) error {
	// TODO (AF): Synchronization
	key := key(area, areaVersion, service, operation)
	if cctx.services[key] == nil {
		logger.Warnf("No interface registered for this operation: %d", key)
		return errors.New("No interface registered for this operation")
	}
	delete(cctx.services, key)
	return nil
}

func (cctx *ClientContext) Close() error {
	// TODO (AF): Close operations and services
	return cctx.Ctx.UnregisterEndPoint(cctx.Uri)
}

// ================================================================================
// Defines Listener interface used by context to route MAL messages

func (cctx *ClientContext) getProvider(stype InteractionType, area UShort, areaVersion UOctet, service UShort, operation UShort) (service, error) {
	key := key(area, areaVersion, service, operation)

	to, ok := cctx.services[key]
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

// TODO (AF): Take in account operations and handlers!!
func (cctx *ClientContext) OnMessage(msg *Message) error {
	if ((msg.InteractionType != MAL_INTERACTIONTYPE_PUBSUB) && (msg.InteractionStage == MAL_IP_STAGE_INIT)) ||
		((msg.InteractionType == MAL_INTERACTIONTYPE_PUBSUB) && ((msg.InteractionStage & 0x1) != 0)) {
		provider, err := cctx.getProvider(msg.InteractionType, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			return err
		}
		switch msg.InteractionType {
		// TODO (AF): We can use msg.InteractionType as selector
		case MAL_INTERACTIONTYPE_SEND:
			sendProvider := provider.(SendProvider)
			transaction := &SendTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			// TODO (AF): use a goroutine
			return sendProvider.OnSend(msg, transaction)
		case MAL_INTERACTIONTYPE_SUBMIT:
			submitProvider := provider.(SubmitProvider)
			transaction := &SubmitTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			// TODO (AF): use a goroutine
			return submitProvider.OnSubmit(msg, transaction)
		case MAL_INTERACTIONTYPE_REQUEST:
			requestProvider := provider.(RequestProvider)
			transaction := &RequestTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			// TODO (AF): use a goroutine
			return requestProvider.OnRequest(msg, transaction)
		case MAL_INTERACTIONTYPE_INVOKE:
			invokeProvider := provider.(InvokeProvider)
			transaction := &InvokeTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			// TODO (AF): use a goroutine
			return invokeProvider.OnInvoke(msg, transaction)
		case MAL_INTERACTIONTYPE_PROGRESS:
			progressProvider := provider.(ProgressProvider)
			transaction := &ProgressTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			// TODO (AF): use a goroutine
			return progressProvider.OnProgress(msg, transaction)
		case MAL_INTERACTIONTYPE_PUBSUB:
			broker := provider.(Broker)
			if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER {
				transaction := &PublisherTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
				// TODO (AF): use a goroutine
				return broker.OnPublishRegister(msg, transaction)
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH {
				transaction := &PublisherTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
				// TODO (AF): use a goroutine
				return broker.OnPublish(msg, transaction)
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER {
				transaction := &PublisherTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
				// TODO (AF): use a goroutine
				return broker.OnPublishDeregister(msg, transaction)
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_REGISTER {
				transaction := &SubscriberTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
				// TODO (AF): use a goroutine
				return broker.OnRegister(msg, transaction)
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_DEREGISTER {
				transaction := &SubscriberTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
				// TODO (AF): use a goroutine
				return broker.OnDeregister(msg, transaction)
			} else {
				// TODO (AF): Log an error, May be wa should not return this error
				return errors.New("Bad interaction stage for PubSub")
			}
		default:
			logger.Warnf("Cannot route message to: %s", *msg.UriTo)
			return nil
		}
	} else {
		// Note (AF): The generated TransactionId is unique for this requesting URI so we
		// can use it as key to retrieve the Operation (This is more restrictive than the
		// MAL API (see section 3.2).
		to, ok := cctx.operations[msg.TransactionId]
		if ok {
			logger.Debugf("onMessage %t", to)
			to.onMessage(msg)
			logger.Debugf("OnMessageMessage transmitted: %s", msg)
		} else {
			logger.Debugf("Cannot route message to: %s?TransactionId=", msg.UriTo, msg.TransactionId)
		}
		return nil
	}
}

// TODO (AF): Take in account operations and handlers!!
func (cctx *ClientContext) OnClose() error {
	logger.Infof("close EndPoint: %s", cctx.Uri)
	for tid, handler := range cctx.operations {
		logger.Debugf("close operation: %d", tid)
		handler.onClose()
	}
	return nil
}
