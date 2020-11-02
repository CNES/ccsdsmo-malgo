/**
 * MIT License
 *
 * Copyright (c) 2018 - 2020 CNES
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
  "github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines ObjectTypeList type

type ObjectTypeList []*ObjectType

var NullObjectTypeList *ObjectTypeList = nil

func NewObjectTypeList(size int) *ObjectTypeList {
  var list ObjectTypeList = ObjectTypeList(make([]*ObjectType, size))
  return &list
}

// ================================================================================
// Defines ObjectTypeList type as an ElementList

func (receiver *ObjectTypeList) Size() int {
  if receiver != nil {
    return len(*receiver)
  }
  return -1
}

func (receiver *ObjectTypeList) GetElementAt(i int) mal.Element {
  if receiver == nil || i >= receiver.Size() {
    return nil
  }
  return (*receiver)[i]
}

func (receiver *ObjectTypeList) AppendElement(element mal.Element) {
  if receiver != nil {
    *receiver = append(*receiver, element.(*ObjectType))
  }
}

// ================================================================================
// Defines ObjectTypeList type as a MAL Composite

func (receiver *ObjectTypeList) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines ObjectTypeList type as a MAL Element

const OBJECTTYPE_LIST_TYPE_SHORT_FORM mal.Integer = -1
const OBJECTTYPE_LIST_SHORT_FORM mal.Long = 0x2000001ffffff

// Registers ObjectTypeList type for polymorphism handling
func init() {
  mal.RegisterMALElement(OBJECTTYPE_LIST_SHORT_FORM, NullObjectTypeList)
}

// Returns the absolute short form of the element type.
func (receiver *ObjectTypeList) GetShortForm() mal.Long {
  return OBJECTTYPE_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ObjectTypeList) GetAreaNumber() mal.UShort {
  return AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ObjectTypeList) GetAreaVersion() mal.UOctet {
  return AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ObjectTypeList) GetServiceNumber() mal.UShort {
    return mal.NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ObjectTypeList) GetTypeShortForm() mal.Integer {
  return OBJECTTYPE_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ObjectTypeList) CreateElement() mal.Element {
  return NewObjectTypeList(0)
}

func (receiver *ObjectTypeList) IsNull() bool {
  return receiver == nil
}

func (receiver *ObjectTypeList) Null() mal.Element {
  return NullObjectTypeList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ObjectTypeList) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(OBJECTTYPE_LIST_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeUInteger(mal.NewUInteger(uint32(len([]*ObjectType(*receiver)))))
  if err != nil {
    return err
  }
  for _, e := range []*ObjectType(*receiver) {
    encoder.EncodeNullableElement(e)
  }
  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ObjectTypeList) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(OBJECTTYPE_LIST_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  size, err := decoder.DecodeUInteger()
  if err != nil {
    return nil, err
  }
  list := ObjectTypeList(make([]*ObjectType, int(*size)))
  for i := 0; i < len(list); i++ {
    elem, err := decoder.DecodeNullableElement(NullObjectType)
    if err != nil {
      return nil, err
    }
    list[i] = elem.(*ObjectType)
  }
  return &list, nil
}
