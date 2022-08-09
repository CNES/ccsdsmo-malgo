package com

import (
	"github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines ObjectIdList type

type ObjectIdList []*ObjectId

var NullObjectIdList *ObjectIdList = nil

func NewObjectIdList(size int) *ObjectIdList {
	var list ObjectIdList = ObjectIdList(make([]*ObjectId, size))
	return &list
}

// ================================================================================
// Defines ObjectIdList type as an ElementList

func (receiver *ObjectIdList) Size() int {
	if receiver != nil {
		return len(*receiver)
	}
	return -1
}

func (receiver *ObjectIdList) GetElementAt(i int) mal.Element {
	if receiver == nil || i >= receiver.Size() {
		return nil
	}
	return (*receiver)[i]
}

func (receiver *ObjectIdList) AppendElement(element mal.Element) {
	if receiver != nil {
		*receiver = append(*receiver, element.(*ObjectId))
	}
}

// ================================================================================
// Defines ObjectIdList type as a MAL Composite

func (receiver *ObjectIdList) Composite() mal.Composite {
	return receiver
}

// ================================================================================
// Defines ObjectIdList type as a MAL Element

const OBJECTID_LIST_TYPE_SHORT_FORM mal.Integer = -3
const OBJECTID_LIST_SHORT_FORM mal.Long = 0x2000001fffffd

// Registers ObjectIdList type for polymorphism handling
func init() {
	mal.RegisterMALElement(OBJECTID_LIST_SHORT_FORM, NullObjectIdList)
}

// Returns the absolute short form of the element type.
func (receiver *ObjectIdList) GetShortForm() mal.Long {
	return OBJECTID_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ObjectIdList) GetAreaNumber() mal.UShort {
	return AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ObjectIdList) GetAreaVersion() mal.UOctet {
	return AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ObjectIdList) GetServiceNumber() mal.UShort {
	return mal.NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ObjectIdList) GetTypeShortForm() mal.Integer {
	return OBJECTID_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ObjectIdList) CreateElement() mal.Element {
	return NewObjectIdList(0)
}

func (receiver *ObjectIdList) IsNull() bool {
	return receiver == nil
}

func (receiver *ObjectIdList) Null() mal.Element {
	return NullObjectIdList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ObjectIdList) Encode(encoder mal.Encoder) error {
	specific := encoder.LookupSpecific(OBJECTID_LIST_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	err := encoder.EncodeUInteger(mal.NewUInteger(uint32(len([]*ObjectId(*receiver)))))
	if err != nil {
		return err
	}
	for _, e := range []*ObjectId(*receiver) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ObjectIdList) Decode(decoder mal.Decoder) (mal.Element, error) {
	specific := decoder.LookupSpecific(OBJECTID_LIST_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := ObjectIdList(make([]*ObjectId, int(*size)))
	for i := 0; i < len(list); i++ {
		elem, err := decoder.DecodeNullableElement(NullObjectId)
		if err != nil {
			return nil, err
		}
		list[i] = elem.(*ObjectId)
	}
	return &list, nil
}
