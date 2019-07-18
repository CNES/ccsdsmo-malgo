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
// Defines MAL UpdateHeaderList type
// ################################################################################

type UpdateHeaderList []*UpdateHeader

var (
	NullUpdateHeaderList *UpdateHeaderList = nil
)

func NewUpdateHeaderList(size int) *UpdateHeaderList {
	var list UpdateHeaderList = UpdateHeaderList(make([]*UpdateHeader, size))
	return &list
}

// ================================================================================
// Defines MAL UpdateHeaderList type as an ElementList

func (list *UpdateHeaderList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

func (list *UpdateHeaderList) GetElementAt(i int) Element {
	if list != nil {
		if i < list.Size() {
			return (*list)[i]
		}
		return nil
	}
	return nil
}

func (list *UpdateHeaderList) AppendElement(element Element) {
	if list != nil {
		*list = append(*list, element.(*UpdateHeader))
	}
}

// ================================================================================
// Defines MAL UpdateHeaderList type as a MAL Composite

func (list *UpdateHeaderList) Composite() Composite {
	return list
}

// ================================================================================
// Defines MAL UpdateHeaderList type as a MAL Element

const MAL_UPDATE_HEADER_LIST_TYPE_SHORT_FORM Integer = -0x1A
const MAL_UPDATE_HEADER_LIST_SHORT_FORM Long = 0x1000001FFFFE6

// Registers MAL UpdateHeaderList type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_UPDATE_HEADER_LIST_SHORT_FORM, NullUpdateHeaderList)
}

// Returns the absolute short form of the element type.
func (*UpdateHeaderList) GetShortForm() Long {
	return MAL_UPDATE_HEADER_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*UpdateHeaderList) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*UpdateHeaderList) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*UpdateHeaderList) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*UpdateHeaderList) GetTypeShortForm() Integer {
	//	return MAL_UPDATE_HEADER_TYPE_SHORT_FORM & 0x01FFFF00
	return MAL_UPDATE_HEADER_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *UpdateHeaderList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*UpdateHeader(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*UpdateHeader(*list) {
		err = encoder.EncodeNullableElement(e)
		if err != nil {
			return err
		}
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *UpdateHeaderList) Decode(decoder Decoder) (Element, error) {
	return DecodeUpdateHeaderList(decoder)
}

// Decodes an instance of UpdateHeaderList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded UpdateHeaderList instance.
func DecodeUpdateHeaderList(decoder Decoder) (*UpdateHeaderList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := UpdateHeaderList(make([]*UpdateHeader, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullUpdateHeader)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*UpdateHeader)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *UpdateHeaderList) CreateElement() Element {
	return NewUpdateHeaderList(0)
}

func (list *UpdateHeaderList) IsNull() bool {
	return list == nil
}

func (*UpdateHeaderList) Null() Element {
	return NullUpdateHeaderList
}
