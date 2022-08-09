package archive

import (
	"github.com/CNES/ccsdsmo-malgo/com"
	"github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines ExpressionOperatorList type

type ExpressionOperatorList []*ExpressionOperator

var NullExpressionOperatorList *ExpressionOperatorList = nil

func NewExpressionOperatorList(size int) *ExpressionOperatorList {
	var list ExpressionOperatorList = ExpressionOperatorList(make([]*ExpressionOperator, size))
	return &list
}

// ================================================================================
// Defines ExpressionOperatorList type as an ElementList

func (receiver *ExpressionOperatorList) Size() int {
	if receiver != nil {
		return len(*receiver)
	}
	return -1
}

func (receiver *ExpressionOperatorList) GetElementAt(i int) mal.Element {
	if receiver == nil || i >= receiver.Size() {
		return nil
	}
	return (*receiver)[i]
}

func (receiver *ExpressionOperatorList) AppendElement(element mal.Element) {
	if receiver != nil {
		*receiver = append(*receiver, element.(*ExpressionOperator))
	}
}

// ================================================================================
// Defines ExpressionOperatorList type as a MAL Composite

func (receiver *ExpressionOperatorList) Composite() mal.Composite {
	return receiver
}

// ================================================================================
// Defines ExpressionOperatorList type as a MAL Element

const EXPRESSIONOPERATOR_LIST_TYPE_SHORT_FORM mal.Integer = -5
const EXPRESSIONOPERATOR_LIST_SHORT_FORM mal.Long = 0x2000201fffffb

// Registers ExpressionOperatorList type for polymorphism handling
func init() {
	mal.RegisterMALElement(EXPRESSIONOPERATOR_LIST_SHORT_FORM, NullExpressionOperatorList)
}

// Returns the absolute short form of the element type.
func (receiver *ExpressionOperatorList) GetShortForm() mal.Long {
	return EXPRESSIONOPERATOR_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ExpressionOperatorList) GetAreaNumber() mal.UShort {
	return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ExpressionOperatorList) GetAreaVersion() mal.UOctet {
	return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ExpressionOperatorList) GetServiceNumber() mal.UShort {
	return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ExpressionOperatorList) GetTypeShortForm() mal.Integer {
	return EXPRESSIONOPERATOR_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ExpressionOperatorList) CreateElement() mal.Element {
	return NewExpressionOperatorList(0)
}

func (receiver *ExpressionOperatorList) IsNull() bool {
	return receiver == nil
}

func (receiver *ExpressionOperatorList) Null() mal.Element {
	return NullExpressionOperatorList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ExpressionOperatorList) Encode(encoder mal.Encoder) error {
	specific := encoder.LookupSpecific(EXPRESSIONOPERATOR_LIST_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	err := encoder.EncodeUInteger(mal.NewUInteger(uint32(len([]*ExpressionOperator(*receiver)))))
	if err != nil {
		return err
	}
	for _, e := range []*ExpressionOperator(*receiver) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ExpressionOperatorList) Decode(decoder mal.Decoder) (mal.Element, error) {
	specific := decoder.LookupSpecific(EXPRESSIONOPERATOR_LIST_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := ExpressionOperatorList(make([]*ExpressionOperator, int(*size)))
	for i := 0; i < len(list); i++ {
		elem, err := decoder.DecodeNullableElement(NullExpressionOperator)
		if err != nil {
			return nil, err
		}
		list[i] = elem.(*ExpressionOperator)
	}
	return &list, nil
}
