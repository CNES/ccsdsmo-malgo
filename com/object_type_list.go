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
// Defines COM ObjectTypeList type
// ################################################################################

type ObjectTypeList []*ObjectType

var (
	NullObjectTypeList *ObjectTypeList = nil
)

func NewObjectTypeList(size int) *ObjectTypeList {
	var list ObjectTypeList = ObjectTypeList(make([]*ObjectType, size))
	return &list
}

// ================================================================================
// Defines COM ObjectTypeList type as an ElementList

func (list *ObjectTypeList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

func (list *ObjectTypeList) GetElementAt(i int) Element {
	if list != nil {
		if i < list.Size() {
			return (*list)[i]
		}
		return nil
	}
	return nil
}

func (list *ObjectTypeList) AppendElement(element Element) {
	if list != nil {
		*list = append(*list, element.(*ObjectType))
	}
}

// ================================================================================
// Defines COM ObjectTypeList type as a MAL Composite

func (list *ObjectTypeList) Composite() Composite {
	return list
}

// ================================================================================
// Defines COM ObjectTypeList type as a MAL Element

const COM_OBJECT_TYPE_LIST_TYPE_SHORT_FORM Integer = -0x1
const COM_OBJECT_TYPE_LIST_SHORT_FORM Long = 0x2000001FFFFFF

// Registers COM ObjectTypeList type for polymorpsism handling
func init() {
	RegisterMALElement(COM_OBJECT_TYPE_LIST_SHORT_FORM, NullObjectTypeList)
}

// Returns the absolute short form of the element type.
func (*ObjectTypeList) GetShortForm() Long {
	return COM_OBJECT_TYPE_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*ObjectTypeList) GetAreaNumber() UShort {
	return COM_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*ObjectTypeList) GetAreaVersion() UOctet {
	return COM_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*ObjectTypeList) GetServiceNumber() UShort {
	return NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*ObjectTypeList) GetTypeShortForm() Integer {
	//	return MAL_ENTITY_REQUEST_TYPE_SHORT_FORM & 0x01FFFF00
	return COM_OBJECT_TYPE_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *ObjectTypeList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*ObjectType(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*ObjectType(*list) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *ObjectTypeList) Decode(decoder Decoder) (Element, error) {
	return DecodeObjectTypeList(decoder)
}

// Decodes an instance of ObjectTypeList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded ObjectTypeList instance.
func DecodeObjectTypeList(decoder Decoder) (*ObjectTypeList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := ObjectTypeList(make([]*ObjectType, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullObjectType)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*ObjectType)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *ObjectTypeList) CreateElement() Element {
	return NewObjectTypeList(0)
}

func (list *ObjectTypeList) IsNull() bool {
	return list == nil
}

func (*ObjectTypeList) Null() Element {
	return NullObjectTypeList
}
