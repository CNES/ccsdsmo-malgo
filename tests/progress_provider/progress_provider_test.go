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
package progress_provider

import (
	"fmt"
	. "github.com/CNES/ccsdsmo-malgo/mal"
	. "github.com/CNES/ccsdsmo-malgo/mal/api"
	_ "github.com/CNES/ccsdsmo-malgo/mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

type TestProgressProvider struct {
	ctx   *Context
	cctx  *ClientContext
	uri   *URI
	nbmsg int
}

// Test TCP transport Progress Interaction using the high level API, this test can be
// used to test interoperability with MAL/C implementation.
//
// docker run -v F:\Users\freyssin\Workspaces\CNES-C\malc2:/opt/malc -it malc-opensuse-shared /bin/bash
// cd malc
// sh ./bin/update_known_projects.sh /usr/local/bin/zproject_known_projects.xml
// cd test/progress_consumer/
// ./genmake check
//
// launches the MAL/C progress_provider the the MAL/GO progress_consumer
func newTestProgressProvider() (*TestProgressProvider, error) {
	ctx, err := NewContext("maltcp://127.0.0.1:6666")
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, "progress_provider/provider")
	if err != nil {
		return nil, err
	}

	provider := &TestProgressProvider{ctx, cctx, cctx.Uri, 0}

	// Handler1
	progressHandler := func(msg *Message, t Transaction) error {
		provider.nbmsg += 1
		if msg != nil {
			transaction := t.(ProgressTransaction)
			par, _ := msg.DecodeLastParameter(NullStringList, false)
			fmt.Printf("\t$$$$$ progressHandler1 receive: %v\n", par)
			transaction.Ack(nil, false)
			for i := 0; i < 10; i++ {
				body := transaction.NewBody()
				body.EncodeLastParameter(NewString(fmt.Sprintf("messsage1.#%d", i)), false)
				transaction.Update(body, false)
				time.Sleep(1 * time.Second)
			}
			body := transaction.NewBody()
			body.EncodeLastParameter(NewString("last message1"), false)
			transaction.Reply(body, false)
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}
	// Registers Progress handler
	cctx.RegisterProgressHandler(200, 1, 1, 5, progressHandler)

	return provider, nil
}

func (provider *TestProgressProvider) close() {
	provider.ctx.Close()
}

// Test TCP transport Progress Interaction using the high level API
func TestProgress(t *testing.T) {
	provider, err := newTestProgressProvider()
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}
	defer provider.close()

	time.Sleep(120 * time.Second)
}
