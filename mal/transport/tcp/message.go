/**
 * MIT License
 *
 * Copyright (c) 2019 CNES
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

type TCPBody struct {
	factory EncodingFactory
	encoder Encoder
	decoder Decoder
	content []byte
}

func NewTCPBody(buf []byte, writeable bool) *TCPBody {
	body := new(TCPBody)
	body.factory = binary.FixedBinaryEncodingFactory
	body.content = buf
	body.Reset(writeable)
	return body
}

func (body *TCPBody) getEncodedContent() []byte {
	if body.encoder == nil {
		// TODO (AF): Needed for body built directly from []byte. Useful ?
		// Normally body.content is always equal to body.encoder.Body()
		return body.content
	} else {
		return body.encoder.Body()
	}
}

func (body *TCPBody) Reset(writeable bool) {
	if writeable {
		body.decoder = nil
		body.encoder = body.factory.NewEncoder(body.content)
	} else {
		body.decoder = body.factory.NewDecoder(body.content)
		body.encoder = nil
	}
}

func (body *TCPBody) SetEncodingFactory(factory EncodingFactory) {
	body.factory = factory
}

func (body *TCPBody) DecodeParameter(element Element) (Element, error) {
	return body.decoder.DecodeNullableElement(element)
}

func (body *TCPBody) DecodeLastParameter(element Element, abstract bool) (Element, error) {
	if abstract {
		return body.decoder.DecodeNullableAbstractElement()
	} else {
		return body.decoder.DecodeNullableElement(element)
	}
}

func (body *TCPBody) EncodeParameter(element Element) error {
	return body.encoder.EncodeNullableElement(element)
}

func (body *TCPBody) EncodeLastParameter(element Element, abstract bool) error {
	if abstract {
		return body.encoder.EncodeNullableAbstractElement(element)
	} else {
		return body.encoder.EncodeNullableElement(element)
	}
}
