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
// Defines MAL UpdateHeader type
// ################################################################################

type UpdateHeader struct {
	Timestamp  Time
	SourceURI  URI
	UpdateType UpdateType
	Key        EntityKey
}

var (
	NullUpdateHeader *UpdateHeader = nil
)

func NewUpdateHeader() *UpdateHeader {
	return new(UpdateHeader)
}

// ================================================================================
// Defines MAL UpdateHeader type as a MAL Composite

func (u *UpdateHeader) Composite() Composite {
	return u
}

// ================================================================================
// Defines MAL UpdateHeader type as a MAL Element

const MAL_UPDATE_HEADER_TYPE_SHORT_FORM Integer = 0x1A
const MAL_UPDATE_HEADER_SHORT_FORM Long = 0x100000100001A

// Returns the absolute short form of the element type.
func (*UpdateHeader) GetShortForm() Long {
	return MAL_UPDATE_HEADER_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*UpdateHeader) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*UpdateHeader) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*UpdateHeader) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*UpdateHeader) GetTypeShortForm() Integer {
	return MAL_UPDATE_HEADER_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (u *UpdateHeader) Encode(encoder Encoder) error {
	err := encoder.EncodeTime(&u.Timestamp)
	if err != nil {
		return err
	}
	err = encoder.EncodeURI(&u.SourceURI)
	if err != nil {
		return err
	}
	err = encoder.EncodeSmallEnum(uint8(u.UpdateType))
	if err != nil {
		return err
	}
	return u.Key.Encode(encoder)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (u *UpdateHeader) Decode(decoder Decoder) (Element, error) {
	return DecodeUpdateHeader(decoder)
}

// Decodes an instance of UpdateHeader using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded EntityRequest instance.
func DecodeUpdateHeader(decoder Decoder) (*UpdateHeader, error) {
	timestamp, err := decoder.DecodeTime()
	if err != nil {
		return nil, err
	}
	sourceURI, err := decoder.DecodeURI()
	if err != nil {
		return nil, err
	}
	updateType, err := decoder.DecodeSmallEnum()
	if err != nil {
		return nil, err
	}
	key, err := DecodeEntityKey(decoder)
	if err != nil {
		return nil, err
	}
	var updateHeader = UpdateHeader{
		Timestamp:  *timestamp,
		SourceURI:  *sourceURI,
		UpdateType: UpdateType(updateType),
		Key:        *key,
	}
	return &updateHeader, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (u *UpdateHeader) CreateElement() Element {
	// TODO (AF):
	//	return new(UpdateHeader)
	return NewUpdateHeader()
}

func (u *UpdateHeader) IsNull() bool {
	return u == nil
}

func (*UpdateHeader) Null() Element {
	return NullUpdateHeader
}
