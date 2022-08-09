package com

import (
	"github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines InstanceBooleanPairList type

type InstanceBooleanPairList []*InstanceBooleanPair

var NullInstanceBooleanPairList *InstanceBooleanPairList = nil

func NewInstanceBooleanPairList(size int) *InstanceBooleanPairList {
	var list InstanceBooleanPairList = InstanceBooleanPairList(make([]*InstanceBooleanPair, size))
	return &list
}

// ================================================================================
// Defines InstanceBooleanPairList type as an ElementList

func (receiver *InstanceBooleanPairList) Size() int {
	if receiver != nil {
		return len(*receiver)
	}
	return -1
}

func (receiver *InstanceBooleanPairList) GetElementAt(i int) mal.Element {
	if receiver == nil || i >= receiver.Size() {
		return nil
	}
	return (*receiver)[i]
}

func (receiver *InstanceBooleanPairList) AppendElement(element mal.Element) {
	if receiver != nil {
		*receiver = append(*receiver, element.(*InstanceBooleanPair))
	}
}

// ================================================================================
// Defines InstanceBooleanPairList type as a MAL Composite

func (receiver *InstanceBooleanPairList) Composite() mal.Composite {
	return receiver
}

// ================================================================================
// Defines InstanceBooleanPairList type as a MAL Element

const INSTANCEBOOLEANPAIR_LIST_TYPE_SHORT_FORM mal.Integer = -5
const INSTANCEBOOLEANPAIR_LIST_SHORT_FORM mal.Long = 0x2000001fffffb

// Registers InstanceBooleanPairList type for polymorphism handling
func init() {
	mal.RegisterMALElement(INSTANCEBOOLEANPAIR_LIST_SHORT_FORM, NullInstanceBooleanPairList)
}

// Returns the absolute short form of the element type.
func (receiver *InstanceBooleanPairList) GetShortForm() mal.Long {
	return INSTANCEBOOLEANPAIR_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *InstanceBooleanPairList) GetAreaNumber() mal.UShort {
	return AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *InstanceBooleanPairList) GetAreaVersion() mal.UOctet {
	return AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *InstanceBooleanPairList) GetServiceNumber() mal.UShort {
	return mal.NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *InstanceBooleanPairList) GetTypeShortForm() mal.Integer {
	return INSTANCEBOOLEANPAIR_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *InstanceBooleanPairList) CreateElement() mal.Element {
	return NewInstanceBooleanPairList(0)
}

func (receiver *InstanceBooleanPairList) IsNull() bool {
	return receiver == nil
}

func (receiver *InstanceBooleanPairList) Null() mal.Element {
	return NullInstanceBooleanPairList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *InstanceBooleanPairList) Encode(encoder mal.Encoder) error {
	specific := encoder.LookupSpecific(INSTANCEBOOLEANPAIR_LIST_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	err := encoder.EncodeUInteger(mal.NewUInteger(uint32(len([]*InstanceBooleanPair(*receiver)))))
	if err != nil {
		return err
	}
	for _, e := range []*InstanceBooleanPair(*receiver) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *InstanceBooleanPairList) Decode(decoder mal.Decoder) (mal.Element, error) {
	specific := decoder.LookupSpecific(INSTANCEBOOLEANPAIR_LIST_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := InstanceBooleanPairList(make([]*InstanceBooleanPair, int(*size)))
	for i := 0; i < len(list); i++ {
		elem, err := decoder.DecodeNullableElement(NullInstanceBooleanPair)
		if err != nil {
			return nil, err
		}
		list[i] = elem.(*InstanceBooleanPair)
	}
	return &list, nil
}
