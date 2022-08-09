package mal

// Defines SessionTypeList type

type SessionTypeList []*SessionType

var NullSessionTypeList *SessionTypeList = nil

func NewSessionTypeList(size int) *SessionTypeList {
	var list SessionTypeList = SessionTypeList(make([]*SessionType, size))
	return &list
}

// ================================================================================
// Defines SessionTypeList type as an ElementList

func (receiver *SessionTypeList) Size() int {
	if receiver != nil {
		return len(*receiver)
	}
	return -1
}

func (receiver *SessionTypeList) GetElementAt(i int) Element {
	if receiver == nil || i >= receiver.Size() {
		return nil
	}
	return (*receiver)[i]
}

func (receiver *SessionTypeList) AppendElement(element Element) {
	if receiver != nil {
		*receiver = append(*receiver, element.(*SessionType))
	}
}

// ================================================================================
// Defines SessionTypeList type as a MAL Composite

func (receiver *SessionTypeList) Composite() Composite {
	return receiver
}

// ================================================================================
// Defines SessionTypeList type as a MAL Element

//const SESSIONTYPE_LIST_TYPE_SHORT_FORM Integer = -20
//const SESSIONTYPE_LIST_SHORT_FORM Long = 0x65000001ffffec

// Registers SessionTypeList type for polymorphism handling
func init() {
	RegisterMALElement(SESSIONTYPE_LIST_SHORT_FORM, NullSessionTypeList)
}

// Returns the absolute short form of the element type.
func (receiver *SessionTypeList) GetShortForm() Long {
	return SESSIONTYPE_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *SessionTypeList) GetAreaNumber() UShort {
	return AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *SessionTypeList) GetAreaVersion() UOctet {
	return AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *SessionTypeList) GetServiceNumber() UShort {
	return NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *SessionTypeList) GetTypeShortForm() Integer {
	return SESSIONTYPE_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *SessionTypeList) CreateElement() Element {
	return NewSessionTypeList(0)
}

func (receiver *SessionTypeList) IsNull() bool {
	return receiver == nil
}

func (receiver *SessionTypeList) Null() Element {
	return NullSessionTypeList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *SessionTypeList) Encode(encoder Encoder) error {
	specific := encoder.LookupSpecific(SESSIONTYPE_LIST_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*SessionType(*receiver)))))
	if err != nil {
		return err
	}
	for _, e := range []*SessionType(*receiver) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *SessionTypeList) Decode(decoder Decoder) (Element, error) {
	specific := decoder.LookupSpecific(SESSIONTYPE_LIST_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := SessionTypeList(make([]*SessionType, int(*size)))
	for i := 0; i < len(list); i++ {
		elem, err := decoder.DecodeNullableElement(NullSessionType)
		if err != nil {
			return nil, err
		}
		list[i] = elem.(*SessionType)
	}
	return &list, nil
}
