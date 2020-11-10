package archive

import (
  "github.com/CNES/ccsdsmo-malgo/mal"
  "github.com/CNES/ccsdsmo-malgo/com"
)

// Defines ArchiveQuery type

type ArchiveQuery struct {
  Domain *mal.IdentifierList
  Network *mal.Identifier
  Provider *mal.URI
  Related mal.Long
  Source *com.ObjectId
  StartTime *mal.FineTime
  EndTime *mal.FineTime
  SortOrder *mal.Boolean
  SortFieldName *mal.String
}

var (
  NullArchiveQuery *ArchiveQuery = nil
)
func NewArchiveQuery() *ArchiveQuery {
  return new(ArchiveQuery)
}

// ================================================================================
// Defines ArchiveQuery type as a MAL Composite

func (receiver *ArchiveQuery) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines ArchiveQuery type as a MAL Element

const ARCHIVEQUERY_TYPE_SHORT_FORM mal.Integer = 2
const ARCHIVEQUERY_SHORT_FORM mal.Long = 0x2000201000002

// Registers ArchiveQuery type for polymorphism handling
func init() {
  mal.RegisterMALElement(ARCHIVEQUERY_SHORT_FORM, NullArchiveQuery)
}

// Returns the absolute short form of the element type.
func (receiver *ArchiveQuery) GetShortForm() mal.Long {
  return ARCHIVEQUERY_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ArchiveQuery) GetAreaNumber() mal.UShort {
  return com.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ArchiveQuery) GetAreaVersion() mal.UOctet {
  return com.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ArchiveQuery) GetServiceNumber() mal.UShort {
    return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ArchiveQuery) GetTypeShortForm() mal.Integer {
  return ARCHIVEQUERY_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ArchiveQuery) CreateElement() mal.Element {
  return new(ArchiveQuery)
}

func (receiver *ArchiveQuery) IsNull() bool {
  return receiver == nil
}

func (receiver *ArchiveQuery) Null() mal.Element {
  return NullArchiveQuery
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ArchiveQuery) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(ARCHIVEQUERY_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeNullableElement(receiver.Domain)
  if err != nil {
    return err
  }
  err = encoder.EncodeNullableIdentifier(receiver.Network)
  if err != nil {
    return err
  }
  err = encoder.EncodeNullableURI(receiver.Provider)
  if err != nil {
    return err
  }
  err = encoder.EncodeLong(&receiver.Related)
  if err != nil {
    return err
  }
  err = encoder.EncodeNullableElement(receiver.Source)
  if err != nil {
    return err
  }
  err = encoder.EncodeNullableFineTime(receiver.StartTime)
  if err != nil {
    return err
  }
  err = encoder.EncodeNullableFineTime(receiver.EndTime)
  if err != nil {
    return err
  }
  err = encoder.EncodeNullableBoolean(receiver.SortOrder)
  if err != nil {
    return err
  }
  err = encoder.EncodeNullableString(receiver.SortFieldName)
  if err != nil {
    return err
  }

  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ArchiveQuery) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(ARCHIVEQUERY_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  Domain, err := decoder.DecodeNullableElement(mal.NullIdentifierList)
  if err != nil {
    return nil, err
  }
  Network, err := decoder.DecodeNullableIdentifier()
  if err != nil {
    return nil, err
  }
  Provider, err := decoder.DecodeNullableURI()
  if err != nil {
    return nil, err
  }
  Related, err := decoder.DecodeLong()
  if err != nil {
    return nil, err
  }
  Source, err := decoder.DecodeNullableElement(com.NullObjectId)
  if err != nil {
    return nil, err
  }
  StartTime, err := decoder.DecodeNullableFineTime()
  if err != nil {
    return nil, err
  }
  EndTime, err := decoder.DecodeNullableFineTime()
  if err != nil {
    return nil, err
  }
  SortOrder, err := decoder.DecodeNullableBoolean()
  if err != nil {
    return nil, err
  }
  SortFieldName, err := decoder.DecodeNullableString()
  if err != nil {
    return nil, err
  }

  var composite = ArchiveQuery {
    Domain: Domain.(*mal.IdentifierList),
    Network: Network,
    Provider: Provider,
    Related: *Related,
    Source: Source.(*com.ObjectId),
    StartTime: StartTime,
    EndTime: EndTime,
    SortOrder: SortOrder,
    SortFieldName: SortFieldName,
  }
  return &composite, nil
}
