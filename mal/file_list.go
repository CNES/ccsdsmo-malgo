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
// Defines MAL FileList type
// ################################################################################

type FileList []*File

var (
	NullFileList *FileList = nil
)

func NewFileList(size int) *FileList {
	var list FileList = FileList(make([]*File, size))
	return &list
}

// ================================================================================
// Defines MAL FileList type as an ElementList

func (list *FileList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

func (list *FileList) GetElementAt(i int) Element {
	if list != nil {
		if i < list.Size() {
			return (*list)[i]
		}
		return nil
	}
	return nil
}

func (list *FileList) AppendElement(element Element) {
	if list != nil {
		*list = append(*list, element.(*File))
	}
}

// ================================================================================
// Defines MAL FileList type as a MAL Composite

func (list *FileList) Composite() Composite {
	return list
}

// ================================================================================
// Defines MAL FileList type as a MAL Element

const MAL_FILE_LIST_TYPE_SHORT_FORM Integer = -0x1E
const MAL_FILE_LIST_SHORT_FORM Long = 0x1000001FFFFE2

// Registers MAL FileList type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_FILE_LIST_SHORT_FORM, NullFileList)
}

// Returns the absolute short form of the element type.
func (*FileList) GetShortForm() Long {
	return MAL_FILE_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*FileList) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*FileList) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*FileList) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*FileList) GetTypeShortForm() Integer {
	//	return MAL_ENTITY_REQUEST_TYPE_SHORT_FORM & 0x01FFFF00
	return MAL_FILE_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *FileList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*File(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*File(*list) {
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
func (list *FileList) Decode(decoder Decoder) (Element, error) {
	return DecodeFileList(decoder)
}

// Decodes an instance of FileList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded FileList instance.
func DecodeFileList(decoder Decoder) (*FileList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := FileList(make([]*File, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullFile)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*File)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *FileList) CreateElement() Element {
	return NewFileList(0)
}

func (list *FileList) IsNull() bool {
	return list == nil
}

func (*FileList) Null() Element {
	return NullFileList
}
