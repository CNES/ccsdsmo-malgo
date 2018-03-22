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
// Defines MAL UIntegerList type
// ################################################################################

type UIntegerList []*UInteger

var (
	NullUIntegerList *UIntegerList = nil
)

func NewUIntegerList(size int) *UIntegerList {
	var list UIntegerList = UIntegerList(make([]*UInteger, size))
	return &list
}

// ================================================================================
// Defines MAL UIntegerList type as an ElementList

func (list *UIntegerList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

// ================================================================================
// Defines MAL UIntegerList type as a MAL Composite

func (list *UIntegerList) Composite() Composite {
	return list
}

// ================================================================================
// Defines MAL UIntegerList type as a MAL Element

const MAL_UINTEGER_LIST_TYPE_SHORT_FORM Integer = -0x0C
const MAL_UINTEGER_LIST_SHORT_FORM Long = 0x1000001FFFFF4

// Registers MAL UOctet type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_UINTEGER_LIST_SHORT_FORM, NullUIntegerList)
}

// Returns the absolute short form of the element type.
func (*UIntegerList) GetShortForm() Long {
	return MAL_UINTEGER_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*UIntegerList) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*UIntegerList) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*UIntegerList) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*UIntegerList) GetTypeShortForm() Integer {
	return MAL_UINTEGER_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *UIntegerList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*UInteger(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*UInteger(*list) {
		encoder.EncodeNullableUInteger(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *UIntegerList) Decode(decoder Decoder) (Element, error) {
	return DecodeUIntegerList(decoder)
}

// Decodes an instance of UIntegerList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded UIntegerList instance.
func DecodeUIntegerList(decoder Decoder) (*UIntegerList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := UIntegerList(make([]*UInteger, int(*size)))
	for i := 0; i < len(list); i++ {
		list[i], err = decoder.DecodeNullableUInteger()
		if err != nil {
			return nil, err
		}
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *UIntegerList) CreateElement() Element {
	return NewUIntegerList(0)
}

func (list *UIntegerList) IsNull() bool {
	return list == nil
}

func (*UIntegerList) Null() Element {
	return NullUIntegerList
}
