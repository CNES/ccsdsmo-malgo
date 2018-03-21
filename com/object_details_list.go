/**
 * MIT License
 *
 * Copyright (c) 2018 CNES
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
package com

import (
	. "github.com/ccsdsmo/malgo/mal"
)

// ################################################################################
// Defines COM ObjectDetailsList type
// ################################################################################

type ObjectDetailsList []*ObjectDetails

var (
	NullObjectDetailsList *ObjectDetailsList = nil
)

func NewObjectDetailsList(size int) *ObjectDetailsList {
	var list ObjectDetailsList = ObjectDetailsList(make([]*ObjectDetails, size))
	return &list
}

// ================================================================================
// Defines COM ObjectDetailsList type as an ElementList

func (list *ObjectDetailsList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

// ================================================================================
// Defines COM ObjectDetailsList type as a MAL Element

const COM_OBJECT_DETAILS_LIST_TYPE_SHORT_FORM Integer = -0x4
const COM_OBJECT_DETAILS_LIST_SHORT_FORM Long = 0x2000001FFFFFC

// Registers COM ObjectDetailsList type for polymorpsism handling
func init() {
	RegisterMALElement(COM_OBJECT_DETAILS_LIST_SHORT_FORM, NullObjectDetailsList)
}

// Returns the absolute short form of the element type.
func (*ObjectDetailsList) GetShortForm() Long {
	return COM_OBJECT_DETAILS_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*ObjectDetailsList) GetAreaNumber() UShort {
	return COM_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*ObjectDetailsList) GetAreaVersion() UOctet {
	return COM_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*ObjectDetailsList) GetServiceNumber() UShort {
	return NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*ObjectDetailsList) GetTypeShortForm() Integer {
	//	return MAL_ENTITY_REQUEST_TYPE_SHORT_FORM & 0x01FFFF00
	return COM_OBJECT_DETAILS_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *ObjectDetailsList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*ObjectDetails(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*ObjectDetails(*list) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *ObjectDetailsList) Decode(decoder Decoder) (Element, error) {
	return DecodeObjectDetailsList(decoder)
}

// Decodes an instance of ObjectDetailsList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded ObjectDetailsList instance.
func DecodeObjectDetailsList(decoder Decoder) (*ObjectDetailsList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := ObjectDetailsList(make([]*ObjectDetails, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullObjectDetails)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*ObjectDetails)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *ObjectDetailsList) CreateElement() Element {
	return NewObjectDetailsList(0)
}

func (list *ObjectDetailsList) IsNull() bool {
	return list == nil
}

func (*ObjectDetailsList) Null() Element {
	return NullObjectDetailsList
}
