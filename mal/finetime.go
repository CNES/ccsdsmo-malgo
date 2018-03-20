/**
 * MIT License
 *
 * Copyright (c) 2017 - 2018 CNES
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
	"time"
)

// ################################################################################
// Defines MAL FineTime type
// Time in picoseconds since January 1, 1970 UTC (Unix Time)
// ################################################################################

type FineTime time.Time

var (
	NullFineTime *FineTime = nil
)

func NewFineTime(t time.Time) *FineTime {
	var val FineTime = FineTime(t)
	return &val
}

func FineTimeNow() *FineTime {
	var val FineTime = FineTime(time.Now())
	return &val
}

// ================================================================================
// Defines MAL FineTime type as a MAL Attribute

func (t *FineTime) attribute() Attribute {
	return t
}

// ================================================================================
// Defines MAL FineTime type as a MAL Element

const MAL_FINETIME_TYPE_SHORT_FORM Integer = 0x11
const MAL_FINETIME_SHORT_FORM Long = 0x1000001000011

// Registers MAL FineTime type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_FINETIME_SHORT_FORM, NullFineTime)
}

// Returns the absolute short form of the element type.
func (*FineTime) GetShortForm() Long {
	return MAL_FINETIME_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*FineTime) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*FineTime) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*FineTime) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*FineTime) GetTypeShortForm() Integer {
	return MAL_FINETIME_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (t *FineTime) Encode(encoder Encoder) error {
	return encoder.EncodeFineTime(t)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (t *FineTime) Decode(decoder Decoder) (Element, error) {
	return decoder.DecodeFineTime()
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (t *FineTime) CreateElement() Element {
	return NewFineTime(time.Now())
}

func (t *FineTime) IsNull() bool {
	if t == nil {
		return true
	} else {
		return false
	}
}

func (*FineTime) Null() Element {
	return NullFineTime
}
