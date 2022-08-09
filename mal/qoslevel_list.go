package mal

// Defines QoSLevelList type

type QoSLevelList []*QoSLevel

var NullQoSLevelList *QoSLevelList = nil

func NewQoSLevelList(size int) *QoSLevelList {
	var list QoSLevelList = QoSLevelList(make([]*QoSLevel, size))
	return &list
}

// ================================================================================
// Defines QoSLevelList type as an ElementList

func (receiver *QoSLevelList) Size() int {
	if receiver != nil {
		return len(*receiver)
	}
	return -1
}

func (receiver *QoSLevelList) GetElementAt(i int) Element {
	if receiver == nil || i >= receiver.Size() {
		return nil
	}
	return (*receiver)[i]
}

func (receiver *QoSLevelList) AppendElement(element Element) {
	if receiver != nil {
		*receiver = append(*receiver, element.(*QoSLevel))
	}
}

// ================================================================================
// Defines QoSLevelList type as a MAL Composite

func (receiver *QoSLevelList) Composite() Composite {
	return receiver
}

// ================================================================================
// Defines QoSLevelList type as a MAL Element

//const QOSLEVEL_LIST_TYPE_SHORT_FORM Integer = -21
//const QOSLEVEL_LIST_SHORT_FORM Long = 0x65000001ffffeb

// Registers QoSLevelList type for polymorphism handling
func init() {
	RegisterMALElement(QOSLEVEL_LIST_SHORT_FORM, NullQoSLevelList)
}

// Returns the absolute short form of the element type.
func (receiver *QoSLevelList) GetShortForm() Long {
	return QOSLEVEL_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *QoSLevelList) GetAreaNumber() UShort {
	return AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *QoSLevelList) GetAreaVersion() UOctet {
	return AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *QoSLevelList) GetServiceNumber() UShort {
	return NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *QoSLevelList) GetTypeShortForm() Integer {
	return QOSLEVEL_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *QoSLevelList) CreateElement() Element {
	return NewQoSLevelList(0)
}

func (receiver *QoSLevelList) IsNull() bool {
	return receiver == nil
}

func (receiver *QoSLevelList) Null() Element {
	return NullQoSLevelList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *QoSLevelList) Encode(encoder Encoder) error {
	specific := encoder.LookupSpecific(QOSLEVEL_LIST_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*QoSLevel(*receiver)))))
	if err != nil {
		return err
	}
	for _, e := range []*QoSLevel(*receiver) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *QoSLevelList) Decode(decoder Decoder) (Element, error) {
	specific := decoder.LookupSpecific(QOSLEVEL_LIST_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := QoSLevelList(make([]*QoSLevel, int(*size)))
	for i := 0; i < len(list); i++ {
		elem, err := decoder.DecodeNullableElement(NullQoSLevel)
		if err != nil {
			return nil, err
		}
		list[i] = elem.(*QoSLevel)
	}
	return &list, nil
}
