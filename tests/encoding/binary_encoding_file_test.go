package encoding

import (
	"bytes"
	. "github.com/ccsdsmo/malgo/mal"
	"github.com/ccsdsmo/malgo/mal/encoding/binary"
	"github.com/ccsdsmo/malgo/mal/encoding/splitbinary"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	BOOL1     Boolean  = false
	UOCTET1   UOctet   = 0
	UOCTET2   UOctet   = 255
	OCTET1    Octet    = -128
	OCTET2    Octet    = 0
	OCTET3    Octet    = 127
	BOOL2     Boolean  = true
	USHORT1   UShort   = 0
	USHORT2   UShort   = 256
	USHORT3   UShort   = 65535
	SHORT1    Short    = -32768
	SHORT2    Short    = -256
	SHORT3    Short    = 0
	SHORT4    Short    = 256
	SHORT5    Short    = 32767
	BOOL3     Boolean  = false
	UINT1     UInteger = 0
	UINT2     UInteger = 256
	UINT3     UInteger = 65536
	UINT4     UInteger = 4294967295
	INT1      Integer  = -2147483648
	INT2      Integer  = -32768
	INT3      Integer  = -256
	INT4      Integer  = 0
	INT5      Integer  = 256
	INT6      Integer  = 32767
	INT7      Integer  = 2147483647
	BOOL4     Boolean  = false
	ULONG1    ULong    = 0
	ULONG2    ULong    = 65536
	ULONG3    ULong    = 4294967295
	LONG1     Long     = -2147483648
	LONG2     Long     = 0
	LONG3     Long     = 2147483647
	BOOL5     Boolean  = false
	FLOAT1    Float    = 1.25E6
	FLOAT2    Float    = -5.8E-2
	BOOL6     Boolean  = false
	DOUBLE1   Double   = 1.25E6
	DOUBLE2   Double   = -5.8E-2
	BOOL7     Boolean  = false
	BLOB1     Blob     = Blob{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	BOOL8     Boolean  = false
	STRING1   String   = String("Hello world")
	BOOL9     Boolean  = true
	TIME1     Time     = Time(time.Unix(int64(1234567), int64(500)))
	BOOL10    Boolean  = true
	FINETIME1 FineTime = FineTime(time.Unix(int64(1234567), int64(500)))
	BOOL11    Boolean  = false
)

// Encodes with FixedBinary and writes in a file
func TestFixedBinaryEncoding(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	encoder := binary.NewBinaryEncoder(buf, false)

	testEncoding(t, encoder)

	buf = encoder.Body()
	os.Remove("gofixedbinary.data")
	err := ioutil.WriteFile("gofixedbinary.data", buf, 0644)
	if err != nil {
		t.Fatalf("Error writing file: %s", err)
	}

	if !deepCompare(t, "gofixedbinary.data", "gofixedbinary.ref") {
		t.Fatalf("TestFixedBinaryEncoding, files differ.")
	}
}

// Reads a file and decodes with FixedBinary
func TestFixedBinaryDecoding(t *testing.T) {
	buf, err := ioutil.ReadFile("gofixedbinary.ref")
	if err != nil {
		t.Fatalf("Error reading file: %s", err)
	}
	decoder := binary.NewBinaryDecoder(buf, false)

	testDecoding(t, decoder)
}

// Encodes with FixedBinary and writes in a file
func TestVarintBinaryEncoding(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	encoder := binary.NewBinaryEncoder(buf, true)

	testEncoding(t, encoder)

	buf = encoder.Body()
	os.Remove("govarintbinary.data")
	err := ioutil.WriteFile("govarintbinary.data", buf, 0644)
	if err != nil {
		t.Fatalf("Error writing file: %s", err)
	}

	if !deepCompare(t, "govarintbinary.data", "govarintbinary.ref") {
		t.Fatalf("TestVarintBinaryEncoding, files differ.")
	}
}

// Reads a file and decodes with FixedBinary
func TestVarintBinaryDecoding(t *testing.T) {
	buf, err := ioutil.ReadFile("govarintbinary.ref")
	if err != nil {
		t.Fatalf("Error reading file: %s", err)
	}
	decoder := binary.NewBinaryDecoder(buf, true)

	testDecoding(t, decoder)
}

// Encodes with SplitBinary and writes in a file
func TestSplitBinaryEncoding(t *testing.T) {
	var length uint32 = 8192
	buf := make([]byte, 0, length)
	bf := make([]byte, 0, 8)
	encoder := splitbinary.NewSplitBinaryEncoder(buf, bf)

	testEncoding(t, encoder)

	buf = encoder.Body()
	os.Remove("gosplitbinary.data")
	err := ioutil.WriteFile("gosplitbinary.data", buf, 0644)
	if err != nil {
		t.Fatalf("Error writing file: %s", err)
	}

	if !deepCompare(t, "gosplitbinary.data", "gosplitbinary.ref") {
		t.Fatalf("TestSplitBinaryEncoding, files differ.")
	}
}

// Reads a file and decodes with SplitBinary
func TestSplitBinaryDecoding(t *testing.T) {
	buf, err := ioutil.ReadFile("gosplitbinary.ref")
	if err != nil {
		t.Fatalf("Error reading file: %s", err)
	}
	decoder := splitbinary.NewSplitBinaryDecoder(buf)

	testDecoding(t, decoder)
}

func testEncoding(t *testing.T, encoder Encoder) {
	writeBoolean(t, encoder, BOOL1)
	writeUOctet(t, encoder, UOCTET1)
	writeUOctet(t, encoder, UOCTET2)
	writeOctet(t, encoder, OCTET1)
	writeOctet(t, encoder, OCTET2)
	writeOctet(t, encoder, OCTET3)
	writeBoolean(t, encoder, BOOL2)
	writeUShort(t, encoder, USHORT1)
	writeUShort(t, encoder, USHORT2)
	writeUShort(t, encoder, USHORT3)
	writeShort(t, encoder, SHORT1)
	writeShort(t, encoder, SHORT2)
	writeShort(t, encoder, SHORT3)
	writeShort(t, encoder, SHORT4)
	writeShort(t, encoder, SHORT5)
	writeBoolean(t, encoder, BOOL3)
	writeUInteger(t, encoder, UINT1)
	writeUInteger(t, encoder, UINT2)
	writeUInteger(t, encoder, UINT3)
	writeUInteger(t, encoder, UINT4)
	writeInteger(t, encoder, INT1)
	writeInteger(t, encoder, INT2)
	writeInteger(t, encoder, INT3)
	writeInteger(t, encoder, INT4)
	writeInteger(t, encoder, INT5)
	writeInteger(t, encoder, INT6)
	writeInteger(t, encoder, INT7)
	writeBoolean(t, encoder, BOOL4)
	writeULong(t, encoder, ULONG1)
	writeULong(t, encoder, ULONG2)
	writeULong(t, encoder, ULONG3)
	writeLong(t, encoder, LONG1)
	writeLong(t, encoder, LONG2)
	writeLong(t, encoder, LONG3)
	writeBoolean(t, encoder, BOOL5)
	writeFloat(t, encoder, FLOAT1)
	writeFloat(t, encoder, FLOAT2)
	writeBoolean(t, encoder, BOOL6)
	writeDouble(t, encoder, DOUBLE1)
	writeDouble(t, encoder, DOUBLE2)
	writeBoolean(t, encoder, BOOL7)
	writeBlob(t, encoder, BLOB1)
	writeBoolean(t, encoder, BOOL8)
	writeString(t, encoder, STRING1)
	writeBoolean(t, encoder, BOOL9)
	writeTime(t, encoder, TIME1)
	writeBoolean(t, encoder, BOOL10)
	writeFineTime(t, encoder, FINETIME1)
	writeBoolean(t, encoder, BOOL11)

}

func testDecoding(t *testing.T, decoder Decoder) {
	testBoolean(t, decoder, BOOL1)
	testUOctet(t, decoder, UOCTET1)
	testUOctet(t, decoder, UOCTET2)
	testOctet(t, decoder, OCTET1)
	testOctet(t, decoder, OCTET2)
	testOctet(t, decoder, OCTET3)
	testBoolean(t, decoder, BOOL2)
	testUShort(t, decoder, USHORT1)
	testUShort(t, decoder, USHORT2)
	testUShort(t, decoder, USHORT3)
	testShort(t, decoder, SHORT1)
	testShort(t, decoder, SHORT2)
	testShort(t, decoder, SHORT3)
	testShort(t, decoder, SHORT4)
	testShort(t, decoder, SHORT5)
	testBoolean(t, decoder, BOOL3)
	testUInteger(t, decoder, UINT1)
	testUInteger(t, decoder, UINT2)
	testUInteger(t, decoder, UINT3)
	testUInteger(t, decoder, UINT4)
	testInteger(t, decoder, INT1)
	testInteger(t, decoder, INT2)
	testInteger(t, decoder, INT3)
	testInteger(t, decoder, INT4)
	testInteger(t, decoder, INT5)
	testInteger(t, decoder, INT6)
	testInteger(t, decoder, INT7)
	testBoolean(t, decoder, BOOL4)
	testULong(t, decoder, ULONG1)
	testULong(t, decoder, ULONG2)
	testULong(t, decoder, ULONG3)
	testLong(t, decoder, LONG1)
	testLong(t, decoder, LONG2)
	testLong(t, decoder, LONG3)
	testBoolean(t, decoder, BOOL5)
	testFloat(t, decoder, FLOAT1)
	testFloat(t, decoder, FLOAT2)
	testBoolean(t, decoder, BOOL6)
	testDouble(t, decoder, DOUBLE1)
	testDouble(t, decoder, DOUBLE2)
	testBoolean(t, decoder, BOOL7)
	testBlob(t, decoder, BLOB1)
	testBoolean(t, decoder, BOOL8)
	testString(t, decoder, STRING1)
	testBoolean(t, decoder, BOOL9)
	testTime(t, decoder, TIME1)
	testBoolean(t, decoder, BOOL10)
	testFineTime(t, decoder, FINETIME1)
	testBoolean(t, decoder, BOOL11)

}

func writeBoolean(t *testing.T, encoder Encoder, b Boolean) {
	err := encoder.EncodeBoolean(&b)
	if err != nil {
		t.Fatalf("Error during encode: %v", b)
	}
}

func testBoolean(t *testing.T, decoder Decoder, ref Boolean) {
	b, err := decoder.DecodeBoolean()
	if err != nil {
		t.Fatalf("Error during decode: %v", ref)
	}
	if *b != ref {
		t.Fatalf("Bad decoded value: %v", *b)
	}
}

func writeUOctet(t *testing.T, encoder Encoder, o UOctet) {
	err := encoder.EncodeUOctet(&o)
	if err != nil {
		t.Fatalf("Error during encode: %d", o)
	}
}

func testUOctet(t *testing.T, decoder Decoder, ref UOctet) {
	o, err := decoder.DecodeUOctet()
	if err != nil {
		t.Fatalf("Error during decode: %d", ref)
	}
	if *o != ref {
		t.Fatalf("Bad decoded value: %d", o)
	}
}

func writeOctet(t *testing.T, encoder Encoder, o Octet) {
	err := encoder.EncodeOctet(&o)
	if err != nil {
		t.Fatalf("Error during encode: %d", o)
	}
}

func testOctet(t *testing.T, decoder Decoder, ref Octet) {
	o, err := decoder.DecodeOctet()
	if err != nil {
		t.Fatalf("Error during decode: %d", ref)
	}
	if *o != ref {
		t.Fatalf("Bad decoded value: %d", o)
	}
}

func writeUShort(t *testing.T, encoder Encoder, o UShort) {
	err := encoder.EncodeUShort(&o)
	if err != nil {
		t.Fatalf("Error during encode: %d", o)
	}
}

func testUShort(t *testing.T, decoder Decoder, ref UShort) {
	o, err := decoder.DecodeUShort()
	if err != nil {
		t.Fatalf("Error during decode: %d", ref)
	}
	if *o != ref {
		t.Fatalf("Bad decoded value: %d", o)
	}
}

func writeShort(t *testing.T, encoder Encoder, o Short) {
	err := encoder.EncodeShort(&o)
	if err != nil {
		t.Fatalf("Error during encode: %d", o)
	}
}

func testShort(t *testing.T, decoder Decoder, ref Short) {
	o, err := decoder.DecodeShort()
	if err != nil {
		t.Fatalf("Error during decode: %d", ref)
	}
	if *o != ref {
		t.Fatalf("Bad decoded value: %d", o)
	}
}

func writeUInteger(t *testing.T, encoder Encoder, o UInteger) {
	err := encoder.EncodeUInteger(&o)
	if err != nil {
		t.Fatalf("Error during encode: %d", o)
	}
}

func testUInteger(t *testing.T, decoder Decoder, ref UInteger) {
	o, err := decoder.DecodeUInteger()
	if err != nil {
		t.Fatalf("Error during decode: %d", ref)
	}
	if *o != ref {
		t.Fatalf("Bad decoded value: %d", o)
	}
}

func writeInteger(t *testing.T, encoder Encoder, o Integer) {
	err := encoder.EncodeInteger(&o)
	if err != nil {
		t.Fatalf("Error during encode: %d", o)
	}
}

func testInteger(t *testing.T, decoder Decoder, ref Integer) {
	o, err := decoder.DecodeInteger()
	if err != nil {
		t.Fatalf("Error during decode: %d", ref)
	}
	if *o != ref {
		t.Fatalf("Bad decoded value: %d", o)
	}
}

func writeULong(t *testing.T, encoder Encoder, o ULong) {
	err := encoder.EncodeULong(&o)
	if err != nil {
		t.Fatalf("Error during encode: %d", o)
	}
}

func testULong(t *testing.T, decoder Decoder, ref ULong) {
	o, err := decoder.DecodeULong()
	if err != nil {
		t.Fatalf("Error during decode: %d", ref)
	}
	if *o != ref {
		t.Fatalf("Bad decoded value: %d", o)
	}
}

func writeLong(t *testing.T, encoder Encoder, o Long) {
	err := encoder.EncodeLong(&o)
	if err != nil {
		t.Fatalf("Error during encode: %d", o)
	}
}

func testLong(t *testing.T, decoder Decoder, ref Long) {
	o, err := decoder.DecodeLong()
	if err != nil {
		t.Fatalf("Error during decode: %d", ref)
	}
	if *o != ref {
		t.Fatalf("Bad decoded value: %d", o)
	}
}

func writeFloat(t *testing.T, encoder Encoder, o Float) {
	err := encoder.EncodeFloat(&o)
	if err != nil {
		t.Fatalf("Error during encode: %f", o)
	}
}

func testFloat(t *testing.T, decoder Decoder, ref Float) {
	o, err := decoder.DecodeFloat()
	if err != nil {
		t.Fatalf("Error during decode: %f", ref)
	}
	if *o != ref {
		t.Errorf("Bad decoded value: %f != %f", *o, ref)
	}
}

func writeDouble(t *testing.T, encoder Encoder, o Double) {
	err := encoder.EncodeDouble(&o)
	if err != nil {
		t.Fatalf("Error during encode: %f", o)
	}
}

func testDouble(t *testing.T, decoder Decoder, ref Double) {
	o, err := decoder.DecodeDouble()
	if err != nil {
		t.Fatalf("Error during decode: %f", ref)
	}
	if *o != ref {
		t.Errorf("Bad decoded value: %f != %f", *o, ref)
	}
}

func writeBlob(t *testing.T, encoder Encoder, o Blob) {
	err := encoder.EncodeBlob(&o)
	if err != nil {
		t.Fatalf("Error during encode: %v", o)
	}
}

func testBlob(t *testing.T, decoder Decoder, ref Blob) {
	o, err := decoder.DecodeBlob()
	if err != nil {
		t.Fatalf("Error during decode: %v", ref)
	}
	if !reflect.DeepEqual(o, &ref) {
		t.Errorf("Bad decoding, got: %v, want: %v", o, ref)
	}
}

func writeString(t *testing.T, encoder Encoder, o String) {
	err := encoder.EncodeString(&o)
	if err != nil {
		t.Fatalf("Error during encode: %s", o)
	}
}

func testString(t *testing.T, decoder Decoder, ref String) {
	o, err := decoder.DecodeString()
	if err != nil {
		t.Fatalf("Error during decode: %s", ref)
	}
	if *o != ref {
		t.Errorf("Bad decoding, got: %s, want: %s", *o, ref)
	}
}

func writeTime(t *testing.T, encoder Encoder, o Time) {
	err := encoder.EncodeTime(&o)
	if err != nil {
		t.Fatalf("Error during encode: %v", o)
	}
}

func testTime(t *testing.T, decoder Decoder, ref Time) {
	o, err := decoder.DecodeTime()
	if err != nil {
		t.Fatalf("Error during decode: %v", ref)
	}
	if (time.Time(*o).UnixNano() / 1000000) != (time.Time(ref).UnixNano() / 1000000) {
		t.Errorf("Bad decoding, got: %v, want: %v", *o, ref)
	}
}

func writeFineTime(t *testing.T, encoder Encoder, o FineTime) {
	err := encoder.EncodeFineTime(&o)
	if err != nil {
		t.Fatalf("Error during encode: %v", o)
	}
}

func testFineTime(t *testing.T, decoder Decoder, ref FineTime) {
	o, err := decoder.DecodeFineTime()
	if err != nil {
		t.Fatalf("Error during decode: %v", ref)
	}
	if (time.Time(*o).UnixNano() / 1000000) != (time.Time(ref).UnixNano() / 1000000) {
		t.Errorf("Bad decoding, got: %v, want: %v", *o, ref)
	}
}

func deepCompare(t *testing.T, file1, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		t.Errorf("Expected file does not exist: %s", file1)
		return false
	}

	f2, err := os.Open(file2)
	if err != nil {
		t.Errorf("Actual file does not exist: %s", file2)
		return false
	}

	for {
		b1 := make([]byte, 8192)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, 8192)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				t.Errorf("Files size differ")
				return false
			} else {
				t.Errorf("Error reading files")
				return false
			}
		}

		if !bytes.Equal(b1, b2) {
			t.Errorf("Files differ")
			return false
		}
	}
}
