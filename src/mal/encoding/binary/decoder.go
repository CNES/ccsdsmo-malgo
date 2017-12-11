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
package binary

import (
	. "mal"
	"math"
	"time"
)

type BinaryDecoder struct {
	GenDecoder
	buffer Buffer
}

// Creates a new decoder using a slice containing binary data to decode.
func NewBinaryDecoder(buf []byte) *BinaryDecoder {
	decoder := &BinaryDecoder{
		buffer: &BinaryBuffer{
			offset: 0,
			buf:    buf,
		},
	}
	decoder.GenDecoder.Decoder = decoder
	return decoder
}

func (decoder *BinaryDecoder) Read() (byte, error) {
	return decoder.buffer.Read()
}

func (decoder *BinaryDecoder) ReadUInt32() (uint32, error) {
	return decoder.buffer.Read32()
}

func (decoder *BinaryDecoder) ReadBody() ([]byte, error) {
	return decoder.buffer.Remaining()
}

// ================================================================================
// Implements Decoder interface

func (decoder *BinaryDecoder) IsNull() (bool, error) {
	b, err := decoder.buffer.Read()
	if err != nil {
		return false, err
	}
	// Note: the encoded value corresponds to the presence flag.
	return (b == FALSE), nil
}

// Decodes the short form of an attribute.
// @return The short form of the attribute.
func (decoder *BinaryDecoder) DecodeAttributeType() (Integer, error) {
	b, err := decoder.buffer.Read()
	if err != nil {
		return -1, err
	}
	return Integer(int(b)), nil
}

// Decodes a Boolean.
// @return The decoded Boolean.
func (decoder *BinaryDecoder) DecodeBoolean() (*Boolean, error) {
	b, err := decoder.buffer.Read()
	if err != nil {
		return new(Boolean), err
	}
	return NewBoolean(b == TRUE), nil
}

// Decodes a Float.
// @return The decoded Float.
func (decoder *BinaryDecoder) DecodeFloat() (*Float, error) {
	f, err := decoder.buffer.Read32()
	if err != nil {
		return nil, err
	}
	return NewFloat(math.Float32frombits(f)), nil
}

// Decodes a Double.
// @return The decoded Double.
func (decoder *BinaryDecoder) DecodeDouble() (*Double, error) {
	d, err := decoder.buffer.Read64()
	if err != nil {
		return nil, err
	}
	return NewDouble(math.Float64frombits(d)), nil
}

// Decodes an Octet.
// @return The decoded Octet.
func (decoder *BinaryDecoder) DecodeOctet() (*Octet, error) {
	o, err := decoder.buffer.Read()
	if err != nil {
		return nil, err
	}
	return NewOctet(int8(o)), nil
}

// Decodes a UOctet.
// @return The decoded UOctet.
func (decoder *BinaryDecoder) DecodeUOctet() (*UOctet, error) {
	o, err := decoder.buffer.Read()
	if err != nil {
		return nil, err
	}
	return NewUOctet(uint8(o)), nil
}

// Decodes a Short.
// @return The decoded Short.
func (decoder *BinaryDecoder) DecodeShort() (*Short, error) {
	s, err := decoder.buffer.Read16()
	if err != nil {
		return nil, err
	}
	return NewShort(int16(s)), nil
}

// Decodes a UShort.
// @return The decoded UShort.
func (decoder *BinaryDecoder) DecodeUShort() (*UShort, error) {
	s, err := decoder.buffer.Read16()
	if err != nil {
		return nil, err
	}
	return NewUShort(uint16(s)), nil
}

// Decodes an Integer.
// @return The decoded Integer.
func (decoder *BinaryDecoder) DecodeInteger() (*Integer, error) {
	i, err := decoder.buffer.Read32()
	if err != nil {
		return nil, err
	}
	return NewInteger(int32(i)), nil
}

// Decodes a UInteger.
// @return The decoded UInteger.
func (decoder *BinaryDecoder) DecodeUInteger() (*UInteger, error) {
	i, err := decoder.buffer.Read32()
	if err != nil {
		return nil, err
	}
	return NewUInteger(i), nil
}

// Decodes a Long.
// @return The decoded Long.
func (decoder *BinaryDecoder) DecodeLong() (*Long, error) {
	l, err := decoder.buffer.Read64()
	if err != nil {
		return nil, err
	}
	return NewLong(int64(l)), nil
}

// Decodes a ULong.
// @return The decoded ULong.
func (decoder *BinaryDecoder) DecodeULong() (*ULong, error) {
	l, err := decoder.buffer.Read64()
	if err != nil {
		return nil, err
	}
	return NewULong(l), nil
}

// TODO (AF): Declares this method in Decoder interface then implements
// DecodeString, DecodeIdentifier and DecodeURI in GenDecoder.
func (decoder *BinaryDecoder) readString() (string, error) {
	length, err := decoder.buffer.Read32()
	if err != nil {
		return "", err
	}
	// TODO (AF): We may avoid a data copy getting bytes directly in source buffer.
	buf := make([]byte, length)
	err = decoder.buffer.ReadBytes(buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// Decodes a String.
// @return The decoded String.
func (decoder *BinaryDecoder) DecodeString() (*String, error) {
	str, err := decoder.readString()
	if err != nil {
		return nil, err
	}
	return NewString(str), nil
}

// Decodes a Blob.
// @return The decoded Blob.
func (decoder *BinaryDecoder) DecodeBlob() (*Blob, error) {
	length, err := decoder.buffer.Read32()
	if err != nil {
		return nil, err
	}
	// TODO (AF): We may avoid a data copy getting bytes directly in source buffer.
	buf := Blob(make([]byte, length))
	err = decoder.buffer.ReadBytes(buf)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

// Decodes an Identifier.
// @return The decoded Identifier.
func (decoder *BinaryDecoder) DecodeIdentifier() (*Identifier, error) {
	str, err := decoder.readString()
	if err != nil {
		return nil, err
	}
	return NewIdentifier(str), nil
}

// Decodes a Duration.
// @return The decoded Duration.
func (decoder *BinaryDecoder) DecodeDuration() (*Duration, error) {
	d, err := decoder.buffer.Read64()
	if err != nil {
		return nil, err
	}
	return NewDuration(math.Float64frombits(d)), nil
}

// Decodes a Time.
// @return The decoded Time.
func (decoder *BinaryDecoder) DecodeTime() (*Time, error) {
	days, err := decoder.buffer.Read16()
	if err != nil {
		return nil, err
	}
	millis, err := decoder.buffer.Read32()
	if err != nil {
		return nil, err
	}

	timestamp := -NANOS_FROM_CCSDS_TO_UNIX_EPOCH
	timestamp += (int64(days) * NANOS_IN_DAY)
	timestamp += (int64(millis) * 1000000)

	t := time.Unix(timestamp/1000000000, timestamp%1000000000)
	return NewTime(t), nil
}

// Decodes a FineTime.
// @return The decoded FineTime.
func (decoder *BinaryDecoder) DecodeFineTime() (*FineTime, error) {
	days, err := decoder.buffer.Read16()
	if err != nil {
		return nil, err
	}
	millis, err := decoder.buffer.Read32()
	if err != nil {
		return nil, err
	}
	picos, err := decoder.buffer.Read32()
	if err != nil {
		return nil, err
	}

	timestamp := -NANOS_FROM_CCSDS_TO_UNIX_EPOCH
	timestamp += int64(days) * NANOS_IN_DAY
	timestamp += int64(millis) * 1000000
	timestamp += int64(picos / 1000)
	t := time.Unix(timestamp/1000000000, timestamp%1000000000)

	return NewFineTime(t), nil
}

// Decodes a URI.
// @return The decoded URI.
func (decoder *BinaryDecoder) DecodeURI() (*URI, error) {
	str, err := decoder.readString()
	if err != nil {
		return nil, err
	}
	return NewURI(str), nil
}

// TODO (AF): Handling of enumeration

func (decoder *BinaryDecoder) DecodeSmallEnum() (uint8, error) {
	o, err := decoder.buffer.Read()
	if err != nil {
		return 0, err
	}
	return o, nil
}

func (decoder *BinaryDecoder) DecodeMediumEnum() (uint16, error) {
	s, err := decoder.buffer.Read16()
	if err != nil {
		return 0, err
	}
	return s, nil
}

func (decoder *BinaryDecoder) DecodelargeEnum() (uint32, error) {
	i, err := decoder.buffer.Read32()
	if err != nil {
		return 0, err
	}
	return i, nil
}
