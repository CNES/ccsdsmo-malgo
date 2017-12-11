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
// Defines MAL UInteger type
// ################################################################################

type UInteger uint32

const (
	UINTEGER_MIN UInteger = 0
	UINTEGER_MAX UInteger = 4294967295
)

var (
	NullUInteger *UInteger = nil
)

func NewUInteger(i uint32) *UInteger {
	var val UInteger = UInteger(i)
	return &val
}

// ================================================================================
// Defines MAL UInteger type as a MAL Attribute

func (i *UInteger) attribute() Attribute {
	return i
}

// ================================================================================
// Defines MAL UInteger type as a MAL Element

const MAL_UINTEGER_TYPE_SHORT_FORM Integer = 0x0C
const MAL_UINTEGER_SHORT_FORM Long = 0x100000100000C

// Returns the absolute short form of the element type.
func (*UInteger) GetShortForm() Long {
	return MAL_UINTEGER_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*UInteger) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*UInteger) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*UInteger) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*UInteger) GetTypeShortForm() Integer {
	return MAL_UINTEGER_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (i *UInteger) Encode(encoder Encoder) error {
	return encoder.EncodeUInteger(i)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (i *UInteger) Decode(decoder Decoder) (Element, error) {
	return decoder.DecodeUInteger()
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (i *UInteger) CreateElement() Element {
	return NewUInteger(0)
}

func (i *UInteger) IsNull() bool {
	if i == nil {
		return true
	} else {
		return false
	}
}

func (*UInteger) Null() Element {
	return NullUInteger
}
