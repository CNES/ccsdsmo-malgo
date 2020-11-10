package archive

import (
  "github.com/CNES/ccsdsmo-malgo/mal"
  "github.com/CNES/ccsdsmo-malgo/com"
)

// Defines ArchiveQueryList type

type ArchiveQueryList []*ArchiveQuery

var NullArchiveQueryList *ArchiveQueryList = nil

func NewArchiveQueryList(size int) *ArchiveQueryList {
  var list ArchiveQueryList = ArchiveQueryList(make([]*ArchiveQuery, size))
  return &list
}

// ================================================================================
// Defines ArchiveQueryList type as an ElementList

func (receiver *ArchiveQueryList) Size() int {
  if receiver != nil {
    return len(*receiver)
  }
  return -1
}

func (receiver *ArchiveQueryList) GetElementAt(i int) mal.Element {
  if receiver == nil || i >= receiver.Size() {
    return nil
  }
  return (*receiver)[i]
}

func (receiver *ArchiveQueryList) AppendElement(element mal.Element) {
  if receiver != nil {
    *receiver = append(*receiver, element.(*ArchiveQuery))
  }
}

// ================================================================================
// Defines ArchiveQueryList type as a MAL Composite

func (receiver *ArchiveQueryList) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines ArchiveQueryList type as a MAL Element

const ARCHIVEQUERY_LIST_TYPE_SHORT_FORM mal.Integer = -2
const ARCHIVEQUERY_LIST_SHORT_FORM mal.Long = 0x2000201fffffe

// Registers ArchiveQueryList type for polymorphism handling
func init() {
  mal.RegisterMALElement(ARCHIVEQUERY_LIST_SHORT_FORM, NullArchiveQueryList)
}

// Returns the absolute short form of the element type.
func (receiver *ArchiveQueryList) GetShortForm() mal.Long {
  return ARCHIVEQUERY_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ArchiveQueryList) GetAreaNumber() mal.UShort {
  return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ArchiveQueryList) GetAreaVersion() mal.UOctet {
  return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ArchiveQueryList) GetServiceNumber() mal.UShort {
    return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ArchiveQueryList) GetTypeShortForm() mal.Integer {
  return ARCHIVEQUERY_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ArchiveQueryList) CreateElement() mal.Element {
  return NewArchiveQueryList(0)
}

func (receiver *ArchiveQueryList) IsNull() bool {
  return receiver == nil
}

func (receiver *ArchiveQueryList) Null() mal.Element {
  return NullArchiveQueryList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ArchiveQueryList) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(ARCHIVEQUERY_LIST_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeUInteger(mal.NewUInteger(uint32(len([]*ArchiveQuery(*receiver)))))
  if err != nil {
    return err
  }
  for _, e := range []*ArchiveQuery(*receiver) {
    encoder.EncodeNullableElement(e)
  }
  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ArchiveQueryList) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(ARCHIVEQUERY_LIST_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  size, err := decoder.DecodeUInteger()
  if err != nil {
    return nil, err
  }
  list := ArchiveQueryList(make([]*ArchiveQuery, int(*size)))
  for i := 0; i < len(list); i++ {
    elem, err := decoder.DecodeNullableElement(NullArchiveQuery)
    if err != nil {
      return nil, err
    }
    list[i] = elem.(*ArchiveQuery)
  }
  return &list, nil
}
