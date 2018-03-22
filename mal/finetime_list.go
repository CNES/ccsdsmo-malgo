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
// Defines MAL FineTimeList type
// ################################################################################

type FineTimeList []*FineTime

var (
	NullFineTimeList *FineTimeList = nil
)

func NewFineTimeList(size int) *FineTimeList {
	var list FineTimeList = FineTimeList(make([]*FineTime, size))
	return &list
}

// ================================================================================
// Defines MAL FineTimeList type as an ElementList

func (list *FineTimeList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

// ================================================================================
// Defines MAL FineTimeList type as a MAL Composite

func (list *FineTimeList) Composite() Composite {
	return list
}

// ================================================================================
// Defines MAL FineTime type as a MAL Element

const MAL_FINETIME_LIST_TYPE_SHORT_FORM Integer = -0x11
const MAL_FINETIME_LIST_SHORT_FORM Long = 0x1000001FFFFEF

// Registers MAL BlobList type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_FINETIME_LIST_SHORT_FORM, NullFineTimeList)
}

// Returns the absolute short form of the element type.
func (*FineTimeList) GetShortForm() Long {
	return MAL_FINETIME_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*FineTimeList) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*FineTimeList) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*FineTimeList) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*FineTimeList) GetTypeShortForm() Integer {
	//	return MAL_FINETIME_TYPE_SHORT_FORM & 0x01FFFF00
	return MAL_FINETIME_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *FineTimeList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*FineTime(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*FineTime(*list) {
		encoder.EncodeNullableFineTime(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *FineTimeList) Decode(decoder Decoder) (Element, error) {
	return DecodeFineTimeList(decoder)
}

// Decodes an instance of BooleanList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded BooleanList instance.
func DecodeFineTimeList(decoder Decoder) (*FineTimeList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := FineTimeList(make([]*FineTime, int(*size)))
	for i := 0; i < len(list); i++ {
		list[i], err = decoder.DecodeNullableFineTime()
		if err != nil {
			return nil, err
		}
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *FineTimeList) CreateElement() Element {
	return NewFineTimeList(0)
}

func (list *FineTimeList) IsNull() bool {
	return list == nil
}

func (*FineTimeList) Null() Element {
	return NullFineTimeList
}
