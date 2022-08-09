package archive

import (
	"github.com/CNES/ccsdsmo-malgo/com"
	"github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines ArchiveDetails type

type ArchiveDetails struct {
	InstId    mal.Long
	Details   com.ObjectDetails
	Network   *mal.Identifier
	Timestamp *mal.FineTime
	Provider  *mal.URI
}

var (
	NullArchiveDetails *ArchiveDetails = nil
)

func NewArchiveDetails() *ArchiveDetails {
	return new(ArchiveDetails)
}

// ================================================================================
// Defines ArchiveDetails type as a MAL Composite

func (receiver *ArchiveDetails) Composite() mal.Composite {
	return receiver
}

// ================================================================================
// Defines ArchiveDetails type as a MAL Element

const ARCHIVEDETAILS_TYPE_SHORT_FORM mal.Integer = 1
const ARCHIVEDETAILS_SHORT_FORM mal.Long = 0x2000201000001

// Registers ArchiveDetails type for polymorphism handling
func init() {
	mal.RegisterMALElement(ARCHIVEDETAILS_SHORT_FORM, NullArchiveDetails)
}

// Returns the absolute short form of the element type.
func (receiver *ArchiveDetails) GetShortForm() mal.Long {
	return ARCHIVEDETAILS_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ArchiveDetails) GetAreaNumber() mal.UShort {
	return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ArchiveDetails) GetAreaVersion() mal.UOctet {
	return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ArchiveDetails) GetServiceNumber() mal.UShort {
	return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ArchiveDetails) GetTypeShortForm() mal.Integer {
	return ARCHIVEDETAILS_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ArchiveDetails) CreateElement() mal.Element {
	return new(ArchiveDetails)
}

func (receiver *ArchiveDetails) IsNull() bool {
	return receiver == nil
}

func (receiver *ArchiveDetails) Null() mal.Element {
	return NullArchiveDetails
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ArchiveDetails) Encode(encoder mal.Encoder) error {
	specific := encoder.LookupSpecific(ARCHIVEDETAILS_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	err := encoder.EncodeLong(&receiver.InstId)
	if err != nil {
		return err
	}
	err = encoder.EncodeElement(&receiver.Details)
	if err != nil {
		return err
	}
	err = encoder.EncodeNullableIdentifier(receiver.Network)
	if err != nil {
		return err
	}
	err = encoder.EncodeNullableFineTime(receiver.Timestamp)
	if err != nil {
		return err
	}
	err = encoder.EncodeNullableURI(receiver.Provider)
	if err != nil {
		return err
	}

	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ArchiveDetails) Decode(decoder mal.Decoder) (mal.Element, error) {
	specific := decoder.LookupSpecific(ARCHIVEDETAILS_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	InstId, err := decoder.DecodeLong()
	if err != nil {
		return nil, err
	}
	Details, err := decoder.DecodeElement(com.NullObjectDetails)
	if err != nil {
		return nil, err
	}
	Network, err := decoder.DecodeNullableIdentifier()
	if err != nil {
		return nil, err
	}
	Timestamp, err := decoder.DecodeNullableFineTime()
	if err != nil {
		return nil, err
	}
	Provider, err := decoder.DecodeNullableURI()
	if err != nil {
		return nil, err
	}

	var composite = ArchiveDetails{
		InstId:    *InstId,
		Details:   *Details.(*com.ObjectDetails),
		Network:   Network,
		Timestamp: Timestamp,
		Provider:  Provider,
	}
	return &composite, nil
}
