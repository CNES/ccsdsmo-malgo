package com

import (
	"github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines ObjectKeyList type

type ObjectKeyList []*ObjectKey

var NullObjectKeyList *ObjectKeyList = nil

func NewObjectKeyList(size int) *ObjectKeyList {
	var list ObjectKeyList = ObjectKeyList(make([]*ObjectKey, size))
	return &list
}

// ================================================================================
// Defines ObjectKeyList type as an ElementList

func (receiver *ObjectKeyList) Size() int {
	if receiver != nil {
		return len(*receiver)
	}
	return -1
}

func (receiver *ObjectKeyList) GetElementAt(i int) mal.Element {
	if receiver == nil || i >= receiver.Size() {
		return nil
	}
	return (*receiver)[i]
}

func (receiver *ObjectKeyList) AppendElement(element mal.Element) {
	if receiver != nil {
		*receiver = append(*receiver, element.(*ObjectKey))
	}
}

// ================================================================================
// Defines ObjectKeyList type as a MAL Composite

func (receiver *ObjectKeyList) Composite() mal.Composite {
	return receiver
}

// ================================================================================
// Defines ObjectKeyList type as a MAL Element

const OBJECTKEY_LIST_TYPE_SHORT_FORM mal.Integer = -2
const OBJECTKEY_LIST_SHORT_FORM mal.Long = 0x2000001fffffe

// Registers ObjectKeyList type for polymorphism handling
func init() {
	mal.RegisterMALElement(OBJECTKEY_LIST_SHORT_FORM, NullObjectKeyList)
}

// Returns the absolute short form of the element type.
func (receiver *ObjectKeyList) GetShortForm() mal.Long {
	return OBJECTKEY_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ObjectKeyList) GetAreaNumber() mal.UShort {
	return AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ObjectKeyList) GetAreaVersion() mal.UOctet {
	return AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ObjectKeyList) GetServiceNumber() mal.UShort {
	return mal.NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ObjectKeyList) GetTypeShortForm() mal.Integer {
	return OBJECTKEY_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ObjectKeyList) CreateElement() mal.Element {
	return NewObjectKeyList(0)
}

func (receiver *ObjectKeyList) IsNull() bool {
	return receiver == nil
}

func (receiver *ObjectKeyList) Null() mal.Element {
	return NullObjectKeyList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ObjectKeyList) Encode(encoder mal.Encoder) error {
	specific := encoder.LookupSpecific(OBJECTKEY_LIST_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	err := encoder.EncodeUInteger(mal.NewUInteger(uint32(len([]*ObjectKey(*receiver)))))
	if err != nil {
		return err
	}
	for _, e := range []*ObjectKey(*receiver) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ObjectKeyList) Decode(decoder mal.Decoder) (mal.Element, error) {
	specific := decoder.LookupSpecific(OBJECTKEY_LIST_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := ObjectKeyList(make([]*ObjectKey, int(*size)))
	for i := 0; i < len(list); i++ {
		elem, err := decoder.DecodeNullableElement(NullObjectKey)
		if err != nil {
			return nil, err
		}
		list[i] = elem.(*ObjectKey)
	}
	return &list, nil
}
