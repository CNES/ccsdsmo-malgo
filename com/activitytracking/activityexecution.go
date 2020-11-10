package activitytracking

import (
  "github.com/CNES/ccsdsmo-malgo/mal"
  "github.com/CNES/ccsdsmo-malgo/com"
)

// Defines ActivityExecution type

type ActivityExecution struct {
  Success mal.Boolean
  ExecutionStage mal.UInteger
  StageCount mal.UInteger
}

var (
  NullActivityExecution *ActivityExecution = nil
)
func NewActivityExecution() *ActivityExecution {
  return new(ActivityExecution)
}

// ================================================================================
// Defines ActivityExecution type as a MAL Composite

func (receiver *ActivityExecution) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines ActivityExecution type as a MAL Element

const ACTIVITYEXECUTION_TYPE_SHORT_FORM mal.Integer = 3
const ACTIVITYEXECUTION_SHORT_FORM mal.Long = 0x2000301000003

// Registers ActivityExecution type for polymorphism handling
func init() {
  mal.RegisterMALElement(ACTIVITYEXECUTION_SHORT_FORM, NullActivityExecution)
}

// Returns the absolute short form of the element type.
func (receiver *ActivityExecution) GetShortForm() mal.Long {
  return ACTIVITYEXECUTION_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ActivityExecution) GetAreaNumber() mal.UShort {
  return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ActivityExecution) GetAreaVersion() mal.UOctet {
  return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ActivityExecution) GetServiceNumber() mal.UShort {
    return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ActivityExecution) GetTypeShortForm() mal.Integer {
  return ACTIVITYEXECUTION_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ActivityExecution) CreateElement() mal.Element {
  return new(ActivityExecution)
}

func (receiver *ActivityExecution) IsNull() bool {
  return receiver == nil
}

func (receiver *ActivityExecution) Null() mal.Element {
  return NullActivityExecution
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ActivityExecution) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(ACTIVITYEXECUTION_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeBoolean(&receiver.Success)
  if err != nil {
    return err
  }
  err = encoder.EncodeUInteger(&receiver.ExecutionStage)
  if err != nil {
    return err
  }
  err = encoder.EncodeUInteger(&receiver.StageCount)
  if err != nil {
    return err
  }

  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ActivityExecution) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(ACTIVITYEXECUTION_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  Success, err := decoder.DecodeBoolean()
  if err != nil {
    return nil, err
  }
  ExecutionStage, err := decoder.DecodeUInteger()
  if err != nil {
    return nil, err
  }
  StageCount, err := decoder.DecodeUInteger()
  if err != nil {
    return nil, err
  }

  var composite = ActivityExecution {
    Success: *Success,
    ExecutionStage: *ExecutionStage,
    StageCount: *StageCount,
  }
  return &composite, nil
}
