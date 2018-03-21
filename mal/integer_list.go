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
// Defines MAL IntegerList type
// ################################################################################

type IntegerList []*Integer

var (
	NullIntegerList *IntegerList = nil
)

func NewIntegerList(size int) *IntegerList {
	var list IntegerList = IntegerList(make([]*Integer, size))
	return &list
}

// ================================================================================
// Defines MAL IntegerList type as an ElementList

func (list *IntegerList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

// ================================================================================
// Defines MAL IntegerList type as a MAL Element

const MAL_INTEGER_LIST_TYPE_SHORT_FORM Integer = -0x0B
const MAL_INTEGER_LIST_SHORT_FORM Long = 0x1000001FFFFF5

// Registers MAL IntegerList type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_INTEGER_LIST_SHORT_FORM, NullIntegerList)
}

// Returns the absolute short form of the element type.
func (*IntegerList) GetShortForm() Long {
	return MAL_INTEGER_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*IntegerList) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*IntegerList) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*IntegerList) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*IntegerList) GetTypeShortForm() Integer {
	//	return MAL_INTEGER_TYPE_SHORT_FORM & 0x01FFFF00
	return MAL_INTEGER_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *IntegerList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*Integer(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*Integer(*list) {
		encoder.EncodeNullableInteger(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *IntegerList) Decode(decoder Decoder) (Element, error) {
	return DecodeIntegerList(decoder)
}

// Decodes an instance of IntegerList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded IntegerList instance.
func DecodeIntegerList(decoder Decoder) (*IntegerList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := IntegerList(make([]*Integer, int(*size)))
	for i := 0; i < len(list); i++ {
		list[i], err = decoder.DecodeNullableInteger()
		if err != nil {
			return nil, err
		}
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *IntegerList) CreateElement() Element {
	return NewIntegerList(0)
}

func (list *IntegerList) IsNull() bool {
	return list == nil
}

func (*IntegerList) Null() Element {
	return NullIntegerList
}
