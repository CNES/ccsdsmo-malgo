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

// ################################################################################
// Defines MAL ShortList type
// ################################################################################

type ShortList []*Short

var (
	NullShortList *ShortList = nil
)

func NewShortList(size int) *ShortList {
	var list ShortList = ShortList(make([]*Short, size))
	return &list
}

// ================================================================================
// Defines MAL ShortList type as an ElementList

func (list *ShortList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

func (list *ShortList) GetElementAt(i int) Element {
	if list != nil {
		if i < list.Size() {
			return (*list)[i]
		}
		return nil
	}
	return nil
}

func (list *ShortList) AppendElement(element Element) {
	if list != nil {
		*list = append(*list, element.(*Short))
	}
}

// ================================================================================
// Defines MAL ShortList type as a MAL Composite

func (list *ShortList) Composite() Composite {
	return list
}

// ================================================================================
// Defines MAL ShortList type as a MAL Element

const MAL_SHORT_LIST_TYPE_SHORT_FORM Integer = -0x09
const MAL_SHORT_LIST_SHORT_FORM Long = 0x1000001FFFFF7

// Registers MAL ShortList type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_SHORT_LIST_SHORT_FORM, NullShortList)
}

// Returns the absolute short form of the element type.
func (*ShortList) GetShortForm() Long {
	return MAL_SHORT_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*ShortList) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*ShortList) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*ShortList) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*ShortList) GetTypeShortForm() Integer {
	//	return MAL_SHORT_TYPE_SHORT_FORM & 0x01FFFF00
	return MAL_SHORT_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *ShortList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*Short(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*Short(*list) {
		encoder.EncodeNullableShort(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *ShortList) Decode(decoder Decoder) (Element, error) {
	return DecodeShortList(decoder)
}

// Decodes an instance of ShortList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded ShortList instance.
func DecodeShortList(decoder Decoder) (*ShortList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := ShortList(make([]*Short, int(*size)))
	for i := 0; i < len(list); i++ {
		list[i], err = decoder.DecodeNullableShort()
		if err != nil {
			return nil, err
		}
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *ShortList) CreateElement() Element {
	return NewShortList(0)
}

func (list *ShortList) IsNull() bool {
	return list == nil
}

func (*ShortList) Null() Element {
	return NullShortList
}
