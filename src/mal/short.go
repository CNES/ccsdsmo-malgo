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
// Defines MAL Short type
// ################################################################################

type Short int16

const (
	SHORT_MIN Short = -32768
	SHORT_MAX Short = 32767
)

var (
	NullShort *Short = nil
)

func NewShort(i int16) *Short {
	var val Short = Short(i)
	return &val
}

// ================================================================================
// Defines MAL Short type as a MAL Attribute

func (s *Short) attribute() Attribute {
	return s
}

// ================================================================================
// Defines MAL Short type as a MAL Element

const MAL_SHORT_TYPE_SHORT_FORM Integer = 0x09
const MAL_SHORT_SHORT_FORM Long = 0x1000001000009

// Returns the absolute short form of the element type.
func (*Short) GetShortForm() Long {
	return MAL_SHORT_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*Short) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*Short) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*Short) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*Short) GetTypeShortForm() Integer {
	return MAL_SHORT_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (s *Short) Encode(encoder Encoder) error {
	return encoder.EncodeShort(s)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (s *Short) Decode(decoder Decoder) (Element, error) {
	return decoder.DecodeShort()
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (s *Short) CreateElement() Element {
	return NewShort(0)
}

func (s *Short) IsNull() bool {
	if s == nil {
		return true
	} else {
		return false
	}
}

func (*Short) Null() Element {
	return NullShort
}
