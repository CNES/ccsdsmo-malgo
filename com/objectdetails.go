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

// Defines ObjectDetails type

type ObjectDetails struct {
  Related *mal.Long
  Source *ObjectId
}

var (
  NullObjectDetails *ObjectDetails = nil
)
func NewObjectDetails() *ObjectDetails {
  return new(ObjectDetails)
}

// ================================================================================
// Defines ObjectDetails type as a MAL Composite

func (receiver *ObjectDetails) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines ObjectDetails type as a MAL Element

const OBJECTDETAILS_TYPE_SHORT_FORM mal.Integer = 4
const OBJECTDETAILS_SHORT_FORM mal.Long = 0x2000001000004

// Registers ObjectDetails type for polymorphism handling
func init() {
  mal.RegisterMALElement(OBJECTDETAILS_SHORT_FORM, NullObjectDetails)
}

// Returns the absolute short form of the element type.
func (receiver *ObjectDetails) GetShortForm() mal.Long {
  return OBJECTDETAILS_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ObjectDetails) GetAreaNumber() mal.UShort {
  return AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ObjectDetails) GetAreaVersion() mal.UOctet {
  return AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ObjectDetails) GetServiceNumber() mal.UShort {
    return mal.NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ObjectDetails) GetTypeShortForm() mal.Integer {
  return OBJECTDETAILS_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ObjectDetails) CreateElement() mal.Element {
  return new(ObjectDetails)
}

func (receiver *ObjectDetails) IsNull() bool {
  return receiver == nil
}

func (receiver *ObjectDetails) Null() mal.Element {
  return NullObjectDetails
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ObjectDetails) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(OBJECTDETAILS_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeNullableLong(receiver.Related)
  if err != nil {
    return err
  }
  err = encoder.EncodeNullableElement(receiver.Source)
  if err != nil {
    return err
  }

  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ObjectDetails) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(OBJECTDETAILS_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  Related, err := decoder.DecodeNullableLong()
  if err != nil {
    return nil, err
  }
  Source, err := decoder.DecodeNullableElement(NullObjectId)
  if err != nil {
    return nil, err
  }

  var composite = ObjectDetails {
    Related: Related,
    Source: Source.(*ObjectId),
  }
  return &composite, nil
}
