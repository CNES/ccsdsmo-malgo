#Tests MAL/GO

This document includes the execution trace of all the tests defined in the malgo project.
It reflects an execution on the 3rd release of malgo (tag MALGO_0_3).
In order to force the exll tests are done with the GOCACHE environment variable set of off.

## Mode normal

```
[malgo@localhost src]$ GOCACHE=off go test github.com/CNES/ccsdsmo-malgo/mal/encoding/binary
ok      github.com/CNES/ccsdsmo-malgo/mal/encoding/binary       2.617s
[malgo@localhost src]$ GOCACHE=off go test github.com/CNES/ccsdsmo-malgo/mal/encoding/splitbinary
ok      github.com/CNES/ccsdsmo-malgo/mal/encoding/splitbinary  2.665s
[malgo@localhost src]$ GOCACHE=off go test github.com/CNES/ccsdsmo-malgo/mal/transport/invm
ok      github.com/CNES/ccsdsmo-malgo/mal/transport/invm        2.507s
[malgo@localhost src]$ GOCACHE=off go test github.com/CNES/ccsdsmo-malgo/mal/transport/tcp
ok      github.com/CNES/ccsdsmo-malgo/mal/transport/tcp 3.277s
[malgo@localhost src]$ GOCACHE=off go test github.com/CNES/ccsdsmo-malgo/mal/api
ok      github.com/CNES/ccsdsmo-malgo/mal/api   5.329s
[malgo@localhost src]$ GOCACHE=off go test github.com/CNES/ccsdsmo-malgo/tests/encoding
ok      github.com/CNES/ccsdsmo-malgo/tests/encoding    0.004s
[malgo@localhost src]$ GOCACHE=off go test github.com/CNES/ccsdsmo-malgo/tests/issue1
ok      github.com/CNES/ccsdsmo-malgo/tests/issue1      0.508s
```

## Mode verbose

### Encoding binary
Tests specific to the binary encoding.

```
[malgo@localhost src]$ GOCACHE=off go test -v github.com/CNES/ccsdsmo-malgo/mal/encoding/binary
=== RUN   TestByte
--- PASS: TestByte (0.00s)
=== RUN   TestUInt32
--- PASS: TestUInt32 (0.01s)
=== RUN   TestUOctet
--- PASS: TestUOctet (0.00s)
=== RUN   TestOctet
--- PASS: TestOctet (0.00s)
=== RUN   TestUShort
--- PASS: TestUShort (0.01s)
=== RUN   TestShort
--- PASS: TestShort (0.01s)
=== RUN   TestUInteger
--- PASS: TestUInteger (0.02s)
=== RUN   TestNullableUInteger
--- PASS: TestNullableUInteger (0.00s)
=== RUN   TestInteger
--- PASS: TestInteger (0.01s)
=== RUN   TestULong
--- PASS: TestULong (0.02s)
=== RUN   TestLong
--- PASS: TestLong (0.02s)
=== RUN   TestFloat
--- PASS: TestFloat (0.02s)
=== RUN   TestDouble
--- PASS: TestDouble (0.02s)
=== RUN   TestBlob
--- PASS: TestBlob (1.73s)
=== RUN   TestString
--- PASS: TestString (0.62s)
=== RUN   TestNullableString
--- PASS: TestNullableString (0.00s)
=== RUN   TestTime
--- PASS: TestTime (0.04s)
=== RUN   TestFineTime
--- PASS: TestFineTime (0.04s)
=== RUN   TestAttribute
--- PASS: TestAttribute (0.00s)
=== RUN   TestNullableAttribute
--- PASS: TestNullableAttribute (0.00s)
=== RUN   TestElement
--- PASS: TestElement (0.00s)
=== RUN   TestNullableElement
--- PASS: TestNullableElement (0.00s)
    binary_test.go:829: Encode:  0xc000073215
    binary_test.go:829: Encode:  <nil>
    binary_test.go:829: Encode:  0xc000073216
    binary_test.go:829: Encode:  0xc000073218
    binary_test.go:829: Encode:  <nil>
    binary_test.go:829: Encode:  <nil>
    binary_test.go:829: Encode:  0xc000073230
    binary_test.go:829: Encode:  0xc000073238
    binary_test.go:829: Encode:  0xc000073240
    binary_test.go:829: Encode:  <nil>
    binary_test.go:829: Encode:  Identifier(BOOM)
    binary_test.go:829: Encode:  0xc000073248
    binary_test.go:829: Encode:  0xc000073249
    binary_test.go:829: Encode:  0xc00007324a
    binary_test.go:829: Encode:  0xc00007324c
    binary_test.go:829: Encode:  0xc000073250
    binary_test.go:829: Encode:  <nil>
    binary_test.go:829: Encode:  <nil>
    binary_test.go:829: Encode:  <nil>
    binary_test.go:829: Encode:  0xc000073254
    binary_test.go:829: Encode:  0xc000073258
    binary_test.go:829: Encode:  0xc000073260
    binary_test.go:829: Encode:  0xc004b2dcc0
    binary_test.go:829: Encode:  0xc000073268
    binary_test.go:829: Encode:  0xc000073270
    binary_test.go:829: Encode:  0xc004b2dcd0
    binary_test.go:829: Encode:  0xc000073271
    binary_test.go:829: Encode:  <nil>
    binary_test.go:829: Encode:  <nil>
    binary_test.go:829: Encode:  <nil>
=== RUN   TestBlobAttribute
--- PASS: TestBlobAttribute (0.00s)
=== RUN   TestTimeAttribute
--- PASS: TestTimeAttribute (0.00s)
=== RUN   TestBooleanList
--- PASS: TestBooleanList (0.00s)
=== RUN   TestOctetList
--- PASS: TestOctetList (0.00s)
=== RUN   TestUOctetList
--- PASS: TestUOctetList (0.00s)
=== RUN   TestShortList
--- PASS: TestShortList (0.00s)
=== RUN   TestUShortList
--- PASS: TestUShortList (0.00s)
=== RUN   TestIntegerList
--- PASS: TestIntegerList (0.00s)
=== RUN   TestUIntegerList
--- PASS: TestUIntegerList (0.00s)
=== RUN   TestLongList
--- PASS: TestLongList (0.00s)
=== RUN   TestULongList
--- PASS: TestULongList (0.00s)
=== RUN   TestFloatList
--- PASS: TestFloatList (0.00s)
=== RUN   TestTimeList
--- PASS: TestTimeList (0.00s)
=== RUN   TestFineTimeList
--- PASS: TestFineTimeList (0.00s)
=== RUN   TestIdentifierList
--- PASS: TestIdentifierList (0.00s)
=== RUN   TestEntityKey
--- PASS: TestEntityKey (0.00s)
=== RUN   TestEntityKeyList
--- PASS: TestEntityKeyList (0.00s)
=== RUN   TestEntityRequest
--- PASS: TestEntityRequest (0.00s)
=== RUN   TestEntityRequestList
--- PASS: TestEntityRequestList (0.00s)
=== RUN   TestAbstractElement
--- PASS: TestAbstractElement (0.00s)
=== RUN   TestNullableAbstractElement
--- PASS: TestNullableAbstractElement (0.00s)
=== RUN   TestBrokerEncoding1
--- PASS: TestBrokerEncoding1 (0.00s)
=== RUN   TestBrokerEncoding2
--- PASS: TestBrokerEncoding2 (0.00s)
PASS
ok      github.com/CNES/ccsdsmo-malgo/mal/encoding/binary       2.600s
```

### Encoding splitbinary
Tests specific to the splitbinary encoding.

```
[malgo@localhost src]$ GOCACHE=off go test -v github.com/CNES/ccsdsmo-malgo/mal/encoding/splitbinary
=== RUN   Test1
--- PASS: Test1 (0.00s)
    splitbinary_test.go:108: SplitBinaryBuffer=[30 30 59] 0 [26 0] 5 9
    splitbinary_test.go:111: SplitBinaryBuffer=[1 26 30 30 59]
    splitbinary_test.go:118: decode:  15
    splitbinary_test.go:124: decode:  false
    splitbinary_test.go:130: decode:  true
    splitbinary_test.go:136: decode:  false
    splitbinary_test.go:142: decode:  30
=== RUN   TestUOctet
--- PASS: TestUOctet (0.00s)
=== RUN   TestOctet
--- PASS: TestOctet (0.00s)
=== RUN   TestUShort
--- PASS: TestUShort (0.01s)
=== RUN   TestShort
--- PASS: TestShort (0.01s)
=== RUN   TestUInteger
--- PASS: TestUInteger (0.02s)
=== RUN   TestNullableUInteger
--- PASS: TestNullableUInteger (0.00s)
=== RUN   TestInteger
--- PASS: TestInteger (0.02s)
=== RUN   TestULong
--- PASS: TestULong (0.03s)
=== RUN   TestLong
--- PASS: TestLong (0.03s)
=== RUN   TestFloat
--- PASS: TestFloat (0.02s)
=== RUN   TestDouble
--- PASS: TestDouble (0.03s)
=== RUN   TestBlob
--- PASS: TestBlob (1.73s)
=== RUN   TestString
--- PASS: TestString (0.66s)
=== RUN   TestNullableString
--- PASS: TestNullableString (0.00s)
=== RUN   TestTime
--- PASS: TestTime (0.04s)
=== RUN   TestFineTime
--- PASS: TestFineTime (0.04s)
=== RUN   TestAttribute
--- PASS: TestAttribute (0.00s)
=== RUN   TestNullableAttribute
--- PASS: TestNullableAttribute (0.00s)
=== RUN   TestElement
--- PASS: TestElement (0.00s)
=== RUN   TestNullableElement
--- PASS: TestNullableElement (0.00s)
    splitbinary_test.go:910: Encode:  0xc000de9958
    splitbinary_test.go:910: Encode:  <nil>
    splitbinary_test.go:910: Encode:  0xc000de9959
    splitbinary_test.go:910: Encode:  0xc000de9960
    splitbinary_test.go:910: Encode:  <nil>
    splitbinary_test.go:910: Encode:  <nil>
    splitbinary_test.go:910: Encode:  0xc000de9968
    splitbinary_test.go:910: Encode:  0xc000de9970
    splitbinary_test.go:910: Encode:  0xc000de9978
    splitbinary_test.go:910: Encode:  <nil>
    splitbinary_test.go:910: Encode:  Identifier(BOOM)
    splitbinary_test.go:910: Encode:  0xc000de9980
    splitbinary_test.go:910: Encode:  0xc000de9981
    splitbinary_test.go:910: Encode:  0xc000de9982
    splitbinary_test.go:910: Encode:  0xc000de9984
    splitbinary_test.go:910: Encode:  0xc000de9988
    splitbinary_test.go:910: Encode:  <nil>
    splitbinary_test.go:910: Encode:  <nil>
    splitbinary_test.go:910: Encode:  <nil>
    splitbinary_test.go:910: Encode:  0xc000de998c
    splitbinary_test.go:910: Encode:  0xc000de9990
    splitbinary_test.go:910: Encode:  0xc000de9998
    splitbinary_test.go:910: Encode:  0xc004532850
    splitbinary_test.go:910: Encode:  0xc000de99a0
    splitbinary_test.go:910: Encode:  0xc000de99a8
    splitbinary_test.go:910: Encode:  0xc004532860
    splitbinary_test.go:910: Encode:  0xc000de99a9
    splitbinary_test.go:910: Encode:  <nil>
    splitbinary_test.go:910: Encode:  <nil>
    splitbinary_test.go:910: Encode:  <nil>
=== RUN   TestBlobAttribute
--- PASS: TestBlobAttribute (0.00s)
=== RUN   TestTimeAttribute
--- PASS: TestTimeAttribute (0.00s)
=== RUN   TestBooleanList
--- PASS: TestBooleanList (0.00s)
=== RUN   TestOctetList
--- PASS: TestOctetList (0.00s)
=== RUN   TestUOctetList
--- PASS: TestUOctetList (0.00s)
=== RUN   TestShortList
--- PASS: TestShortList (0.00s)
=== RUN   TestUShortList
--- PASS: TestUShortList (0.00s)
=== RUN   TestIntegerList
--- PASS: TestIntegerList (0.00s)
=== RUN   TestUIntegerList
--- PASS: TestUIntegerList (0.00s)
=== RUN   TestLongList
--- PASS: TestLongList (0.00s)
=== RUN   TestULongList
--- PASS: TestULongList (0.00s)
=== RUN   TestFloatList
--- PASS: TestFloatList (0.00s)
=== RUN   TestTimeList
--- PASS: TestTimeList (0.00s)
=== RUN   TestFineTimeList
--- PASS: TestFineTimeList (0.00s)
=== RUN   TestIdentifierList
--- PASS: TestIdentifierList (0.00s)
=== RUN   TestEntityKey
--- PASS: TestEntityKey (0.00s)
=== RUN   TestEntityKeyList
--- PASS: TestEntityKeyList (0.00s)
=== RUN   TestEntityRequest
--- PASS: TestEntityRequest (0.00s)
=== RUN   TestEntityRequestList
--- PASS: TestEntityRequestList (0.00s)
PASS
ok      github.com/CNES/ccsdsmo-malgo/mal/encoding/splitbinary  2.682s
```

### Transport InVM
Tests specific to the InVM specific transport.

```
[malgo@localhost src]$ GOCACHE=off go test -v github.com/CNES/ccsdsmo-malgo/mal/transport/invm
=== RUN   TestLocal1
Registered consumer:  0xc0000667d0
Registered provider:  0xc0000667e0
consumer:  &{0xc0000ba050 0xc0000667d0 0xc0000761e0 0} <nil>
receive:  message1 ,  <nil>
receive:  message2 ,  <nil>
--- PASS: TestLocal1 (1.26s)
=== RUN   TestLocal2
Registered consumer:  invm://local1/consumer
Registered provider:  invm://local2/provider
consumer:  &{0xc0000ba230 0xc0000668f0 0xc000076300 0} <nil>
receive:  message1 ,  <nil>
receive:  message2 ,  <nil>
--- PASS: TestLocal2 (1.25s)
PASS
ok      github.com/CNES/ccsdsmo-malgo/mal/transport/invm        2.519s
```

### Transport TCP
Tests specific to the TCP standardized transport.

```
[malgo@localhost src]$ GOCACHE=off go test -v  github.com/CNES/ccsdsmo-malgo/mal/transport/tcp
=== RUN   TestMessage
[1 0 0 0 8 109 101 115 115 97 103 101 49]
message1
--- PASS: TestMessage (0.00s)
=== RUN   TestTCP1
Registered consumer:  0xc00007ebc0
Registered provider:  0xc00007ebd0
consumer:  &{0xc0000dc0f0 0xc00007ebc0 0xc00008e2a0 0} <nil>
receive:  message1 ,  <nil>
receive:  message2 ,  <nil>
2019-09-25 09:00:57 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16000), connection closed
2019-09-25 09:00:57 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:45254), connection closed
--- PASS: TestTCP1 (2.01s)
=== RUN   TestTCP2
Registered consumer:  maltcp://127.0.0.1:16001/consumer
Registered provider:  maltcp://127.0.0.1:16002/provider
consumer:  &{0xc0000ea0a0 0xc0000281a0 0xc000064180 0} <nil>
receive:  message1 ,  <nil>
receive:  message2 ,  <nil>
2019-09-25 09:00:58 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16002), connection closed
2019-09-25 09:00:58 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:47592), connection closed
--- PASS: TestTCP2 (1.25s)
PASS
ok      github.com/CNES/ccsdsmo-malgo/mal/transport/tcp 3.271s
```

### API
Tests related to the MAL/GO high-level API (SEND, SUBMIT, REQUEST, INVOKE, PROGRESS and PUBSUB MAL interaction). 

```
[malgo@localhost src]$ GOCACHE=off go test -v github.com/CNES/ccsdsmo-malgo/mal/api
=== RUN   TestSend
        $$$$$ sendHandler receive:  message1 <nil>
        $$$$$ sendHandler receive:  message2 <nil>
2019-09-25 09:01:07 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16001), connection closed
2019-09-25 09:01:07 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:49710), connection closed
--- PASS: TestSend (0.52s)
=== RUN   TestSubmit
        $$$$$ submitHandler receive:  message1 <nil>
        &&&&& Submit1: OK
        $$$$$ submitHandler receive:  message2 <nil>
        &&&&& Submit2: OK
2019-09-25 09:01:07 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16001), connection closed
2019-09-25 09:01:07 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:47598), connection closed
2019-09-25 09:01:07 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:49712), connection closed
2019-09-25 09:01:07 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16002), connection closed
--- PASS: TestSubmit (0.26s)
=== RUN   TestRequest
        $$$$$ requestHandler receive:  message1 <nil>
        &&&&& Request1: OK,  reply message
        $$$$$ requestHandler receive:  message2 <nil>
        &&&&& Request2: OK,  reply message
2019-09-25 09:01:07 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16001), connection closed
2019-09-25 09:01:07 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:47602), connection closed
2019-09-25 09:01:07 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:49716), connection closed
2019-09-25 09:01:07 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16002), connection closed
--- PASS: TestRequest (0.26s)
=== RUN   TestInvoke
        $$$$$ invokeProvider receive:  message1 <nil>
        &&&&& Invoke1: OK,  message1
        $$$$$ invokeProvider receive:  message2 <nil>
        &&&&& Invoke2: OK,  message2
2019-09-25 09:01:08 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16001), connection closed
2019-09-25 09:01:08 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:47606), connection closed
2019-09-25 09:01:08 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:49720), connection closed
2019-09-25 09:01:08 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16002), connection closed
--- PASS: TestInvoke (0.77s)
=== RUN   TestProgress
        $$$$$ progressHandler1 receive:  message1 <nil>
        &&&&& Progress1: OK
        &&&&& Progress1: Update ->  messsage1.#0
        &&&&& Progress1: Update ->  messsage1.#1
        &&&&& Progress1: Update ->  messsage1.#2
        &&&&& Progress1: Update ->  messsage1.#3
        &&&&& Progress1: Update ->  messsage1.#4
        &&&&& Progress1: Update ->  messsage1.#5
        &&&&& Progress1: Update ->  messsage1.#6
        &&&&& Progress1: Update ->  messsage1.#7
        &&&&& Progress1: Update ->  messsage1.#8
        &&&&& Progress1: Update ->  messsage1.#9
        &&&&& Progress1: Response ->  last message
        $$$$$ progressHandler2 receive:  message2 <nil>
        &&&&& Progress2: OK
        &&&&& Progress2: Update ->  messsage2.#0
        &&&&& Progress2: Update ->  messsage2.#1
        &&&&& Progress2: Update ->  messsage2.#2
        &&&&& Progress2: Update ->  messsage2.#3
        &&&&& Progress2: Update ->  messsage2.#4
        &&&&& Progress2: Response ->  last message2
2019-09-25 09:01:08 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16001), connection closed
2019-09-25 09:01:08 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:47610), connection closed
2019-09-25 09:01:08 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:49724), connection closed
2019-09-25 09:01:08 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16002), connection closed
--- PASS: TestProgress (0.35s)
=== RUN   TestPubSub
        ##########
        # OnPublishRegister:  register#1 <nil>
        ##########
        # OnRegister:  register#2 <nil>
        ##########
        # OnPublish:  publish#1 <nil>
        ##########
        # OnPublish: Notify
        ##########
        # OnPublish:  publish#2 <nil>
        ##########
        # OnPublish: Notify
        &&&&& Subscriber notified: OK,  publish#1
        &&&&& Subscriber notified: OK,  publish#2
        ##########
        # OnPublishDeregister:  deregister#1 <nil>
        ##########
        # OnDeregister:  deregister#2 <nil>
2019-09-25 09:01:09 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:49734), connection closed
2019-09-25 09:01:09 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16003), connection closed
2019-09-25 09:01:09 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:47614), connection closed
2019-09-25 09:01:09 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:56620), connection closed
2019-09-25 09:01:09 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16002), connection closed
2019-09-25 09:01:09 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:56624), connection closed
2019-09-25 09:01:09 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16001), connection closed
--- PASS: TestPubSub (0.51s)
=== RUN   TestReset
2019-09-25 09:01:09 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16003), connection closed
        $$$$$ submitHandler receive:  message1 <nil>
        &&&&& Submit1: OK
        $$$$$ submitHandler receive:  message2 <nil>
        &&&&& Submit2: OK
2019-09-25 09:01:09 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:47622), connection closed
2019-09-25 09:01:09 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16001), connection closed
2019-09-25 09:01:09 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:49736), connection closed
2019-09-25 09:01:09 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16002), connection closed
--- PASS: TestReset (0.50s)
=== RUN   TestNested1
        &&&&& Provider1 receive:  message1 <nil>
        $$$$$ Provider2 receive:  message from provider1 <nil>
        $$$$$ Provider2 ack sent: OK
        &&&&& Nested Invoke1: Ack from provider2
        $$$$$ Provider2 reply sent: OK
        &&&&& Nested Invoke1: OK,  message from provider1 <nil>
        &&&&& Invoke1: OK,  message1
        &&&&& Provider1 receive:  message2 <nil>
        $$$$$ Provider2 receive:  message from provider1 <nil>
        $$$$$ Provider2 ack sent: OK
        &&&&& Nested Invoke1: Ack from provider2
        $$$$$ Provider2 reply sent: OK
        &&&&& Nested Invoke1: OK,  message from provider1 <nil>
        &&&&& Invoke2: OK,  message2
2019-09-25 09:01:10 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16001), connection closed
2019-09-25 09:01:10 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:47626), connection closed
2019-09-25 09:01:10 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16003), connection closed
2019-09-25 09:01:10 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:49746), connection closed
2019-09-25 09:01:10 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:49740), connection closed
2019-09-25 09:01:10 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16002), connection closed
2019-09-25 09:01:10 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:56636), connection closed
2019-09-25 09:01:10 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16001), connection closed
--- PASS: TestNested1 (0.77s)
=== RUN   TestNested2
        $$$$$ Provider1 receive:  message1 <nil>
        $$$$$ Provider2 receive:  message from provider1 <nil>
        $$$$$ Provider2 ack sent: OK
        $$$$$ Provider2 reply sent: OK
        &&&&& Nested Invoke1: OK,  message from provider1 <nil>
        &&&&& Invoke1: OK,  message1
        $$$$$ Provider1 receive:  message2 <nil>
        $$$$$ Provider2 receive:  message from provider1 <nil>
        $$$$$ Provider2 ack sent: OK
        $$$$$ Provider2 reply sent: OK
        &&&&& Nested Invoke1: OK,  message from provider1 <nil>
        &&&&& Invoke2: OK,  message2
2019-09-25 09:01:11 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:47634), connection closed
2019-09-25 09:01:11 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16001), connection closed
2019-09-25 09:01:11 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:49748), connection closed
2019-09-25 09:01:11 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16002), connection closed
2019-09-25 09:01:11 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16001), connection closed
2019-09-25 09:01:11 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:49752), connection closed
--- PASS: TestNested2 (0.77s)
=== RUN   TestNested3Provider
        $$$$$ Provider1 receive:  message1 <nil>
        $$$$$ Provider2 receive:  message from provider1 <nil>
        $$$$$ Provider2 ack sent: OK
        $$$$$ Provider2 reply sent: OK
        &&&&& Nested Invoke1: OK,  message from provider1 <nil>
        &&&&& Invoke1: OK,  message1
        $$$$$ Provider1 receive:  message2 <nil>
        $$$$$ Provider2 receive:  message from provider1 <nil>
        $$$$$ Provider2 ack sent: OK
        $$$$$ Provider2 reply sent: OK
        &&&&& Nested Invoke1: OK,  message from provider1 <nil>
        &&&&& Invoke2: OK,  message2
2019-09-25 09:01:12 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16001), connection closed
2019-09-25 09:01:12 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:47640), connection closed
2019-09-25 09:01:12 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:49754), connection closed
2019-09-25 09:01:12 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16002), connection closed
2019-09-25 09:01:12 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16001), connection closed
2019-09-25 09:01:12 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:49758), connection closed
--- PASS: TestNested3Provider (0.75s)
PASS
ok      github.com/CNES/ccsdsmo-malgo/mal/api   5.455s
```

### Encoding
Tests allowing to verify the interoperability of the encoding implementation. These tests are common to Java and C implementation.

```
[malgo@localhost src]$ GOCACHE=off go test -v github.com/CNES/ccsdsmo-malgo/tests/encoding
=== RUN   TestFixedBinaryEncoding
--- PASS: TestFixedBinaryEncoding (0.00s)
=== RUN   TestFixedBinaryDecoding
--- PASS: TestFixedBinaryDecoding (0.00s)
=== RUN   TestVarintBinaryEncoding
--- PASS: TestVarintBinaryEncoding (0.00s)
=== RUN   TestVarintBinaryDecoding
--- PASS: TestVarintBinaryDecoding (0.00s)
=== RUN   TestSplitBinaryEncoding
--- PASS: TestSplitBinaryEncoding (0.00s)
=== RUN   TestSplitBinaryDecoding
--- PASS: TestSplitBinaryDecoding (0.00s)
PASS
ok      github.com/CNES/ccsdsmo-malgo/tests/encoding    0.004s
```

### Transport Issue1
This test allows to verify the correction about an issue on TCP transport implementation.

```
[malgo@localhost src]$ GOCACHE=off go test -v github.com/CNES/ccsdsmo-malgo/tests/issue1
=== RUN   Test1
        $$$$$ sendHandler receive:  message1 <nil>
        $$$$$ sendHandler receive:  message2 <nil>
        $$$$$ submitHandler receive:  message1 <nil>
        &&&&& Submit1: OK
        $$$$$ submitHandler receive:  message2 <nil>
        &&&&& Submit2: OK
2019-09-25 09:01:24 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:47646), connection closed
2019-09-25 09:01:24 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16001), connection closed
2019-09-25 09:01:24 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:49760), connection closed
2019-09-25 09:01:24 WARNING mal.transport.tcp tcp.go:293 TCPTransport.readMessage(127.0.0.1:16002), connection closed
--- PASS: Test1 (0.51s)
PASS
ok      github.com/CNES/ccsdsmo-malgo/tests/issue1      0.515s
```
