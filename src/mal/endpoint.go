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
package mal

import (
	"errors"
	"sync/atomic"
)

const (
	// Defines the default length of end-point channel.
	// Be careful: Depending of the application logic if there is not enough slots in the
	// end-point channel there may be a blocking of the underlying MAL context thread.
	dflt_endpoint_channel_length uint = 10
)

//type EndPoint interface {
//	NewEndPoint(uri *URI, ch chan Message) EndPoint, error
//	Send(msg *Message) error
//	Recv() Message, error
//}

// This concept is part of a low level MAL API, it should not be used by classic users.
type EndPoint struct {
	Ctx *Context
	Uri *URI
	ch  chan *Message
	tid uint64
}

// Creates a new end-point in related MAL context for specified service.
// If the message channel is not specified (nil) it is created.
//
// @param ctx		the underlying MAL context
// @param service	the identity of service to register
// @param ch		the message channel to receive message
func NewEndPoint(ctx *Context, service string, ch chan *Message) (*EndPoint, error) {
	if ch == nil {
		ch = make(chan *Message, dflt_endpoint_channel_length)
	}
	uri := ctx.NewURI(service)
	endpoint := &EndPoint{ctx, uri, ch, 0}
	err := ctx.RegisterEndPoint(uri, endpoint)
	if err != nil {
		return nil, err
	}
	return endpoint, nil
}

func (endpoint *EndPoint) Send(msg *Message) error {
	// Verify that the urifrom is the endpoint URI.
	if msg.UriFrom == nil {
		msg.UriFrom = endpoint.Uri
	} else if *msg.UriFrom != *endpoint.Uri {
		logger.Warnf("EndPoint.Send, bad urifrom field force endpoint URI")
		msg.UriFrom = endpoint.Uri
	}
	return endpoint.Ctx.Send(msg)
}

func (endpoint *EndPoint) Recv() (*Message, error) {
	msg, ok := <-endpoint.ch
	if ok {
		return msg, nil
	} else {
		return nil, errors.New("MAL context closed")
	}
}

// Gets next TransactionId for this end-point.
func (endpoint *EndPoint) TransactionId() ULong {
	return ULong(atomic.AddUint64(&endpoint.tid, 1))
}

// Closes this end-point.
func (endpoint *EndPoint) Close() error {
	return endpoint.Ctx.UnregisterEndPoint(endpoint.Uri)
}

// ================================================================================
// Defines Listener interface used by context to route MAL messages

func (endpoint *EndPoint) OnMessage(msg *Message) {
	endpoint.ch <- msg
}

func (endpoint *EndPoint) OnClose() error {
	// TODO (AF): ?
	return nil
}
