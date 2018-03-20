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
// Defines COM ObjectKeyList type
// ################################################################################

type ObjectKeyList []*ObjectKey

var (
	NullObjectKeyList *ObjectKeyList = nil
)

func NewObjectKeyList(size int) *ObjectKeyList {
	var list ObjectKeyList = ObjectKeyList(make([]*ObjectKey, size))
	return &list
}

// ================================================================================
// Defines COM ObjectKeyList type as a MAL Element

const COM_OBJECT_KEY_LIST_TYPE_SHORT_FORM Integer = -0x2
const COM_OBJECT_KEY_LIST_SHORT_FORM Long = 0x2000001FFFFFE

// Registers COM ObjectKeyList type for polymorpsism handling
func init() {
	RegisterMALElement(COM_OBJECT_KEY_LIST_SHORT_FORM, NullObjectKeyList)
}

// Returns the absolute short form of the element type.
func (*ObjectKeyList) GetShortForm() Long {
	return COM_OBJECT_KEY_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*ObjectKeyList) GetAreaNumber() UShort {
	return COM_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*ObjectKeyList) GetAreaVersion() UOctet {
	return COM_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*ObjectKeyList) GetServiceNumber() UShort {
	return NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*ObjectKeyList) GetTypeShortForm() Integer {
	//	return MAL_ENTITY_REQUEST_TYPE_SHORT_FORM & 0x01FFFF00
	return COM_OBJECT_KEY_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *ObjectKeyList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*ObjectKey(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*ObjectKey(*list) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *ObjectKeyList) Decode(decoder Decoder) (Element, error) {
	return DecodeObjectKeyList(decoder)
}

// Decodes an instance of ObjectKeyList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded ObjectKeyList instance.
func DecodeObjectKeyList(decoder Decoder) (*ObjectKeyList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := ObjectKeyList(make([]*ObjectKey, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullObjectKey)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*ObjectKey)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *ObjectKeyList) CreateElement() Element {
	return NewObjectKeyList(0)
}

func (list *ObjectKeyList) IsNull() bool {
	return list == nil
}

func (*ObjectKeyList) Null() Element {
	return NullObjectKeyList
}
