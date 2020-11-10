package mal

import ()

// Defines UpdateTypeList type

type UpdateTypeList []*UpdateType

var NullUpdateTypeList *UpdateTypeList = nil

func NewUpdateTypeList(size int) *UpdateTypeList {
  var list UpdateTypeList = UpdateTypeList(make([]*UpdateType, size))
  return &list
}

// ================================================================================
// Defines UpdateTypeList type as an ElementList

func (receiver *UpdateTypeList) Size() int {
  if receiver != nil {
    return len(*receiver)
  }
  return -1
}

func (receiver *UpdateTypeList) GetElementAt(i int) Element {
  if receiver == nil || i >= receiver.Size() {
    return nil
  }
  return (*receiver)[i]
}

func (receiver *UpdateTypeList) AppendElement(element Element) {
  if receiver != nil {
    *receiver = append(*receiver, element.(*UpdateType))
  }
}

// ================================================================================
// Defines UpdateTypeList type as a MAL Composite

func (receiver *UpdateTypeList) Composite() Composite {
  return receiver
}

// ================================================================================
// Defines UpdateTypeList type as a MAL Element

//const UPDATETYPE_LIST_TYPE_SHORT_FORM Integer = -22
//const UPDATETYPE_LIST_SHORT_FORM Long = 0x65000001ffffea

// Registers UpdateTypeList type for polymorphism handling
func init() {
  RegisterMALElement(UPDATETYPE_LIST_SHORT_FORM, NullUpdateTypeList)
}

// Returns the absolute short form of the element type.
func (receiver *UpdateTypeList) GetShortForm() Long {
  return UPDATETYPE_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *UpdateTypeList) GetAreaNumber() UShort {
  return AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *UpdateTypeList) GetAreaVersion() UOctet {
  return AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *UpdateTypeList) GetServiceNumber() UShort {
    return NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *UpdateTypeList) GetTypeShortForm() Integer {
  return UPDATETYPE_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *UpdateTypeList) CreateElement() Element {
  return NewUpdateTypeList(0)
}

func (receiver *UpdateTypeList) IsNull() bool {
  return receiver == nil
}

func (receiver *UpdateTypeList) Null() Element {
  return NullUpdateTypeList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *UpdateTypeList) Encode(encoder Encoder) error {
  specific := encoder.LookupSpecific(UPDATETYPE_LIST_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*UpdateType(*receiver)))))
  if err != nil {
    return err
  }
  for _, e := range []*UpdateType(*receiver) {
    encoder.EncodeNullableElement(e)
  }
  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *UpdateTypeList) Decode(decoder Decoder) (Element, error) {
  specific := decoder.LookupSpecific(UPDATETYPE_LIST_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  size, err := decoder.DecodeUInteger()
  if err != nil {
    return nil, err
  }
  list := UpdateTypeList(make([]*UpdateType, int(*size)))
  for i := 0; i < len(list); i++ {
    elem, err := decoder.DecodeNullableElement(NullUpdateType)
    if err != nil {
      return nil, err
    }
    list[i] = elem.(*UpdateType)
  }
  return &list, nil
}
