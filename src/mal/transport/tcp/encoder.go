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
	"mal/encoding/binary"
)

func (transport *TCPTransport) encode(msg *Message) ([]byte, error) {
	// TODO (AF): calculates the size of the encoded message
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	encoder := binary.NewBinaryEncoder(buf, false)

	sdu, err := encodeSDU(msg.InteractionType, msg.InteractionStage)
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}
	err = encoder.Write(sdu | (transport.version << 5))
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	err = encoder.EncodeUShort(&msg.ServiceArea)
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	err = encoder.EncodeUShort(&msg.Service)
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	err = encoder.EncodeUShort(&msg.Operation)
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	err = encoder.EncodeUOctet(&msg.AreaVersion)
	if err != nil {
		// TODO (AF): handle error
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
		// TODO (AF): handle error
		return nil, err
	}

	err = encoder.EncodeULong(&msg.TransactionId)
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	err = encoder.Write(transport.flags)
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	err = encoder.EncodeUOctet(&msg.EncodingId)
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	// Skips variable length field
	err = encoder.WriteUInt32(0)
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	if transport.sourceFlag {
		err = encoder.EncodeURI(msg.UriFrom)
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
	} else {
		// TODO (AF): Optimized mapping, writes only URI service part
	}

	if transport.destinatioFlag {
		err = encoder.EncodeURI(msg.UriTo)
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
	} else {
		// TODO (AF): Optimized mapping, writes only URI service part
	}

	if transport.priorityFlag {
		err = encoder.EncodeUInteger(&msg.Priority)
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
	}

	if transport.timestampFlag {
		err = encoder.EncodeTime(&msg.Timestamp)
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
	}

	if transport.networkZoneFlag {
		err = encoder.EncodeIdentifier(&msg.NetworkZone)
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
	}

	if transport.sessionNameFlag {
		err = encoder.EncodeIdentifier(&msg.SessionName)
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
	}

	if transport.domainFlag {
		err = msg.Domain.Encode(encoder)
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
	}

	if transport.authenticationIdFlag {
		err = encoder.EncodeBlob(&msg.AuthenticationId)
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
	}

	err = encoder.WriteBody(msg.Body)

	return encoder.Body(), nil
}
