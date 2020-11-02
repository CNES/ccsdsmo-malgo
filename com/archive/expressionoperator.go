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


// Defines ExpressionOperator type

type ExpressionOperator uint32
const (
  EXPRESSIONOPERATOR_EQUAL_OVAL = iota
  EXPRESSIONOPERATOR_EQUAL_NVAL = 1
  EXPRESSIONOPERATOR_DIFFER_OVAL
  EXPRESSIONOPERATOR_DIFFER_NVAL = 2
  EXPRESSIONOPERATOR_GREATER_OVAL
  EXPRESSIONOPERATOR_GREATER_NVAL = 3
  EXPRESSIONOPERATOR_GREATER_OR_EQUAL_OVAL
  EXPRESSIONOPERATOR_GREATER_OR_EQUAL_NVAL = 4
  EXPRESSIONOPERATOR_LESS_OVAL
  EXPRESSIONOPERATOR_LESS_NVAL = 5
  EXPRESSIONOPERATOR_LESS_OR_EQUAL_OVAL
  EXPRESSIONOPERATOR_LESS_OR_EQUAL_NVAL = 6
  EXPRESSIONOPERATOR_CONTAINS_OVAL
  EXPRESSIONOPERATOR_CONTAINS_NVAL = 7
  EXPRESSIONOPERATOR_ICONTAINS_OVAL
  EXPRESSIONOPERATOR_ICONTAINS_NVAL = 8
)

// Conversion table OVAL->NVAL
var nvalTable = []uint32 {
  EXPRESSIONOPERATOR_EQUAL_NVAL,
  EXPRESSIONOPERATOR_DIFFER_NVAL,
  EXPRESSIONOPERATOR_GREATER_NVAL,
  EXPRESSIONOPERATOR_GREATER_OR_EQUAL_NVAL,
  EXPRESSIONOPERATOR_LESS_NVAL,
  EXPRESSIONOPERATOR_LESS_OR_EQUAL_NVAL,
  EXPRESSIONOPERATOR_CONTAINS_NVAL,
  EXPRESSIONOPERATOR_ICONTAINS_NVAL,
}

var (
  EXPRESSIONOPERATOR_EQUAL = ExpressionOperator(EXPRESSIONOPERATOR_EQUAL_OVAL)
  EXPRESSIONOPERATOR_DIFFER = ExpressionOperator(EXPRESSIONOPERATOR_DIFFER_OVAL)
  EXPRESSIONOPERATOR_GREATER = ExpressionOperator(EXPRESSIONOPERATOR_GREATER_OVAL)
  EXPRESSIONOPERATOR_GREATER_OR_EQUAL = ExpressionOperator(EXPRESSIONOPERATOR_GREATER_OR_EQUAL_OVAL)
  EXPRESSIONOPERATOR_LESS = ExpressionOperator(EXPRESSIONOPERATOR_LESS_OVAL)
  EXPRESSIONOPERATOR_LESS_OR_EQUAL = ExpressionOperator(EXPRESSIONOPERATOR_LESS_OR_EQUAL_OVAL)
  EXPRESSIONOPERATOR_CONTAINS = ExpressionOperator(EXPRESSIONOPERATOR_CONTAINS_OVAL)
  EXPRESSIONOPERATOR_ICONTAINS = ExpressionOperator(EXPRESSIONOPERATOR_ICONTAINS_OVAL)
)

var NullExpressionOperator *ExpressionOperator = nil
func NewExpressionOperator(i uint32) *ExpressionOperator {
  var val ExpressionOperator = ExpressionOperator(i)
  return &val
}

// ================================================================================
// Defines ExpressionOperator type as a MAL Element

const EXPRESSIONOPERATOR_TYPE_SHORT_FORM mal.Integer = 5
const EXPRESSIONOPERATOR_SHORT_FORM mal.Long = 0x2000201000005

// Registers ExpressionOperator type for polymorphism handling
func init() {
  mal.RegisterMALElement(EXPRESSIONOPERATOR_SHORT_FORM, NullExpressionOperator)
}

// Returns the absolute short form of the element type.
func (receiver *ExpressionOperator) GetShortForm() mal.Long {
  return EXPRESSIONOPERATOR_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ExpressionOperator) GetAreaNumber() mal.UShort {
  return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ExpressionOperator) GetAreaVersion() mal.UOctet {
  return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ExpressionOperator) GetServiceNumber() mal.UShort {
    return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ExpressionOperator) GetTypeShortForm() mal.Integer {
  return EXPRESSIONOPERATOR_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ExpressionOperator) CreateElement() mal.Element {
  return NewExpressionOperator(0)
}

func (receiver *ExpressionOperator) IsNull() bool {
  return receiver == nil
}

func (receiver *ExpressionOperator) Null() mal.Element {
  return NullExpressionOperator
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ExpressionOperator) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(EXPRESSIONOPERATOR_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  value := mal.NewUOctet(uint8(uint32(*receiver)))
  return encoder.EncodeUOctet(value)
}


// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ExpressionOperator) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(EXPRESSIONOPERATOR_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  elem, err := decoder.DecodeUOctet()
  if err != nil {
    return receiver.Null(), err
  }
  value := ExpressionOperator(uint32(uint8(*elem)))
  return &value, nil
}

