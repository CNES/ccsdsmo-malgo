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
package mal

import ()

// ################################################################################
// Defines MAL UShort type
// ################################################################################

type UShort uint16

const (
	USHORT_MIN UShort = 0
	USHORT_MAX UShort = 65535
)

var (
	NullUShort *UShort = nil
)

func NewUShort(i uint16) *UShort {
	var val UShort = UShort(i)
	return &val
}

// ================================================================================
// Defines MAL Short type as a MAL Attribute

func (s *UShort) attribute() Attribute {
	return s
}

// ================================================================================
// Defines MAL UShort type as a MAL Element

const MAL_USHORT_TYPE_SHORT_FORM Integer = 0x0A
const MAL_USHORT_SHORT_FORM Long = 0x100000100000A

// Registers MAL UShort type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_USHORT_SHORT_FORM, NullUShort)
}

// Returns the absolute short form of the element type.
func (*UShort) GetShortForm() Long {
	return MAL_USHORT_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*UShort) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*UShort) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*UShort) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*UShort) GetTypeShortForm() Integer {
	return MAL_USHORT_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (s *UShort) Encode(encoder Encoder) error {
	return encoder.EncodeUShort(s)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (s *UShort) Decode(decoder Decoder) (Element, error) {
	return decoder.DecodeUShort()
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (s *UShort) CreateElement() Element {
	return NewUShort(0)
}

func (s *UShort) IsNull() bool {
	if s == nil {
		return true
	} else {
		return false
	}
}

func (*UShort) Null() Element {
	return NullUShort
}
