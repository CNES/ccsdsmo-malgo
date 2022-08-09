package mal

// Defines InteractionTypeList type

type InteractionTypeList []*InteractionType

var NullInteractionTypeList *InteractionTypeList = nil

func NewInteractionTypeList(size int) *InteractionTypeList {
	var list InteractionTypeList = InteractionTypeList(make([]*InteractionType, size))
	return &list
}

// ================================================================================
// Defines InteractionTypeList type as an ElementList

func (receiver *InteractionTypeList) Size() int {
	if receiver != nil {
		return len(*receiver)
	}
	return -1
}

func (receiver *InteractionTypeList) GetElementAt(i int) Element {
	if receiver == nil || i >= receiver.Size() {
		return nil
	}
	return (*receiver)[i]
}

func (receiver *InteractionTypeList) AppendElement(element Element) {
	if receiver != nil {
		*receiver = append(*receiver, element.(*InteractionType))
	}
}

// ================================================================================
// Defines InteractionTypeList type as a MAL Composite

func (receiver *InteractionTypeList) Composite() Composite {
	return receiver
}

// ================================================================================
// Defines InteractionTypeList type as a MAL Element

//const INTERACTIONTYPE_LIST_TYPE_SHORT_FORM Integer = -19
//const INTERACTIONTYPE_LIST_SHORT_FORM Long = 0x65000001ffffed

// Registers InteractionTypeList type for polymorphism handling
func init() {
	RegisterMALElement(INTERACTIONTYPE_LIST_SHORT_FORM, NullInteractionTypeList)
}

// Returns the absolute short form of the element type.
func (receiver *InteractionTypeList) GetShortForm() Long {
	return INTERACTIONTYPE_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *InteractionTypeList) GetAreaNumber() UShort {
	return AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *InteractionTypeList) GetAreaVersion() UOctet {
	return AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *InteractionTypeList) GetServiceNumber() UShort {
	return NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *InteractionTypeList) GetTypeShortForm() Integer {
	return INTERACTIONTYPE_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *InteractionTypeList) CreateElement() Element {
	return NewInteractionTypeList(0)
}

func (receiver *InteractionTypeList) IsNull() bool {
	return receiver == nil
}

func (receiver *InteractionTypeList) Null() Element {
	return NullInteractionTypeList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *InteractionTypeList) Encode(encoder Encoder) error {
	specific := encoder.LookupSpecific(INTERACTIONTYPE_LIST_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*InteractionType(*receiver)))))
	if err != nil {
		return err
	}
	for _, e := range []*InteractionType(*receiver) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *InteractionTypeList) Decode(decoder Decoder) (Element, error) {
	specific := decoder.LookupSpecific(INTERACTIONTYPE_LIST_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := InteractionTypeList(make([]*InteractionType, int(*size)))
	for i := 0; i < len(list); i++ {
		elem, err := decoder.DecodeNullableElement(NullInteractionType)
		if err != nil {
			return nil, err
		}
		list[i] = elem.(*InteractionType)
	}
	return &list, nil
}
