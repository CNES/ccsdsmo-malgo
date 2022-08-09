package archive

import (
	"github.com/CNES/ccsdsmo-malgo/com"
	"github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines CompositeFilterList type

type CompositeFilterList []*CompositeFilter

var NullCompositeFilterList *CompositeFilterList = nil

func NewCompositeFilterList(size int) *CompositeFilterList {
	var list CompositeFilterList = CompositeFilterList(make([]*CompositeFilter, size))
	return &list
}

// ================================================================================
// Defines CompositeFilterList type as an ElementList

func (receiver *CompositeFilterList) Size() int {
	if receiver != nil {
		return len(*receiver)
	}
	return -1
}

func (receiver *CompositeFilterList) GetElementAt(i int) mal.Element {
	if receiver == nil || i >= receiver.Size() {
		return nil
	}
	return (*receiver)[i]
}

func (receiver *CompositeFilterList) AppendElement(element mal.Element) {
	if receiver != nil {
		*receiver = append(*receiver, element.(*CompositeFilter))
	}
}

// ================================================================================
// Defines CompositeFilterList type as a MAL Composite

func (receiver *CompositeFilterList) Composite() mal.Composite {
	return receiver
}

// ================================================================================
// Defines CompositeFilterList type as a MAL Element

const COMPOSITEFILTER_LIST_TYPE_SHORT_FORM mal.Integer = -3
const COMPOSITEFILTER_LIST_SHORT_FORM mal.Long = 0x2000201fffffd

// Registers CompositeFilterList type for polymorphism handling
func init() {
	mal.RegisterMALElement(COMPOSITEFILTER_LIST_SHORT_FORM, NullCompositeFilterList)
}

// Returns the absolute short form of the element type.
func (receiver *CompositeFilterList) GetShortForm() mal.Long {
	return COMPOSITEFILTER_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *CompositeFilterList) GetAreaNumber() mal.UShort {
	return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *CompositeFilterList) GetAreaVersion() mal.UOctet {
	return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *CompositeFilterList) GetServiceNumber() mal.UShort {
	return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *CompositeFilterList) GetTypeShortForm() mal.Integer {
	return COMPOSITEFILTER_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *CompositeFilterList) CreateElement() mal.Element {
	return NewCompositeFilterList(0)
}

func (receiver *CompositeFilterList) IsNull() bool {
	return receiver == nil
}

func (receiver *CompositeFilterList) Null() mal.Element {
	return NullCompositeFilterList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *CompositeFilterList) Encode(encoder mal.Encoder) error {
	specific := encoder.LookupSpecific(COMPOSITEFILTER_LIST_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	err := encoder.EncodeUInteger(mal.NewUInteger(uint32(len([]*CompositeFilter(*receiver)))))
	if err != nil {
		return err
	}
	for _, e := range []*CompositeFilter(*receiver) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *CompositeFilterList) Decode(decoder mal.Decoder) (mal.Element, error) {
	specific := decoder.LookupSpecific(COMPOSITEFILTER_LIST_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := CompositeFilterList(make([]*CompositeFilter, int(*size)))
	for i := 0; i < len(list); i++ {
		elem, err := decoder.DecodeNullableElement(NullCompositeFilter)
		if err != nil {
			return nil, err
		}
		list[i] = elem.(*CompositeFilter)
	}
	return &list, nil
}
