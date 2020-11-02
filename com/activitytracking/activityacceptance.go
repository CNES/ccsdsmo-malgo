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

// Defines ActivityAcceptance type

type ActivityAcceptance struct {
  Success mal.Boolean
}

var (
  NullActivityAcceptance *ActivityAcceptance = nil
)
func NewActivityAcceptance() *ActivityAcceptance {
  return new(ActivityAcceptance)
}

// ================================================================================
// Defines ActivityAcceptance type as a MAL Composite

func (receiver *ActivityAcceptance) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines ActivityAcceptance type as a MAL Element

const ACTIVITYACCEPTANCE_TYPE_SHORT_FORM mal.Integer = 2
const ACTIVITYACCEPTANCE_SHORT_FORM mal.Long = 0x2000301000002

// Registers ActivityAcceptance type for polymorphism handling
func init() {
  mal.RegisterMALElement(ACTIVITYACCEPTANCE_SHORT_FORM, NullActivityAcceptance)
}

// Returns the absolute short form of the element type.
func (receiver *ActivityAcceptance) GetShortForm() mal.Long {
  return ACTIVITYACCEPTANCE_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ActivityAcceptance) GetAreaNumber() mal.UShort {
  return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ActivityAcceptance) GetAreaVersion() mal.UOctet {
  return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ActivityAcceptance) GetServiceNumber() mal.UShort {
    return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ActivityAcceptance) GetTypeShortForm() mal.Integer {
  return ACTIVITYACCEPTANCE_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ActivityAcceptance) CreateElement() mal.Element {
  return new(ActivityAcceptance)
}

func (receiver *ActivityAcceptance) IsNull() bool {
  return receiver == nil
}

func (receiver *ActivityAcceptance) Null() mal.Element {
  return NullActivityAcceptance
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ActivityAcceptance) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(ACTIVITYACCEPTANCE_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeBoolean(&receiver.Success)
  if err != nil {
    return err
  }

  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ActivityAcceptance) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(ACTIVITYACCEPTANCE_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  Success, err := decoder.DecodeBoolean()
  if err != nil {
    return nil, err
  }

  var composite = ActivityAcceptance {
    Success: *Success,
  }
  return &composite, nil
}
