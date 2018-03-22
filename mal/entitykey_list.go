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
// Defines MAL EntityKeyList type
// ################################################################################

type EntityKeyList []*EntityKey

var (
	NullEntityKeyList *EntityKeyList = nil
)

func NewEntityKeyList(size int) *EntityKeyList {
	var list EntityKeyList = EntityKeyList(make([]*EntityKey, size))
	return &list
}

// ================================================================================
// Defines MAL EntityKeyList type as an ElementList

func (list *EntityKeyList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

// ================================================================================
// Defines MAL EntityKeyList type as a MAL Composite

func (list *EntityKeyList) Composite() Composite {
	return list
}

// ================================================================================
// Defines MAL EntityKeyList type as a MAL Element

const MAL_ENTITY_KEY_LIST_TYPE_SHORT_FORM Integer = -0x19
const MAL_ENTITY_KEY_LIST_SHORT_FORM Long = 0x1000001FFFF19

// Registers MAL EntityKeyList type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_ENTITY_KEY_LIST_SHORT_FORM, NullEntityKeyList)
}

// Returns the absolute short form of the element type.
func (*EntityKeyList) GetShortForm() Long {
	return MAL_ENTITY_KEY_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*EntityKeyList) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*EntityKeyList) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*EntityKeyList) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*EntityKeyList) GetTypeShortForm() Integer {
	//	return MAL_ENTITY_KEY_TYPE_SHORT_FORM & 0x01FFFF00
	return MAL_ENTITY_KEY_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *EntityKeyList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*EntityKey(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*EntityKey(*list) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *EntityKeyList) Decode(decoder Decoder) (Element, error) {
	return DecodeEntityKeyList(decoder)
}

// Decodes an instance of EntityKeyList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded EntityKeyList instance.
func DecodeEntityKeyList(decoder Decoder) (*EntityKeyList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := EntityKeyList(make([]*EntityKey, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullEntityKey)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*EntityKey)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *EntityKeyList) CreateElement() Element {
	return NewEntityKeyList(0)
}

func (list *EntityKeyList) IsNull() bool {
	return list == nil
}

func (*EntityKeyList) Null() Element {
	return NullEntityKeyList
}
