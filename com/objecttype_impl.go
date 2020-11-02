/**
 * MIT License
 *
 * Copyright (c) 2018 - 2020 CNES
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
	"errors"
	"fmt"
	"github.com/CNES/ccsdsmo-malgo/mal"
)

// ================================================================================
// A COM object holds a body made of a MAL object
// The type of the MAL object is defined in the specification of the COM object type.
// It must be made known to the system with a call to RegisterMALBodyType

var comTypesMap = make(map[ObjectType]mal.Long)

func (t *ObjectType) RegisterMALBodyType(shortForm mal.Long) error {
	if t == nil {
		return errors.New("Unexpected null type in RegisterMALBodyType")
	}
	if t.Area == 0 || t.Version == 0 || t.Number == 0 {
		return errors.New("Unexpected null type field in RegisterMALBodyType")
	}
	val := comTypesMap[*t]
	if val != 0 {
		return errors.New("A value has already been registered for type")
	}
	// check the MAL type corresponding to shortForm is known to the system
	_, error := mal.LookupMALElement(shortForm)
	if error != nil {
		return error
	}
	comTypesMap[*t] = shortForm
	return nil
}

func (t *ObjectType) GetMALBodyType() mal.Long {
	if t == nil {
		return 0
	}
	return comTypesMap[*t]
}

func (t *ObjectType) GetMALBodyListType() mal.Long {
	if t == nil {
		return 0
	}
	shortForm := comTypesMap[*t]
	if shortForm == 0 {
		return 0
	}
	// this short form should never represent a list type
	numberPart := mal.ULong(shortForm & 0xFFFFFF)
	numberListPart := mal.ULong((-mal.Long(numberPart)) & 0xFFFFFF)
	return mal.Long((mal.ULong(shortForm) & 0xFFFFFFFFFF000000) | numberListPart)
}

// ================================================================================
// Implements Stringer interface

func (t *ObjectType) String() string {
	return fmt.Sprintf("ObjectType(%d, %d, %d, %d)", t.Area, t.Service, t.Version, t.Number)
}
