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
// Defines MAL DurationList type
// ################################################################################

type DurationList []*Duration

var (
	NullDurationList *IntegerList = nil
)

func NewDurationList(size int) *DurationList {
	var list DurationList = DurationList(make([]*Duration, size))
	return &list
}

// ================================================================================
// Defines MAL DurationList type as an ElementList

func (list *DurationList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

// ================================================================================
// Defines MAL DurationList type as a MAL Composite

func (list *DurationList) Composite() Composite {
	return list
}

// ================================================================================
// Defines MAL DurationList type as a MAL Element

const MAL_DURATION_LIST_TYPE_SHORT_FORM Integer = -0x03
const MAL_DURATION_LIST_SHORT_FORM Long = 0x1000001FFFFFD

// Registers MAL DurationList type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_DURATION_LIST_SHORT_FORM, NullDurationList)
}

// Returns the absolute short form of the element type.
func (*DurationList) GetShortForm() Long {
	return MAL_DURATION_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*DurationList) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*DurationList) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*DurationList) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*DurationList) GetTypeShortForm() Integer {
	//	return MAL_DURATION_TYPE_SHORT_FORM & 0x01FFFF00
	return MAL_DURATION_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *DurationList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*Duration(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*Duration(*list) {
		encoder.EncodeNullableDuration(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *DurationList) Decode(decoder Decoder) (Element, error) {
	return DecodeDurationList(decoder)
}

// Decodes an instance of DurationList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded DurationList instance.
func DecodeDurationList(decoder Decoder) (*DurationList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := DurationList(make([]*Duration, int(*size)))
	for i := 0; i < len(list); i++ {
		list[i], err = decoder.DecodeNullableDuration()
		if err != nil {
			return nil, err
		}
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *DurationList) CreateElement() Element {
	return NewDurationList(0)
}

func (list *DurationList) IsNull() bool {
	return list == nil
}

func (*DurationList) Null() Element {
	return NullDurationList
}
