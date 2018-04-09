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
// Defines MAL UShortList type
// ################################################################################

type UShortList []*UShort

var (
	NullUShortList *UShortList = nil
)

func NewUShortList(size int) *UShortList {
	var list UShortList = UShortList(make([]*UShort, size))
	return &list
}

// ================================================================================
// Defines MAL UShortList type as an ElementList

func (list *UShortList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

func (list *UShortList) GetElementAt(i int) Element {
	if list != nil {
		if i <= list.Size() {
			return (*list)[i]
		}
		return nil
	}
	return nil
}

// ================================================================================
// Defines MAL UShortList type as a MAL Composite

func (list *UShortList) Composite() Composite {
	return list
}

// ================================================================================
// Defines MAL UShortList type as a MAL Element

const MAL_USHORT_LIST_TYPE_SHORT_FORM Integer = -0x0A
const MAL_USHORT_LIST_SHORT_FORM Long = 0x1000001FFFFF6

// Registers MAL UShortList type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_USHORT_LIST_SHORT_FORM, NullUShortList)
}

// Returns the absolute short form of the element type.
func (*UShortList) GetShortForm() Long {
	return MAL_USHORT_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*UShortList) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*UShortList) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*UShortList) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Return the relative short form of the element type.
func (*UShortList) GetTypeShortForm() Integer {
	//	return MAL_USHORT_TYPE_SHORT_FORM & 0x01FFFF00
	return MAL_USHORT_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *UShortList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*UShort(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*UShort(*list) {
		encoder.EncodeNullableUShort(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *UShortList) Decode(decoder Decoder) (Element, error) {
	return DecodeUShortList(decoder)
}

// Decodes an instance of UShortList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded UShortList instance.
func DecodeUShortList(decoder Decoder) (*UShortList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := UShortList(make([]*UShort, int(*size)))
	for i := 0; i < len(list); i++ {
		list[i], err = decoder.DecodeNullableUShort()
		if err != nil {
			return nil, err
		}
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *UShortList) CreateElement() Element {
	return NewUShortList(0)
}

func (list *UShortList) IsNull() bool {
	return list == nil
}

func (*UShortList) Null() Element {
	return NullUShortList
}
