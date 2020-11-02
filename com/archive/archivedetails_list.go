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
package archive

import (
  "github.com/CNES/ccsdsmo-malgo/mal"
  "github.com/CNES/ccsdsmo-malgo/com"
)

// Defines ArchiveDetailsList type

type ArchiveDetailsList []*ArchiveDetails

var NullArchiveDetailsList *ArchiveDetailsList = nil

func NewArchiveDetailsList(size int) *ArchiveDetailsList {
  var list ArchiveDetailsList = ArchiveDetailsList(make([]*ArchiveDetails, size))
  return &list
}

// ================================================================================
// Defines ArchiveDetailsList type as an ElementList

func (receiver *ArchiveDetailsList) Size() int {
  if receiver != nil {
    return len(*receiver)
  }
  return -1
}

func (receiver *ArchiveDetailsList) GetElementAt(i int) mal.Element {
  if receiver == nil || i >= receiver.Size() {
    return nil
  }
  return (*receiver)[i]
}

func (receiver *ArchiveDetailsList) AppendElement(element mal.Element) {
  if receiver != nil {
    *receiver = append(*receiver, element.(*ArchiveDetails))
  }
}

// ================================================================================
// Defines ArchiveDetailsList type as a MAL Composite

func (receiver *ArchiveDetailsList) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines ArchiveDetailsList type as a MAL Element

const ARCHIVEDETAILS_LIST_TYPE_SHORT_FORM mal.Integer = -1
const ARCHIVEDETAILS_LIST_SHORT_FORM mal.Long = 0x2000201ffffff

// Registers ArchiveDetailsList type for polymorphism handling
func init() {
  mal.RegisterMALElement(ARCHIVEDETAILS_LIST_SHORT_FORM, NullArchiveDetailsList)
}

// Returns the absolute short form of the element type.
func (receiver *ArchiveDetailsList) GetShortForm() mal.Long {
  return ARCHIVEDETAILS_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ArchiveDetailsList) GetAreaNumber() mal.UShort {
  return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ArchiveDetailsList) GetAreaVersion() mal.UOctet {
  return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ArchiveDetailsList) GetServiceNumber() mal.UShort {
    return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ArchiveDetailsList) GetTypeShortForm() mal.Integer {
  return ARCHIVEDETAILS_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ArchiveDetailsList) CreateElement() mal.Element {
  return NewArchiveDetailsList(0)
}

func (receiver *ArchiveDetailsList) IsNull() bool {
  return receiver == nil
}

func (receiver *ArchiveDetailsList) Null() mal.Element {
  return NullArchiveDetailsList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ArchiveDetailsList) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(ARCHIVEDETAILS_LIST_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeUInteger(mal.NewUInteger(uint32(len([]*ArchiveDetails(*receiver)))))
  if err != nil {
    return err
  }
  for _, e := range []*ArchiveDetails(*receiver) {
    encoder.EncodeNullableElement(e)
  }
  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ArchiveDetailsList) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(ARCHIVEDETAILS_LIST_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  size, err := decoder.DecodeUInteger()
  if err != nil {
    return nil, err
  }
  list := ArchiveDetailsList(make([]*ArchiveDetails, int(*size)))
  for i := 0; i < len(list); i++ {
    elem, err := decoder.DecodeNullableElement(NullArchiveDetails)
    if err != nil {
      return nil, err
    }
    list[i] = elem.(*ArchiveDetails)
  }
  return &list, nil
}
