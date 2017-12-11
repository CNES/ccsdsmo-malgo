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

type Pair struct {
	First  Attribute
	Second Attribute
}

var (
	NullPair *Pair = nil
)

func NewPair() *Pair {
	return new(Pair)
}

// ================================================================================
// Defines MAL Pair type as a MAL Composite

func (pair *Pair) composite() Composite {
	return pair
}

// ================================================================================
// Defines MAL Pair type as a MAL Element

const MAL_PAIR_TYPE_SHORT_FORM Integer = 0x1C
const MAL_PAIR_SHORT_FORM Long = 0x100000100001C

// Returns the absolute short form of the element type.
func (*Pair) GetShortForm() Long {
	return MAL_PAIR_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*Pair) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*Pair) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*Pair) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*Pair) GetTypeShortForm() Integer {
	return MAL_PAIR_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (pair *Pair) Encode(encoder Encoder) error {
	err := encoder.EncodeNullableAttribute(pair.First)
	if err != nil {
		return err
	}
	return encoder.EncodeNullableAttribute(pair.Second)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (pair *Pair) Decode(decoder Decoder) (Element, error) {
	return DecodePair(decoder)
}

// Decodes an instance of Pair using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded EntityRequest instance.
func DecodePair(decoder Decoder) (*Pair, error) {
	first, err := decoder.DecodeNullableAttribute()
	if err != nil {
		return nil, err
	}
	second, err := decoder.DecodeNullableAttribute()
	if err != nil {
		return nil, err
	}
	var pair = Pair{
		first, second,
	}
	return &pair, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (pair *Pair) CreateElement() Element {
	// TODO (AF):
	//	return new(Pair)
	return NewPair()
}

func (pair *Pair) IsNull() bool {
	return pair == nil
}

func (*Pair) Null() Element {
	return NullPair
}
