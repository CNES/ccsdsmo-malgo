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
// Defines MAL String type
// ################################################################################

type String string

var (
	NullString *String = nil
)

func NewString(s string) *String {
	var val String = String(s)
	return &val
}

// ================================================================================
// Defines MAL Short String as a MAL Attribute

func (s *String) attribute() Attribute {
	return s
}

// ================================================================================
// Defines MAL String type as a MAL Element

const MAL_STRING_TYPE_SHORT_FORM Integer = 0x0F
const MAL_STRING_SHORT_FORM Long = 0x100000100000F

// Registers MAL String type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_STRING_SHORT_FORM, NullString)
}

// Returns the absolute short form of the element type.
func (*String) GetShortForm() Long {
	return MAL_STRING_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*String) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*String) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*String) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*String) GetTypeShortForm() Integer {
	return MAL_STRING_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (str *String) Encode(encoder Encoder) error {
	return encoder.EncodeString(str)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (str *String) Decode(decoder Decoder) (Element, error) {
	return decoder.DecodeString()
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (str *String) CreateElement() Element {
	return NewString("")
}

func (str *String) IsNull() bool {
	if str == nil {
		return true
	} else {
		return false
	}
}

func (*String) Null() Element {
	return NullString
}
