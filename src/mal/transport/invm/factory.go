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
	. "github.com/ccsdsmo/malgo/src/mal"
	"net/url"
)

type InVMTransportFactory struct {
}

var contexts map[string]*InVMTransport = make(map[string]*InVMTransport)

func init() {
	RegisterTransportFactory("invm", new(InVMTransportFactory))
}

func (*InVMTransportFactory) NewTransport(u *url.URL, ctx TransportCallback) (Transport, *URI, error) {
	// Builds base URI from URL
	base := url.URL{Scheme: u.Scheme, Host: u.Host}
	uri := URI(base.String())

	logger.Infof("InVMTransportFactory.InVMTransportFactory: registers %s", uri)

	// Gets parameters from URL
	params := u.Query()

	transport := &InVMTransport{
		uri:    uri,
		ctx:    ctx,
		params: params,
	}

	// Registers the MAL context
	if contexts[string(uri)] != nil {
		logger.Warnf("InVMTransportFactory.InVMTransportFactory: MAL context already registered ", uri)
		return nil, NullURI, errors.New("MAL context already registered")
	}
	contexts[string(uri)] = transport

	return transport, &transport.uri, nil
}
