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
package mal

import ()

// ################################################################################
// Defines MAL UOctet type
// ################################################################################

type UOctet uint8

const (
	UOCTET_MIN UOctet = 0
	UOCTET_MAX UOctet = 255
)

var (
	NullUOctet *UOctet = nil
)

func NewUOctet(i uint8) *UOctet {
	var val UOctet = UOctet(i)
	return &val
}

// ================================================================================
// Defines MAL UOctet type as a MAL Attribute

func (o *UOctet) attribute() Attribute {
	return o
}

// ================================================================================
// Defines MAL UOctet type as a MAL Element

const MAL_UOCTET_TYPE_SHORT_FORM Integer = 0x08
const MAL_UOCTET_SHORT_FORM Long = 0x1000001000008

// Returns the absolute short form of the element type.
func (*UOctet) GetShortForm() Long {
	return MAL_UOCTET_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*UOctet) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*UOctet) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*UOctet) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*UOctet) GetTypeShortForm() Integer {
	return MAL_UOCTET_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (o *UOctet) Encode(encoder Encoder) error {
	return encoder.EncodeUOctet(o)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (o *UOctet) Decode(decoder Decoder) (Element, error) {
	return decoder.DecodeUOctet()
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (o *UOctet) CreateElement() Element {
	return NewUOctet(0)
}

func (o *UOctet) IsNull() bool {
	if o == nil {
		return true
	} else {
		return false
	}
}

func (*UOctet) Null() Element {
	return NullUOctet
}
