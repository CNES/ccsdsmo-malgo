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

// ================================================================================
// Encoding interface, implemented by specific encoding technology.
type Encoder interface {
	Body() []byte

	EncodeNull() error

	EncodeNotNull() error

	// Note (AF): The EncodeNullable methods for attribute are no really needed as nullable
	// attribute can be encode with EncodeNullableElement (EncodeNullableAttribute method
	// are also useless).

	// Encodes a non-null Boolean.
	// @param att The Boolean to encode.
	EncodeBoolean(att *Boolean) error

	// Encodes a Boolean that may be null
	// @param att The Boolean to encode.
	EncodeNullableBoolean(att *Boolean) error

	// Encodes a non-null Float.
	// @param att The Float to encode.
	EncodeFloat(att *Float) error

	// Encodes a Float that may be null
	// @param att The Float to encode.
	EncodeNullableFloat(att *Float) error

	// Encodes a non-null Double.
	// @param att The Double to encode.
	EncodeDouble(att *Double) error

	// Encodes a Double that may be null
	// @param att The Double to encode.
	EncodeNullableDouble(att *Double) error

	// Encodes a non-null Octet.
	// @param att The Octet to encode.
	EncodeOctet(att *Octet) error

	// Encodes an Octet that may be null
	// @param att The Octet to encode.
	EncodeNullableOctet(att *Octet) error

	// Encodes a non-null UOctet.
	// @param att The UOctet to encode.
	EncodeUOctet(att *UOctet) error

	// Encodes a UOctet that may be null
	// @param att The UOctet to encode.
	EncodeNullableUOctet(att *UOctet) error

	// Encodes a non-null Short.
	// @param att The Short to encode.
	EncodeShort(att *Short) error

	// Encodes a Short that may be null
	// @param att The Short to encode.
	EncodeNullableShort(att *Short) error

	// Encodes a non-null UShort.
	// @param att The UShort to encode.
	EncodeUShort(att *UShort) error

	// Encodes a UShort that may be null
	// @param att The UShort to encode.
	EncodeNullableUShort(att *UShort) error

	// Encodes a non-null Integer.
	// @param att The Integer to encode.
	EncodeInteger(att *Integer) error

	// Encodes an Integer that may be null
	// @param att The Integer to encode.
	EncodeNullableInteger(att *Integer) error

	// Encodes a non-null UInteger.
	// @param att The UInteger to encode.
	EncodeUInteger(att *UInteger) error

	// Encodes a UInteger that may be null
	// @param att The UInteger to encode.
	EncodeNullableUInteger(att *UInteger) error

	// Encodes a non-null Long.
	// @param att The Long to encode.
	EncodeLong(att *Long) error

	// Encodes a Long that may be null
	// @param att The Long to encode.
	EncodeNullableLong(att *Long) error

	// Encodes a non-null ULong.
	// @param att The ULong to encode.
	EncodeULong(att *ULong) error

	// Encodes a ULong that may be null
	// @param att The ULong to encode.
	EncodeNullableULong(att *ULong) error

	// Encodes a non-null String.
	// @param att The String to encode.
	EncodeString(att *String) error

	// Encodes a String that may be null
	// @param att The String to encode.
	EncodeNullableString(att *String) error

	// Encodes a non-null Blob.
	// @param att The Blob to encode.
	EncodeBlob(att *Blob) error

	// Encodes a Blob that may be null
	// @param att The Blob to encode.
	EncodeNullableBlob(att *Blob) error

	// Encodes a non-null Duration.
	// @param att The Duration to encode.
	EncodeDuration(att *Duration) error

	// Encodes a Duration that may be null
	// @param att The Duration to encode.
	EncodeNullableDuration(att *Duration) error

	// Encodes a non-null FineTime.
	// @param att The FineTime to encode.
	EncodeFineTime(att *FineTime) error

	// Encodes a FineTime that may be null
	// @param att The FineTime to encode.
	EncodeNullableFineTime(att *FineTime) error

	// Encodes a non-null Identifier.
	// @param att The Identifier to encode.
	EncodeIdentifier(att *Identifier) error

	// Encodes an Identifier that may be null
	// @param att The Identifier to encode.
	EncodeNullableIdentifier(att *Identifier) error

	// Encodes a non-null Time.
	// @param att The Time to encode.
	EncodeTime(att *Time) error

	// Encodes a Time that may be null
	// @param att The Time to encode.
	EncodeNullableTime(att *Time) error

	// Encodes a non-null URI.
	// @param att The URI to encode.
	// @throws IllegalArgumentException If the argument is null.
	EncodeURI(att *URI) error

	// Encodes a URI that may be null
	// @param att The URI to encode.
	EncodeNullableURI(att *URI) error

	// TODO (AF): Handling of enumeration

	EncodeSmallEnum(ordinal uint8) error
	EncodeMediumEnum(ordinal uint16) error
	EncodelargeEnum(ordinal uint32) error

	// Encodes a non-null Element.
	// @param element The Element to encode.
	EncodeElement(element Element) error

	// Encodes an Element that may be null
	// @param element The Element to encode.
	EncodeNullableElement(element Element) error

	// Encodes a non-null abstract Element (use for polymorphism).
	// @param element The Element to encode.
	EncodeAbstractElement(element Element) error

	// Encodes an abstract Element that may be null (use for polymorphism).
	// @param element The Element to encode.
	EncodeNullableAbstractElement(element Element) error

	// Encodes the short form of an attribute.
	EncodeAttributeType(Integer) error

	// Encodes a non-null Attribute.
	// @param att The Attribute to encode.
	EncodeAttribute(att Attribute) error

	// Encodes an Attribute that may be null
	// @param att The Attribute to encode.
	EncodeNullableAttribute(att Attribute) error

	// Encodes a list of Element given as a slice of Element.
	// Should only use to encode List< <<Update Value Type>> > in Broker.
	// @param list The list of Element to encode
	EncodeElementList(list []Element) error

	EncodeList(list ElementList) error
	EncodeNullableList(list ElementList) error

	// Gets a specific encoder for the specified type
	LookupSpecific(shortForm Long) SpecificEncoder
}

type SpecificEncoder func(element Element, encoder Encoder) error

type GenEncoder struct {
	Self Encoder

	// Registry for specific encoding functions
	Registry map[int64]SpecificEncoder
}

// TODO (AF): Move all nullable from binary decoder..

// Encodes a Boolean that may be null
// @param att The Boolean to encode.
func (encoder *GenEncoder) EncodeNullableBoolean(att *Boolean) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeBoolean(att)
	}
}

// Encodes a Float that may be null
// @param att The Float to encode.
func (encoder *GenEncoder) EncodeNullableFloat(att *Float) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeFloat(att)
	}
}

// Encodes a Double that may be null
// @param att The Double to encode.
func (encoder *GenEncoder) EncodeNullableDouble(att *Double) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeDouble(att)
	}
}

// Encodes an Octet that may be null
// @param att The Octet to encode.
func (encoder *GenEncoder) EncodeNullableOctet(att *Octet) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeOctet(att)
	}
}

// Encodes a UOctet that may be null
// @param att The UOctet to encode.
func (encoder *GenEncoder) EncodeNullableUOctet(att *UOctet) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeUOctet(att)
	}
}

// Encodes a Short that may be null
// @param att The Short to encode.
func (encoder *GenEncoder) EncodeNullableShort(att *Short) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeShort(att)
	}
}

// Encodes a UShort that may be null
// @param att The UShort to encode.
func (encoder *GenEncoder) EncodeNullableUShort(att *UShort) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeUShort(att)
	}
}

// Encodes an Integer that may be null
// @param att The Integer to encode.
func (encoder *GenEncoder) EncodeNullableInteger(att *Integer) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeInteger(att)
	}
}

// Encodes a UInteger that may be null
// @param att The UInteger to encode.
func (encoder *GenEncoder) EncodeNullableUInteger(att *UInteger) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeUInteger(att)
	}
}

// Encodes a Long that may be null
// @param att The Long to encode.
func (encoder *GenEncoder) EncodeNullableLong(att *Long) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeLong(att)
	}
}

// Encodes a ULong that may be null
// @param att The ULong to encode.
func (encoder *GenEncoder) EncodeNullableULong(att *ULong) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeULong(att)
	}
}

// Encodes a String that may be null
// @param att The String to encode.
func (encoder *GenEncoder) EncodeNullableString(str *String) error {
	if str == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeString(str)
	}
}

// Encodes a Blob that may be null
// @param att The Blob to encode.
func (encoder *GenEncoder) EncodeNullableBlob(blob *Blob) error {
	if blob == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeBlob(blob)
	}
}

// Encodes an Identifier that may be null
// @param att The Identifier to encode.
func (encoder *GenEncoder) EncodeNullableIdentifier(att *Identifier) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeIdentifier(att)
	}
}

// Encodes a Duration that may be null
// @param att The Duration to encode.
func (encoder *GenEncoder) EncodeNullableDuration(att *Duration) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeDuration(att)
	}
}

// Encodes a Time that may be null
// @param att The Time to encode.
func (encoder *GenEncoder) EncodeNullableTime(att *Time) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeTime(att)
	}
}

// Encodes a FineTime that may be null
// @param att The FineTime to encode.
func (encoder *GenEncoder) EncodeNullableFineTime(att *FineTime) error {
	if att == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeFineTime(att)
	}
}

// Encodes a URI that may be null
// @param att The URI to encode.
func (encoder *GenEncoder) EncodeNullableURI(uri *URI) error {
	if uri == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeURI(uri)
	}
}

// Encodes a non-null Element.
// @param element The Element to encode.
func (encoder *GenEncoder) EncodeElement(element Element) error {
	return element.Encode(encoder.Self)
}

// Encodes an Element that may be null
// @param element The Element to encode.
func (encoder *GenEncoder) EncodeNullableElement(element Element) error {
	if element == element.Null() { // element.IsNull()
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return element.Encode(encoder.Self)
	}
}

// Encodes a non-null abstract Element (use for polymorphism).
// @param element The Element to encode.
func (encoder *GenEncoder) EncodeAbstractElement(element Element) error {
	shortForm := element.GetShortForm()
	err := encoder.Self.EncodeLong(&shortForm)
	if err != nil {
		return err
	}
	return element.Encode(encoder.Self)
}

// Encodes an abstract Element that may be null (use for polymorphism).
// @param element The Element to encode.
func (encoder *GenEncoder) EncodeNullableAbstractElement(element Element) error {
	if element == element.Null() { // element.IsNull()
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeAbstractElement(element)
	}
}

// Encodes a non-null Attribute.
// @param att The Attribute to encode.
func (encoder *GenEncoder) EncodeAttribute(att Attribute) error {
	err := encoder.Self.EncodeAttributeType(att.GetTypeShortForm())
	if err != nil {
		return err
	}
	return att.Encode(encoder.Self)
}

// Encodes an Attribute that may be null
// @param att The Attribute to encode.
func (encoder *GenEncoder) EncodeNullableAttribute(att Attribute) error {
	if att == att.Null() { // att.IsNull()
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeAttribute(att)
	}
}

// Encodes a list of Element given as a slice of Element.
// Should only use to encode List< <<Update Value Type>> > in Broker.
// @param list The list of Element to encode
func (encoder *GenEncoder) EncodeElementList(list []Element) error {
	// Computes the short form of the list.
	shortForm := list[0].GetShortForm()
	x := -int32((int64(shortForm) & 0x00FFFFFF) | 0xFF000000)
	y := (int64(shortForm) & mask) | int64(x)
	shortForm = Long(y)
	err := encoder.Self.EncodeLong(&shortForm)
	if err != nil {
		return err
	}
	err = encoder.Self.EncodeUInteger(NewUInteger(uint32(len(list))))
	if err != nil {
		return err
	}
	for _, e := range list {
		encoder.Self.EncodeNullableElement(e)
	}
	return nil
}

// Note (AF): This code below provides a generic view of list as Element slices.
// It must be enhanced in order to offer somes details (verification that all
// ELement  implements the same concrete type for exemple). It needs also to add
// methods encoding list as generic types (writing the corresponding short form).
// Finally it should not replace the actual list view, but rather offers an
// alternative to developper.

//func (encoder *GenEncoder) EncodeList(list []Element) error {
//	err := encoder.Self.EncodeUInteger(NewUInteger(uint32(len(list))))
//	if err != nil {
//		return err
//	}
//	for _, e := range list {
//		encoder.Self.EncodeNullableElement(e)
//	}
//	return nil
//}
//
//func (encoder *GenEncoder) EncodeNullableList(list []Element) error {
//	if list == nil {
//		return encoder.Self.EncodeNull()
//	} else {
//		err := encoder.Self.EncodeNotNull()
//		if err != nil {
//			return err
//		}
//		return encoder.EncodeList(list)
//	}
//}

// Use the ElementList interface instead of the []Element representation

func (encoder *GenEncoder) EncodeList(list ElementList) error {
	err := encoder.Self.EncodeUInteger(NewUInteger(uint32(list.Size())))
	if err != nil {
		return err
	}
	for i := 0; i < list.Size(); i++ {
		encoder.Self.EncodeNullableElement(list.GetElementAt(i))
	}
	return nil
}

func (encoder *GenEncoder) EncodeNullableList(list ElementList) error {
	if list == nil {
		return encoder.Self.EncodeNull()
	} else {
		err := encoder.Self.EncodeNotNull()
		if err != nil {
			return err
		}
		return encoder.Self.EncodeList(list)
	}
}

// Functions allowing to handle specific encoders

func NewEncoderRegistry() map[int64]SpecificEncoder {
	return make(map[int64]SpecificEncoder)
}

func RegisterSpecificEncoder(registry map[int64]SpecificEncoder, shortForm Long, specific SpecificEncoder) error {
	rlogger.Debugf("EncoderRegistry.RegisterSpecific: %x", (int64)(shortForm))
	_, ok := registry[(int64)(shortForm)]
	if ok {
		rlogger.Errorf("EncoderRegistry.RegisterSpecific: %x already registered", (int64)(shortForm))
		return errors.New("EncoderRegistry.RegisterSpecific: already registered")
	}
	registry[(int64)(shortForm)] = specific
	return nil
}

func (encoder *GenEncoder) LookupSpecific(shortForm Long) SpecificEncoder {
	if encoder.Registry == nil {
		return nil
	}
	return LookupSpecificEncoder(encoder.Registry, shortForm)
}

func LookupSpecificEncoder(registry map[int64]SpecificEncoder, shortForm Long) SpecificEncoder {
	rlogger.Debugf("EncoderRegistry.LookupSpecific: %x", (int64)(shortForm))
	specific, ok := registry[(int64)(shortForm)]
	if !ok {
		rlogger.Debugf("EncoderRegistry.LookupSpecific: unknown %x element", (int64)(shortForm))
		return nil
	}
	return specific
}

func DeregisterSpecificEncoder(registry map[int64]SpecificEncoder, shortForm Long) error {
	rlogger.Debugf("EncoderRegistry.DeregisterSpecific: %x", (int64)(shortForm))
	_, ok := registry[(int64)(shortForm)]
	if !ok {
		rlogger.Errorf("EncoderRegistry.DeregisterSpecific: %x not registered", (int64)(shortForm))
		return errors.New("EncoderRegistry.DeregisterSpecific: not registered")
	}
	delete(registry, (int64)(shortForm))
	return nil
}
