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
// Defines MAL File type
// ################################################################################

type File struct {
	Name             Identifier
	MimeType         *String
	CreationDate     *Time
	ModificationDate *Time
	Size             *ULong
	Content          *Blob
	MetaData         *NamedValueList
}

var (
	NullFile *File = nil
)

func NewFile() *File {
	return new(File)
}

// ================================================================================
// Defines MAL File type as a MAL Composite

func (file *File) composite() Composite {
	return file
}

// ================================================================================
// Defines MAL File type as a MAL Element

const MAL_FILE_TYPE_SHORT_FORM Integer = 0x1E
const MAL_FILE_SHORT_FORM Long = 0x100000100001E

// Returns the absolute short form of the element type.
func (*File) GetShortForm() Long {
	return MAL_FILE_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*File) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*File) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*File) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*File) GetTypeShortForm() Integer {
	return MAL_FILE_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (file *File) Encode(encoder Encoder) error {
	err := encoder.EncodeIdentifier(&file.Name)
	if err != nil {
		return err
	}
	err = encoder.EncodeNullableString(file.MimeType)
	if err != nil {
		return err
	}
	err = encoder.EncodeNullableTime(file.CreationDate)
	if err != nil {
		return err
	}
	err = encoder.EncodeNullableTime(file.ModificationDate)
	if err != nil {
		return err
	}
	err = encoder.EncodeNullableULong(file.Size)
	if err != nil {
		return err
	}
	err = encoder.EncodeNullableBlob(file.Content)
	if err != nil {
		return err
	}
	return encoder.EncodeNullableElement(file.MetaData)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (file *File) Decode(decoder Decoder) (Element, error) {
	return DecodeFile(decoder)
}

// Decodes an instance of File using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded EntityRequest instance.
func DecodeFile(decoder Decoder) (*File, error) {
	name, err := decoder.DecodeIdentifier()
	if err != nil {
		return nil, err
	}
	mimeType, err := decoder.DecodeNullableString()
	if err != nil {
		return nil, err
	}
	creationDate, err := decoder.DecodeNullableTime()
	if err != nil {
		return nil, err
	}
	modificationDate, err := decoder.DecodeNullableTime()
	if err != nil {
		return nil, err
	}
	size, err := decoder.DecodeNullableULong()
	if err != nil {
		return nil, err
	}
	content, err := decoder.DecodeNullableBlob()
	if err != nil {
		return nil, err
	}
	var metaData *NamedValueList
	element, err := decoder.DecodeNullableElement(metaData)
	metaData = element.(*NamedValueList)
	if err != nil {
		return nil, err
	}
	var file = File{
		*name, mimeType, creationDate, modificationDate, size, content, metaData,
	}
	return &file, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (file *File) CreateElement() Element {
	// TODO (AF):
	//	return new(File)
	return NewFile()
}

func (file *File) IsNull() bool {
	return file == nil
}

func (*File) Null() Element {
	return NullFile
}
