package activitytracking

import (
  "github.com/CNES/ccsdsmo-malgo/mal"
  "github.com/CNES/ccsdsmo-malgo/com"
)

// Defines OperationActivity type

type OperationActivity struct {
  InteractionType mal.InteractionType
}

var (
  NullOperationActivity *OperationActivity = nil
)
func NewOperationActivity() *OperationActivity {
  return new(OperationActivity)
}

// ================================================================================
// Defines OperationActivity type as a MAL Composite

func (receiver *OperationActivity) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines OperationActivity type as a MAL Element

const OPERATIONACTIVITY_TYPE_SHORT_FORM mal.Integer = 4
const OPERATIONACTIVITY_SHORT_FORM mal.Long = 0x2000301000004

// Registers OperationActivity type for polymorphism handling
func init() {
  mal.RegisterMALElement(OPERATIONACTIVITY_SHORT_FORM, NullOperationActivity)
}

// Returns the absolute short form of the element type.
func (receiver *OperationActivity) GetShortForm() mal.Long {
  return OPERATIONACTIVITY_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *OperationActivity) GetAreaNumber() mal.UShort {
  return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *OperationActivity) GetAreaVersion() mal.UOctet {
  return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *OperationActivity) GetServiceNumber() mal.UShort {
    return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *OperationActivity) GetTypeShortForm() mal.Integer {
  return OPERATIONACTIVITY_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *OperationActivity) CreateElement() mal.Element {
  return new(OperationActivity)
}

func (receiver *OperationActivity) IsNull() bool {
  return receiver == nil
}

func (receiver *OperationActivity) Null() mal.Element {
  return NullOperationActivity
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *OperationActivity) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(OPERATIONACTIVITY_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeElement(&receiver.InteractionType)
  if err != nil {
    return err
  }

  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *OperationActivity) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(OPERATIONACTIVITY_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  InteractionType, err := decoder.DecodeElement(mal.NullInteractionType)
  if err != nil {
    return nil, err
  }

  var composite = OperationActivity {
    InteractionType: *InteractionType.(*mal.InteractionType),
  }
  return &composite, nil
}
