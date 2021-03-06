package com

import (
  "github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines ObjectKey type

type ObjectKey struct {
  Domain mal.IdentifierList
  InstId mal.Long
}

var (
  NullObjectKey *ObjectKey = nil
)
func NewObjectKey() *ObjectKey {
  return new(ObjectKey)
}

// ================================================================================
// Defines ObjectKey type as a MAL Composite

func (receiver *ObjectKey) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines ObjectKey type as a MAL Element

const OBJECTKEY_TYPE_SHORT_FORM mal.Integer = 2
const OBJECTKEY_SHORT_FORM mal.Long = 0x2000001000002

// Registers ObjectKey type for polymorphism handling
func init() {
  mal.RegisterMALElement(OBJECTKEY_SHORT_FORM, NullObjectKey)
}

// Returns the absolute short form of the element type.
func (receiver *ObjectKey) GetShortForm() mal.Long {
  return OBJECTKEY_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ObjectKey) GetAreaNumber() mal.UShort {
  return AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ObjectKey) GetAreaVersion() mal.UOctet {
  return AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ObjectKey) GetServiceNumber() mal.UShort {
    return mal.NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ObjectKey) GetTypeShortForm() mal.Integer {
  return OBJECTKEY_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ObjectKey) CreateElement() mal.Element {
  return new(ObjectKey)
}

func (receiver *ObjectKey) IsNull() bool {
  return receiver == nil
}

func (receiver *ObjectKey) Null() mal.Element {
  return NullObjectKey
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ObjectKey) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(OBJECTKEY_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeElement(&receiver.Domain)
  if err != nil {
    return err
  }
  err = encoder.EncodeLong(&receiver.InstId)
  if err != nil {
    return err
  }

  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ObjectKey) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(OBJECTKEY_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  Domain, err := decoder.DecodeElement(mal.NullIdentifierList)
  if err != nil {
    return nil, err
  }
  InstId, err := decoder.DecodeLong()
  if err != nil {
    return nil, err
  }

  var composite = ObjectKey {
    Domain: *Domain.(*mal.IdentifierList),
    InstId: *InstId,
  }
  return &composite, nil
}
