/**
 * MIT License
 *
 * Copyright (c) 2017 - 2019 CNES
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
	"errors"
	. "github.com/CNES/ccsdsmo-malgo/mal"
	"math"
	"time"
)

type BinaryDecoder struct {
	GenDecoder
	Varint bool
	In     Buffer
}

// Creates a new decoder using a slice containing binary data to decode.
func NewBinaryDecoder(buf []byte, varint bool) *BinaryDecoder {
	decoder := &BinaryDecoder{
		Varint: varint,
		In: &BinaryBuffer{
			Offset: 0,
			Buf:    buf,
		},
	}
	decoder.GenDecoder.Self = decoder
	return decoder
}

// ================================================================================
// These methods are not part of the Encoder interface, they are needed to encode
// message header

func (decoder *BinaryDecoder) Read() (byte, error) {
	return decoder.In.Read()
}

func (decoder *BinaryDecoder) ReadUInt32() (uint32, error) {
	return decoder.In.Read32()
}

// Returns the part of buffer that still needs to be decoded
func (decoder *BinaryDecoder) Remaining() []byte {
	return decoder.In.Remaining()
}

// Reads an unsigned varint as defined in 5.25 section of the specification.
//func (decoder *BinaryDecoder) readUVarInt() (uint64, error) {
//	var value uint64 = 0
//	var i uint = 0
//	for {
//		b, err := decoder.In.Read()
//		if err != nil {
//			return value, err
//		}
//
//		value |= uint64(b&0x7F) << i
//		if (b & 0x80) == 0 {
//			break
//		}
//		i += 7
//	}
//
//	return value, nil
//}

// ================================================================================
// Implements Decoder interface

func (decoder *BinaryDecoder) IsNull() (bool, error) {
	// TODO (AF): It should be simple to have a method from buffer that reads the
	// presence flag directly (or defining a IsPresent method in decoder).
	b, err := decoder.In.ReadFlag()
	if err != nil {
		return false, err
	}
	// Note: the encoded value corresponds to the presence flag.
	return (b == false), nil
}

// Decodes the short form of an attribute.
// @return The short form of the attribute.
func (decoder *BinaryDecoder) DecodeAttributeType() (Integer, error) {
	b, err := decoder.In.Read()
	if err != nil {
		return -1, err
	}
	return Integer(int(b)), nil
}

// Decodes a Boolean.
// @return The decoded Boolean.
func (decoder *BinaryDecoder) DecodeBoolean() (*Boolean, error) {
	b, err := decoder.In.ReadFlag()
	if err != nil {
		// TODO (AF): Should be nil?
		return new(Boolean), err
	}
	return NewBoolean(b), nil
}

// Decodes a Float.
// @return The decoded Float.
func (decoder *BinaryDecoder) DecodeFloat() (*Float, error) {
	if decoder.Varint {
		value, err := decoder.In.ReadUVarInt()
		if err != nil {
			return nil, err
		}
		if (value & 0xFFFFFFFF00000000) != 0 {
			return nil, errors.New("Error decoding varint integer: " + string(value))
		}
		var res int32 = 0
		if (value & 1) != 0 {
			res = -1
		}
		res = res ^ int32(value>>1)
		return NewFloat(math.Float32frombits(uint32(res))), nil
	} else {
		f, err := decoder.In.Read32()
		if err != nil {
			return nil, err
		}
		return NewFloat(math.Float32frombits(f)), nil
	}
}

// Decodes a Double.
// @return The decoded Double.
func (decoder *BinaryDecoder) DecodeDouble() (*Double, error) {
	if decoder.Varint {
		value, err := decoder.In.ReadUVarInt()
		if err != nil {
			return nil, err
		}
		var res int64 = 0
		if (value & 1) != 0 {
			res = -1
		}
		res = res ^ int64(value>>1)
		return NewDouble(math.Float64frombits(uint64(res))), nil
	} else {
		d, err := decoder.In.Read64()
		if err != nil {
			return nil, err
		}
		return NewDouble(math.Float64frombits(d)), nil
	}
}

// Decodes an Octet.
// @return The decoded Octet.
func (decoder *BinaryDecoder) DecodeOctet() (*Octet, error) {
	o, err := decoder.In.Read()
	if err != nil {
		return nil, err
	}
	return NewOctet(int8(o)), nil
}

// Decodes a UOctet.
// @return The decoded UOctet.
func (decoder *BinaryDecoder) DecodeUOctet() (*UOctet, error) {
	o, err := decoder.In.Read()
	if err != nil {
		return nil, err
	}
	return NewUOctet(uint8(o)), nil
}

// Decodes a Short.
// @return The decoded Short.
func (decoder *BinaryDecoder) DecodeShort() (*Short, error) {
	if decoder.Varint {
		value, err := decoder.In.ReadUVarInt()
		if err != nil {
			return nil, err
		}
		if (value & 0xFFFFFFFFFFFE0000) != 0 {
			return nil, errors.New("Error decoding varint short: " + string(value))
		}
		var res int16 = 0
		if (value & 1) != 0 {
			res = -1
		}
		res = res ^ int16(value>>1)
		return NewShort(res), nil
	} else {
		s, err := decoder.In.Read16()
		if err != nil {
			return nil, err
		}
		return NewShort(int16(s)), nil
	}
}

// Decodes a UShort.
// @return The decoded UShort.
func (decoder *BinaryDecoder) DecodeUShort() (*UShort, error) {
	if decoder.Varint {
		value, err := decoder.In.ReadUVarInt()
		if err != nil {
			return nil, err
		}
		if (value & 0xFFFFFFFFFFFF0000) != 0 {
			return nil, errors.New("Error decoding varint short: " + string(value))
		}
		return NewUShort(uint16(value)), nil
	} else {
		s, err := decoder.In.Read16()
		if err != nil {
			return nil, err
		}
		return NewUShort(s), nil
	}
}

// Decodes an Integer.
// @return The decoded Integer.
func (decoder *BinaryDecoder) DecodeInteger() (*Integer, error) {
	if decoder.Varint {
		value, err := decoder.In.ReadUVarInt()
		if err != nil {
			return nil, err
		}
		if (value & 0xFFFFFFFE00000000) != 0 {
			return nil, errors.New("Error decoding varint integer: " + string(value))
		}
		var res int32 = 0
		if (value & 1) != 0 {
			res = -1
		}
		res = res ^ int32(value>>1)
		return NewInteger(res), nil
	} else {
		i, err := decoder.In.Read32()
		if err != nil {
			return nil, err
		}
		return NewInteger(int32(i)), nil
	}
}

// Decodes a UInteger.
// @return The decoded UInteger.
func (decoder *BinaryDecoder) DecodeUInteger() (*UInteger, error) {
	if decoder.Varint {
		value, err := decoder.In.ReadUVarInt()
		if err != nil {
			return nil, err
		}
		if (value & 0xFFFFFFFF00000000) != 0 {
			return nil, errors.New("Error decoding varint integer: " + string(value))
		}
		return NewUInteger(uint32(value)), nil
	} else {
		i, err := decoder.In.Read32()
		if err != nil {
			return nil, err
		}
		return NewUInteger(i), nil
	}
}

// Decodes a Long.
// @return The decoded Long.
func (decoder *BinaryDecoder) DecodeLong() (*Long, error) {
	if decoder.Varint {
		value, err := decoder.In.ReadUVarInt()
		if err != nil {
			return nil, err
		}
		var res int64 = 0
		if (value & 1) != 0 {
			res = -1
		}
		res = res ^ int64(value>>1)
		return NewLong(res), nil
	} else {
		l, err := decoder.In.Read64()
		if err != nil {
			return nil, err
		}
		return NewLong(int64(l)), nil
	}
}

// Decodes a ULong.
// @return The decoded ULong.
func (decoder *BinaryDecoder) DecodeULong() (*ULong, error) {
	if decoder.Varint {
		value, err := decoder.In.ReadUVarInt()
		if err != nil {
			return nil, err
		}
		return NewULong(value), nil
	} else {
		l, err := decoder.In.Read64()
		if err != nil {
			return nil, err
		}
		return NewULong(l), nil
	}
}

// TODO (AF): Declares this method in Decoder interface then implements
// DecodeString, DecodeIdentifier and DecodeURI in GenDecoder.
// Implements a readBuf method (see Encoder.encodeBuf) and uses it in readString
// and decodeBlob.
func (decoder *BinaryDecoder) readString() (string, error) {
	var length uint32
	if decoder.Varint {
		value, err := decoder.In.ReadUVarInt()
		if err != nil {
			return "", err
		}
		if (value & 0xFFFFFFFF00000000) != 0 {
			return "", errors.New("Error decoding varint integer: " + string(value))
		}
		length = uint32(value)
	} else {
		value, err := decoder.In.Read32()
		if err != nil {
			return "", err
		}
		length = value
	}
	// TODO (AF): We may avoid a data copy getting bytes directly in source buffer.
	buf := make([]byte, length)
	err := decoder.In.ReadBytes(buf)
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

// Decodes an Identifier.
// @return The decoded Identifier.
func (decoder *BinaryDecoder) DecodeIdentifier() (*Identifier, error) {
	str, err := decoder.readString()
	if err != nil {
		return nil, err
	}
	return NewIdentifier(str), nil
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

// Decodes a Blob.
// @return The decoded Blob.
func (decoder *BinaryDecoder) DecodeBlob() (*Blob, error) {
	var length uint32
	if decoder.Varint {
		value, err := decoder.In.ReadUVarInt()
		if err != nil {
			return nil, err
		}
		if (value & 0xFFFFFFFF00000000) != 0 {
			return nil, errors.New("Error decoding varint integer: " + string(value))
		}
		length = uint32(value)
	} else {
		value, err := decoder.In.Read32()
		if err != nil {
			return nil, err
		}
		length = value
	}
	// TODO (AF): We may avoid a data copy getting bytes directly in source buffer.
	buf := Blob(make([]byte, length))
	err := decoder.In.ReadBytes(buf)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

// Decodes a Duration.
// @return The decoded Duration.
func (decoder *BinaryDecoder) DecodeDuration() (*Duration, error) {
	d, err := decoder.In.Read64()
	if err != nil {
		return nil, err
	}
	return NewDuration(math.Float64frombits(d)), nil
}

// Decodes a Time.
// @return The decoded Time.
func (decoder *BinaryDecoder) DecodeTime() (*Time, error) {
	days, err := decoder.In.Read16()
	if err != nil {
		return nil, err
	}
	millis, err := decoder.In.Read32()
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
	days, err := decoder.In.Read16()
	if err != nil {
		return nil, err
	}
	millis, err := decoder.In.Read32()
	if err != nil {
		return nil, err
	}
	picos, err := decoder.In.Read32()
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

// TODO (AF): Handling of enumeration

func (decoder *BinaryDecoder) DecodeSmallEnum() (uint8, error) {
	o, err := decoder.In.Read()
	if err != nil {
		return 0, err
	}
	return o, nil
}

func (decoder *BinaryDecoder) DecodeMediumEnum() (uint16, error) {
	s, err := decoder.In.Read16()
	if err != nil {
		return 0, err
	}
	return s, nil
}

func (decoder *BinaryDecoder) DecodelargeEnum() (uint32, error) {
	i, err := decoder.In.Read32()
	if err != nil {
		return 0, err
	}
	return i, nil
}
