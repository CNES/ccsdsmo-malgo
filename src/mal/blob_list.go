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
// Defines MAL BlobList type
// ################################################################################

type BlobList []*Blob

var (
	NullBlobList *BlobList = nil
)

func NewBlobList(size int) *BlobList {
	var blob BlobList = BlobList(make([]*Blob, size))
	return &blob
}

// ================================================================================
// Defines MAL BlobList type as a MAL Element

const MAL_BLOB_LIST_TYPE_SHORT_FORM Integer = -0x01
const MAL_BLOB_LIST_SHORT_FORM Long = 0x1000001FFFFFF

// Returns the absolute short form of the element type.
func (*BlobList) GetShortForm() Long {
	return MAL_BLOB_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*BlobList) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*BlobList) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*BlobList) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*BlobList) GetTypeShortForm() Integer {
	//	return MAL_BLOB_TYPE_SHORT_FORM & 0x01FFFF00
	return MAL_BLOB_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *BlobList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*Blob(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*Blob(*list) {
		encoder.EncodeNullableBlob(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *BlobList) Decode(decoder Decoder) (Element, error) {
	return DecodeBlobList(decoder)
}

// Decodes an instance of BlobList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded BlobList instance.
func DecodeBlobList(decoder Decoder) (*BlobList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := BlobList(make([]*Blob, int(*size)))
	for i := 0; i < len(list); i++ {
		list[i], err = decoder.DecodeNullableBlob()
		if err != nil {
			return nil, err
		}
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *BlobList) CreateElement() Element {
	return NewBlobList(0)
}

func (list *BlobList) IsNull() bool {
	if list == nil {
		return true
	} else {
		return false
	}
}

func (*BlobList) Null() Element {
	return NullBlobList
}
