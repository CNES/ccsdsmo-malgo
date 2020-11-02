/**
 * MIT License
 *
 * Copyright (c) 2020 CNES
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
package activitytracking

import (
  "github.com/CNES/ccsdsmo-malgo/mal"
  "github.com/CNES/ccsdsmo-malgo/com"
)

// Defines ActivityTransferList type

type ActivityTransferList []*ActivityTransfer

var NullActivityTransferList *ActivityTransferList = nil

func NewActivityTransferList(size int) *ActivityTransferList {
  var list ActivityTransferList = ActivityTransferList(make([]*ActivityTransfer, size))
  return &list
}

// ================================================================================
// Defines ActivityTransferList type as an ElementList

func (receiver *ActivityTransferList) Size() int {
  if receiver != nil {
    return len(*receiver)
  }
  return -1
}

func (receiver *ActivityTransferList) GetElementAt(i int) mal.Element {
  if receiver == nil || i >= receiver.Size() {
    return nil
  }
  return (*receiver)[i]
}

func (receiver *ActivityTransferList) AppendElement(element mal.Element) {
  if receiver != nil {
    *receiver = append(*receiver, element.(*ActivityTransfer))
  }
}

// ================================================================================
// Defines ActivityTransferList type as a MAL Composite

func (receiver *ActivityTransferList) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines ActivityTransferList type as a MAL Element

const ACTIVITYTRANSFER_LIST_TYPE_SHORT_FORM mal.Integer = -1
const ACTIVITYTRANSFER_LIST_SHORT_FORM mal.Long = 0x2000301ffffff

// Registers ActivityTransferList type for polymorphism handling
func init() {
  mal.RegisterMALElement(ACTIVITYTRANSFER_LIST_SHORT_FORM, NullActivityTransferList)
}

// Returns the absolute short form of the element type.
func (receiver *ActivityTransferList) GetShortForm() mal.Long {
  return ACTIVITYTRANSFER_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ActivityTransferList) GetAreaNumber() mal.UShort {
  return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ActivityTransferList) GetAreaVersion() mal.UOctet {
  return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ActivityTransferList) GetServiceNumber() mal.UShort {
    return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ActivityTransferList) GetTypeShortForm() mal.Integer {
  return ACTIVITYTRANSFER_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ActivityTransferList) CreateElement() mal.Element {
  return NewActivityTransferList(0)
}

func (receiver *ActivityTransferList) IsNull() bool {
  return receiver == nil
}

func (receiver *ActivityTransferList) Null() mal.Element {
  return NullActivityTransferList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ActivityTransferList) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(ACTIVITYTRANSFER_LIST_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeUInteger(mal.NewUInteger(uint32(len([]*ActivityTransfer(*receiver)))))
  if err != nil {
    return err
  }
  for _, e := range []*ActivityTransfer(*receiver) {
    encoder.EncodeNullableElement(e)
  }
  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ActivityTransferList) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(ACTIVITYTRANSFER_LIST_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  size, err := decoder.DecodeUInteger()
  if err != nil {
    return nil, err
  }
  list := ActivityTransferList(make([]*ActivityTransfer, int(*size)))
  for i := 0; i < len(list); i++ {
    elem, err := decoder.DecodeNullableElement(NullActivityTransfer)
    if err != nil {
      return nil, err
    }
    list[i] = elem.(*ActivityTransfer)
  }
  return &list, nil
}
