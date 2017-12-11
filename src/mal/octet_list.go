/**
 * MIT License
 *
 * Copyright (c) 2017 CNES
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
// Defines MAL OctetList type
// ################################################################################

type OctetList []*Octet

var (
	NullOctetList *OctetList = nil
)

func NewOctetList(size int) *OctetList {
	var list OctetList = OctetList(make([]*Octet, size))
	return &list
}

// ================================================================================
// Defines MAL OctetList type as a MAL Element

const MAL_OCTET_LIST_TYPE_SHORT_FORM Integer = -0x07
const MAL_OCTET_LIST_SHORT_FORM Long = 0x1000001FFFFF9

// Returns the absolute short form of the element type.
func (*OctetList) GetShortForm() Long {
	return MAL_OCTET_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*OctetList) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*OctetList) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*OctetList) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*OctetList) GetTypeShortForm() Integer {
	//	return MAL_OCTET_TYPE_SHORT_FORM & 0x01FFFF00
	return MAL_OCTET_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *OctetList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*Octet(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*Octet(*list) {
		encoder.EncodeNullableOctet(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *OctetList) Decode(decoder Decoder) (Element, error) {
	return DecodeOctetList(decoder)
}

// Decodes an instance of OctetList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded OctetList instance.
func DecodeOctetList(decoder Decoder) (*OctetList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := OctetList(make([]*Octet, int(*size)))
	for i := 0; i < len(list); i++ {
		list[i], err = decoder.DecodeNullableOctet()
		if err != nil {
			return nil, err
		}
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (*OctetList) CreateElement() Element {
	return NewOctetList(0)
}

func (list *OctetList) IsNull() bool {
	return list == nil
}

func (*OctetList) Null() Element {
	return NullOctetList
}
