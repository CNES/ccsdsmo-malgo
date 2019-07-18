/**
 * MIT License
 *
 * Copyright (c) 2017 - 2019 CNES
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
package mal

import (
	"errors"
)

// Decoding interface, implemented by specific decoding technology.
type Decoder interface {
	// Returns true is the next value is NULL.
	IsNull() (bool, error)

	// Note (AF): The DecodeNullable methods for attribute are no really needed as nullable
	// attribute can be decode with DecodeNullableElement (DecodeNullableAttribute method
	// are also useless).

	// Decodes a Boolean.
	// @return The decoded Boolean.
	DecodeBoolean() (*Boolean, error)

	// Decodes a Boolean that may be null.
	// @return The decoded Boolean or null.
	DecodeNullableBoolean() (*Boolean, error)

	// Decodes a Float.
	// @return The decoded Float.
	DecodeFloat() (*Float, error)

	// Decodes a Float that may be null.
	// @return The decoded Float or null.
	DecodeNullableFloat() (*Float, error)

	// Decodes a Double.
	// @return The decoded Double.
	DecodeDouble() (*Double, error)

	// Decodes a Double that may be null.
	// @return The decoded Double or null.
	DecodeNullableDouble() (*Double, error)

	// Decodes an Octet.
	// @return The decoded Octet.
	DecodeOctet() (*Octet, error)

	// Decodes an Octet that may be null.
	// @return The decoded Octet or null.
	DecodeNullableOctet() (*Octet, error)

	// Decodes a UOctet.
	// @return The decoded UOctet.
	DecodeUOctet() (*UOctet, error)

	// Decodes a UOctet that may be null.
	// @return The decoded UOctet or null.
	DecodeNullableUOctet() (*UOctet, error)

	// Decodes a Short.
	// @return The decoded Short.
	DecodeShort() (*Short, error)

	// Decodes a Short that may be null.
	// @return The decoded Short or null.
	DecodeNullableShort() (*Short, error)

	// Decodes a UShort.
	// @return The decoded UShort.
	DecodeUShort() (*UShort, error)

	// Decodes a UShort that may be null.
	// @return The decoded UShort or null.
	DecodeNullableUShort() (*UShort, error)

	// Decodes an Integer.
	// @return The decoded Integer.
	DecodeInteger() (*Integer, error)

	// Decodes an Integer that may be null.
	// @return The decoded Integer or null.
	DecodeNullableInteger() (*Integer, error)

	// Decodes a UInteger.
	// @return The decoded UInteger.
	DecodeUInteger() (*UInteger, error)

	// Decodes a UInteger that may be null.
	// @return The decoded UInteger or null.
	DecodeNullableUInteger() (*UInteger, error)

	// Decodes a Long.
	// @return The decoded Long.
	DecodeLong() (*Long, error)

	// Decodes a Long that may be null.
	// @return The decoded Long or null.
	DecodeNullableLong() (*Long, error)

	// Decodes a ULong.
	// @return The decoded ULong.
	DecodeULong() (*ULong, error)

	// Decodes a ULong that may be null.
	// @return The decoded ULong or null.
	DecodeNullableULong() (*ULong, error)

	// Decodes a String.
	// @return The decoded String.
	DecodeString() (*String, error)

	// Decodes a String that may be null.
	// @return The decoded String or null.
	DecodeNullableString() (*String, error)

	// Decodes a Blob.
	// @return The decoded Blob.
	DecodeBlob() (*Blob, error)

	// Decodes a Blob that may be null.
	// @return The decoded Blob or null.
	DecodeNullableBlob() (*Blob, error)

	// Decodes a Duration.
	// @return The decoded Duration.
	DecodeDuration() (*Duration, error)

	// Decodes a Duration that may be null.
	// @return The decoded Duration or null.
	DecodeNullableDuration() (*Duration, error)

	// Decodes a FineTime.
	// @return The decoded FineTime.
	DecodeFineTime() (*FineTime, error)

	// Decodes a FineTime that may be null.
	// @return The decoded FineTime or null.
	DecodeNullableFineTime() (*FineTime, error)

	// Decodes an Identifier.
	// @return The decoded Identifier.
	DecodeIdentifier() (*Identifier, error)

	// Decodes an Identifier that may be null.
	// @return The decoded Identifier or null.
	DecodeNullableIdentifier() (*Identifier, error)

	// Decodes a Time.
	// @return The decoded Time.
	DecodeTime() (*Time, error)

	// Decodes a Time that may be null.
	// @return The decoded Time or null.
	DecodeNullableTime() (*Time, error)

	// Decodes a URI.
	// @return The decoded URI.
	DecodeURI() (*URI, error)

	// Decodes a URI that may be null.
	// @return The decoded URI or null.
	DecodeNullableURI() (*URI, error)

	// TODO (AF): Handling of enumeration

	DecodeSmallEnum() (uint8, error)
	DecodeMediumEnum() (uint16, error)
	DecodelargeEnum() (uint32, error)

	// Decodes an Element.
	// @param element An instance of the element to decode.
	// @return The decoded Element.
	DecodeElement(element Element) (Element, error)

	// Decodes an Element that may be null.
	// @param element An instance of the element to decode.
	// @return The decoded Element or null.
	DecodeNullableElement(element Element) (Element, error)

	// Decodes an abstract Element using polymorphism.
	// @return The decoded Element.
	DecodeAbstractElement() (Element, error)

	// Decodes an abstract Element that may be null using polymorphism.
	// @return The decoded Element or null.
	DecodeNullableAbstractElement() (Element, error)

	// Decodes the short form of an attribute.
	// @return The short form of the attribute.
	DecodeAttributeType() (Integer, error)

	// Decodes an Attribute.
	// @return The decoded Attribute.
	DecodeAttribute() (Attribute, error)

	// Decodes an Attribute that may be null.
	// @return The decoded Attribute or null.
	DecodeNullableAttribute() (Attribute, error)

	// Decodes a list of Element as a slice of Element.
	// Should only use to decode List< <<Update Value Type>> > in Broker.
	// @return The decoded list ofElement
	DecodeElementList() ([]Element, error)

	// Gets a specific decoder for the specified type
	LookupSpecific(shortForm Long) SpecificDecoder
}

type SpecificDecoder func(decoder Decoder) (Element, error)

type GenDecoder struct {
	Self Decoder

	// Registry for specific decoding functions
	Registry map[int64]SpecificDecoder
}

// Decodes a Boolean that may be null.
// @return The decoded Boolean or null.
func (decoder *GenDecoder) DecodeNullableBoolean() (*Boolean, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeBoolean()
	}
}

// Decodes a Float that may be null.
// @return The decoded Float or null.
func (decoder *GenDecoder) DecodeNullableFloat() (*Float, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeFloat()
	}
}

// Decodes a Double that may be null.
// @return The decoded Double or null.
func (decoder *GenDecoder) DecodeNullableDouble() (*Double, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeDouble()
	}
}

// Decodes an Octet that may be null.
// @return The decoded Octet or null.
func (decoder *GenDecoder) DecodeNullableOctet() (*Octet, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeOctet()
	}
}

// Decodes a UOctet that may be null.
// @return The decoded UOctet or null.
func (decoder *GenDecoder) DecodeNullableUOctet() (*UOctet, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeUOctet()
	}
}

// Decodes a Short that may be null.
// @return The decoded Short or null.
func (decoder *GenDecoder) DecodeNullableShort() (*Short, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeShort()
	}
}

// Decodes a UShort that may be null.
// @return The decoded UShort or null.
func (decoder *GenDecoder) DecodeNullableUShort() (*UShort, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeUShort()
	}
}

// Decodes an Integer that may be null.
// @return The decoded Integer or null.
func (decoder *GenDecoder) DecodeNullableInteger() (*Integer, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeInteger()
	}
}

// Decodes a UInteger that may be null.
// @return The decoded UInteger or null.
func (decoder *GenDecoder) DecodeNullableUInteger() (*UInteger, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeUInteger()
	}
}

// Decodes a Long that may be null.
// @return The decoded Long or null.
func (decoder *GenDecoder) DecodeNullableLong() (*Long, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeLong()
	}
}

// Decodes a ULong that may be null.
// @return The decoded ULong or null.
func (decoder *GenDecoder) DecodeNullableULong() (*ULong, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeULong()
	}
}

// Decodes a String that may be null.
// @return The decoded String or null.
func (decoder *GenDecoder) DecodeNullableString() (*String, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeString()
	}
}

// Decodes a Blob that may be null.
// @return The decoded Blob or null.
func (decoder *GenDecoder) DecodeNullableBlob() (*Blob, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeBlob()
	}
}

// Decodes an Identifier that may be null.
// @return The decoded Identifier or null.
func (decoder *GenDecoder) DecodeNullableIdentifier() (*Identifier, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeIdentifier()
	}
}

// Decodes a Duration that may be null.
// @return The decoded Duration or null.
func (decoder *GenDecoder) DecodeNullableDuration() (*Duration, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeDuration()
	}
}

// Decodes a Time that may be null.
// @return The decoded Time or null.
func (decoder *GenDecoder) DecodeNullableTime() (*Time, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeTime()
	}
}

// Decodes a FineTime that may be null.
// @return The decoded FineTime or null.
func (decoder *GenDecoder) DecodeNullableFineTime() (*FineTime, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeFineTime()
	}
}

// Decodes a URI that may be null.
// @return The decoded URI or null.
func (decoder *GenDecoder) DecodeNullableURI() (*URI, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeURI()
	}
}

// Decodes an Element.
// @param element An instance of the element to decode.
// @return The decoded Element.
func (decoder *GenDecoder) DecodeElement(element Element) (Element, error) {
	return element.Decode(decoder.Self)
}

// Decodes an Element that may be null.
// @param element An instance of the element to decode.
// @return The decoded Element or null.
func (decoder *GenDecoder) DecodeNullableElement(element Element) (Element, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return element.Null(), nil
	} else {
		return element.Decode(decoder.Self)
	}
}

// Decodes an abstract Element using polymorphism.
// @return The decoded Element.
func (decoder *GenDecoder) DecodeAbstractElement() (Element, error) {
	shortForm, err := decoder.Self.DecodeLong()
	if err != nil {
		return nil, err
	}
	element, err := LookupMALElement(*shortForm)
	return element.Decode(decoder.Self)
}

// Decodes an abstract Element that may be null using polymorphism.
// @return The decoded Element or null.
func (decoder *GenDecoder) DecodeNullableAbstractElement() (Element, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return NullElement, nil
	} else {
		return decoder.Self.DecodeAbstractElement()
	}
}

// Decodes an Attribute.
// @return The decoded Attribute.
func (decoder *GenDecoder) DecodeAttribute() (Attribute, error) {
	typeval, err := decoder.Self.DecodeAttributeType()
	if err != nil {
		return nil, err
	}
	switch typeval {
	case MAL_BLOB_TYPE_SHORT_FORM:
		return decoder.Self.DecodeBlob()
	case MAL_BOOLEAN_TYPE_SHORT_FORM:
		return decoder.Self.DecodeBoolean()
	case MAL_DURATION_TYPE_SHORT_FORM:
		return decoder.Self.DecodeDuration()
	case MAL_FLOAT_TYPE_SHORT_FORM:
		return decoder.Self.DecodeFloat()
	case MAL_DOUBLE_TYPE_SHORT_FORM:
		return decoder.Self.DecodeDouble()
	case MAL_IDENTIFIER_TYPE_SHORT_FORM:
		return decoder.Self.DecodeIdentifier()
	case MAL_OCTET_TYPE_SHORT_FORM:
		return decoder.Self.DecodeOctet()
	case MAL_UOCTET_TYPE_SHORT_FORM:
		return decoder.Self.DecodeUOctet()
	case MAL_SHORT_TYPE_SHORT_FORM:
		return decoder.Self.DecodeShort()
	case MAL_USHORT_TYPE_SHORT_FORM:
		return decoder.Self.DecodeUShort()
	case MAL_INTEGER_TYPE_SHORT_FORM:
		return decoder.Self.DecodeInteger()
	case MAL_UINTEGER_TYPE_SHORT_FORM:
		return decoder.Self.DecodeUInteger()
	case MAL_LONG_TYPE_SHORT_FORM:
		return decoder.Self.DecodeLong()
	case MAL_ULONG_TYPE_SHORT_FORM:
		return decoder.Self.DecodeULong()
	case MAL_STRING_TYPE_SHORT_FORM:
		return decoder.Self.DecodeString()
	case MAL_TIME_TYPE_SHORT_FORM:
		return decoder.Self.DecodeTime()
	case MAL_FINETIME_TYPE_SHORT_FORM:
		return decoder.Self.DecodeFineTime()
	case MAL_URI_TYPE_SHORT_FORM:
		return decoder.Self.DecodeURI()
	default:
		return nil, errors.New("Unknow attribute: " + string(typeval))
	}
}

// Decodes an Attribute that may be null.
// @return The decoded Attribute or null.
func (decoder *GenDecoder) DecodeNullableAttribute() (Attribute, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.Self.DecodeAttribute()
	}
}

const mask int64 = -16777216 // 0xFFFFFFFFFF000000

// Decodes a list of Element as a slice of Element.
// Only use to decode List< <<Update Value Type>> > in Broker.
// @return The decoded list ofElement
func (decoder *GenDecoder) DecodeElementList() ([]Element, error) {
	shortForm, err := decoder.Self.DecodeLong()
	if err != nil {
		return nil, err
	}
	// This is the short form of the list, so we computes the short
	// form of the list entry.
	x := -int32((int64(*shortForm) & 0x00FFFFFF) | 0xFF000000)
	y := (int64(*shortForm) & mask) | int64(x)
	shortForm = NewLong(y)
	element, err := LookupMALElement(*shortForm)
	size, err := decoder.Self.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := make([]Element, int(*size))
	for i := 0; i < len(list); i++ {
		list[i], err = decoder.Self.DecodeNullableElement(element)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

// Note (AF): see corresponding comment in encoder about generic view for lists.

//func (decoder *GenDecoder) DecodeList(element Element) ([]Element, error) {
//	size, err := decoder.Self.DecodeUInteger()
//	if err != nil {
//		return nil, err
//	}
//	list := make([]Element, int(*size))
//	for i := 0; i < len(list); i++ {
//		list[i], err = decoder.Self.DecodeNullableElement(element)
//		if err != nil {
//			return nil, err
//		}
//	}
//	return list, nil
//}
//
//func (decoder *GenDecoder) DecodeNullableList(element Element) ([]Element, error) {
//	null, err := decoder.Self.IsNull()
//	if err != nil {
//		return nil, err
//	}
//	if null {
//		return nil, nil
//	} else {
//		return decoder.DecodeList(element)
//	}
//}

// Use interface ElementList instead of []Element
func (decoder *GenDecoder) DecodeList(element Element) (ElementList, error) {
	size, err := decoder.Self.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	seed, err := LookupMALElement(GetListShortForm(element))
	list := seed.CreateElement().(ElementList)
	for i := 0; i < int(*size); i++ {
		item, err := decoder.Self.DecodeNullableElement(element)
		if err != nil {
			return nil, err
		}
		list.AppendElement(item)
	}
	return list, nil
}

func (decoder *GenDecoder) DecodeNullableList(element Element) (ElementList, error) {
	null, err := decoder.Self.IsNull()
	if err != nil {
		return nil, err
	}
	if null {
		return nil, nil
	} else {
		return decoder.DecodeList(element)
	}
}

// Functions allowing to handle specific decoders

func NewDecoderRegistry() map[int64]SpecificDecoder {
	return make(map[int64]SpecificDecoder)
}

func RegisterSpecificDecoder(registry map[int64]SpecificDecoder, shortForm Long, specific SpecificDecoder) error {
	rlogger.Debugf("DecoderRegistry.RegisterSpecific: %x", (int64)(shortForm))
	_, ok := registry[(int64)(shortForm)]
	if ok {
		rlogger.Errorf("DecoderRegistry.RegisterSpecific: %x already registered", (int64)(shortForm))
		return errors.New("DecoderRegistry.RegisterSpecific: already registered")
	}
	registry[(int64)(shortForm)] = specific
	return nil
}

func (decoder *GenDecoder) LookupSpecific(shortForm Long) SpecificDecoder {
	if decoder.Registry == nil {
		return nil
	}
	return LookupSpecificDecoder(decoder.Registry, shortForm)
}

func LookupSpecificDecoder(registry map[int64]SpecificDecoder, shortForm Long) SpecificDecoder {
	rlogger.Debugf("DecoderRegistry.LookupSpecific: %x", (int64)(shortForm))
	specific, ok := registry[(int64)(shortForm)]
	if !ok {
		rlogger.Debugf("DecoderRegistry.LookupSpecific: unknown %x element", (int64)(shortForm))
		return nil
	}
	return specific
}

func DeregisterSpecificDecoder(registry map[int64]SpecificDecoder, shortForm Long) error {
	rlogger.Debugf("DecoderRegistry.DeregisterSpecific: %x", (int64)(shortForm))
	_, ok := registry[(int64)(shortForm)]
	if !ok {
		rlogger.Errorf("DecoderRegistry.DeregisterSpecific: %x not registered", (int64)(shortForm))
		return errors.New("DecoderRegistry.DeregisterSpecific: not registered")
	}
	delete(registry, (int64)(shortForm))
	return nil
}
