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
package invm

import (
	"errors"
	"fmt"
	. "mal"
	"net/url"
)

type InVMTransport struct {
	uri    URI
	ctx    TransportCallback
	params map[string][]string
}

func (*InVMTransport) Transmit(msg *Message) error {
	u, err := url.Parse(string(*msg.UriTo))
	if err != nil {
		fmt.Println("Cannot route message, urito=", msg.UriTo)
		return err
	}
	urito := url.URL{Scheme: u.Scheme, Host: u.Host}
	fmt.Println("Forward to: ", urito)
	transport, ok := contexts[urito.String()]
	if ok {
		fmt.Println("Transmit: ", msg)
		return transport.ctx.Receive(msg)
	}
	return errors.New("Cannot route message, urito=" + string(*msg.UriTo))
}

func (transport *InVMTransport) TransmitMultiple(msgs ...*Message) error {
	return transport.ctx.ReceiveMultiple(msgs...)
}

func (transport *InVMTransport) Close() error {
	fmt.Println("close transport", transport.uri)
	delete(contexts, string(transport.uri))
	return nil
}
