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
// Defines MAL NamedValue type
// ################################################################################

type NamedValue struct {
	Name  *Identifier
	Value Attribute
}

var (
	NullNamedValue *NamedValue = nil
)

func NewNamedValue() *NamedValue {
	return new(NamedValue)
}

// ================================================================================
// Defines MAL NamedValue type as a MAL Composite

func (pair *NamedValue) Composite() Composite {
	return pair
}

// ================================================================================
// Defines MAL NamedValue type as a MAL Element

const MAL_NAMED_VALUE_TYPE_SHORT_FORM Integer = 0x1D
const MAL_NAMED_VALUE_SHORT_FORM Long = 0x100000100001D

// Registers MAL NamedValue type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_NAMED_VALUE_SHORT_FORM, NullNamedValue)
}

// Returns the absolute short form of the element type.
func (*NamedValue) GetShortForm() Long {
	return MAL_NAMED_VALUE_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*NamedValue) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*NamedValue) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*NamedValue) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*NamedValue) GetTypeShortForm() Integer {
	return MAL_NAMED_VALUE_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (pair *NamedValue) Encode(encoder Encoder) error {
	err := encoder.EncodeNullableIdentifier(pair.Name)
	if err != nil {
		return err
	}
	return encoder.EncodeNullableAttribute(pair.Value)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (pair *NamedValue) Decode(decoder Decoder) (Element, error) {
	return DecodeNamedValue(decoder)
}

// Decodes an instance of IdBooleanPair using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded EntityRequest instance.
func DecodeNamedValue(decoder Decoder) (*NamedValue, error) {
	name, err := decoder.DecodeNullableIdentifier()
	if err != nil {
		return nil, err
	}
	value, err := decoder.DecodeAttribute()
	if err != nil {
		return nil, err
	}
	var pair = NamedValue{
		name, value,
	}
	return &pair, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (pair *NamedValue) CreateElement() Element {
	// TODO (AF):
	//	return new(NamedValue)
	return NewNamedValue()
}

func (pair *NamedValue) IsNull() bool {
	return pair == nil
}

func (*NamedValue) Null() Element {
	return NullNamedValue
}
