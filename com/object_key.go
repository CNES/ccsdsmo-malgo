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
// Defines COM ObjectKey type
// ################################################################################

/**
 * The ObjectKey structure combines a domain and an object instance identifier such that
 * it identifies the instance of an object for a specific domain.
 */
type ObjectKey struct {
	// The domain of the object instance.
	Domain IdentifierList
	// The unique identifier of the object instance. Must not be '0' for values as this is the wildcard.
	InstId Long
}

var (
	NullObjectKey *ObjectKey = nil
)

func NewObjectKey() *ObjectKey {
	return new(ObjectKey)
}

// ================================================================================
// Defines COM ObjectKey type as a MAL Composite

func (key *ObjectKey) Composite() Composite {
	return key
}

// ================================================================================
// Defines COM ObjectKey type as a MAL Element

const COM_OBJECT_KEY_TYPE_SHORT_FORM Integer = 0x02
const COM_OBJECT_KEY_SHORT_FORM Long = 0x2000001000002

// Registers COM ObjectKey type for polymorpsism handling
func init() {
	RegisterMALElement(COM_OBJECT_KEY_SHORT_FORM, NullObjectKey)
}

// Returns the absolute short form of the element type.
func (*ObjectKey) GetShortForm() Long {
	return COM_OBJECT_KEY_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*ObjectKey) GetAreaNumber() UShort {
	return COM_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*ObjectKey) GetAreaVersion() UOctet {
	return COM_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*ObjectKey) GetServiceNumber() UShort {
	return NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*ObjectKey) GetTypeShortForm() Integer {
	return COM_OBJECT_KEY_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (key *ObjectKey) Encode(encoder Encoder) error {
	err := encoder.EncodeElement(&key.Domain)
	if err != nil {
		return err
	}
	return encoder.EncodeLong(&key.InstId)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (key *ObjectKey) Decode(decoder Decoder) (Element, error) {
	return DecodeObjectKey(decoder)
}

// Decodes an instance of EntityKey using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded EntityKey instance.
func DecodeObjectKey(decoder Decoder) (*ObjectKey, error) {
	domain, err := decoder.DecodeElement(NullIdentifierList)
	if err != nil {
		return nil, err
	}
	instId, err := decoder.DecodeLong()
	if err != nil {
		return nil, err
	}
	var key = ObjectKey{Domain: *domain.(*IdentifierList), InstId: *instId}
	return &key, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (key *ObjectKey) CreateElement() Element {
	return new(ObjectKey)
}

func (key *ObjectKey) IsNull() bool {
	return key == nil
}

func (key *ObjectKey) Null() Element {
	return NullObjectKey
}
