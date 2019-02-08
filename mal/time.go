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
	"time"
)

// ################################################################################
// Defines MAL Time type
// Time in milliseconds since January 1, 1970 UTC (Unix Time)
// ################################################################################

const (
	NANOS_FROM_CCSDS_TO_UNIX_EPOCH int64 = 378691200000000000
	NANOS_IN_DAY                   int64 = 86400000000000
)

type Time time.Time

var (
	NullTime *Time = nil
)

func NewTime(t time.Time) *Time {
	var val Time = Time(t)
	return &val
}

func TimeNow() *Time {
	var val Time = Time(time.Now())
	return &val
}

func (t *Time) UnixNano() int64 {
	return time.Time(*t).UnixNano()
}

func FromUnixNano(t int64) *Time {
	var val Time = Time(time.Unix(t/1000000000, t%1000000000))
	return &val
}

// ================================================================================
// Defines MAL Time type as a MAL Attribute

func (t *Time) attribute() Attribute {
	return t
}

// ================================================================================
// Defines MAL Time type as a MAL Element

const MAL_TIME_TYPE_SHORT_FORM Integer = 0x10
const MAL_TIME_SHORT_FORM Long = 0x1000001000010

// Registers MAL Time type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_TIME_SHORT_FORM, NullTime)
}

// Returns the absolute short form of the element type.
func (*Time) GetShortForm() Long {
	return MAL_TIME_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*Time) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*Time) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*Time) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*Time) GetTypeShortForm() Integer {
	return MAL_TIME_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (t *Time) Encode(encoder Encoder) error {
	return encoder.EncodeTime(t)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (t *Time) Decode(decoder Decoder) (Element, error) {
	return decoder.DecodeTime()
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (t *Time) CreateElement() Element {
	return NewTime(time.Now())
}

func (t *Time) IsNull() bool {
	if t == nil {
		return true
	} else {
		return false
	}
}

func (*Time) Null() Element {
	return NullTime
}
