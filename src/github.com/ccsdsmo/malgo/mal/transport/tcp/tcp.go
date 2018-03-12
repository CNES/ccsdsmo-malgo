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
	. "github.com/ccsdsmo/malgo/mal"
	"github.com/ccsdsmo/malgo/mal/debug"
	"io"
	"net"
	"net/url"
	"strconv"
)

const (
	// Name of preoperty allowing to fix the underlying protocol: currently tcp, tcp4 or tcp6.
	// By default use tcp.
	NETWORK_PROPERTY string = "network"

	VARIABLE_LENGTH_OFFSET uint32 = 19
	FIXED_HEADER_LENGTH    uint32 = 23
)

var (
	logger debug.Logger = debug.GetLogger("mal.transport.tcp")
)

type TCPTransport struct {
	uri    URI
	ctx    TransportCallback
	params map[string][]string

	version byte

	network string
	address string
	port    uint16

	running bool

	// Channel for outgoing messages.
	ch   chan *Message
	ends chan bool

	listen net.Listener
	conns  map[string]net.Conn

	optimizeURI bool

	sourceFlag           bool
	sourceId             *URI
	destinationFlag      bool
	destinationId        *URI
	priorityFlag         bool
	priority             UInteger
	timestampFlag        bool
	networkZoneFlag      bool
	networkZone          Identifier
	sessionNameFlag      bool
	sessionName          Identifier
	domainFlag           bool
	domain               IdentifierList
	authenticationIdFlag bool
	authenticationId     Blob

	flags byte

	dfltPriority         UInteger
	dfltNetworkZone      Identifier
	dfltSessionName      Identifier
	dfltAuthenticationId Blob
	dfltDomain           IdentifierList
}

// Initializes the MAL/TCP context.
func (transport *TCPTransport) init() error {
	transport.running = false

	transport.version = 1

	// TODO (AF): Configure flags
	transport.flags = 0
	// Note: Should be always true
	transport.sourceFlag = true
	if transport.sourceFlag {
		transport.flags |= (1 << 7)
	}
	// Note: Should be always true
	transport.destinationFlag = true
	if transport.destinationFlag {
		transport.flags |= (1 << 6)
	}
	transport.priorityFlag = true
	if transport.priorityFlag {
		transport.flags |= (1 << 5)
	}
	transport.timestampFlag = true
	if transport.timestampFlag {
		transport.flags |= (1 << 4)
	}
	transport.networkZoneFlag = true
	if transport.networkZoneFlag {
		transport.flags |= (1 << 3)
	}
	transport.sessionNameFlag = true
	if transport.sessionNameFlag {
		transport.flags |= (1 << 2)
	}
	transport.domainFlag = true
	if transport.domainFlag {
		transport.flags |= (1 << 1)
	}
	transport.authenticationIdFlag = true
	if transport.authenticationIdFlag {
		transport.flags |= 1
	}

	// Get protocol: tcp, tcp4 or tcp6.
	if p := transport.params[NETWORK_PROPERTY]; p != nil {
		transport.network = p[0]
	} else {
		transport.network = "tcp"
	}

	transport.conns = make(map[string]net.Conn)
	// TODO (AF): Fix length of channel
	transport.ch = make(chan *Message, 10)
	transport.ends = make(chan bool)

	return nil
}

// Starts the MAL/TCP context.
func (transport *TCPTransport) start() error {
	// If the host in the address parameter is empty or a literal unspecified IP address,
	// Listen listens on all available unicast and anycast IP addresses of the local system.
	// To only use IPv4, use "tcp4" a network parameter.
	listen, err := net.Listen(transport.network, ":"+strconv.Itoa(int(transport.port)))
	if err != nil {
		logger.Errorf("TCPTransport.start, cannot create listen socket: %s", err.Error())
		return err
	}

	transport.running = true

	transport.listen = listen
	go transport.handleConn(listen)
	// Note: May be we have to create multiples threads to handle outgoing messages.
	go transport.handleOut()

	return nil
}

func (transport *TCPTransport) handleConn(listen net.Listener) {
	for {
		cnx, err := listen.Accept()
		if err != nil {
			// If closing don't log an error.
			if transport.running {
				logger.Errorf("TCPTransport.handleConn, error accepting connection: %s", err.Error())
			}
			break
		}
		logger.Infof("TCPTransport.handleConn, accept connexion from %s", cnx.RemoteAddr())
		go transport.handleIn(cnx)
	}
	logger.Infof("TCPTransport.HandleConn exited")
}

// A utility function which tests if an error returned from a TCPConnection or
// TCPListener is actually an EOF. In some edge cases this which should be treated
// as EOFs are not returned as one.
func isEOF(err error) bool {
	if err == io.EOF {
		return true
	} else if oerr, ok := err.(*net.OpError); ok {
		/* this hack happens because the error is returned when the
		 * network socket is closing and instead of returning a
		 * io.EOF it returns this error.New(...) struct. */
		if oerr.Err.Error() == "use of closed network connection" {
			return true
		}
	} else {
		if err.Error() == "use of closed network connection" {
			return true
		}
	}
	return false
}

func (transport *TCPTransport) handleIn(cnx net.Conn) {
	// Registers the new connection
	uri := cnx.RemoteAddr().String()
	transport.conns[uri] = cnx

	for transport.running {
		logger.Debugf("TCPTransport.HandleIn(%s), wait for message.", cnx.RemoteAddr())
		msg, err := transport.readMessage(cnx)

		if err != nil {
			if isEOF(err) {
				break
			} else {
				continue
			}
		}
		logger.Debugf("TCPTransport.HandleIn(%s), receives message: %s", cnx.RemoteAddr(), msg)
		if msg != nil {
			transport.ctx.Receive(msg)
		}
	}
	// Closes the connection
	cnx.Close()
	// Removes connection from transport.conns
	delete(transport.conns, uri)
	logger.Infof("TCPTransport.HandleIn(%s) exited: %s", cnx.RemoteAddr(), cnx.RemoteAddr())
}

func (transport *TCPTransport) readMessage(cnx net.Conn) (*Message, error) {
	var buf []byte = make([]byte, FIXED_HEADER_LENGTH)

	// Reads the fixed part of MAL message header
	for offset := 0; offset < int(FIXED_HEADER_LENGTH); {
		logger.Debugf("TCPTransport.readMessage(%s), waiting message: ..", cnx.RemoteAddr())
		nb, err := cnx.Read(buf[offset:])
		if err != nil {
			if !isEOF(err) {
				logger.Errorf("TCPTransport.readMessage(%s), error reading fixed part of message: %s", cnx.RemoteAddr(), err.Error())
			} else {
				logger.Warnf("TCPTransport.readMessage(%s), connection closed", cnx.RemoteAddr())
			}
			return nil, err
		}
		offset += nb
	}

	// Get the variable length of message
	length := FIXED_HEADER_LENGTH +
		uint32(buf[VARIABLE_LENGTH_OFFSET+3]) | uint32(buf[VARIABLE_LENGTH_OFFSET+2])<<8 |
		uint32(buf[VARIABLE_LENGTH_OFFSET+1])<<16 | uint32(buf[VARIABLE_LENGTH_OFFSET])<<24
	logger.Debugf("Reads message header, length: %d", length)

	// Allocate a new buffer and copy the fixed part of MAL message header
	var newbuf []byte = make([]byte, length)
	copy(newbuf, buf)

	// Reads fully the message
	for offset := int(FIXED_HEADER_LENGTH); offset < len(newbuf); {
		nb, err := cnx.Read(newbuf[offset:])
		if err != nil {
			logger.Errorf("TCPTransport.readMessage(%s), error reading message: %s", cnx.RemoteAddr(), err.Error())
			return nil, err
		}
		offset += nb
		logger.Debugf("Reads: %d", offset)
	}

	// Decodes the message
	msg, err := transport.decode(newbuf, cnx.RemoteAddr().String())
	if err != nil {
		logger.Errorf("TCPTransport.readMessage(%s), error receiving message: %s", err.Error())
		return nil, err
	}
	logger.Debugf("TCPTransport.readMessage(%s), receives %s from %s to %s", cnx.RemoteAddr(), msg, *msg.UriFrom, *msg.UriTo)

	return msg, nil
}

func (transport *TCPTransport) handleOut() {
	var msg *Message
	var nbtry uint
	for {
		logger.Debugf("TCPTransport.handleOut, wait message..")
		if msg == nil {
			msg, _ = <-transport.ch
			nbtry = 0
		}
		if msg != nil {
			logger.Debugf("TCPTransport.handleOut, get Message %+v", *msg)
			u, err := url.Parse(string(*msg.UriTo))
			if err != nil {
				logger.Errorf("TCPTransport.handleOut, cannot route message to %s", *msg.UriTo)
				continue
			}
			urito := u.Host

			cnx, ok := transport.conns[urito]
			if !ok {
				logger.Debugf("TCPTransport.handleOut, creates connection to %s", urito)
				cnx, err = net.Dial("tcp", urito)
				if err != nil {
					logger.Errorf("TCPTransport.handleOut, cannot creates connection to %s: %s", urito, err.Error())
					// TODO (AF): Handles the faulty message, forwards it to error listener
					continue
				}
				// Registers the created connection
				transport.conns[urito] = cnx
				go transport.handleIn(cnx)
			}
			logger.Debugf("TCPTransport.handleOut, send message to %s", *msg.UriTo)
			err = transport.writeMessage(cnx, msg)
			if err != nil {
				nbtry += 1
				logger.Debugf("TCPTransport.handleOut, error sending message: %s", err.Error())
				// Closes the connection to retrieve a clean state
				cnx.Close()
				// Removes connection from transport.conns
				delete(transport.conns, urito)
				// try to send anew the message
				if nbtry < 3 {
					continue
				} else {
					// TODO (AF): Handles the faulty message, forwards it to error listener
				}
			}
			msg = nil
		} else {
			logger.Infof("TCPTransport.handleOut, ends")
			transport.ends <- true
		}
	}
	logger.Debugf("TCPTransport.handleOut exited")
}

func write32(value uint32, buf []byte) {
	buf[0] = byte(value >> 24)
	buf[1] = byte(value >> 16)
	buf[2] = byte(value >> 8)
	buf[3] = byte(value >> 0)
}

func (transport *TCPTransport) writeMessage(cnx net.Conn, msg *Message) error {
	buf, err := transport.encode(msg)
	if err != nil {
		// TODO (AF): Logging
		return err
	}
	logger.Debugf("Writes message: %d", len(buf))
	write32(uint32(len(buf))-FIXED_HEADER_LENGTH, buf[VARIABLE_LENGTH_OFFSET:VARIABLE_LENGTH_OFFSET+4])
	logger.Debugf("Message transmitted: ", buf)
	_, err = cnx.Write(buf)
	if err != nil {
		logger.Errorf("Transport.writeMessage, cannot send to %s", cnx.RemoteAddr())
		return err
	}
	return nil
}

func (transport *TCPTransport) Transmit(msg *Message) error {
	logger.Debugf("Transmit: %+v", *msg)
	transport.ch <- msg
	logger.Debugf("Transmited")
	return nil
}

func (transport *TCPTransport) TransmitMultiple(msgs ...*Message) error {
	for _, msg := range msgs {
		err := transport.Transmit(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (transport *TCPTransport) Close() error {
	transport.running = false
	close(transport.ch)
	transport.listen.Close()
	for id, cnx := range transport.conns {
		logger.Debugf("Transport.Close, close connection: %s", id)
		cnx.Close()
	}
	// TODO (AF):
	return nil
}
