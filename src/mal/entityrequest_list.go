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
// Defines MAL EntityRequestList type
// ################################################################################

type EntityRequestList []*EntityRequest

var (
	NullEntityRequestList *EntityRequestList = nil
)

func NewEntityRequestList(size int) *EntityRequestList {
	var list EntityRequestList = EntityRequestList(make([]*EntityRequest, size))
	return &list
}

// ================================================================================
// Defines MAL EntityRequestList type as a MAL Element

const MAL_ENTITY_REQUEST_LIST_TYPE_SHORT_FORM Integer = -0x18
const MAL_ENTITY_REQUEST_LIST_SHORT_FORM Long = 0x1000001FFFF18

// Returns the absolute short form of the element type.
func (*EntityRequestList) GetShortForm() Long {
	return MAL_ENTITY_REQUEST_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*EntityRequestList) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*EntityRequestList) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*EntityRequestList) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*EntityRequestList) GetTypeShortForm() Integer {
	//	return MAL_ENTITY_REQUEST_TYPE_SHORT_FORM & 0x01FFFF00
	return MAL_ENTITY_REQUEST_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *EntityRequestList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*EntityRequest(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*EntityRequest(*list) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *EntityRequestList) Decode(decoder Decoder) (Element, error) {
	return DecodeEntityRequestList(decoder)
}

// Decodes an instance of EntityRequestList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded EntityRequestList instance.
func DecodeEntityRequestList(decoder Decoder) (*EntityRequestList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := EntityRequestList(make([]*EntityRequest, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullEntityRequest)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*EntityRequest)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *EntityRequestList) CreateElement() Element {
	return NewEntityRequestList(0)
}

func (list *EntityRequestList) IsNull() bool {
	return list == nil
}

func (*EntityRequestList) Null() Element {
	return NullEntityRequestList
}
