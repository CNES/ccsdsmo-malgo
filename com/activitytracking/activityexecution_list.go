package activitytracking

import (
	"github.com/CNES/ccsdsmo-malgo/com"
	"github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines ActivityExecutionList type

type ActivityExecutionList []*ActivityExecution

var NullActivityExecutionList *ActivityExecutionList = nil

func NewActivityExecutionList(size int) *ActivityExecutionList {
	var list ActivityExecutionList = ActivityExecutionList(make([]*ActivityExecution, size))
	return &list
}

// ================================================================================
// Defines ActivityExecutionList type as an ElementList

func (receiver *ActivityExecutionList) Size() int {
	if receiver != nil {
		return len(*receiver)
	}
	return -1
}

func (receiver *ActivityExecutionList) GetElementAt(i int) mal.Element {
	if receiver == nil || i >= receiver.Size() {
		return nil
	}
	return (*receiver)[i]
}

func (receiver *ActivityExecutionList) AppendElement(element mal.Element) {
	if receiver != nil {
		*receiver = append(*receiver, element.(*ActivityExecution))
	}
}

// ================================================================================
// Defines ActivityExecutionList type as a MAL Composite

func (receiver *ActivityExecutionList) Composite() mal.Composite {
	return receiver
}

// ================================================================================
// Defines ActivityExecutionList type as a MAL Element

const ACTIVITYEXECUTION_LIST_TYPE_SHORT_FORM mal.Integer = -3
const ACTIVITYEXECUTION_LIST_SHORT_FORM mal.Long = 0x2000301fffffd

// Registers ActivityExecutionList type for polymorphism handling
func init() {
	mal.RegisterMALElement(ACTIVITYEXECUTION_LIST_SHORT_FORM, NullActivityExecutionList)
}

// Returns the absolute short form of the element type.
func (receiver *ActivityExecutionList) GetShortForm() mal.Long {
	return ACTIVITYEXECUTION_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ActivityExecutionList) GetAreaNumber() mal.UShort {
	return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ActivityExecutionList) GetAreaVersion() mal.UOctet {
	return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ActivityExecutionList) GetServiceNumber() mal.UShort {
	return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ActivityExecutionList) GetTypeShortForm() mal.Integer {
	return ACTIVITYEXECUTION_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ActivityExecutionList) CreateElement() mal.Element {
	return NewActivityExecutionList(0)
}

func (receiver *ActivityExecutionList) IsNull() bool {
	return receiver == nil
}

func (receiver *ActivityExecutionList) Null() mal.Element {
	return NullActivityExecutionList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ActivityExecutionList) Encode(encoder mal.Encoder) error {
	specific := encoder.LookupSpecific(ACTIVITYEXECUTION_LIST_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	err := encoder.EncodeUInteger(mal.NewUInteger(uint32(len([]*ActivityExecution(*receiver)))))
	if err != nil {
		return err
	}
	for _, e := range []*ActivityExecution(*receiver) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ActivityExecutionList) Decode(decoder mal.Decoder) (mal.Element, error) {
	specific := decoder.LookupSpecific(ACTIVITYEXECUTION_LIST_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := ActivityExecutionList(make([]*ActivityExecution, int(*size)))
	for i := 0; i < len(list); i++ {
		elem, err := decoder.DecodeNullableElement(NullActivityExecution)
		if err != nil {
			return nil, err
		}
		list[i] = elem.(*ActivityExecution)
	}
	return &list, nil
}
