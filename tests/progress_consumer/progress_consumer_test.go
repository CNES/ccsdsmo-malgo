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
package progress_consumer

import (
	"fmt"
	. "github.com/CNES/ccsdsmo-malgo/mal"
	. "github.com/CNES/ccsdsmo-malgo/mal/api"
	"github.com/CNES/ccsdsmo-malgo/mal/encoding/binary"
	_ "github.com/CNES/ccsdsmo-malgo/mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
)

const (
	provider_url = "maltcp://127.0.0.1:6660/progress_provider/provider"
	consumer_url = "maltcp://127.0.0.1:16002"
)

// Test TCP transport Progress Interaction using the high level API, this test can be
// used to test interoperability with MAL/C implementation.
//
// docker run -v F:\Users\freyssin\Workspaces\CNES-C\malc2:/opt/malc -p 6660:6666 -it malc-opensuse-shared /bin/bash
// cd malc
// sh ./bin/update_known_projects.sh /usr/local/bin/zproject_known_projects.xml
// cd test/progress_provider/
// ./genmake check
//
// launches the MAL/C progress_provider the the MAL/GO progress_consumer
func TestProgressConsumer(t *testing.T) {
	consumer_ctx, err := NewContext(consumer_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	defer consumer_ctx.Close()

	consumer, err := NewClientContext(consumer_ctx, "consumer")
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
		return
	}

	nbmsg := 0

	var providerUri *URI = NewURI(provider_url)
	op1 := consumer.NewProgressOperation(providerUri, 200, 1, 1, 5)

	list := NewStringList(5)
	([]*String)(*list)[0] = NewString("list-element-0")
	([]*String)(*list)[1] = NewString("list-element-1")
	([]*String)(*list)[2] = NewString("list-element-2")
	([]*String)(*list)[3] = NewString("list-element-3")
	([]*String)(*list)[4] = NewString("list-element-4")
	buf := make([]byte, 0, 8192)
	encoder := binary.NewBinaryEncoder(buf, false)
	encoder.EncodeNullableElement(list)

	op1.Progress(encoder.Body())
	fmt.Println("\t&&&&& Progress1: OK")

	updt, err := op1.GetUpdate()
	if err != nil {
		t.Error(err)
	}
	for updt != nil {
		nbmsg += 1
		fmt.Println("\t&&&&& Progress1: Update -> ", updt.Body)
		updt, err = op1.GetUpdate()
		if err != nil {
			t.Error(err)
		}
	}
	rep, err := op1.GetResponse()
	if err != nil {
		t.Error(err)
	}
	nbmsg += 1

	decoder := binary.NewBinaryDecoder(rep.Body, false)
	decoder.DecodeNullableElement(list)
	fmt.Println("\t&&&&& Progress1: Response -> ", *([]*String)(*list)[0])

	if nbmsg != 11 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	}
}
