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
	"fmt"
)

type Listener interface {
	OnMessage(msg *Message) error
	OnClose() error
}

//type Context interface {
//	newURI(id string) Uri_t
//	newEndPoint(uri Uri_t, ch chan) EndPoint, Error
//}

type Context struct {
	uri       URI
	listeners map[URI]Listener
	ch        chan *Message
	ends      chan bool
	transport Transport
}

func NewContext(url string) (*Context, error) {
	listeners := make(map[URI]Listener)
	// TODO (AF): Fix length of channel
	ch := make(chan *Message, 10)
	ends := make(chan bool)
	ctx := &Context{
		listeners: listeners,
		ch:        ch,
		ends:      ends,
	}

	transport, uri, err := NewTransport(url, ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println("Transport created: ", uri)

	ctx.uri = *uri
	ctx.transport = transport

	go ctx.handle()
	return ctx, nil
}

func (ctx *Context) NewURI(id string) *URI {
	// TODO (AF): Verify the uri
	uri := URI(string(ctx.uri) + "/" + id)
	return &uri
}

func (ctx *Context) RegisterEndPoint(uri *URI, listener Listener) error {
	// TODO (AF): Verify if the context is alive
	if listener == nil {
		return errors.New("EndPoint is nil")
	}
	// TODO (AF): Verify the uri
	// Verify that the URI is not already registered
	old := ctx.listeners[*uri]
	if old != nil {
		return errors.New("EndPoint already exists: " + string(*uri))
	}
	// Registers the uri
	ctx.listeners[*uri] = listener
	return nil
}

func (ctx *Context) GetEndPoint(uri *URI) (Listener, error) {
	listener := ctx.listeners[*uri]
	if listener == nil {
		return nil, errors.New("EndPoint doesn't exist: " + string(*uri))
	}
	return listener, nil
}

func (ctx *Context) UnregisterEndPoint(uri *URI) error {
	listener := ctx.listeners[*uri]
	if listener == nil {
		return errors.New("EndPoint doesn't exist: " + string(*uri))
	}
	delete(ctx.listeners, *uri)
	return nil
}

// Routes messages from transport to registered MAL EndPoint
func (ctx *Context) handle() {
	for {
		msg, more := <-ctx.ch
		if more {
			fmt.Println("Receive: ", msg)
			to, ok := ctx.listeners[*msg.UriTo]
			if ok {
				fmt.Printf("%t\n", to)
				to.OnMessage(msg)
				fmt.Println("Message transmitted: ", msg)
			} else {
				fmt.Println("Cannot route message to: ", *msg.UriTo)
			}
		} else {
			fmt.Println("Context ends: ", msg)
			ctx.ends <- true
		}
	}
}

func (ctx *Context) Close() error {
	for uri, listener := range ctx.listeners {
		fmt.Println("close EndPoint ", uri)
		listener.OnClose()
	}

	return ctx.transport.Close()
	// TODO (AF):
	//	close(ctx.ch)
	//	<-ctx.ends
}

// ================================================================================
// TransportCallback interface

func (ctx *Context) Send(msg *Message) error {
	return ctx.transport.Transmit(msg)
}

//func (ctx *Context) Send(msg *Message) error {
//	if ctx.ch != nil {
//		ctx.ch <- *msg
//		return nil
//	}
//	return errors.New("MAL context closed")
//}

func (ctx *Context) Receive(msg *Message) error {
	ctx.ch <- msg
	return nil
}

func (ctx *Context) ReceiveMultiple(msgs ...*Message) error {
	for _, msg := range msgs {
		ctx.Receive(msg)
	}
	return nil
}
