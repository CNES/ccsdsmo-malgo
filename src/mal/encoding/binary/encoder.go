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
	"errors"
	. "mal"
	"math"
	"time"
)

type BinaryEncoder struct {
	GenEncoder
	Varint bool
	Out    Buffer
}

// TODO (AF): To remove
//func (encoder *BinaryEncoder) InitBinaryEncoder(varint bool) {
//	encoder.GenEncoder.Encoder = encoder
//	encoder.Varint = varint
//}

// Creates a new encoder using a slice with sufficient capacity to encode datas.
// If the slice is not empty the encoded datas are append afterwards.
func NewBinaryEncoder(buf []byte, varint bool) *BinaryEncoder {
	encoder := &BinaryEncoder{
		Varint: varint,
		Out: &BinaryBuffer{
			Offset: 0,
			Buf:    buf,
		},
	}
	encoder.GenEncoder.Encoder = encoder
	return encoder
}

// Returns a slice containing all encoded data as needed to be sent
func (encoder *BinaryEncoder) Body() []byte {
	return encoder.Out.(*BinaryBuffer).Buf
}

// TODO (AF): No longer needed
// Writes an unsigned varint as defined in 5.25 section of the specification.
//func (encoder *BinaryEncoder) writeUVarInt(value uint64) error {
//	for (value & 0xFFFFFFFFFFFFFF80) != 0 {
//		err := encoder.Out.Write(byte((value & 0x7F) | 0x80))
//		if err != nil {
//			return err
//		}
//		value >>= 7
//	}
//	return encoder.Out.Write(byte(value & 0x7F))
//}

// ================================================================================
// These methods are not part of the Encoder interface, they are needed to encode
// message header

func (encoder *BinaryEncoder) Write(b byte) error {
	return encoder.Out.Write(b)
}

func (encoder *BinaryEncoder) WriteUInt32(i uint32) error {
	return encoder.Out.Write32(i)
}

func (encoder *BinaryEncoder) WriteBody(buf []byte) error {
	return encoder.Out.WriteBytes(buf)
}

// ================================================================================
// Implements Encoder interface

func (encoder *BinaryEncoder) EncodeNull() error {
	// Encode False (presence flag) -> 0
	return encoder.Out.WriteFlag(false)
}

func (encoder *BinaryEncoder) EncodeNotNull() error {
	// Encode True (presence flag) -> 1
	return encoder.Out.WriteFlag(true)
}

// Encodes the short form of an attribute.
func (encoder *BinaryEncoder) EncodeAttributeType(typeval Integer) error {
	return encoder.Out.Write(byte(typeval))
}

// Encodes a non-null Boolean.
// @param att The Boolean to encode.
func (encoder *BinaryEncoder) EncodeBoolean(att *Boolean) error {
	if *att {
		return encoder.Out.WriteFlag(true)
	} else {
		return encoder.Out.WriteFlag(false)
	}
}

// Encodes a non-null Float.
// @param att The Float to encode.
func (encoder *BinaryEncoder) EncodeFloat(att *Float) error {
	val := math.Float32bits(float32(*att))
	return encoder.Out.Write32(val)
}

// Encodes a non-null Double.
// @param att The Double to encode.
func (encoder *BinaryEncoder) EncodeDouble(att *Double) error {
	val := math.Float64bits(float64(*att))
	return encoder.Out.Write64(val)
}

// Encodes a non-null Octet.
// @param att The Octet to encode.
func (encoder *BinaryEncoder) EncodeOctet(att *Octet) error {
	return encoder.Out.Write(byte(*att))
}

// Encodes a non-null UOctet.
// @param att The UOctet to encode.
func (encoder *BinaryEncoder) EncodeUOctet(att *UOctet) error {
	return encoder.Out.Write(byte(*att))
}

// Encodes a non-null Short.
// @param att The Short to encode.
func (encoder *BinaryEncoder) EncodeShort(att *Short) error {
	if encoder.Varint {
		value := int16(*att)
		return encoder.Out.WriteUVarInt(uint64((value<<1)^(value>>15)) & 0xFFFF)
	} else {
		return encoder.Out.Write16(uint16(*att))
	}
}

// Encodes a non-null UShort.
// @param att The UShort to encode.
func (encoder *BinaryEncoder) EncodeUShort(att *UShort) error {
	if encoder.Varint {
		return encoder.Out.WriteUVarInt(uint64(*att))
	} else {
		return encoder.Out.Write16(uint16(*att))
	}
}

// Encodes a non-null Integer.
// @param att The Integer to encode.
func (encoder *BinaryEncoder) EncodeInteger(att *Integer) error {
	if encoder.Varint {
		value := int32(*att)
		return encoder.Out.WriteUVarInt(uint64((value<<1)^(value>>31)) & 0xFFFFFFFF)
	} else {
		return encoder.Out.Write32(uint32(*att))
	}
}

// Encodes a non-null UInteger.
// @param att The UInteger to encode.
func (encoder *BinaryEncoder) EncodeUInteger(att *UInteger) error {
	if encoder.Varint {
		return encoder.Out.WriteUVarInt(uint64(*att))
	} else {
		return encoder.Out.Write32(uint32(*att))
	}
}

// Encodes a non-null Long.
// @param att The Long to encode.
func (encoder *BinaryEncoder) EncodeLong(att *Long) error {
	if encoder.Varint {
		value := int64(*att)
		return encoder.Out.WriteUVarInt(uint64((value << 1) ^ (value >> 63)))
	} else {
		return encoder.Out.Write64(uint64(*att))
	}
}

// Encodes a non-null ULong.
// @param att The ULong to encode.
func (encoder *BinaryEncoder) EncodeULong(att *ULong) error {
	if encoder.Varint {
		return encoder.Out.WriteUVarInt(uint64(*att))
	} else {
		return encoder.Out.Write64(uint64(*att))
	}
}

// Encodes a non-null String.
// @param att The String to encode.
func (encoder *BinaryEncoder) EncodeString(str *String) error {
	buf := []byte(*str)
	err := encoder.Out.Write32(uint32(len(buf)))
	if err != nil {
		return err
	}
	return encoder.Out.WriteBytes(buf)
}

// Encodes a non-null Blob.
// @param att The Blob to encode.
func (encoder *BinaryEncoder) EncodeBlob(blob *Blob) error {
	err := encoder.Out.Write32(uint32(len(*blob)))
	if err != nil {
		return err
	}
	return encoder.Out.WriteBytes([]byte(*blob))
}

// Encodes a non-null Identifier.
// @param att The Identifier to encode.
func (encoder *BinaryEncoder) EncodeIdentifier(id *Identifier) error {
	buf := []byte(*id)
	err := encoder.Out.Write32(uint32(len(buf)))
	if err != nil {
		return err
	}
	return encoder.Out.WriteBytes(buf)
}

// Encodes a non-null Duration.
// @param att The Duration to encode.
func (encoder *BinaryEncoder) EncodeDuration(att *Duration) error {
	val := math.Float64bits(float64(*att))
	return encoder.Out.Write64(val)
}

// Encodes a non-null Time.
// @param att The Time to encode.
func (encoder *BinaryEncoder) EncodeTime(t *Time) error {
	timestamp := int64(time.Time(*t).UnixNano())
	timestamp += NANOS_FROM_CCSDS_TO_UNIX_EPOCH

	days := timestamp / NANOS_IN_DAY
	millis := (timestamp % NANOS_IN_DAY) / 1000000

	if days > 65535 {
		// TODO (AF): Verify if days can be negative in CCSDS Time format
		return errors.New("Cannot encode Time: " + time.Time(*t).String())
	}

	encoder.Out.Write16(uint16(days))
	encoder.Out.Write32(uint32(millis))

	return nil
}

// Encodes a non-null FineTime.
// @param att The FineTime to encode.
func (encoder *BinaryEncoder) EncodeFineTime(t *FineTime) error {
	timestamp := int64(time.Time(*t).UnixNano())
	timestamp += NANOS_FROM_CCSDS_TO_UNIX_EPOCH

	days := timestamp / NANOS_IN_DAY
	millis := (timestamp % NANOS_IN_DAY) / 1000000
	picos := ((timestamp % NANOS_IN_DAY) % 1000000) * 1000

	if days > 65535 {
		// TODO (AF): Verify if days can be negative in CCSDS Time format
		return errors.New("Cannot encode FineTime: " + time.Time(*t).String())
	}

	encoder.Out.Write16(uint16(days))
	encoder.Out.Write32(uint32(millis))
	encoder.Out.Write32(uint32(picos))

	return nil
}

// Encodes a non-null URI.
// @param att The URI to encode.
// @throws IllegalArgumentException If the argument is null.
func (encoder *BinaryEncoder) EncodeURI(uri *URI) error {
	buf := []byte(*uri)
	err := encoder.Out.Write32(uint32(len(buf)))
	if err != nil {
		return err
	}
	return encoder.Out.WriteBytes(buf)
}

// TODO (AF): Handling of enumeration

func (encoder *BinaryEncoder) EncodeSmallEnum(ordinal uint8) error {
	return encoder.Out.Write(ordinal)
}

func (encoder *BinaryEncoder) EncodeMediumEnum(ordinal uint16) error {
	return encoder.Out.Write16(ordinal)
}

func (encoder *BinaryEncoder) EncodelargeEnum(ordinal uint32) error {
	return encoder.Out.Write32(ordinal)
}
