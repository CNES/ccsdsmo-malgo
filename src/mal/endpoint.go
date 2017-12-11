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

//type EndPoint interface {
//	NewEndPoint(uri *URI, ch chan Message) EndPoint, error
//	Send(msg *Message) error
//	Recv() Message, error
//}

// This concept is part of a low level MAL API, it should not be used by classic users.
// @see OperationContext
type EndPoint struct {
	Ctx *Context
	Uri *URI
	ch  chan *Message
	tid uint64
}

func NewEndPoint(ctx *Context, service string, ch chan *Message) (*EndPoint, error) {
	if ch == nil {
		// TODO (AF): Fix length of channel?
		ch = make(chan *Message, 10)
	}
	// TODO (AF): Verify the uri
	uri := ctx.NewURI(service)
	endpoint := &EndPoint{ctx, uri, ch, 0}
	err := ctx.RegisterEndPoint(uri, endpoint)
	if err != nil {
		return nil, err
	}
	return endpoint, nil
}

func (endpoint *EndPoint) Send(msg *Message) error {
	// TODO (AF): We may verify that the urifrom is the endpoint URI, we can also
	// overload it.
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

func (endpoint *EndPoint) TransactionId() ULong {
	return ULong(atomic.AddUint64(&endpoint.tid, 1))
}

func (endpoint *EndPoint) Close() error {
	return endpoint.Ctx.UnregisterEndPoint(endpoint.Uri)
}

// ================================================================================
// Defines Listener interface used by context to route MAL messages

func (endpoint *EndPoint) OnMessage(msg *Message) error {
	endpoint.ch <- msg
	return nil
}

func (endpoint *EndPoint) OnClose() error {
	// TODO (AF): ?
	return nil
}
