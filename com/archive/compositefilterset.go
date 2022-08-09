package archive

import (
	"github.com/CNES/ccsdsmo-malgo/com"
	"github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines CompositeFilterSet type

type CompositeFilterSet struct {
	Filters CompositeFilterList
}

var (
	NullCompositeFilterSet *CompositeFilterSet = nil
)

func NewCompositeFilterSet() *CompositeFilterSet {
	return new(CompositeFilterSet)
}

// ================================================================================
// Defines CompositeFilterSet type as a MAL Composite

func (receiver *CompositeFilterSet) Composite() mal.Composite {
	return receiver
}

// ================================================================================
// Defines CompositeFilterSet type as a MAL Element

const COMPOSITEFILTERSET_TYPE_SHORT_FORM mal.Integer = 4
const COMPOSITEFILTERSET_SHORT_FORM mal.Long = 0x2000201000004

// Registers CompositeFilterSet type for polymorphism handling
func init() {
	mal.RegisterMALElement(COMPOSITEFILTERSET_SHORT_FORM, NullCompositeFilterSet)
}

// Returns the absolute short form of the element type.
func (receiver *CompositeFilterSet) GetShortForm() mal.Long {
	return COMPOSITEFILTERSET_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *CompositeFilterSet) GetAreaNumber() mal.UShort {
	return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *CompositeFilterSet) GetAreaVersion() mal.UOctet {
	return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *CompositeFilterSet) GetServiceNumber() mal.UShort {
	return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *CompositeFilterSet) GetTypeShortForm() mal.Integer {
	return COMPOSITEFILTERSET_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *CompositeFilterSet) CreateElement() mal.Element {
	return new(CompositeFilterSet)
}

func (receiver *CompositeFilterSet) IsNull() bool {
	return receiver == nil
}

func (receiver *CompositeFilterSet) Null() mal.Element {
	return NullCompositeFilterSet
}

// ================================================================================
// Defines CompositeFilterSet type as a QueryFilter
func (receiver *CompositeFilterSet) QueryFilter() QueryFilter {
	return receiver
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *CompositeFilterSet) Encode(encoder mal.Encoder) error {
	specific := encoder.LookupSpecific(COMPOSITEFILTERSET_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	err := encoder.EncodeElement(&receiver.Filters)
	if err != nil {
		return err
	}

	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *CompositeFilterSet) Decode(decoder mal.Decoder) (mal.Element, error) {
	specific := decoder.LookupSpecific(COMPOSITEFILTERSET_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	Filters, err := decoder.DecodeElement(NullCompositeFilterList)
	if err != nil {
		return nil, err
	}

	var composite = CompositeFilterSet{
		Filters: *Filters.(*CompositeFilterList),
	}
	return &composite, nil
}
