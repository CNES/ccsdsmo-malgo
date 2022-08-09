/**
 * MIT License
 *
 * Copyright (c) 2020 CNES
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
/*
 * This file has been automatically generated by fr.cnes.mo:StubGenerator_go
 * It has then be slightly transformed to match the underlying type uint8 originally defined in the mal.
 * The mal should eventually use the standard generated definition of the type.
 */
package mal

import (
	"fmt"
)

// Defines UpdateType type

// the generator would define the type as uint32 instead of uint8
//type UpdateType uint32
const (
	UPDATETYPE_CREATION_OVAL = iota
	UPDATETYPE_CREATION_NVAL = 1
	UPDATETYPE_UPDATE_OVAL
	UPDATETYPE_UPDATE_NVAL = 2
	UPDATETYPE_MODIFICATION_OVAL
	UPDATETYPE_MODIFICATION_NVAL = 3
	UPDATETYPE_DELETION_OVAL
	UPDATETYPE_DELETION_NVAL = 4
)

// Conversion table OVAL->NVAL
var nvalTable_UpdateType = []uint32{
	UPDATETYPE_CREATION_NVAL,
	UPDATETYPE_UPDATE_NVAL,
	UPDATETYPE_MODIFICATION_NVAL,
	UPDATETYPE_DELETION_NVAL,
}

// Conversion map NVAL->OVAL
var ovalMap_UpdateType map[uint32]uint32

var (
	UPDATETYPE_CREATION     = UpdateType(UPDATETYPE_CREATION_NVAL)
	UPDATETYPE_UPDATE       = UpdateType(UPDATETYPE_UPDATE_NVAL)
	UPDATETYPE_MODIFICATION = UpdateType(UPDATETYPE_MODIFICATION_NVAL)
	UPDATETYPE_DELETION     = UpdateType(UPDATETYPE_DELETION_NVAL)
)

var NullUpdateType *UpdateType = nil

func init() {
	ovalMap_UpdateType = make(map[uint32]uint32)
	for oval, nval := range nvalTable_UpdateType {
		ovalMap_UpdateType[nval] = uint32(oval)
	}
}

func (receiver UpdateType) GetNumericValue() uint32 {
	//  return uint32(receiver)
	return uint32(uint8(receiver))
}
func (receiver UpdateType) GetOrdinalValue() (uint32, error) {
	nval := receiver.GetNumericValue()
	oval, ok := ovalMap_UpdateType[nval]
	if !ok {
		return 0, fmt.Errorf("Invalid UpdateType value: %d", nval)
	}
	return oval, nil
}
func UpdateTypeFromNumericValue(nval uint32) (UpdateType, error) {
	_, ok := ovalMap_UpdateType[nval]
	if !ok {
		return UpdateType(0), fmt.Errorf("Invalid numeric value for UpdateType: %v", nval)
	}
	return UpdateType(nval), nil
}
func UpdateTypeFromOrdinalValue(oval uint32) (UpdateType, error) {
	if oval >= uint32(len(nvalTable_UpdateType)) {
		return UpdateType(0), fmt.Errorf("Invalid ordinal value for UpdateType: %v", oval)
	}
	return UpdateType(nvalTable_UpdateType[oval]), nil
}

// ================================================================================
// Defines UpdateType type as a MAL Element

//const UPDATETYPE_TYPE_SHORT_FORM Integer = 22
//const UPDATETYPE_SHORT_FORM Long = 0x65000001000016

// Registers UpdateType type for polymorphism handling
func init() {
	RegisterMALElement(UPDATETYPE_SHORT_FORM, NullUpdateType)
}

// Returns the absolute short form of the element type.
func (receiver *UpdateType) GetShortForm() Long {
	return UPDATETYPE_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *UpdateType) GetAreaNumber() UShort {
	return AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *UpdateType) GetAreaVersion() UOctet {
	return AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *UpdateType) GetServiceNumber() UShort {
	return NULL_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *UpdateType) GetTypeShortForm() Integer {
	return UPDATETYPE_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *UpdateType) CreateElement() Element {
	return NullUpdateType
}

func (receiver *UpdateType) IsNull() bool {
	return receiver == nil
}

func (receiver *UpdateType) Null() Element {
	return NullUpdateType
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *UpdateType) Encode(encoder Encoder) error {
	specific := encoder.LookupSpecific(UPDATETYPE_SHORT_FORM)
	if specific != nil {
		return specific(receiver, encoder)
	}

	oval, err := receiver.GetOrdinalValue()
	if err != nil {
		return err
	}
	value := NewUOctet(uint8(oval))
	return encoder.EncodeUOctet(value)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *UpdateType) Decode(decoder Decoder) (Element, error) {
	specific := decoder.LookupSpecific(UPDATETYPE_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}

	elem, err := decoder.DecodeUOctet()
	if err != nil {
		return receiver.Null(), err
	}
	value, err := UpdateTypeFromOrdinalValue(uint32(uint8(*elem)))
	return &value, err
}