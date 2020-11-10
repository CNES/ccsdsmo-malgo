package archive

import (
  "github.com/CNES/ccsdsmo-malgo/mal"
  "github.com/CNES/ccsdsmo-malgo/com"
)

// Defines CompositeFilter type

type CompositeFilter struct {
  FieldName mal.String
  Type ExpressionOperator
  FieldValue mal.Attribute
}

var (
  NullCompositeFilter *CompositeFilter = nil
)
func NewCompositeFilter() *CompositeFilter {
  return new(CompositeFilter)
}

// ================================================================================
// Defines CompositeFilter type as a MAL Composite

func (receiver *CompositeFilter) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines CompositeFilter type as a MAL Element

const COMPOSITEFILTER_TYPE_SHORT_FORM mal.Integer = 3
const COMPOSITEFILTER_SHORT_FORM mal.Long = 0x2000201000003

// Registers CompositeFilter type for polymorphism handling
func init() {
  mal.RegisterMALElement(COMPOSITEFILTER_SHORT_FORM, NullCompositeFilter)
}

// Returns the absolute short form of the element type.
func (receiver *CompositeFilter) GetShortForm() mal.Long {
  return COMPOSITEFILTER_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *CompositeFilter) GetAreaNumber() mal.UShort {
  return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *CompositeFilter) GetAreaVersion() mal.UOctet {
  return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *CompositeFilter) GetServiceNumber() mal.UShort {
    return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *CompositeFilter) GetTypeShortForm() mal.Integer {
  return COMPOSITEFILTER_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *CompositeFilter) CreateElement() mal.Element {
  return new(CompositeFilter)
}

func (receiver *CompositeFilter) IsNull() bool {
  return receiver == nil
}

func (receiver *CompositeFilter) Null() mal.Element {
  return NullCompositeFilter
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *CompositeFilter) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(COMPOSITEFILTER_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeString(&receiver.FieldName)
  if err != nil {
    return err
  }
  err = encoder.EncodeElement(&receiver.Type)
  if err != nil {
    return err
  }
  err = encoder.EncodeNullableAttribute(receiver.FieldValue)
  if err != nil {
    return err
  }

  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *CompositeFilter) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(COMPOSITEFILTER_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  FieldName, err := decoder.DecodeString()
  if err != nil {
    return nil, err
  }
  Type, err := decoder.DecodeElement(NullExpressionOperator)
  if err != nil {
    return nil, err
  }
  FieldValue, err := decoder.DecodeNullableAttribute()
  if err != nil {
    return nil, err
  }

  var composite = CompositeFilter {
    FieldName: *FieldName,
    Type: *Type.(*ExpressionOperator),
    FieldValue: FieldValue,
  }
  return &composite, nil
}
