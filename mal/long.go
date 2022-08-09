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
// Defines MAL Long type
// ################################################################################

type Long int64

const (
	LONG_MIN Long = -9223372036854775808
	LONG_MAX Long = 9223372036854775807
)

var (
	NullLong *Long = nil
)

func NewLong(i int64) *Long {
	var val Long = Long(i)
	return &val
}

// ================================================================================
// Defines MAL Long type as a MAL Attribute

func (l *Long) attribute() Attribute {
	return l
}

// ================================================================================
// Defines MAL Long type as a MAL Element

const MAL_LONG_TYPE_SHORT_FORM Integer = 0x0D
const MAL_LONG_SHORT_FORM Long = 0x100000100000D

// Registers MAL Long type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_LONG_SHORT_FORM, NullLong)
}

// Returns the absolute short form of the element type.
func (*Long) GetShortForm() Long {
	return MAL_LONG_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*Long) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*Long) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*Long) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*Long) GetTypeShortForm() Integer {
	return MAL_LONG_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (l *Long) Encode(encoder Encoder) error {
	return encoder.EncodeLong(l)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (l *Long) Decode(decoder Decoder) (Element, error) {
	return decoder.DecodeLong()
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (l *Long) CreateElement() Element {
	return NewLong(0)
}

func (l *Long) IsNull() bool {
	if l == nil {
		return true
	} else {
		return false
	}
}

func (*Long) Null() Element {
	return NullLong
}
