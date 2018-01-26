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

// Defines a generic root handler interface
type Handler func(*Message, Transaction) error

// ================================================================================
// SendHandler

// TODO (AF):
//type SendHandler func(*Message, SendTransaction) error

//func (hctx *ProviderContext) RegisterSendHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler SendHandler) error {
func (cctx *ClientContext) RegisterSendHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return cctx.registerHdl(MAL_INTERACTIONTYPE_SEND, area, areaVersion, service, operation, handler)
}

// ================================================================================
// SubmitHandler

// TODO (AF):
//type SubmitHandler func(*Message, SubmitTransaction) error

//func (hctx *ProviderContext) RegisterSendHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler SendHandler) error {
func (cctx *ClientContext) RegisterSubmitHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return cctx.registerHdl(MAL_INTERACTIONTYPE_SUBMIT, area, areaVersion, service, operation, handler)
}

// ================================================================================
// RequestHandler

// TODO (AF):
//type RequestHandler func(*Message, RequestTransaction) error

//func (hctx *ProviderContext) RegisterRequestHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler RequestHandler) error {
func (cctx *ClientContext) RegisterRequestHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return cctx.registerHdl(MAL_INTERACTIONTYPE_REQUEST, area, areaVersion, service, operation, handler)
}

// ================================================================================
// InvokeHandler

// TODO (AF):
//type InvokeHandler func(*Message, InvokeTransaction) error

//func (hctx *ProviderContext) RegisterInvokeHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler InvokeHandler) error {
func (cctx *ClientContext) RegisterInvokeHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return cctx.registerHdl(MAL_INTERACTIONTYPE_INVOKE, area, areaVersion, service, operation, handler)
}

// ================================================================================
// ProgressHandler

// TODO (AF):
//type ProgressHandler func(*Message, ProgressTransaction) error

//func (hctx *ProviderContext) RegisterSendHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler SendHandler) error {
func (cctx *ClientContext) RegisterProgressHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return cctx.registerHdl(MAL_INTERACTIONTYPE_PROGRESS, area, areaVersion, service, operation, handler)
}

// ================================================================================
// BrokerHandler: There is only one handler but 2 transactions type depending of the
// incoming interaction.

// TODO (AF):
//type BrokerHandler func(*Message, BrokerTransaction) error

//func (hctx *ProviderContext) RegisterBrokerHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler BrokerHandler) error {
func (cctx *ClientContext) RegisterBrokerHandler(area UShort, areaVersion UOctet, service UShort, operation UShort, handler Handler) error {
	return cctx.registerHdl(MAL_INTERACTIONTYPE_PUBSUB, area, areaVersion, service, operation, handler)
}

// ================================================================================
// Defines Listener interface used by context to route MAL messages

func (cctx *ClientContext) getHandler(hdltype InteractionType, area UShort, areaVersion UOctet, service UShort, operation UShort) (Handler, error) {
	key := key(area, areaVersion, service, operation)

	to, ok := cctx.hdl[key]
	if ok {
		if to.handlerType == hdltype {
			return to.handler, nil
		} else {
			logger.Errorf("Bad handler type: %d should be %d", to.handlerType, hdltype)
			return nil, errors.New("Bad handler type")
		}
	} else {
		logger.Errorf("MAL handler not registered: %d", key)
		return nil, errors.New("MAL handler not registered")
	}
}
