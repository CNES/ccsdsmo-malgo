package archive

import (
	"github.com/CNES/ccsdsmo-malgo/com"
	"github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines CompositeFilterSetList type

type CompositeFilterSetList []*CompositeFilterSet

var NullCompositeFilterSetList *CompositeFilterSetList = nil

func NewCompositeFilterSetList(size int) *CompositeFilterSetList {
	var list CompositeFilterSetList = CompositeFilterSetList(make([]*CompositeFilterSet, size))
	return &list
}

// ================================================================================
// Defines CompositeFilterSetList type as an ElementList

func (receiver *CompositeFilterSetList) Size() int {
	if receiver != nil {
		return len(*receiver)
	}
	return -1
}

func (receiver *CompositeFilterSetList) GetElementAt(i int) mal.Element {
	if receiver == nil || i >= receiver.Size() {
		return nil
	}
	return (*receiver)[i]
}

func (receiver *CompositeFilterSetList) AppendElement(element mal.Element) {
	if receiver != nil {
		*receiver = append(*receiver, element.(*CompositeFilterSet))
	}
}

// ================================================================================
// Defines CompositeFilterSetList type as a MAL Composite

func (receiver *CompositeFilterSetList) Composite() mal.Composite {
	return receiver
}

// ================================================================================
// Defines CompositeFilterSetList type as a MAL Element

const COMPOSITEFILTERSET_LIST_TYPE_SHORT_FORM mal.Integer = -4
const COMPOSITEFILTERSET_LIST_SHORT_FORM mal.Long = 0x2000201fffffc

// Registers CompositeFilterSetList type for polymorphism handling
func init() {
	mal.RegisterMALElement(COMPOSITEFILTERSET_LIST_SHORT_FORM, NullCompositeFilterSetList)
}

// Returns the absolute short form of the element type.
func (receiver *CompositeFilterSetList) GetShortForm() mal.Long {
	return COMPOSITEFILTERSET_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *CompositeFilterSetList) GetAreaNumber() mal.UShort {
	return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *CompositeFilterSetList) GetAreaVersion() mal.UOctet {
	return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *CompositeFilterSetList) GetServiceNumber() mal.UShort {
	return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *CompositeFilterSetList) GetTypeShortForm() mal.Integer {
	return COMPOSITEFILTERSET_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *CompositeFilterSetList) CreateElement() mal.Element {
	return NewCompositeFilterSetList(0)
}

func (receiver *CompositeFilterSetList) IsNull() bool {
	return receiver == nil
}

func (receiver *CompositeFilterSetList) Null() mal.Element {
	return NullCompositeFilterSetList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *CompositeFilterSetList) Encode(encoder mal.Encoder) error {
	specific := encoder.LookupSpecific(COMPOSITEFILTERSET_LIST_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	err := encoder.EncodeUInteger(mal.NewUInteger(uint32(len([]*CompositeFilterSet(*receiver)))))
	if err != nil {
		return err
	}
	for _, e := range []*CompositeFilterSet(*receiver) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *CompositeFilterSetList) Decode(decoder mal.Decoder) (mal.Element, error) {
	specific := decoder.LookupSpecific(COMPOSITEFILTERSET_LIST_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := CompositeFilterSetList(make([]*CompositeFilterSet, int(*size)))
	for i := 0; i < len(list); i++ {
		elem, err := decoder.DecodeNullableElement(NullCompositeFilterSet)
		if err != nil {
			return nil, err
		}
		list[i] = elem.(*CompositeFilterSet)
	}
	return &list, nil
}

// ================================================================================
// Defines CompositeFilterSetList type as a QueryFilterList
func (receiver *CompositeFilterSetList) QueryFilterList() QueryFilterList {
	return receiver
}
