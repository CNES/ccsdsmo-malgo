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
	. "github.com/CNES/ccsdsmo-malgo/mal"
	"github.com/CNES/ccsdsmo-malgo/mal/encoding/binary"
)

func (transport *TCPTransport) encode(msg *Message) ([]byte, error) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	encoder := binary.NewBinaryEncoder(buf, false)

	sdu, err := encodeSDU(msg.InteractionType, msg.InteractionStage)
	if err != nil {
		logger.Errorf("TCPTransport.encode, cannot encode SDU: %s", err.Error())
		return nil, err
	}
	err = encoder.Write(sdu | (transport.version << 5))
	if err != nil {
		logger.Errorf("TCPTransport.encode, cannot write SDU: %s", err.Error())
		return nil, err
	}

	err = encoder.EncodeUShort(&msg.ServiceArea)
	if err != nil {
		logger.Errorf("TCPTransport.encode, cannot encode ServiceArea: %s", err.Error())
		return nil, err
	}

	err = encoder.EncodeUShort(&msg.Service)
	if err != nil {
		logger.Errorf("TCPTransport.encode, cannot encode Service: %s", err.Error())
		return nil, err
	}

	err = encoder.EncodeUShort(&msg.Operation)
	if err != nil {
		logger.Errorf("TCPTransport.encode, cannot encode Operation: %s", err.Error())
		return nil, err
	}

	err = encoder.EncodeUOctet(&msg.AreaVersion)
	if err != nil {
		logger.Errorf("TCPTransport.encode, cannot encode AreaVersion: %s", err.Error())
		return nil, err
	}

	var b byte = 0
	if msg.IsErrorMessage {
		b = 0x80
	}
	b |= (byte(msg.QoSLevel) << 4)
	b |= byte(msg.Session)
	err = encoder.Write(b)
	if err != nil {
		logger.Errorf("TCPTransport.encode, cannot write Flags: %s", err.Error())
		return nil, err
	}

	err = encoder.EncodeULong(&msg.TransactionId)
	if err != nil {
		logger.Errorf("TCPTransport.encode, cannot encode TransactionId: %s", err.Error())
		return nil, err
	}

	err = encoder.Write(transport.flags)
	if err != nil {
		logger.Errorf("TCPTransport.encode, cannot write Transport flags: %s", err.Error())
		return nil, err
	}

	err = encoder.EncodeUOctet(&msg.EncodingId)
	if err != nil {
		logger.Errorf("TCPTransport.encode, cannot encode EncodingId: %s", err.Error())
		return nil, err
	}

	// Skips variable length field
	err = encoder.WriteUInt32(0)
	if err != nil {
		logger.Errorf("TCPTransport.encode, cannot skip variable length: %s", err.Error())
		return nil, err
	}

	// Remaining data are now encoded in PDU using varint.
	encoder.Varint = true

	if transport.sourceFlag {
		if transport.optimizeURI {
			// Optimized mapping, writes only URI service part
			err = encoder.EncodeString(msg.UriFrom.GetService())
			if err != nil {
				logger.Errorf("TCPTransport.encode, cannot encode URIFrom: %s", err.Error())
				return nil, err
			}
		} else {
			err = encoder.EncodeURI(msg.UriFrom)
			if err != nil {
				logger.Errorf("TCPTransport.encode, cannot encode URIFrom: %s", err.Error())
				return nil, err
			}
		}
	}

	if transport.destinationFlag {
		if transport.optimizeURI {
			// Optimized mapping, writes only URI service part
			err = encoder.EncodeString(msg.UriTo.GetService())
			if err != nil {
				logger.Errorf("TCPTransport.encode, cannot encode URITo: %s", err.Error())
				return nil, err
			}
		} else {
			err = encoder.EncodeURI(msg.UriTo)
			if err != nil {
				logger.Errorf("TCPTransport.encode, cannot encode URITo: %s", err.Error())
				return nil, err
			}
		}
	}

	if transport.priorityFlag {
		err = encoder.EncodeUInteger(&msg.Priority)
		if err != nil {
			logger.Errorf("TCPTransport.encode, cannot encode Priority: %s", err.Error())
			return nil, err
		}
	}

	if transport.timestampFlag {
		err = encoder.EncodeTime(&msg.Timestamp)
		if err != nil {
			logger.Errorf("TCPTransport.encode, cannot encode TimeStamp: %s", err.Error())
			return nil, err
		}
	}

	if transport.networkZoneFlag {
		err = encoder.EncodeIdentifier(&msg.NetworkZone)
		if err != nil {
			logger.Errorf("TCPTransport.encode, cannot encode NetworkZone: %s", err.Error())
			return nil, err
		}
	}

	if transport.sessionNameFlag {
		err = encoder.EncodeIdentifier(&msg.SessionName)
		if err != nil {
			logger.Errorf("TCPTransport.encode, cannot encode SessionName: %s", err.Error())
			return nil, err
		}
	}

	if transport.domainFlag {
		err = msg.Domain.Encode(encoder)
		if err != nil {
			logger.Errorf("TCPTransport.encode, cannot encode Domain: %s", err.Error())
			return nil, err
		}
	}

	if transport.authenticationIdFlag {
		err = encoder.EncodeBlob(&msg.AuthenticationId)
		if err != nil {
			logger.Errorf("TCPTransport.encode, cannot encode AuthenticationId: %s", err.Error())
			return nil, err
		}
	}

	err = encoder.WriteBody(msg.Body)
	if err != nil {
		logger.Errorf("TCPTransport.encode, cannot write body: %s", err.Error())
		return nil, err
	}

	return encoder.Out.(*binary.BinaryBuffer).Buf, nil
}
