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
package tcp

import (
	"errors"
	. "mal"
	"mal/encoding/binary"
	"strings"
)

func (transport *TCPTransport) decode(buf []byte, from string) (*Message, error) {
	decoder := binary.NewBinaryDecoder(buf, false)

	b, err := decoder.Read()
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}
	if ((b >> 5) & 0x07) != transport.version {
		return nil, errors.New("MAL/TCP Version Unknown")
	}
	sdu := b & 0x1F

	interactionType, interactionStage, err := decodeSDU(sdu)
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	serviceArea, err := decoder.DecodeUShort()
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	service, err := decoder.DecodeUShort()
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	operation, err := decoder.DecodeUShort()
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	areaVersion, err := decoder.DecodeUOctet()
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	b, err = decoder.Read()
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}
	isError := ((b >> 7) & 0x01) == binary.TRUE
	qos := (b >> 4) & 0x07
	session := b & 0xF

	transactionId, err := decoder.DecodeULong()
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	b, err = decoder.Read()
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}
	source_flag := ((b >> 7) & 0x01) == binary.TRUE
	destination_flag := ((b >> 6) & 0x01) == binary.TRUE
	priority_flag := ((b >> 5) & 0x01) == binary.TRUE
	timestamp_flag := ((b >> 4) & 0x01) == binary.TRUE
	network_zone_flag := ((b >> 3) & 0x01) == binary.TRUE
	session_name_flag := ((b >> 2) & 0x01) == binary.TRUE
	domain_flag := ((b >> 1) & 0x01) == binary.TRUE
	authentication_id_flag := (b & 0x01) == binary.TRUE

	encodingId, err := decoder.DecodeUOctet()
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	// Skips variable length field
	_, err = decoder.ReadUInt32()
	if err != nil {
		// TODO (AF): handle error
		return nil, err
	}

	var urifrom *URI = nil
	if source_flag {
		urifrom, err = decoder.DecodeURI()
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
		if !strings.HasPrefix(string(*urifrom), MALTCP) {
			var uri URI = URI(MALTCP + from + string(*urifrom))
			urifrom = &uri
		}
	} else {
		var uri URI = URI(MALTCP + from)
		urifrom = &uri
	}

	var urito *URI = nil
	if destination_flag {
		urito, err = decoder.DecodeURI()
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
		if !strings.HasPrefix(string(*urito), MALTCP) {
			var uri URI = URI(string(transport.uri) + string(*urito))
			urito = &uri
		}
	} else {
		var uri URI = transport.uri
		urito = &uri
	}

	var priority *UInteger = nil
	if priority_flag {
		priority, err = decoder.DecodeUInteger()
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
	} else {
		priority = &transport.dfltPriority
	}

	var timestamp *Time = nil
	if timestamp_flag {
		timestamp, err = decoder.DecodeTime()
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
	} else {
		timestamp = TimeNow()
	}

	var networkZone *Identifier = nil
	if network_zone_flag {
		networkZone, err = decoder.DecodeIdentifier()
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
	} else {
		networkZone = &transport.dfltNetworkZone
	}

	var sessionName *Identifier = nil
	if session_name_flag {
		sessionName, err = decoder.DecodeIdentifier()
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
	} else {
		sessionName = &transport.dfltSessionName
	}

	var domain *IdentifierList = nil
	if domain_flag {
		domain, err = DecodeIdentifierList(decoder)
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
	} else {
		domain = &transport.dfltDomain
	}

	var authenticationId *Blob = nil
	if authentication_id_flag {
		authenticationId, err = decoder.DecodeBlob()
		if err != nil {
			// TODO (AF): handle error
			return nil, err
		}
	} else {
		// Makes a copy to avoid modification of default value
		authenticationId = &transport.dfltAuthenticationId
	}

	// The remaining part of the buffer corresponds to the body part
	// of the message.
	body, err := decoder.Remaining()

	var msg *Message = &Message{
		UriFrom:          urifrom,
		UriTo:            urito,
		AuthenticationId: *authenticationId,
		EncodingId:       *encodingId,
		Timestamp:        *timestamp,
		QoSLevel:         QoSLevel(qos),
		Priority:         *priority,
		Domain:           *domain,
		NetworkZone:      *networkZone,
		Session:          SessionType(session),
		SessionName:      *sessionName,
		InteractionType:  interactionType,
		InteractionStage: interactionStage,
		TransactionId:    *transactionId,
		ServiceArea:      *serviceArea,
		Service:          *service,
		Operation:        *operation,
		AreaVersion:      *areaVersion,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}

	return msg, nil
}
