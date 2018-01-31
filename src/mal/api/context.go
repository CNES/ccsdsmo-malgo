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

// Defines a generic handler interface for providers
type ProviderHandler func(*Message, Transaction) error

type OperationHandler interface {
	onMessage(msg *Message)
	onClose()
}

type pDesc struct {
	stype InteractionType
	// TODO (AF): Not really needed, these fields are included in the correponding key.
	area        UShort
	areaVersion UOctet
	service     UShort
	operation   UShort
	handler     ProviderHandler
}

type ClientContext struct {
	Ctx        *Context
	Uri        *URI
	operations map[ULong]OperationHandler
	handlers   map[uint64](*pDesc)
	txcounter  uint64
}

func NewClientContext(ctx *Context, service string) (*ClientContext, error) {
	// TODO (AF): Verify the uri
	uri := ctx.NewURI(service)
	operations := make(map[ULong]OperationHandler)
	handlers := make(map[uint64](*pDesc))
	cctx := &ClientContext{ctx, uri, operations, handlers, 0}
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

func (cctx *ClientContext) registerProviderHandler(hdltype InteractionType, area UShort, areaVersion UOctet, service UShort, operation UShort, handler ProviderHandler) error {
	// TODO (AF): Synchronization
	key := key(area, areaVersion, service, operation)
	old := cctx.handlers[key]

	if old != nil {
		logger.Errorf("MAL handler already registered: %d", key)
		return errors.New("MAL handler already registered")
	} else {
		logger.Debugf("MAL handler registered: %d", key)
	}

	var desc = &pDesc{
		stype:       hdltype,
		area:        area,
		areaVersion: areaVersion,
		service:     service,
		operation:   operation,
		handler:     handler,
	}

	cctx.handlers[key] = desc
	return nil
}

func (cctx *ClientContext) deregisterProviderHandler(hdltype InteractionType, area UShort, areaVersion UOctet, service UShort, operation UShort) error {
	// TODO (AF): Synchronization
	key := key(area, areaVersion, service, operation)
	if cctx.handlers[key] == nil {
		logger.Warnf("No interface registered for this operation: %d", key)
		return errors.New("No interface registered for this operation")
	}
	delete(cctx.handlers, key)
	return nil
}

func (cctx *ClientContext) Close() error {
	// TODO (AF): Close operations and services
	return cctx.Ctx.UnregisterEndPoint(cctx.Uri)
}

// ================================================================================
// Defines Listener interface used by context to route MAL messages

func (cctx *ClientContext) getProviderHandler(stype InteractionType, area UShort, areaVersion UOctet, service UShort, operation UShort) (ProviderHandler, error) {
	key := key(area, areaVersion, service, operation)

	to, ok := cctx.handlers[key]
	if ok {
		if to.stype == stype {
			return to.handler, nil
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
func (cctx *ClientContext) OnMessage(msg *Message) {
	if ((msg.InteractionType != MAL_INTERACTIONTYPE_PUBSUB) && (msg.InteractionStage == MAL_IP_STAGE_INIT)) ||
		((msg.InteractionType == MAL_INTERACTIONTYPE_PUBSUB) && ((msg.InteractionStage & 0x1) != 0)) {
		handler, err := cctx.getProviderHandler(msg.InteractionType, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			// TODO (AF): Log an error? Adds an error listener?
			logger.Errorf("Cannot route message: %tv", msg)
		}
		var transaction Transaction
		switch msg.InteractionType {
		case MAL_INTERACTIONTYPE_SEND:
			transaction = &SendTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		case MAL_INTERACTIONTYPE_SUBMIT:
			transaction = &SubmitTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		case MAL_INTERACTIONTYPE_REQUEST:
			transaction = &RequestTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		case MAL_INTERACTIONTYPE_INVOKE:
			transaction = &InvokeTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		case MAL_INTERACTIONTYPE_PROGRESS:
			transaction = &ProgressTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
		case MAL_INTERACTIONTYPE_PUBSUB:
			if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER {
				transaction = &PublisherTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH {
				transaction = &PublisherTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER {
				transaction = &PublisherTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_REGISTER {
				transaction = &SubscriberTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_DEREGISTER {
				transaction = &SubscriberTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			} else {
				// TODO (AF): Log an error? Adds an error listener?
				logger.Errorf("Unknown interaction stage for PubSub: %tv", msg)
			}
		default:
			// TODO (AF): Log an error? Adds an error listener?
			logger.Errorf("Unknown interaction type: %s", msg)
			return
		}
		go handler(msg, transaction)
	} else {
		// Note (AF): The generated TransactionId is unique for this requesting URI so we
		// can use it as key to retrieve the Operation (This is more restrictive than the
		// MAL API (see section 3.2).
		to, ok := cctx.operations[msg.TransactionId]
		if ok {
			logger.Debugf("Operation.onMessage %t", to)
			// There is no need to call a go routine as this code is not blocking.
			to.onMessage(msg)
			logger.Debugf("OnMessageMessage handled: %s", msg)
		} else {
			logger.Errorf("Unknown TransactionID, cannot route message: %tv", msg)
		}
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
