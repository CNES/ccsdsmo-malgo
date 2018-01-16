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

type Buffer interface {
	Write(v byte) error
	Write16(v uint16) error
	Write32(v uint32) error
	Write64(v uint64) error
	WriteBytes(v []byte) error
	// Writes an unsigned varint as defined in 5.25 section of the specification.
	WriteUVarInt(value uint64) error
	WriteFlag(f bool) error
	// Reset the buffer allowing to reuse it
	Reset(write bool)
	Read() (byte, error)
	Read16() (uint16, error)
	Read32() (uint32, error)
	Read64() (uint64, error)
	ReadBytes(buf []byte) error
	// Reads an unsigned varint as defined in 5.25 section of the specification.
	ReadUVarInt() (uint64, error)
	ReadFlag() (bool, error)
	// Returns the part of buffer that still needs to be decoded
	Remaining() ([]byte, error)
}
