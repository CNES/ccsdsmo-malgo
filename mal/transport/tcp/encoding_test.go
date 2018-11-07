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
package tcp

import (
	. "github.com/ccsdsmo/malgo/mal"
	"testing"
	"time"
)

func TestMessage(t *testing.T) {
	from := URI("maltcp://192.168.1.80:12345/Service1")
	to := URI("maltcp://192.168.1.81:54321/Service2")
	msg1 := &Message{
		UriFrom:          &from,
		UriTo:            &to,
		Timestamp:        *TimeNow(),
		Body:             []byte("message1"),
		InteractionType:  MAL_INTERACTIONTYPE_PROGRESS,
		InteractionStage: MAL_IP_STAGE_PROGRESS_UPDATE,
		Domain: IdentifierList([]*Identifier{
			NewIdentifier("DOMAIN1"),
			NewIdentifier("DOMAIN2"),
			nil,
			NewIdentifier("Domain4"),
		}),
	}

	transport1 := &TCPTransport{
		uri:     URI("maltcp://192.168.1.80:12345"),
		version: 1,

		sourceFlag:           true,
		destinationFlag:      true,
		priorityFlag:         true,
		timestampFlag:        true,
		networkZoneFlag:      true,
		sessionNameFlag:      true,
		domainFlag:           true,
		authenticationIdFlag: true,

		flags: 0xFF,
	}

	buf, err := transport1.encode(msg1)
	if err != nil {
		t.Fatalf("Error during encode: %s", err)
	}

	transport2 := &TCPTransport{
		uri:     URI("maltcp://192.168.1.81:54321"),
		version: 1,

		sourceFlag:           true,
		destinationFlag:      true,
		priorityFlag:         true,
		timestampFlag:        true,
		networkZoneFlag:      true,
		sessionNameFlag:      true,
		domainFlag:           true,
		authenticationIdFlag: true,

		flags: 0xFF,
	}

	msg2, err := transport2.decode(buf, "192.168.1.80:12345")
	if err != nil {
		t.Fatalf("Error during encode: %s", err)
	}
	// NOTE (AF): Be careful the DeepEqual method is costly :-(
	if !messageEqual(t, msg1, msg2) {
		t.Errorf("Bad decoding, got: %v, want: %v", msg1, msg2)
	}
}

func messageEqual(t *testing.T, msg1 *Message, msg2 *Message) bool {
	if *msg1.UriFrom != *msg2.UriFrom {
		t.Logf("UriFrom different, got %s, expect %s\n", *msg2.UriFrom, *msg1.UriFrom)
		return false
	}
	if *msg1.UriTo != *msg2.UriTo {
		t.Logf("UriTo different, got %s, expect %s\n", *msg2.UriTo, *msg1.UriTo)
		return false
	}
	if (time.Time(msg1.Timestamp).UnixNano() / 1000000) != (time.Time(msg2.Timestamp).UnixNano() / 1000000) {
		t.Logf("Timestamp different, got %v, expect %v\n", msg2.Timestamp, msg1.Timestamp)
		return false
	}
	if msg1.InteractionType != msg2.InteractionType {
		t.Logf("InteractionType different, got %v, expect %v\n", msg2.InteractionType, msg1.InteractionType)
		return false
	}
	if msg1.InteractionStage != msg2.InteractionStage {
		t.Logf("InteractionStage different, got %v, expect %v\n", msg2.InteractionStage, msg1.InteractionStage)
		return false
	}
	return true
}
