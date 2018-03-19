/**
 * MIT License
 *
 * Copyright (c) 2018 CNES
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
package com

import (
	. "github.com/ccsdsmo/malgo/mal"
)

// ################################################################################
// Defines COM ObjectType type
// ################################################################################

/**
 * The ObjectType structure uniquely identifies the type of an object. It is the combination
 * of the area number, service number, area version, and service object type number. The combined
 * parts are able to fit inside a MAL::Long (for implementations that prefer to index on a single
 * numeric field rather than a structure).
 */
type ObjectType struct {
	// Area Number where the object type is defined. Must not be '0' for values as this is the wildcard.
	Area UShort
	// Service Number of the service where the object type is defined. Must not be '0' for values as this is the wildcard.
	Service UShort
	// Area Version of the service where the object type is defined. Must not be '0' for values as this is the wildcard
	Version UOctet
	// The service specific object number. Must not be '0' for values as this is the wildcard.
	Number UShort
}

var (
	NullObjectType *ObjectType = nil
)

func NewObjectType() *ObjectType {
	return new(ObjectType)
}

// ================================================================================
// Defines COM ObjectType type as a MAL Composite

func (t *ObjectType) Composite() Composite {
	return t
}

// ================================================================================
// Defines COM ObjectType type as a MAL Element

const COM_OBJECT_TYPE_TYPE_SHORT_FORM Integer = 0x01
const COM_OBJECT_TYPE_SHORT_FORM Long = 0x2000001000001

// Returns the absolute short form of the element type.
func (*ObjectType) GetShortForm() Long {
	return COM_OBJECT_TYPE_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*ObjectType) GetAreaNumber() UShort {
	return COM_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*ObjectType) GetAreaVersion() UOctet {
	return COM_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*ObjectType) GetServiceNumber() UShort {
	return NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*ObjectType) GetTypeShortForm() Integer {
	return COM_OBJECT_TYPE_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (t *ObjectType) Encode(encoder Encoder) error {
	err := encoder.EncodeUShort(&t.Area)
	if err != nil {
		return err
	}
	err = encoder.EncodeUShort(&t.Service)
	if err != nil {
		return err
	}
	err = encoder.EncodeUOctet(&t.Version)
	if err != nil {
		return err
	}
	return encoder.EncodeUShort(&t.Number)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (t *ObjectType) Decode(decoder Decoder) (Element, error) {
	return DecodeObjectType(decoder)
}

// Decodes an instance of EntityKey using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded EntityKey instance.
func DecodeObjectType(decoder Decoder) (*ObjectType, error) {
	area, err := decoder.DecodeUShort()
	if err != nil {
		return nil, err
	}
	service, err := decoder.DecodeUShort()
	if err != nil {
		return nil, err
	}
	version, err := decoder.DecodeUOctet()
	if err != nil {
		return nil, err
	}
	number, err := decoder.DecodeUShort()
	if err != nil {
		return nil, err
	}
	var t = ObjectType{Area: *area, Service: *service, Version: *version, Number: *number}
	return &t, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (t *ObjectType) CreateElement() Element {
	return new(ObjectType)
}

func (t *ObjectType) IsNull() bool {
	return t == nil
}

func (*ObjectType) Null() Element {
	return NullObjectType
}
