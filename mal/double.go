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
// Defines MAL Double type
// ################################################################################

type Double float64

var (
	NullDouble *Double = nil
)

func NewDouble(d float64) *Double {
	var val Double = Double(d)
	return &val
}

// ================================================================================
// Defines MAL Double type as a MAL Attribute

func (d *Double) attribute() Attribute {
	return d
}

// ================================================================================
// Defines MAL Double type as a MAL Element

const MAL_DOUBLE_TYPE_SHORT_FORM Integer = 0x05
const MAL_DOUBLE_SHORT_FORM Long = 0x1000001000005

// Registers MAL Double type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_DOUBLE_SHORT_FORM, NullDouble)
}

// Returns the absolute short form of the element type.
func (*Double) GetShortForm() Long {
	return MAL_DOUBLE_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*Double) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*Double) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*Double) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*Double) GetTypeShortForm() Integer {
	return MAL_DOUBLE_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (d *Double) Encode(encoder Encoder) error {
	return encoder.EncodeDouble(d)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (d *Double) Decode(decoder Decoder) (Element, error) {
	return decoder.DecodeDouble()
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (d *Double) CreateElement() Element {
	return NewDouble(0)
}

func (i *Double) IsNull() bool {
	if i == nil {
		return true
	} else {
		return false
	}
}

func (*Double) Null() Element {
	return NullDouble
}
