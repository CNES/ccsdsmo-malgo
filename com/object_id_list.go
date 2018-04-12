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
// Defines COM ObjectIdList type
// ################################################################################

type ObjectIdList []*ObjectId

var (
	NullObjectIdList *ObjectIdList = nil
)

func NewObjectIdList(size int) *ObjectIdList {
	var list ObjectIdList = ObjectIdList(make([]*ObjectId, size))
	return &list
}

// ================================================================================
// Defines COM ObjectIdList type as an ElementList

func (list *ObjectIdList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

func (list *ObjectIdList) GetElementAt(i int) Element {
	if list != nil {
		if i < list.Size() {
			return (*list)[i]
		}
		return nil
	}
	return nil
}

func (list *ObjectIdList) AppendElement(element Element) {
	if list != nil {
		*list = append(*list, element.(*ObjectId))
	}
}

// ================================================================================
// Defines MAL ObjectIdList type as a MAL Composite

func (list *ObjectIdList) Composite() Composite {
	return list
}

// ================================================================================
// Defines COM ObjectIdList type as a MAL Element

const COM_OBJECT_ID_LIST_TYPE_SHORT_FORM Integer = -0x3
const COM_OBJECT_ID_LIST_SHORT_FORM Long = 0x2000001FFFFFD

// Registers COM ObjectIdList type for polymorpsism handling
func init() {
	RegisterMALElement(COM_OBJECT_ID_LIST_SHORT_FORM, NullObjectIdList)
}

// Returns the absolute short form of the element type.
func (*ObjectIdList) GetShortForm() Long {
	return COM_OBJECT_ID_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*ObjectIdList) GetAreaNumber() UShort {
	return COM_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*ObjectIdList) GetAreaVersion() UOctet {
	return COM_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*ObjectIdList) GetServiceNumber() UShort {
	return NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*ObjectIdList) GetTypeShortForm() Integer {
	//	return MAL_ENTITY_REQUEST_TYPE_SHORT_FORM & 0x01FFFF00
	return COM_OBJECT_ID_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *ObjectIdList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*ObjectId(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*ObjectId(*list) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *ObjectIdList) Decode(decoder Decoder) (Element, error) {
	return DecodeObjectIdList(decoder)
}

// Decodes an instance of ObjectIdList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded ObjectIdList instance.
func DecodeObjectIdList(decoder Decoder) (*ObjectIdList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := ObjectIdList(make([]*ObjectId, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullObjectId)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*ObjectId)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *ObjectIdList) CreateElement() Element {
	return NewObjectIdList(0)
}

func (list *ObjectIdList) IsNull() bool {
	return list == nil
}

func (*ObjectIdList) Null() Element {
	return NullObjectIdList
}
