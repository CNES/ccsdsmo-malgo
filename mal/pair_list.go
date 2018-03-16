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
// Defines MAL PairList type
// ################################################################################

type PairList []*Pair

var (
	NullPairList *PairList = nil
)

func NewPairList(size int) *PairList {
	var list PairList = PairList(make([]*Pair, size))
	return &list
}

// ================================================================================
// Defines MAL PairList type as a MAL Element

const MAL_PAIR_LIST_TYPE_SHORT_FORM Integer = -0x1C
const MAL_PAIR_LIST_SHORT_FORM Long = 0x1000001FFFF1C

// Returns the absolute short form of the element type.
func (*PairList) GetShortForm() Long {
	return MAL_PAIR_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*PairList) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*PairList) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*PairList) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*PairList) GetTypeShortForm() Integer {
	//	return MAL_PAIR_TYPE_SHORT_FORM & 0x01FFFF00
	return MAL_PAIR_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *PairList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*Pair(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*Pair(*list) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *PairList) Decode(decoder Decoder) (Element, error) {
	return DecodePairList(decoder)
}

// Decodes an instance of PairList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded PairList instance.
func DecodePairList(decoder Decoder) (*PairList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := PairList(make([]*Pair, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullPair)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*Pair)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *PairList) CreateElement() Element {
	return NewPairList(0)
}

func (list *PairList) IsNull() bool {
	return list == nil
}

func (*PairList) Null() Element {
	return NullPairList
}
