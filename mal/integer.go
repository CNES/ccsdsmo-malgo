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

// ################################################################################
// Defines MAL Integer type
// ################################################################################

type Integer int32

const (
	INTEGER_MIN Integer = -2147483648
	INTEGER_MAX Integer = 2147483647
)

var (
	NullInteger *Integer = nil
)

func NewInteger(i int32) *Integer {
	var val Integer = Integer(i)
	return &val
}

// ================================================================================
// Defines MAL Integer type as a MAL Attribute

func (i *Integer) attribute() Attribute {
	return i
}

// ================================================================================
// Defines MAL Integer type as a MAL Element

const MAL_INTEGER_TYPE_SHORT_FORM Integer = 0x0B
const MAL_INTEGER_SHORT_FORM Long = 0x100000100000B

// Registers MAL Integer type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_INTEGER_SHORT_FORM, NullInteger)
}

// Returns the absolute short form of the element type.
func (*Integer) GetShortForm() Long {
	return MAL_INTEGER_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*Integer) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*Integer) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*Integer) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*Integer) GetTypeShortForm() Integer {
	return MAL_INTEGER_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (i *Integer) Encode(encoder Encoder) error {
	return encoder.EncodeInteger(i)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (i *Integer) Decode(decoder Decoder) (Element, error) {
	return decoder.DecodeInteger()
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (i *Integer) CreateElement() Element {
	return NewInteger(0)
}

func (i *Integer) IsNull() bool {
	if i == nil {
		return true
	} else {
		return false
	}
}

func (*Integer) Null() Element {
	return NullInteger
}
