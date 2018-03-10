package malactor_test

import (
	"fmt"
	. "github.com/ccsdsmo/malgo/mal"
	_ "github.com/ccsdsmo/malgo/mal/transport/invm"
	_ "github.com/ccsdsmo/malgo/mal/transport/tcp"
	. "malactor"
	"testing"
	"time"
)

var (
	invm string = "invm://local"
)

func TestLocal(t *testing.T) {
	ctx, err := NewContext(invm)
	if err != nil {
		t.Fatal("Error creating context, ", err)
	}

	consumer, err := NewEndPoint(ctx, "consumer", nil)
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
	}
	fmt.Println("Registered consumer: ", consumer.Uri)

	routing := new(Routing)
	//	provider := new(MyProvider)
	//	routing.registerProviderProgress(..., provider)
	actor, err := NewActor(ctx, "provider", routing)
	fmt.Println("Registered provider: ", actor)
	actor.Start()

	endpoint, err := ctx.GetEndPoint(ctx.NewURI("consumer"))
	fmt.Println("consumer: ", endpoint, err)

	time.Sleep(1 * time.Second)
	msg1 := &Message{
		UriFrom: consumer.Uri,
		UriTo:   actor.Uri(),
		Body:    []byte("message1"),
	}
	consumer.Send(msg1)
	fmt.Println("Message sent: ", string(msg1.Body))

	time.Sleep(250 * time.Millisecond)
	msg2 := &Message{
		UriFrom: consumer.Uri,
		UriTo:   actor.Uri(),
		Body:    []byte("message2"),
	}
	consumer.Send(msg2)
	fmt.Println("Message sent: ", string(msg2.Body))

	time.Sleep(250 * time.Millisecond)
	actor.Stop()

	time.Sleep(250 * time.Millisecond)
	ctx.Close()

	time.Sleep(250 * time.Millisecond)
}
