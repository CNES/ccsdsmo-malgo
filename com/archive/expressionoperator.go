package archive

import (
  "fmt"
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
var nvalTable_ExpressionOperator = []uint32 {
  EXPRESSIONOPERATOR_EQUAL_NVAL,
  EXPRESSIONOPERATOR_DIFFER_NVAL,
  EXPRESSIONOPERATOR_GREATER_NVAL,
  EXPRESSIONOPERATOR_GREATER_OR_EQUAL_NVAL,
  EXPRESSIONOPERATOR_LESS_NVAL,
  EXPRESSIONOPERATOR_LESS_OR_EQUAL_NVAL,
  EXPRESSIONOPERATOR_CONTAINS_NVAL,
  EXPRESSIONOPERATOR_ICONTAINS_NVAL,
}

// Conversion map NVAL->OVAL
var ovalMap_ExpressionOperator map[uint32]uint32

var (
  EXPRESSIONOPERATOR_EQUAL = ExpressionOperator(EXPRESSIONOPERATOR_EQUAL_NVAL)
  EXPRESSIONOPERATOR_DIFFER = ExpressionOperator(EXPRESSIONOPERATOR_DIFFER_NVAL)
  EXPRESSIONOPERATOR_GREATER = ExpressionOperator(EXPRESSIONOPERATOR_GREATER_NVAL)
  EXPRESSIONOPERATOR_GREATER_OR_EQUAL = ExpressionOperator(EXPRESSIONOPERATOR_GREATER_OR_EQUAL_NVAL)
  EXPRESSIONOPERATOR_LESS = ExpressionOperator(EXPRESSIONOPERATOR_LESS_NVAL)
  EXPRESSIONOPERATOR_LESS_OR_EQUAL = ExpressionOperator(EXPRESSIONOPERATOR_LESS_OR_EQUAL_NVAL)
  EXPRESSIONOPERATOR_CONTAINS = ExpressionOperator(EXPRESSIONOPERATOR_CONTAINS_NVAL)
  EXPRESSIONOPERATOR_ICONTAINS = ExpressionOperator(EXPRESSIONOPERATOR_ICONTAINS_NVAL)
)

var NullExpressionOperator *ExpressionOperator = nil

func init() {
  ovalMap_ExpressionOperator = make(map[uint32]uint32)
  for oval, nval := range nvalTable_<Enum> {
    ovalMap_ExpressionOperator[nval] = uint32(oval)
  }
}

func (receiver ExpressionOperator) GetNumericValue() uint32 {
  return uint32(receiver)
}
func (receiver ExpressionOperator) GetOrdinalValue() uint32 {
  nvalue := receiver.GetNumericValue()
  return ovalMap_ExpressionOperator[nvalue]
}
func ExpressionOperatorFromNumericValue(nval uint32) (ExpressionOperator, error) {
  _, ok := ovalMap_ExpressionOperator[nval]
  if !ok {
    return ExpressionOperator(0), fmt.Errorf("Invalid numeric value for ExpressionOperator: %v", nval)
  }
  return ExpressionOperator(nval), nil
}
func ExpressionOperatorFromOrdinalValue(oval uint32) (ExpressionOperator, error) {
  if oval >= uint32(len(nvalTable_ExpressionOperator)) {
    return ExpressionOperator(0), fmt.Errorf("Invalid ordinal value for ExpressionOperator: %v", oval)
  }
  return ExpressionOperator(nvalTable_ExpressionOperator[oval]), nil
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
  return NullExpressionOperator
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

  value := mal.NewUOctet(codedNativeType(receiver.GetOrdinalValue()))
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
  value, err := ExpressionOperatorFromOrdinalValue(uint32(uint8(*elem)))
  return &value, err
}

