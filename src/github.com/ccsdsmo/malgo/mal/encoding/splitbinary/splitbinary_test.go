/**
 * MIT License
 *
 * Copyright (c) 2018 CNES
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
package splitbinary_test

import (
	. "github.com/ccsdsmo/malgo/mal"
	"github.com/ccsdsmo/malgo/mal/encoding/splitbinary"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	i := NewInteger(15)
	err := encoder.EncodeInteger(i)
	if err != nil {
		t.Fatalf("Error during encode: 15")
	}
	b := NewBoolean(false)
	err = encoder.EncodeBoolean(b)
	if err != nil {
		t.Fatalf("Error during encode: false")
	}
	b = NewBoolean(true)
	err = encoder.EncodeBoolean(b)
	if err != nil {
		t.Fatalf("Error during encode: true")
	}
	b = NewBoolean(false)
	err = encoder.EncodeBoolean(b)
	if err != nil {
		t.Fatalf("Error during encode: false")
	}
	b = NewBoolean(true)
	err = encoder.EncodeBoolean(b)
	if err != nil {
		t.Fatalf("Error during encode: true")
	}
	b = NewBoolean(true)
	err = encoder.EncodeBoolean(b)
	if err != nil {
		t.Fatalf("Error during encode: true")
	}
	b = NewBoolean(false)
	err = encoder.EncodeBoolean(b)
	if err != nil {
		t.Fatalf("Error during encode: false")
	}
	b = NewBoolean(false)
	err = encoder.EncodeBoolean(b)
	if err != nil {
		t.Fatalf("Error during encode: false")
	}
	b = NewBoolean(false)
	err = encoder.EncodeBoolean(b)
	if err != nil {
		t.Fatalf("Error during encode: false")
	}
	b = NewBoolean(false)
	err = encoder.EncodeBoolean(b)
	if err != nil {
		t.Fatalf("Error during encode: false")
	}
	s := NewUShort(30)
	err = encoder.EncodeUShort(s)
	if err != nil {
		t.Fatalf("Error during encode: %d", s)
	}
	l := NewLong(-30)
	err = encoder.EncodeLong(l)
	if err != nil {
		t.Fatalf("Error during encode: %d", l)
	}

	var buffer *splitbinary.SplitBinaryBuffer = encoder.Out.(*splitbinary.SplitBinaryBuffer)
	t.Logf("SplitBinaryBuffer=", buffer.Buf, buffer.Offset, buffer.Bitfield, buffer.Bitfield_len, buffer.Bitfield_idx)
	//	decoder := splitbinary.NewSplitBinaryDecoder2(buffer.Buf, buffer.Bitfield, buffer.Bitfield_len)
	buf = encoder.Body()
	t.Logf("SplitBinaryBuffer=", buf)
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	i, err = decoder.DecodeInteger()
	if err != nil {
		t.Fatalf("Error during decode")
	}
	t.Log("decode: ", *i)

	b, err = decoder.DecodeBoolean()
	if err != nil {
		t.Fatalf("Error during decode")
	}
	t.Log("decode: ", *b)

	b, err = decoder.DecodeBoolean()
	if err != nil {
		t.Fatalf("Error during decode")
	}
	t.Log("decode: ", *b)

	b, err = decoder.DecodeBoolean()
	if err != nil {
		t.Fatalf("Error during decode")
	}
	t.Log("decode: ", *b)

	s, err = decoder.DecodeUShort()
	if err != nil {
		t.Fatalf("Error during decode")
	}
	t.Log("decode: ", *s)
}

func TestUOctet(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var o UOctet = 0
	for {
		err := encoder.EncodeUOctet(&o)
		if err != nil {
			t.Fatalf("Error during encode: %d", o)
		}
		o += 1
		if o == 0 {
			break
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i := 0; i < 256; i++ {
		o, err := decoder.DecodeUOctet()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if *o != UOctet(i) {
			t.Errorf("Bad decoding, got: %d, want: %d", o, i)
		}
	}
}

// TODO (AF): Test NullableUOctet

func TestOctet(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var o Octet = -128
	for {
		err := encoder.EncodeOctet(&o)
		if err != nil {
			t.Fatalf("Error during encode: %d", o)
		}
		o += 1
		if o == -128 {
			break
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i := -128; i < 128; i++ {
		o, err := decoder.DecodeOctet()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if *o != Octet(i) {
			t.Errorf("Bad decoding, got: %d, want: %d", o, i)
		}
	}
}

// TODO (AF): Test NullableOctet

func TestUShort(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var s UShort = 0
	for {
		err := encoder.EncodeUShort(&s)
		if err != nil {
			t.Fatalf("Error during encode: %d", s)
		}
		s += 1
		if s == 0 {
			break
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i := 0; i < 65536; i++ {
		s, err := decoder.DecodeUShort()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if *s != UShort(i) {
			t.Errorf("Bad decoding, got: %d, want: %d", s, i)
		}
	}
}

// TODO (AF): Test NullableUShort

func TestShort(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var s Short = -32768
	for {
		err := encoder.EncodeShort(&s)
		if err != nil {
			t.Fatalf("Error during encode: %d", s)
		}
		s += 1
		if s == -32768 {
			break
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i := -32768; i < 32768; i++ {
		s, err := decoder.DecodeShort()
		if err != nil {
			t.Fatalf("Error during decode (%s): %d", err, i)
		}
		if *s != Short(i) {
			t.Errorf("Bad decoding, got: %d, want: %d", s, i)
		}
	}
}

// TODO (AF): Test NullableShort

func TestUInteger(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list [65536]UInteger
	for i := 0; i < len(list); i++ {
		list[i] = UInteger(rand.Uint32())
	}
	for _, x := range list {
		err := encoder.EncodeUInteger(&x)
		if err != nil {
			t.Fatalf("Error during encode: %d", x)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeUInteger()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if x != *y {
			t.Errorf("Bad decoding, got: %d, want: %d", x, *y)
		}
	}
}

func TestNullableUInteger(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list = []*UInteger{
		nil,
		NewUInteger(rand.Uint32()),
		nil,
		NewUInteger(rand.Uint32()),
		NewUInteger(rand.Uint32()),
		nil,
		nil,
		NewUInteger(rand.Uint32()),
		NewUInteger(rand.Uint32()),
		nil,
	}

	for _, x := range list {
		err := encoder.EncodeNullableUInteger(x)
		if err != nil {
			t.Fatalf("Error during encode: %d", x)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeNullableUInteger()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if x == nil {
			if y != nil {
				t.Errorf("Bad decoding #%d, got: %d, want: nil", i, *y)
			}
			continue
		}
		if *x != *y {
			t.Errorf("Bad decoding, got: %d, want: %d", *x, *y)
		}
	}
}

func TestInteger(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list [65536]Integer
	for i := 0; i < len(list); i++ {
		list[i] = Integer(rand.Uint32())
	}
	for _, x := range list {
		err := encoder.EncodeInteger(&x)
		if err != nil {
			t.Fatalf("Error during encode: %d", x)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeInteger()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if x != *y {
			t.Errorf("Bad decoding, got: %d, want: %d", x, *y)
		}
	}
}

// TODO (AF): Test NullableInteger

func TestULong(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list [65536]ULong
	for i := 0; i < len(list); i++ {
		list[i] = ULong(rand.Uint64())
	}
	for _, x := range list {
		err := encoder.EncodeULong(&x)
		if err != nil {
			t.Fatalf("Error during encode: %d", x)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeULong()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if x != *y {
			t.Errorf("Bad decoding, got: %d, want: %d", x, *y)
		}
	}
}

// TODO (AF): Test NullableULong

func TestLong(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list [65536]Long
	for i := 0; i < len(list); i++ {
		list[i] = Long(rand.Uint64())
	}
	for _, x := range list {
		err := encoder.EncodeLong(&x)
		if err != nil {
			t.Fatalf("Error during encode: %d", x)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeLong()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if x != *y {
			t.Errorf("Bad decoding, got: %d, want: %d", x, *y)
		}
	}
}

// TODO (AF): Test NullableLong

func TestFloat(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list [65536]Float
	for i := 0; i < len(list); i++ {
		list[i] = Float(2 * (rand.Float32() - 0.5) * 3.4E+38)
	}
	for _, x := range list {
		err := encoder.EncodeFloat(&x)
		if err != nil {
			t.Fatalf("Error during encode: %d", x)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeFloat()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if x != *y {
			t.Errorf("Bad decoding, got: %d, want: %d", x, *y)
		}
	}
}

// TODO (AF): Test NullableFloat

func TestDouble(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list [65536]Double
	for i := 0; i < len(list); i++ {
		list[i] = Double(rand.NormFloat64())
	}
	for _, x := range list {
		err := encoder.EncodeDouble(&x)
		if err != nil {
			t.Fatalf("Error during encode: %d", x)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeDouble()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if x != *y {
			t.Errorf("Bad decoding, got: %d, want: %d", x, *y)
		}
	}
}

// TODO (AF): Test NullableDouble

func TestBlob(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list [16384]Blob
	for i := 0; i < len(list); i++ {
		n := rand.Intn(512)
		b := make([]byte, n)
		rand.Read(b)
		list[i] = Blob(b)
	}
	for _, x := range list {
		err := encoder.EncodeBlob(&x)
		if err != nil {
			t.Fatalf("Error during encode: %d", x)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeBlob()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		// NOTE (AF): Be careful the DeepEqual method is costly :-(
		if !reflect.DeepEqual(x, *y) {
			t.Errorf("Bad decoding, got: %d, want: %d", x, y)
			return
		}
	}
}

// TODO (AF): Test NullableBlob

func TestString(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list [65536]String
	for i := 0; i < len(list); i++ {
		n := rand.Intn(1024)
		b := make([]byte, n)
		rand.Read(b)
		list[i] = String(string(b))
	}
	for _, x := range list {
		err := encoder.EncodeString(&x)
		if err != nil {
			t.Fatalf("Error during encode: %d", x)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeString()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if x != *y {
			t.Errorf("Bad decoding, got: %d, want: %d", x, *y)
		}
	}
}

func TestNullableString(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list = []*String{
		nil,
		NewString("abcde"),
		nil,
		NewString("fghi"),
		NewString("jklmnop"),
		nil,
		nil,
		NewString("rstuvw"),
		NewString("xyz"),
		nil,
	}
	for _, x := range list {
		err := encoder.EncodeNullableString(x)
		if err != nil {
			t.Fatalf("Error during encode: %d", x)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeNullableString()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if x == nil {
			if y != nil {
				t.Errorf("Bad decoding #%d, got: %d, want: nil", i, *y)
			}
			continue
		}

		if *x != *y {
			t.Errorf("Bad decoding, got: %d, want: %d", x, *y)
		}
	}
}

func TestTime(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list [65536]Time
	list[0] = *TimeNow()
	for i := 1; i < len(list); i++ {
		list[i] = Time(time.Unix(int64(rand.Intn(2429913600)), int64(rand.Intn(999999999))))
	}
	for _, x := range list {
		err := encoder.EncodeTime(&x)
		if err != nil {
			t.Fatalf("Error during encode: %s", err)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeTime()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if (time.Time(x).UnixNano() / 1000000) != (time.Time(*y).UnixNano() / 1000000) {
			t.Errorf("Bad decoding, got: %d, want: %d", *y, x)
		}
	}
}

// TODO (AF): Test NullableTime

func TestFineTime(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list [65536]FineTime
	list[0] = *FineTimeNow()
	for i := 1; i < len(list); i++ {
		list[i] = FineTime(time.Unix(int64(rand.Intn(2429913600)), int64(rand.Intn(999999999))))
	}
	for _, x := range list {
		err := encoder.EncodeFineTime(&x)
		if err != nil {
			t.Fatalf("Error during encode: %s", err)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeFineTime()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if (time.Time(x).UnixNano()) != (time.Time(*y).UnixNano()) {
			t.Errorf("Bad decoding, got: %d, want: %d", *y, x)
		}
	}
}

// TODO (AF): Test NullableFineTime

// TODO (AF): Test Duration, Identifier, URI, .. and Nullable associated.

func TestAttribute(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list = []Attribute{
		NewBoolean(false),
		NewBoolean(true),
		NewDuration(12.34E56),
		NewDuration(1000),
		NewFloat(123.45E12),
		NewDouble(12345.67E89),
		NewIdentifier("BOOM"),
		NewOctet(-127),
		NewUOctet(255),
		NewShort(-32767),
		NewUShort(65535),
		NewInteger(123456789),
		NewUInteger(987654321),
		NewLong(-1),
		NewULong(9223372036854775807),
		NewString("Hello world"),
		NewDouble(9.99999E99),
		NewOctet(0),
		NewURI("http://www.scalagent.com"),
		NewUOctet(0),
	}
	for _, x := range list {
		err := encoder.EncodeAttribute(x)
		if err != nil {
			t.Fatalf("Error during encode: %s", err)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeAttribute()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		//		if x != y {
		//			t.Errorf("Bad decoding, got: %t, want: %t", y, x)
		//		}
		if reflect.DeepEqual(&x, y) {
			t.Errorf("Bad decoding, got: %t, want: %t", y, x)
		}
	}
}

func TestNullableAttribute(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list = []Attribute{
		NewBoolean(false),
		NullBoolean,
		NewBoolean(true),
		NewDuration(12.34E56),
		NullDuration,
		NewDuration(1000),
		NewFloat(123.45E12),
		NewDouble(12345.67E89),
		NullUInteger,
		NewIdentifier("BOOM"),
		NewOctet(-127),
		NewUOctet(255),
		NewShort(-32767),
		NewUShort(65535),
		NewInteger(123456789),
		NullInteger,
		NullOctet,
		NewUInteger(987654321),
		NewLong(-1),
		NewULong(9223372036854775807),
		NewString("Hello world"),
		NullString,
		NewDouble(9.99999E99),
		NewOctet(0),
		NewURI("http://www.scalagent.com"),
		NewUOctet(0),
	}
	for _, x := range list {
		err := encoder.EncodeNullableAttribute(x)
		if err != nil {
			t.Fatalf("Error during encode: %s", err)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeNullableAttribute()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		//		if x != y {
		//			t.Errorf("Bad decoding, got: %t, want: %t", y, x)
		//		}
		if reflect.DeepEqual(&x, y) {
			t.Errorf("Bad decoding, got: %t, want: %t", y, x)
		}
	}
}

func TestElement(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list = []Element{
		NewBoolean(false),
		NewBoolean(true),
		NewDuration(12.34E56),
		NewDuration(1000),
		NewFloat(123.45E12),
		NewDouble(12345.67E89),
		NewIdentifier("BOOM"),
		NewOctet(-127),
		NewUOctet(255),
		NewShort(-32767),
		NewUShort(65535),
		NewInteger(123456789),
		NewUInteger(987654321),
		NewLong(-1),
		NewULong(9223372036854775807),
		NewString("Hello world"),
		NewDouble(9.99999E99),
		NewOctet(0),
		NewURI("http://www.scalagent.com"),
		NewUOctet(0),
	}
	for _, x := range list {
		err := encoder.EncodeElement(x)
		if err != nil {
			t.Fatalf("Error during encode: %s", err)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeElement(x)
		if err != nil {
			t.Fatalf("Error during decode: %d", i)

		}
		if reflect.DeepEqual(&x, y) {
			t.Errorf("Bad decoding, got: %t, want: %t", y, x)
		}
	}
}

func TestNullableElement(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list = []Element{
		NewBoolean(false),
		NullBlob,
		NewBoolean(true),
		NewDuration(12.34E56),
		NullBoolean,
		NullURI,
		NewDuration(1000),
		NewFloat(123.45E12),
		NewDouble(12345.67E89),
		NullString,
		NewIdentifier("BOOM"),
		NewOctet(-127),
		NewUOctet(255),
		NewShort(-32767),
		NewUShort(65535),
		NewInteger(123456789),
		NullOctet,
		NullBlob,
		NullInteger,
		NewUInteger(987654321),
		NewLong(-1),
		NewULong(9223372036854775807),
		NewString("Hello world"),
		NewDouble(9.99999E99),
		NewOctet(0),
		NewURI("http://www.scalagent.com"),
		NewUOctet(0),
		NullBlobList,
		NullStringList,
		NullBlobList,
	}
	for _, x := range list {
		t.Log("Encode: ", x)
		err := encoder.EncodeNullableElement(x)
		if err != nil {
			t.Fatalf("Error during encode: %s", err)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeNullableElement(x)
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if reflect.DeepEqual(&x, y) {
			t.Errorf("Bad decoding, got: %t, want: %t", y, x)
		}
	}
}

func TestBlobAttribute(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list = []Attribute{
		&Blob{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		&Blob{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	}
	for _, x := range list {
		err := encoder.EncodeAttribute(x)
		if err != nil {
			t.Fatalf("Error during encode: %s", err)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeAttribute()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if !reflect.DeepEqual(x, y) {
			t.Errorf("Bad decoding, got: %t, want: %t", y, x)
		}
	}
}

func TestTimeAttribute(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list = []Attribute{
		NewTime(time.Now().Truncate(time.Millisecond)),
		NewTime(time.Unix(int64(rand.Intn(2429913600)), int64(rand.Intn(999))*1000000)),
		NewTime(time.Unix(int64(rand.Intn(2429913600)), int64(rand.Intn(999))*1000000)),
		NewTime(time.Unix(int64(rand.Intn(2429913600)), int64(rand.Intn(999))*1000000)),
		NewTime(time.Unix(int64(rand.Intn(2429913600)), int64(rand.Intn(999))*1000000)),
		NewFineTime(time.Now().Truncate(time.Nanosecond)),
		NewFineTime(time.Unix(int64(rand.Intn(2429913600)), int64(rand.Intn(999999999)))),
		NewFineTime(time.Unix(int64(rand.Intn(2429913600)), int64(rand.Intn(999999999)))),
		NewFineTime(time.Unix(int64(rand.Intn(2429913600)), int64(rand.Intn(999999999)))),
		NewFineTime(time.Unix(int64(rand.Intn(2429913600)), int64(rand.Intn(999999999)))),
	}
	for _, x := range list {
		err := encoder.EncodeAttribute(x)
		if err != nil {
			t.Fatalf("Error during encode: %s", err)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeAttribute()
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if !reflect.DeepEqual(x, y) {
			t.Errorf("Bad decoding #%d, got: %t, want: %t", i, y, x)
		}
	}
}

// TODO (AF): Test Element, Composite, List, etc.

func TestBooleanList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var idlist = []*Boolean{
		NewBoolean(true),
		NewBoolean(false),
		nil,
		NewBoolean(true),
	}
	var x = BooleanList(idlist)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeBooleanList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

func TestOctetList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var idlist = []*Octet{
		NewOctet(-128),
		NewOctet(0),
		nil,
		NewOctet(127),
		NullOctet,
	}
	var x = OctetList(idlist)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeOctetList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

func TestUOctetList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var idlist = []*UOctet{
		NewUOctet(255),
		NewUOctet(0),
		nil,
		NewUOctet(127),
		NullUOctet,
	}
	var x = UOctetList(idlist)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeUOctetList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

func TestShortList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var idlist = []*Short{
		NewShort(-32768),
		NewShort(0),
		nil,
		NewShort(32767),
		NullShort,
	}
	var x = ShortList(idlist)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeShortList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

func TestUShortList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var idlist = []*UShort{
		NewUShort(65535),
		NewUShort(0),
		nil,
		NewUShort(127),
		NullUShort,
	}
	var x = UShortList(idlist)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeUShortList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

func TestIntegerList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var idlist = []*Integer{
		NewInteger(int32(INTEGER_MIN)),
		NewInteger(0),
		nil,
		NewInteger(int32(INTEGER_MAX)),
		NullInteger,
	}
	var x = IntegerList(idlist)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeIntegerList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

func TestUIntegerList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var idlist = []*UInteger{
		NewUInteger(uint32(UINTEGER_MIN)),
		NewUInteger(0),
		nil,
		NewUInteger(uint32(UINTEGER_MAX)),
		NullUInteger,
	}
	var x = UIntegerList(idlist)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeUIntegerList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

func TestLongList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var idlist = []*Long{
		NewLong(int64(LONG_MIN)),
		NewLong(0),
		nil,
		NewLong(int64(LONG_MAX)),
		NullLong,
	}
	var x = LongList(idlist)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeLongList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

func TestULongList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var idlist = []*ULong{
		NewULong(uint64(ULONG_MIN)),
		NewULong(0),
		nil,
		NewULong(uint64(ULONG_MAX)),
		NullULong,
	}
	var x = ULongList(idlist)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeULongList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

func TestFloatList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var idlist = []*Float{
		NewFloat(12.4),
		NewFloat(16.54E12),
		NullFloat,
		NewFloat(-17.99E-15),
		nil,
	}
	var x = FloatList(idlist)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeFloatList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

func TestTimeList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var idlist = []*Time{
		NewTime(time.Now().Truncate(time.Millisecond)),
		NewTime(time.Unix(int64(rand.Intn(2429913600)), int64(rand.Intn(999))*1000000)),
		nil,
		NewTime(time.Unix(int64(rand.Intn(2429913600)), int64(rand.Intn(999))*1000000)),
		NullTime,
	}
	var x = TimeList(idlist)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeTimeList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

func TestFineTimeList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var idlist = []*FineTime{
		NewFineTime(time.Now().Truncate(time.Nanosecond)),
		NewFineTime(time.Unix(int64(rand.Intn(2429913600)), int64(rand.Intn(999999999)))),
		nil,
		NewFineTime(time.Unix(int64(rand.Intn(2429913600)), int64(rand.Intn(999999999)))),
		NullFineTime,
	}
	var x = FineTimeList(idlist)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeFineTimeList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

func TestIdentifierList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var idlist = []*Identifier{
		NewIdentifier("DOMAIN1"),
		NewIdentifier("DOMAIN2"),
		nil,
		NewIdentifier("Domain4"),
		NullIdentifier,
	}
	var x = IdentifierList(idlist)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeIdentifierList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

func TestEntityKey(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list = []*EntityKey{
		&EntityKey{NewIdentifier("abcde"), NewLong(0), NewLong(1), NewLong(2)},
		&EntityKey{NewIdentifier("fghijklmn"), NewLong(3), NewLong(4), NewLong(5)},
		&EntityKey{nil, NewLong(6), NewLong(7), nil},
		&EntityKey{NewIdentifier("opqrst"), nil, nil, NewLong(2)},
		&EntityKey{nil, nil, nil, nil},
	}
	for _, x := range list {
		err := x.Encode(encoder)
		if err != nil {
			t.Fatalf("Error during encode: %s", err)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeElement(list[0])
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if !reflect.DeepEqual(x, y) {
			t.Errorf("Bad decoding, got: %t, want: %t", y, x)
			t.Log(*x.FirstSubKey, *y.(*EntityKey).FirstSubKey)
		}
	}
}

func TestEntityKeyList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	key1 := EntityKey{NewIdentifier("abcde"), NewLong(0), NewLong(1), NewLong(2)}
	key2 := EntityKey{NewIdentifier("fghijklmn"), NewLong(3), NewLong(4), NewLong(5)}
	key3 := EntityKey{nil, NewLong(6), NewLong(7), nil}
	key4 := EntityKey{NewIdentifier("opqrst"), nil, nil, NewLong(2)}
	key5 := EntityKey{nil, nil, nil, nil}

	var list = []*EntityKey{
		&key1,
		&key2,
		nil,
		&key3,
		&key4,
		NullEntityKey,
		nil,
		&key5,
	}

	var x = EntityKeyList(list)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeEntityKeyList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

func TestEntityRequest(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	key1 := EntityKey{NewIdentifier("abcde"), NewLong(0), NewLong(1), NewLong(2)}
	key2 := EntityKey{NewIdentifier("fghijklmn"), NewLong(3), NewLong(4), NewLong(5)}

	var list = []*EntityRequest{
		&EntityRequest{
			NullIdentifierList,
			true,
			false,
			true,
			false,
			EntityKeyList{
				&key2,
			},
		},
		&EntityRequest{
			&IdentifierList{
				NewIdentifier("domain"),
			},
			true,
			false,
			true,
			false,
			EntityKeyList{
				&key1,
			},
		},
		&EntityRequest{
			&IdentifierList{
				NewIdentifier("domain1"),
				NewIdentifier("domain2"),
			},
			false,
			false,
			false,
			false,
			EntityKeyList{
				&key1,
				NullEntityKey,
				&key2,
			},
		},
	}
	for _, x := range list {
		err := x.Encode(encoder)
		if err != nil {
			t.Fatalf("Error during encode: %s", err)
		}
	}

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	for i, x := range list {
		y, err := decoder.DecodeElement(list[0])
		if err != nil {
			t.Fatalf("Error during decode: %d", i)
		}
		if !reflect.DeepEqual(x, y) {
			t.Errorf("Bad decoding, got: %t, want: %t", y, x)
		}
	}
}

func TestEntityRequestList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	key1 := EntityKey{NewIdentifier("abcde"), NewLong(0), NewLong(1), NewLong(2)}
	key2 := EntityKey{NewIdentifier("fghijklmn"), NewLong(3), NewLong(4), NewLong(5)}
	key3 := EntityKey{nil, NewLong(6), NewLong(7), nil}
	key4 := EntityKey{NewIdentifier("opqrst"), nil, nil, NewLong(2)}
	key5 := EntityKey{nil, nil, nil, nil}

	var list = []*EntityRequest{
		&EntityRequest{
			NullIdentifierList,
			true,
			false,
			true,
			false,
			EntityKeyList{
				&key2, &key3, &key5,
			},
		},
		&EntityRequest{
			&IdentifierList{
				NewIdentifier("domain"),
				NullIdentifier,
			},
			true,
			false,
			true,
			false,
			EntityKeyList{
				&key1, NullEntityKey, nil, &key2, &key4,
			},
		},
		&EntityRequest{
			&IdentifierList{
				NewIdentifier("domain1"),
				NewIdentifier("domain2"),
			},
			false,
			false,
			false,
			false,
			EntityKeyList{
				&key1, NullEntityKey, NullEntityKey, &key2,
			},
		},
		&EntityRequest{
			nil,
			false,
			false,
			false,
			false,
			EntityKeyList{
				&key1, NullEntityKey, NullEntityKey, &key2,
			},
		},
	}

	var x = EntityRequestList(list)
	x.Encode(encoder)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := DecodeEntityRequestList(decoder)
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, x)
	}
}

// TODO (AF): File, FileList, IdBooleanPair, IdBooleanPairList, NamedValue, NamedValueList,
// Pair, PairList, Subscription, SubscriptionList, UpdateHeader, UpdateHeaderList.

// Note (AF): see corresponding comment in encoder about generic view for lists.

func TestOctetElementList(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	var list = []Element{
		NewOctet(-128),
		NewOctet(0),
		NullOctet,
		NewOctet(127),
	}
	encoder.EncodeList(list)

	buf = encoder.Body()
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	y, err := decoder.DecodeList(NewOctet(0))
	if err != nil {
		t.Fatalf("Error during decode:", err)
	}
	if !reflect.DeepEqual(list, y) {
		t.Errorf("Bad decoding, got: %t, want: %t", y, list)
	}
}
