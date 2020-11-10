package com

import (
  "github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines ObjectId type

type ObjectId struct {
  Type ObjectType
  Key ObjectKey
}

var (
  NullObjectId *ObjectId = nil
)
func NewObjectId() *ObjectId {
  return new(ObjectId)
}

// ================================================================================
// Defines ObjectId type as a MAL Composite

func (receiver *ObjectId) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines ObjectId type as a MAL Element

const OBJECTID_TYPE_SHORT_FORM mal.Integer = 3
const OBJECTID_SHORT_FORM mal.Long = 0x2000001000003

// Registers ObjectId type for polymorphism handling
func init() {
  mal.RegisterMALElement(OBJECTID_SHORT_FORM, NullObjectId)
}

// Returns the absolute short form of the element type.
func (receiver *ObjectId) GetShortForm() mal.Long {
  return OBJECTID_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ObjectId) GetAreaNumber() mal.UShort {
  return AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ObjectId) GetAreaVersion() mal.UOctet {
  return AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ObjectId) GetServiceNumber() mal.UShort {
    return mal.NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ObjectId) GetTypeShortForm() mal.Integer {
  return OBJECTID_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ObjectId) CreateElement() mal.Element {
  return new(ObjectId)
}

func (receiver *ObjectId) IsNull() bool {
  return receiver == nil
}

func (receiver *ObjectId) Null() mal.Element {
  return NullObjectId
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ObjectId) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(OBJECTID_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeElement(&receiver.Type)
  if err != nil {
    return err
  }
  err = encoder.EncodeElement(&receiver.Key)
  if err != nil {
    return err
  }

  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ObjectId) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(OBJECTID_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  Type, err := decoder.DecodeElement(NullObjectType)
  if err != nil {
    return nil, err
  }
  Key, err := decoder.DecodeElement(NullObjectKey)
  if err != nil {
    return nil, err
  }

  var composite = ObjectId {
    Type: *Type.(*ObjectType),
    Key: *Key.(*ObjectKey),
  }
  return &composite, nil
}
