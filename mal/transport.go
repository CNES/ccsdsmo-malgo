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
	"net/url"
)

type Transport interface {
	// Returns a new Message ready to encode
	NewMessage() *Message
	// Returns a new Body ready to encode
	NewBody() Body
	//	SupportedQoS(qos QoSLevel) bool
	//	SupportedIP(ip InteractionType) bool
	Transmit(msg *Message) error
	TransmitMultiple(msgs ...*Message) error
	Close() error
}

type TransportCallback interface {
	//	Ack()
	//	Error()
	Receive(msg *Message) error
	ReceiveMultiple(msgs ...*Message) error
}

type TransportFactory interface {
	NewTransport(u *url.URL, ctx TransportCallback) (Transport, *URI, error)
}

var transports map[string]TransportFactory = make(map[string]TransportFactory)

func RegisterTransportFactory(name string, factory TransportFactory) {
	logger.Infof("RegisterTransportFactory: %s", name)
	transports[name] = factory
}

func NewTransport(cfgURL string, ctx TransportCallback) (Transport, *URI, error) {
	u, err := url.Parse(cfgURL)
	if err != nil {
		logger.Warnf("NewTransport: cannot parse %s", cfgURL)
		return nil, NULL_URI, errors.New("Bad URL: " + cfgURL)
	}

	factory := transports[u.Scheme]
	if factory != nil {
		return factory.NewTransport(u, ctx)
	}
	logger.Warnf("NewTransport: unknow transport %s", cfgURL)
	return nil, NULL_URI, errors.New("Unknow transport: " + cfgURL)
}
