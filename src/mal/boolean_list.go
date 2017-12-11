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
// Defines MAL BooleanList type
// ################################################################################

type BooleanList []*Boolean

var (
	NullBooleanList *BooleanList = nil
)

func NewBooleanList(size int) *BooleanList {
	var list BooleanList = BooleanList(make([]*Boolean, size))
	return &list
}

// ================================================================================
// Defines MAL BooleanList type as a MAL Element

const MAL_BOOLEAN_LIST_TYPE_SHORT_FORM Integer = -0x02
const MAL_BOOLEAN_LIST_SHORT_FORM Long = 0x1000001FFFFFE

// Returns the absolute short form of the element type.
func (*BooleanList) GetShortForm() Long {
	return MAL_BOOLEAN_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*BooleanList) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*BooleanList) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*BooleanList) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*BooleanList) GetTypeShortForm() Integer {
	//	return MAL_BOOLEAN_TYPE_SHORT_FORM & 0x01FFFF00
	return MAL_BOOLEAN_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *BooleanList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*Boolean(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*Boolean(*list) {
		encoder.EncodeNullableBoolean(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *BooleanList) Decode(decoder Decoder) (Element, error) {
	return DecodeBooleanList(decoder)
}

// Decodes an instance of BooleanList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded BooleanList instance.
func DecodeBooleanList(decoder Decoder) (*BooleanList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := BooleanList(make([]*Boolean, int(*size)))
	for i := 0; i < len(list); i++ {
		list[i], err = decoder.DecodeNullableBoolean()
		if err != nil {
			return nil, err
		}
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *BooleanList) CreateElement() Element {
	return NewBooleanList(0)
}

func (list *BooleanList) IsNull() bool {
	return list == nil
}

func (*BooleanList) Null() Element {
	return NullBooleanList
}
