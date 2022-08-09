package activitytracking

import (
	"github.com/CNES/ccsdsmo-malgo/com"
	"github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines ActivityTransfer type

type ActivityTransfer struct {
	Success          mal.Boolean
	EstimateDuration *mal.Duration
	NextDestination  *mal.URI
}

var (
	NullActivityTransfer *ActivityTransfer = nil
)

func NewActivityTransfer() *ActivityTransfer {
	return new(ActivityTransfer)
}

// ================================================================================
// Defines ActivityTransfer type as a MAL Composite

func (receiver *ActivityTransfer) Composite() mal.Composite {
	return receiver
}

// ================================================================================
// Defines ActivityTransfer type as a MAL Element

const ACTIVITYTRANSFER_TYPE_SHORT_FORM mal.Integer = 1
const ACTIVITYTRANSFER_SHORT_FORM mal.Long = 0x2000301000001

// Registers ActivityTransfer type for polymorphism handling
func init() {
	mal.RegisterMALElement(ACTIVITYTRANSFER_SHORT_FORM, NullActivityTransfer)
}

// Returns the absolute short form of the element type.
func (receiver *ActivityTransfer) GetShortForm() mal.Long {
	return ACTIVITYTRANSFER_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ActivityTransfer) GetAreaNumber() mal.UShort {
	return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ActivityTransfer) GetAreaVersion() mal.UOctet {
	return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ActivityTransfer) GetServiceNumber() mal.UShort {
	return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ActivityTransfer) GetTypeShortForm() mal.Integer {
	return ACTIVITYTRANSFER_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ActivityTransfer) CreateElement() mal.Element {
	return new(ActivityTransfer)
}

func (receiver *ActivityTransfer) IsNull() bool {
	return receiver == nil
}

func (receiver *ActivityTransfer) Null() mal.Element {
	return NullActivityTransfer
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ActivityTransfer) Encode(encoder mal.Encoder) error {
	specific := encoder.LookupSpecific(ACTIVITYTRANSFER_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	err := encoder.EncodeBoolean(&receiver.Success)
	if err != nil {
		return err
	}
	err = encoder.EncodeNullableDuration(receiver.EstimateDuration)
	if err != nil {
		return err
	}
	err = encoder.EncodeNullableURI(receiver.NextDestination)
	if err != nil {
		return err
	}

	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ActivityTransfer) Decode(decoder mal.Decoder) (mal.Element, error) {
	specific := decoder.LookupSpecific(ACTIVITYTRANSFER_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	Success, err := decoder.DecodeBoolean()
	if err != nil {
		return nil, err
	}
	EstimateDuration, err := decoder.DecodeNullableDuration()
	if err != nil {
		return nil, err
	}
	NextDestination, err := decoder.DecodeNullableURI()
	if err != nil {
		return nil, err
	}

	var composite = ActivityTransfer{
		Success:          *Success,
		EstimateDuration: EstimateDuration,
		NextDestination:  NextDestination,
	}
	return &composite, nil
}
