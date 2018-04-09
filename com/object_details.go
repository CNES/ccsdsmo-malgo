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
// Defines COM ObjectDetails type
// ################################################################################

/**
 * The ObjectDetails type is used to hold the extra information associated with an object
 * instance, namely the related and source links.
 */
type ObjectDetails struct {
	// Contains the object instance identifier of a related object (e.g. the ActionDefinition that an
	// Action uses). This is service specific. The ObjectType of the related object is specified in the
	// service specification. The related object must exist in the same domain as this object.
	Related *Long
	// An object which is at the origin of the object creation (e.g. the procedure from which an action
	// was triggered).
	Source *ObjectId
}

var (
	NullObjectDetails *ObjectDetails = nil
)

func NewObjectDetails() *ObjectDetails {
	return new(ObjectDetails)
}

// ================================================================================
// Defines COM ObjectDetails type as a MAL Composite

func (details *ObjectDetails) Composite() Composite {
	return details
}

// ================================================================================
// Defines COM ObjectDetails type as a MAL Element

const COM_OBJECT_DETAILS_TYPE_SHORT_FORM Integer = 0x04
const COM_OBJECT_DETAILS_SHORT_FORM Long = 0x2000001000004

// Registers COM ObjectDetails type for polymorpsism handling
func init() {
	RegisterMALElement(COM_OBJECT_DETAILS_SHORT_FORM, NullObjectDetails)
}

// Returns the absolute short form of the element type.
func (*ObjectDetails) GetShortForm() Long {
	return COM_OBJECT_DETAILS_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*ObjectDetails) GetAreaNumber() UShort {
	return COM_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*ObjectDetails) GetAreaVersion() UOctet {
	return COM_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*ObjectDetails) GetServiceNumber() UShort {
	return NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*ObjectDetails) GetTypeShortForm() Integer {
	return COM_OBJECT_DETAILS_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (details *ObjectDetails) Encode(encoder Encoder) error {
	err := encoder.EncodeNullableLong(details.Related)
	if err != nil {
		return err
	}
	return encoder.EncodeNullableElement(details.Source)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (details *ObjectDetails) Decode(decoder Decoder) (Element, error) {
	return DecodeObjectDetails(decoder)
}

// Decodes an instance of EntityKey using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded EntityKey instance.
func DecodeObjectDetails(decoder Decoder) (*ObjectDetails, error) {
	related, err := decoder.DecodeNullableLong()
	if err != nil {
		return nil, err
	}
	source, err := decoder.DecodeNullableElement(NullObjectId)
	if err != nil {
		return nil, err
	}
	var details = ObjectDetails{
		Related: related,
		Source:  source.(*ObjectId),
	}
	return &details, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (details *ObjectDetails) CreateElement() Element {
	return new(ObjectDetails)
}

func (details *ObjectDetails) IsNull() bool {
	return details == nil
}

func (*ObjectDetails) Null() Element {
	return NullObjectDetails
}
