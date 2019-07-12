/**
 * MIT License
 *
 * Copyright (c) 2017 - 2019 CNES
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
)

type Listener interface {
	// TODO (AF): We should remove error return (useless and unused)
	OnMessage(msg *Message)
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
	// Access Control handler
	achdlr    AccessControl
	errch     chan *MessageError
	transport Transport
}

// Be careful: Depending of the application logic if there is not enough slots in
// the internal channel there may be a blocking of the underlying transport threads.
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

	logger.Infof("NewContext: transport created: %s", *uri)

	ctx.uri = *uri
	ctx.transport = transport

	go ctx.handle()
	return ctx, nil
}

func (ctx *Context) NewMessage() *Message {
	return ctx.transport.NewMessage()
}

// Returns a new Body ready to encode
func (ctx *Context) NewBody() Body {
	return ctx.transport.NewBody()
}

// Note (AF): May be we should provide a non programmatic way to fix the AccessControl handler
// in the MAL context (using a factory as in MAL Java API for example).
func (ctx *Context) SetAccessControl(achdlr AccessControl) {
	ctx.achdlr = achdlr
}

func (ctx *Context) SetErrorChannel(errch chan *MessageError) {
	if ctx.errch != nil {
		close(ctx.errch)
	}
	ctx.errch = errch
}

func (ctx *Context) NewURI(id string) *URI {
	// TODO (AF): Verify the uri
	uri := URI(string(ctx.uri) + "/" + id)
	return &uri
}

// Registers an End-Point with the specified MAL URI.
func (ctx *Context) RegisterEndPoint(uri *URI, listener Listener) error {
	// TODO (AF): Verify if the context is alive
	if listener == nil {
		logger.Warnf("Context.RegisterEndPoint: Cannot not register nil listener for %s", *uri)
		return errors.New("EndPoint is nil")
	}
	// TODO (AF): Verify the uri
	// Verify that the URI is not already registered
	old := ctx.listeners[*uri]
	if old != nil {
		logger.Warnf("Context.RegisterEndPoint: %s already registered", *uri)
		return errors.New("EndPoint already exists: " + string(*uri))
	}
	// Registers the uri
	ctx.listeners[*uri] = listener
	return nil
}

// Gets the End-Point asssociated with the specified MAL URI.
func (ctx *Context) GetEndPoint(uri *URI) (Listener, error) {
	listener := ctx.listeners[*uri]
	if listener == nil {
		logger.Warnf("Context.GetEndPoint: %s not registered", *uri)
		return nil, errors.New("EndPoint doesn't exist: " + string(*uri))
	}
	return listener, nil
}

// Removes the End-Point asssociated with the specified MAL URI.
func (ctx *Context) UnregisterEndPoint(uri *URI) error {
	listener := ctx.listeners[*uri]
	if listener == nil {
		logger.Warnf("Context.UnregisterEndPoint: %s not registered", *uri)
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
			logger.Debugf("Context.handle: receive: %+v for %s", msg, *msg.UriTo)
			to, ok := ctx.listeners[*msg.UriTo]
			if ok {
				to.OnMessage(msg)
				logger.Debugf("Context.handle: Message delivered: %s", msg)
			} else {
				logger.Errorf("Context.handle: Cannot route message to: %s", *msg.UriTo)
			}
		} else {
			logger.Infof("Context.handle: ends: %s", msg)
			ctx.ends <- true
		}
	}
}

// Close the MAL context, closing all registered listeners (end-point, etc.) and transport.
func (ctx *Context) Close() error {
	for uri, listener := range ctx.listeners {
		logger.Infof("Context.Close: %s", uri)
		listener.OnClose()
	}

	return ctx.transport.Close()
	// TODO (AF):
	//	close(ctx.ch)
	//	<-ctx.ends
}

// ================================================================================
// TransportCallback interface

// Method implementing SEND request to the transport layer.
func (ctx *Context) Send(msg *Message) error {
	msg.Timestamp = *TimeNow()
	// TODO (AF): May be we should handle errors internally using the error channel.
	if ctx.achdlr != nil {
		err := ctx.achdlr.check(msg)
		if err != nil {
			return err
		}
	}
	return ctx.transport.Transmit(msg)
}

// Method implementing RECEIVE indication from transport layer.
func (ctx *Context) Receive(msg *Message) error {
	// TODO (AF): This method should not returned errors. Errors should be handled
	// internally using the error channel.
	if ctx.achdlr != nil {
		err := ctx.achdlr.check(msg)
		if err != nil {
			return err
		}
	}
	logger.Debugf("Context.Receive: forward to client %s", *msg.UriTo)
	ctx.ch <- msg
	return nil
}

// Method implementing RECEIVEMULTIPLE indication from transport layer.
func (ctx *Context) ReceiveMultiple(msgs ...*Message) error {
	for _, msg := range msgs {
		ctx.Receive(msg)
	}
	return nil
}
