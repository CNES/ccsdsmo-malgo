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
package api1_test

import (
	"fmt"
	. "mal"
	. "mal/api1"
	_ "mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

const (
	maltcp1 = "maltcp://127.0.0.1:16001"
	maltcp2 = "maltcp://127.0.0.1:16002"
)

// ########## ########## ########## ########## ########## ########## ########## ##########

// Test TCP transport Send Interaction using the high level API
func TestSend(t *testing.T) {
	ctx1, err := NewContext(maltcp1)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	provider, err := NewProviderContext(ctx1, "provider")
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}

	ctx2, err := NewContext(maltcp2)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	consumer, err := NewOperationContext(ctx2, "consumer")
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
		return
	}

	nbmsg := 0

	// Declares handler function here so nbmsg local variable is accessible throught
	// the closure.
	sendHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			fmt.Println("\n\n $$$$$ sendHandler receive: ", string(msg.Body), "\n\n")
			nbmsg += 1
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}

	provider.RegisterSendHandler(200, 1, 1, 1, sendHandler)

	op1, err := consumer.NewSendOperation(200, 1, 1, 1)
	msg1 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message1"),
	}
	op1.Send(msg1)

	op2, err := consumer.NewSendOperation(200, 1, 1, 1)
	msg2 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message2"),
	}
	op2.Send(msg2)

	time.Sleep(250 * time.Millisecond)
	ctx1.Close()
	ctx2.Close()

	if nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	}

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}

// ########## ########## ########## ########## ########## ########## ########## ##########

// Test TCP transport Submit Interaction using the high level API
func TestSubmit(t *testing.T) {
	ctx1, err := NewContext(maltcp1)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	provider, err := NewProviderContext(ctx1, "provider")
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}

	ctx2, err := NewContext(maltcp2)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	consumer, err := NewOperationContext(ctx2, "consumer")
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
		return
	}

	nbmsg := 0

	// Declares handler function here so nbmsg local variable is accessible throught
	// the closure.

	submitHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			fmt.Println("\n\n $$$$$ submitHandler receive: ", string(msg.Body), "\n\n")
			nbmsg += 1
			transaction := t.(SubmitTransaction)
			transaction.Ack(nil)
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}

	provider.RegisterSubmitHandler(200, 1, 1, 1, submitHandler)

	op1, err := consumer.NewSubmitOperation(200, 1, 1, 1)
	msg1 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message1"),
	}
	op1.Submit(msg1)

	fmt.Println("\n\n &&&&& Submit1: OK\n\n")

	op2, err := consumer.NewSubmitOperation(200, 1, 1, 1)
	msg2 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message2"),
	}
	op2.Submit(msg2)

	fmt.Println("\n\n &&&&& Submit2: OK\n\n")

	time.Sleep(250 * time.Millisecond)
	ctx1.Close()
	ctx2.Close()

	if nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	}

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}

// ########## ########## ########## ########## ########## ########## ########## ##########

// Test TCP transport Progress Interaction using the high level API
func TestProgress(t *testing.T) {
	ctx1, err := NewContext(maltcp1)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	provider, err := NewProviderContext(ctx1, "provider")
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}

	ctx2, err := NewContext(maltcp2)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	consumer, err := NewOperationContext(ctx2, "consumer")
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
		return
	}

	nbmsg := 0

	// Declares handler function here so nbmsg local variable is accessible throught
	// the closure.

	progressHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			fmt.Println("\n\n $$$$$ progressHandler receive: ", string(msg.Body), "\n\n")
			transaction := t.(ProgressTransaction)
			transaction.Ack(nil)
			for i := 0; i < 10; i++ {
				transaction.Update([]byte(fmt.Sprintf("messsage#%d", i)), nil)
			}
			transaction.Reply([]byte("last message"), nil)
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}

	provider.RegisterProgressHandler(200, 1, 1, 1, progressHandler)

	op1, err := consumer.NewProgressOperation(200, 1, 1, 1)
	msg1 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message1"),
	}
	op1.Progress(msg1)
	fmt.Println("\n\n &&&&& Progress1: OK\n\n")

	updt, err := op1.GetUpdate()
	if err != nil {
		t.Error(err)
	}
	for updt != nil {
		nbmsg += 1
		fmt.Println("\n\n &&&&& Progress1: Update -> ", string(updt.Body), "\n\n")
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
	fmt.Println("\n\n &&&&& Progress1: Response -> ", string(rep.Body), "\n\n")

	time.Sleep(250 * time.Millisecond)
	ctx1.Close()
	ctx2.Close()

	if nbmsg != 11 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	}

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}
