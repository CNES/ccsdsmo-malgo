package com

import (
	"github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines ObjectDetailsList type

type ObjectDetailsList []*ObjectDetails

var NullObjectDetailsList *ObjectDetailsList = nil

func NewObjectDetailsList(size int) *ObjectDetailsList {
	var list ObjectDetailsList = ObjectDetailsList(make([]*ObjectDetails, size))
	return &list
}

// ================================================================================
// Defines ObjectDetailsList type as an ElementList

func (receiver *ObjectDetailsList) Size() int {
	if receiver != nil {
		return len(*receiver)
	}
	return -1
}

func (receiver *ObjectDetailsList) GetElementAt(i int) mal.Element {
	if receiver == nil || i >= receiver.Size() {
		return nil
	}
	return (*receiver)[i]
}

func (receiver *ObjectDetailsList) AppendElement(element mal.Element) {
	if receiver != nil {
		*receiver = append(*receiver, element.(*ObjectDetails))
	}
}

// ================================================================================
// Defines ObjectDetailsList type as a MAL Composite

func (receiver *ObjectDetailsList) Composite() mal.Composite {
	return receiver
}

// ================================================================================
// Defines ObjectDetailsList type as a MAL Element

const OBJECTDETAILS_LIST_TYPE_SHORT_FORM mal.Integer = -4
const OBJECTDETAILS_LIST_SHORT_FORM mal.Long = 0x2000001fffffc

// Registers ObjectDetailsList type for polymorphism handling
func init() {
	mal.RegisterMALElement(OBJECTDETAILS_LIST_SHORT_FORM, NullObjectDetailsList)
}

// Returns the absolute short form of the element type.
func (receiver *ObjectDetailsList) GetShortForm() mal.Long {
	return OBJECTDETAILS_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ObjectDetailsList) GetAreaNumber() mal.UShort {
	return AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ObjectDetailsList) GetAreaVersion() mal.UOctet {
	return AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ObjectDetailsList) GetServiceNumber() mal.UShort {
	return mal.NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ObjectDetailsList) GetTypeShortForm() mal.Integer {
	return OBJECTDETAILS_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ObjectDetailsList) CreateElement() mal.Element {
	return NewObjectDetailsList(0)
}

func (receiver *ObjectDetailsList) IsNull() bool {
	return receiver == nil
}

func (receiver *ObjectDetailsList) Null() mal.Element {
	return NullObjectDetailsList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ObjectDetailsList) Encode(encoder mal.Encoder) error {
	specific := encoder.LookupSpecific(OBJECTDETAILS_LIST_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	err := encoder.EncodeUInteger(mal.NewUInteger(uint32(len([]*ObjectDetails(*receiver)))))
	if err != nil {
		return err
	}
	for _, e := range []*ObjectDetails(*receiver) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ObjectDetailsList) Decode(decoder mal.Decoder) (mal.Element, error) {
	specific := decoder.LookupSpecific(OBJECTDETAILS_LIST_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := ObjectDetailsList(make([]*ObjectDetails, int(*size)))
	for i := 0; i < len(list); i++ {
		elem, err := decoder.DecodeNullableElement(NullObjectDetails)
		if err != nil {
			return nil, err
		}
		list[i] = elem.(*ObjectDetails)
	}
	return &list, nil
}
