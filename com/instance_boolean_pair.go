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
// Defines COM InstanceBooleanPair type
// ################################################################################

/**
 * Simple pair of an object instance identifier and a Boolean value.
 */
type InstanceBooleanPair struct {
	// The object instance identifier.
	Id Long
	// An associated Boolean value.
	Value Boolean
}

var (
	NullInstanceBooleanPair *InstanceBooleanPair = nil
)

func NewInstanceBooleanPair() *InstanceBooleanPair {
	return new(InstanceBooleanPair)
}

// ================================================================================
// Defines COM InstanceBooleanPair type as a MAL Composite

func (pair *InstanceBooleanPair) Composite() Composite {
	return pair
}

// ================================================================================
// Defines COM InstanceBooleanPair type as a MAL Element

const COM_INSTANCE_BOOLEAN_PAIR_TYPE_SHORT_FORM Integer = 0x05
const COM_INSTANCE_BOOLEAN_PAIR_SHORT_FORM Long = 0x2000001000005

// Registers COM InstanceBooleanPair type for polymorpsism handling
func init() {
	RegisterMALElement(COM_INSTANCE_BOOLEAN_PAIR_SHORT_FORM, NullInstanceBooleanPair)
}

// Returns the absolute short form of the element type.
func (*InstanceBooleanPair) GetShortForm() Long {
	return COM_INSTANCE_BOOLEAN_PAIR_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*InstanceBooleanPair) GetAreaNumber() UShort {
	return COM_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*InstanceBooleanPair) GetAreaVersion() UOctet {
	return COM_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*InstanceBooleanPair) GetServiceNumber() UShort {
	return NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*InstanceBooleanPair) GetTypeShortForm() Integer {
	return COM_INSTANCE_BOOLEAN_PAIR_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (pair *InstanceBooleanPair) Encode(encoder Encoder) error {
	err := encoder.EncodeLong(&pair.Id)
	if err != nil {
		return err
	}
	return encoder.EncodeBoolean(&pair.Value)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (pair *InstanceBooleanPair) Decode(decoder Decoder) (Element, error) {
	return DecodeInstanceBooleanPair(decoder)
}

// Decodes an instance of EntityKey using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded EntityKey instance.
func DecodeInstanceBooleanPair(decoder Decoder) (*InstanceBooleanPair, error) {
	id, err := decoder.DecodeLong()
	if err != nil {
		return nil, err
	}
	value, err := decoder.DecodeBoolean()
	if err != nil {
		return nil, err
	}
	var pair = InstanceBooleanPair{Id: *id, Value: *value}
	return &pair, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (pair *InstanceBooleanPair) CreateElement() Element {
	return new(InstanceBooleanPair)
}

func (pair *InstanceBooleanPair) IsNull() bool {
	return pair == nil
}

func (pair *InstanceBooleanPair) Null() Element {
	return NullInstanceBooleanPair
}
