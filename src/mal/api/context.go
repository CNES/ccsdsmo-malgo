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
	"sync/atomic"
)

type OperationHandler interface {
	onMessage(msg *Message)
	onClose()
}

type handlerDesc struct {
	handlerType InteractionType
	area        UShort
	areaVersion UOctet
	service     UShort
	operation   UShort
	handler     Handler
}

type ClientContext struct {
	Ctx       *Context
	Uri       *URI
	ops       map[ULong]OperationHandler
	hdl       map[uint64](*handlerDesc)
	txcounter uint64
}

func NewClientContext(ctx *Context, service string) (*ClientContext, error) {
	// TODO (AF): Verify the uri
	uri := ctx.NewURI(service)
	ops := make(map[ULong]OperationHandler)
	hdl := make(map[uint64](*handlerDesc))
	cctx := &ClientContext{ctx, uri, ops, hdl, 0}
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
	old := cctx.ops[tid]
	if old != nil {
		logger.Warnf("Handler already registered for this transaction: %d", tid)
		return errors.New("Handler already registered for this transaction")
	}
	cctx.ops[tid] = handler
	return nil
}

func (cctx *ClientContext) deregisterOp(tid ULong) error {
	// TODO (AF): Synchronization
	if cctx.ops[tid] == nil {
		logger.Warnf("No handler registered for this transaction: %d", tid)
		return errors.New("No handler registered for this transaction")
	}
	delete(cctx.ops, tid)
	return nil
}

func (cctx *ClientContext) registerHdl(hdltype InteractionType, area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	key := key(area, areaVersion, service, operation)
	old := cctx.hdl[key]

	if old != nil {
		logger.Errorf("MAL handler already registered: %d", key)
		return errors.New("MAL handler already registered")
	} else {
		logger.Debugf("MAL handler registered: %d", key)
	}

	var desc = &handlerDesc{
		handlerType: hdltype,
		area:        area,
		areaVersion: areaVersion,
		service:     service,
		operation:   operation,
		handler:     handler,
	}

	cctx.hdl[key] = desc
	return nil
}

func (cctx *ClientContext) Close() error {
	// TODO (AF): Close operations and handlers
	return cctx.Ctx.UnregisterEndPoint(cctx.Uri)
}

// ================================================================================
// Defines Listener interface used by context to route MAL messages

// TODO (AF): Take in account operations and handlers!!
func (cctx *ClientContext) OnMessage(msg *Message) error {
	if ((msg.InteractionType != MAL_INTERACTIONTYPE_PUBSUB) && (msg.InteractionStage == MAL_IP_STAGE_INIT)) ||
		((msg.InteractionType == MAL_INTERACTIONTYPE_PUBSUB) && ((msg.InteractionStage & 0x1) != 0)) {
		handler, err := cctx.getHandler(msg.InteractionType, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation)
		if err != nil {
			return err
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
			if (msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER) ||
				(msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH) ||
				(msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER) {
				transaction = &PublisherTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			} else if (msg.InteractionStage == MAL_IP_STAGE_PUBSUB_REGISTER) ||
				(msg.InteractionStage == MAL_IP_STAGE_PUBSUB_DEREGISTER) {
				transaction = &SubscriberTransactionX{TransactionX{cctx.Ctx, cctx.Uri, msg.UriFrom, msg.TransactionId, msg.ServiceArea, msg.AreaVersion, msg.Service, msg.Operation}}
			} else {
				// TODO (AF): Log an error, May be we should not return this error
				return errors.New("Bad interaction stage for PubSub")
			}
		default:
			logger.Debugf("Cannot route message to: %s", *msg.UriTo)
			return nil
		}

		// TODO (AF): use a goroutine
		return handler(msg, transaction)
	} else {
		// Note (AF): The generated TransactionId is unique for this requesting URI so we
		// can use it as key to retrieve the Operation (This is more restrictive than the
		// MAL API (see section 3.2).
		to, ok := cctx.ops[msg.TransactionId]
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
	for tid, handler := range cctx.ops {
		logger.Debugf("close operation: %d", tid)
		handler.onClose()
	}
	return nil
}
