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

import ()

// ################################################################################
// Defines MAL Blob type
// ################################################################################

type Blob []byte

var (
	NullBlob *Blob = nil
)

func (blob Blob) Copy() Blob {
	buf := make([]byte, blob.Capacity())
	copy(buf, []byte(blob))
	return Blob(buf)
}

// Return the content of the Blob
func (blob Blob) Value() []byte {
	return []byte(blob)
}

// Returns the length of the Blob
func (blob Blob) Length() int {
	return len(blob)
}

// Returns the capacity of the Blob
func (blob Blob) Capacity() int {
	return cap(blob)
}

// ================================================================================
// Defines MAL Blob type as a MAL Attribute

func (blob *Blob) attribute() Attribute {
	return blob
}

// ================================================================================
// Defines MAL Blob type as a MAL Element

const MAL_BLOB_TYPE_SHORT_FORM Integer = 0x01
const MAL_BLOB_SHORT_FORM Long = 0x1000001000001

// Registers MAL Blob type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_BLOB_SHORT_FORM, NullBlob)
}

// Returns the absolute short form of the element type.
func (*Blob) GetShortForm() Long {
	return MAL_BLOB_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*Blob) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*Blob) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*Blob) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*Blob) GetTypeShortForm() Integer {
	return MAL_BLOB_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (blob *Blob) Encode(encoder Encoder) error {
	return encoder.EncodeBlob(blob)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (blob *Blob) Decode(decoder Decoder) (Element, error) {
	return decoder.DecodeBlob()
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (blob *Blob) CreateElement() Element {
	return NullBlob
}

func (blob Blob) IsNull() bool {
	return blob == nil
}

func (*Blob) Null() Element {
	return NullBlob
}
