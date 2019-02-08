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
package invm

import (
	"errors"
	. "github.com/CNES/ccsdsmo-malgo/mal"
	"github.com/CNES/ccsdsmo-malgo/mal/debug"
	"net/url"
)

var (
	logger debug.Logger = debug.GetLogger("mal.transport.invm")
)

type InVMTransport struct {
	uri    URI
	ctx    TransportCallback
	params map[string][]string
}

// Returns a new Message ready to encode
func (transport *InVMTransport) NewMessage() *Message {
	msg := &Message{Body: NewInVMBody(make([]byte, 0, 1024), true)}
	return msg
}

// Returns a new Body ready to encode
func (transport *InVMTransport) NewBody() Body {
	return NewInVMBody(make([]byte, 0, 1024), true)
}

func (*InVMTransport) Transmit(msg *Message) error {
	u, err := url.Parse(string(*msg.UriTo))
	if err != nil {
		logger.Errorf("Cannot parse urito=%s, %s", *msg.UriTo, err)
		return err
	}

	// Transform Body to readable
	msg.Body.(*InVMBody).content = msg.Body.(*InVMBody).getEncodedContent()
	msg.Body.Reset(false)

	urito := url.URL{Scheme: u.Scheme, Host: u.Host}
	transport, ok := contexts[urito.String()]
	if ok {
		logger.Debugf("Forward Message%+v to %s", *msg, *msg.UriTo)
		return transport.ctx.Receive(msg)
	}
	logger.Errorf("Cannot route Message%+v to %s", msg, *msg.UriTo)
	return errors.New("Cannot route message to" + string(*msg.UriTo))
}

func (transport *InVMTransport) TransmitMultiple(msgs ...*Message) error {
	return transport.ctx.ReceiveMultiple(msgs...)
}

func (transport *InVMTransport) Close() error {
	logger.Infof("InVMTransport.Close: %s", transport.uri)
	delete(contexts, string(transport.uri))
	return nil
}
