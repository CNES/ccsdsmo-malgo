/**
 * MIT License
 *
 * Copyright (c) 2018 - 2019 CNES
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
	"fmt"
	. "github.com/CNES/ccsdsmo-malgo/mal"
)

// ################################################################################
// Defines COM ObjectId type
// ################################################################################

/**
 * The ObjectId structure combines an object type and an object key such that it identifies
 * the instance and type of an object for a specific domain.
 */
type ObjectId struct {
	// The fully qualified unique identifier of the type.
	Type *ObjectType
	// The combination of the object domain and object instance identifier.
	Key *ObjectKey
}

var (
	NullObjectId *ObjectId = nil
)

func NewObjectId() *ObjectId {
	return new(ObjectId)
}

// ================================================================================
// Defines COM ObjectId type as a MAL Composite

func (id *ObjectId) Composite() Composite {
	return id
}

// ================================================================================
// Defines COM ObjectId type as a MAL Element

const COM_OBJECT_ID_TYPE_SHORT_FORM Integer = 0x03
const COM_OBJECT_ID_SHORT_FORM Long = 0x2000001000003

// Registers COM ObjectId type for polymorpsism handling
func init() {
	RegisterMALElement(COM_OBJECT_ID_SHORT_FORM, NullObjectId)
}

// Returns the absolute short form of the element type.
func (*ObjectId) GetShortForm() Long {
	return COM_OBJECT_ID_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*ObjectId) GetAreaNumber() UShort {
	return COM_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*ObjectId) GetAreaVersion() UOctet {
	return COM_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*ObjectId) GetServiceNumber() UShort {
	return NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*ObjectId) GetTypeShortForm() Integer {
	return COM_OBJECT_ID_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (id *ObjectId) Encode(encoder Encoder) error {
	err := encoder.EncodeElement(id.Type)
	if err != nil {
		return err
	}
	return encoder.EncodeElement(id.Key)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (id *ObjectId) Decode(decoder Decoder) (Element, error) {
	return DecodeObjectId(decoder)
}

// Decodes an instance of EntityKey using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded EntityKey instance.
func DecodeObjectId(decoder Decoder) (*ObjectId, error) {
	t, err := decoder.DecodeElement(NullObjectType)
	if err != nil {
		return nil, err
	}
	k, err := decoder.DecodeElement(NullObjectKey)
	if err != nil {
		return nil, err
	}
	var id = ObjectId{Type: t.(*ObjectType), Key: k.(*ObjectKey)}
	return &id, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (id *ObjectId) CreateElement() Element {
	return new(ObjectId)
}

func (id *ObjectId) IsNull() bool {
	return id == nil
}

func (*ObjectId) Null() Element {
	return NullObjectId
}

// ================================================================================
// Implements Stringer interface

func (id *ObjectId) String() string {
	return fmt.Sprintf("ObjectId(%s, %s)", id.Type, id.Key)
}
