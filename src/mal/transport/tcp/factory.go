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
package tcp

import (
	. "mal"
	"net/url"
	"strconv"
)

const (
	MALTCP     string = "maltcp"
	MALTCP_URI string = "maltcp://"
)

type TCPTransportFactory struct {
}

func init() {
	RegisterTransportFactory(MALTCP, new(TCPTransportFactory))
}

func (*TCPTransportFactory) NewTransport(u *url.URL, ctx TransportCallback) (Transport, *URI, error) {
	// Builds base URI from URL
	base := url.URL{Scheme: u.Scheme, Host: u.Host}
	uri := URI(base.String())

	logger.Infof("TCPTransportFactory.TCPTransportFactory: registers %s", uri)

	// Gets parameters from URL
	params := u.Query()

	// Get the listening address. If this parameter is empty or a literal unspecified
	// IP address, the transport listens on all available unicast and anycast IP addresses
	//of the local system.
	address := u.Hostname()
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		logger.Errorf("TCPTransportFactory.NewTransport: Bad URL, cannot get listening port.")
		return nil, NULL_URI, err
	}

	transport := &TCPTransport{
		uri:     uri,
		ctx:     ctx,
		params:  params,
		address: address,
		port:    uint16(port),
	}

	err = transport.init()
	if err != nil {
		logger.Errorf("TCPTransportFactory.NewTransport: Cannot initialize transport.")
		return nil, NULL_URI, err
	}
	err = transport.start()
	if err != nil {
		logger.Errorf("TCPTransportFactory.NewTransport: Cannot start transport.")
		return nil, NULL_URI, err
	}

	return transport, &transport.uri, nil
}
