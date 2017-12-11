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
// Defines MAL Float type
// ################################################################################

type Float float32

var (
	NullFloat *Float = nil
)

func NewFloat(f float32) *Float {
	var val Float = Float(f)
	return &val
}

// ================================================================================
// Defines MAL Float type as a MAL Attribute

func (f *Float) attribute() Attribute {
	return f
}

// ================================================================================
// Defines MAL Float type as a MAL Element

const MAL_FLOAT_TYPE_SHORT_FORM Integer = 0x04
const MAL_FLOAT_SHORT_FORM Long = 0x1000001000004

// Returns the absolute short form of the element type.
func (*Float) GetShortForm() Long {
	return MAL_FLOAT_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*Float) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*Float) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*Float) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*Float) GetTypeShortForm() Integer {
	return MAL_FLOAT_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (f *Float) Encode(encoder Encoder) error {
	return encoder.EncodeFloat(f)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (f *Float) Decode(decoder Decoder) (Element, error) {
	return decoder.DecodeFloat()
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (f *Float) CreateElement() Element {
	return NewFloat(0)
}

func (i *Float) IsNull() bool {
	if i == nil {
		return true
	} else {
		return false
	}
}

func (*Float) Null() Element {
	return NullFloat
}
