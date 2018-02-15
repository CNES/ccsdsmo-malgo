MAL/GO API
==========


Introduction
============

This GO API represents all MAL concepts, especially the interaction patterns, the MAL message format and the data model. The main goal of this API is to offer a simple and efficient API that makes full use of GO features.

Concepts
========

Most concepts of MAL/GO API are imposed by the MAL specification:

  -	**Message**: MAL Message definition.
  -	**Data types**: Element, Attribute, standard attributes, Composite, base composites and List of these elements.
  - **Access Control**: Access control and security aspects.
  -	**Encoding**: Encoding API.
  -	**Transport**: Transport API.

However, the MAL specification does not define some concepts that are specific to this implementation:

  -	**Context**: MAL communication context.
  -	**End-Point**: Base entity allowing to send and receive MAL messages.
  -	**Operation**: Entity formalizing consumer's request to providers. Each MAL interaction corresponds to a specific operation.
  -	**Handler**, **Service**: Entities formalizing provider's replies to requests. In the first case the behavior of the provider is defined by a unique method (handler), while in the second case it is defined through an interface (service)

Overview of the MAL/GO implementation
=====================================

The MAL/GO implementation shall comprise the following levels:

  - the generic API used by every MAL clients;
  - the data structures and related interfaces;
  - the service consumer API;
  - the service provider API;
  - the broker API.

The MAL/GO implementation consists of several modules:

  -	**mal**: Defines the generic low level MAL API allowing the use of MAL level core concepts. It defines MAL data structures.
  -	**mal/encoding**: Implements the encoding APIs for each encoding format, currently Binary, FixedBinary and SplitBinary.
  -	**mal/transport**: Manages the mapping of MAL concepts to the underlying concepts of transport chosen. Currently two transport are implemented, a local InVM transport and the MAL/TCP standard.
  - **mal/api**: Implements a high level API simplifying the implementation of consumer, provider and broker. It offers the concepts of Operation (consumer side) and Handler (provider and broker side) to handle complex MAL interactions.

MAL data types
--------------

The **MAL Element** is the base type of all data constructs, all types are derived from it. The **Element** GO interface represents this **MAL Element** type, it defines the base interface allowing to encode (resp. decode) each data type.

The **MAL Attribute** is the base type of all attributes of the MAL data model. Attributes are contained within Composites and are used to build complex structures that make the data model. The **Attribute** Go interface represents this **MAL Attribute** type, all GO types defining a MAL attribute inherits from this type. All MAL attributes except Time and FineTime are represented by GO base types.

The **MAL Composite** is the base structure for composite structures that contain a set of elements. The **Composite** GO interface represents this structure, all GO types defining a MAL Composite inherits from this type.

For each GO interface or structure, named X, defining a MAL Element, a XList type shall be defined to handle MAL list of this element.

All data types defined in the MAL standard and the corresponding list are defined in the MAL/GO API.

MAL message
-----------

The **Message** GO structure represents the **MAL Message**. It defines all header data fields, the body is handled as a byte array (encoded representation of this message body). The encoding of the message to (resp. from) the Protocol Data Unit (PDU) is the responsibility of the transport implementation.

MAL context
-----------

The **Context** GO structure is the base entity enabling a client to use the communication functions provided by the MAL layer. Normally it is not be used as
is but through the **EndPoint** structure or the high level API. However you can send MAL messages using the **Context.Send** method, and receive MAL Message
registering a **Listener**.

###Initialization and configuration

A MAL context is created and initialized from a call to **NewContext** function. This function takes in parameter the URI of this context, the corresponding transport will be instantiated and initialized through a factory using this URI (this transport factory needs to be first registered).

```go
func NewContext(url string) (*Context, error)
```

**Note:** Register a transport factory only needs to import for side effect the corresponding package in your program.

###Configuration

There are 2 methods allowing to configure a MAL context:

  - **SetAccessControl** allows to set the Access Control component that checks all incoming and outgoing messages.
  - **SetErrorChannel** allows to set the GO channel for error messages (incoming and outgoing).

```go
func (ctx *Context) SetAccessControl(achdlr AccessControl)
func (ctx *Context) SetErrorChannel(errch chan *MessageError)
```

###Messages sending

The **Context** defines the **Send** request method to transmit a message through the transport component.

```go
func (ctx *Context) Send(msg *Message) error
```

###Messages reception

The **Context** defines the **Receive** indication method to handle incoming messages from the transport component. All incoming messages are pushed in the internal channel if the security check is ok. Then they can be delivered to registered listeners.

```go
func (ctx *Context) Receive(msg *Message) error
```

The **Context** defines method to register (resp. unregister) End-Points, an End-Point is a message listener associated with a specific MAL URI. After registration all messages sent to the corresponding URI are delivered to the End-Point through the onMessage method of its **Listener** interface.
The GetEndPoint method allows to retrieve the End-Point associated with a specific URI.

```go
func (ctx *Context) GetEndPoint(uri *URI) (Listener, error)
func (ctx *Context) RegisterEndPoint(uri *URI, listener Listener) error
func (ctx *Context) UnregisterEndPoint(uri *URI) error
```

###Access control

###Message error

MAL end-point
-------------

The **EndPoint** defines a generic implementation of the **Listener** interface, it is the base entity allowing to send and receive MAL Messages.
It defines a method to send messages **EndPoint.Send**, a blocking method to receive messages **EndPoint.Receive**, and a method to close and unregister
the end-point **EndPoint.Close**. The construction of each MAL message shall be done manually.
 
Each end-point handles an atomic counter allowing to generates the appropriate MAL TransactionId using the **EndPoint.TransactionId** method.
A end-point is created by using the NewEndPoint function.

```go
func NewEndPoint(ctx *Context, service string, ch chan *Message) (*EndPoint, error) {

func (endpoint *EndPoint) TransactionId() ULong
func (endpoint *EndPoint) Send(msg *Message) error
func (endpoint *EndPoint) Recv() (*Message, error)

func (endpoint *EndPoint) Close() error
```

MAL client context
------------------

The **ClientContext** entity is the entry point of the high level API. It corresponds to a unique MAL URI and is registered as an end-point with the
underlying MAL context. The **ClientContext** entity manages the MAL TransactionId. It allows providers to register handlers to process consumer's
requests. It allows consumers to create operations to initiate and manage interaction with providers.

###Initialization and configuration

A **ClientContext** is created and initialized from a call to **NewClientContext** function (package mal/api). This function takes in parameter the
underlying MAL context, and the name of the service (last part of the MAL URI for the corresponding end-point).

```go
func NewClientContext(ctx *Context, service string) (*ClientContext, error)
```

###Registering provider's handler

The **ClientContext** entity defines a set of 6 methods allowing providers to register handlers to process consumer's requests. Each method is dedicated
to a particular MAL interaction, the corresponding handler receives in parameter a **Transaction** entity corresponding to the interaction. This entity
provides methods to handle the interaction.

  - **RegisterSendHandler** allows to register an handler for a Send interaction. When called the handler function receives a **SendTransaction** parameter
  with no method.
  - **RegisterSubmitHandler** allows to register an handler for a Submit interaction. When called the handler function receives a **SubmitTransaction** parameter that provides an **Ack** method allowing to acknowledge the consumer.
  - **RegisterRequestHandler** allows to register an handler for a Request interaction. When called the handler function receives a **RequestTransaction** parameter that provides an **Reply** method allowing to reply to the consumer.
  - **RegisterInvokeHandler** allows to register an handler for a Invoke interaction. When called the handler function receives an **InvokeTransaction** parameter that provides 2 methods. The **Ack** method to acknowledge the the incoming message, and the **Reply** method to reply to the consumer.
  - **RegisterProgressHandler** allows to register an handler for a progress interaction. When called the handler function receives a **SendTransaction** parameter that provides 3 methods. The **Ack** method to acknowledge the the incoming message, the **Update** method to send updates and the **Reply** method to reply to the consumer.
  - **RegisterBrokerHandler** allows to register an handler for a PubSub interaction. When called the handler function receives either a **SubscriberTransaction** or a **PublisherTransaction**. If the interaction comes from the subscriber the transaction parameter is a **SubscriberTransaction**, it allows to acknowledge REGISTER or DEREGISTER messages, and to send NOTIFY to subscriber.If the interaction comes from the publisher the transaction parameter is a **PublisherTransaction**, it allows to acknowledge PUBLISH\_REGISTER or PUBLISH\_DEREGISTER messages. 

The definition of handler interface is:

```go
type ProviderHandler func(*Message, Transaction) error

// Register an handler for Send interaction
func (cctx *ClientContext) RegisterSendHandler(area UShort, areaVersion UOctet, service UShort, operation UShort,
                                               handler ProviderHandler) error
                                               
// Register an handler for Submit interaction
func (cctx *ClientContext) RegisterSubmitHandler(area UShort, areaVersion UOctet, service UShort, operation UShort,
                                                 handler ProviderHandler) error
                                                 
// Register an handler for Request interaction
func (cctx *ClientContext) RegisterRequestHandler(area UShort, areaVersion UOctet, service UShort, operation UShort,
                                                  handler ProviderHandler) error
                                                  
// Register an handler for Invoke interaction
func (cctx *ClientContext) RegisterInvokeHandler(area UShort, areaVersion UOctet, service UShort, operation UShort,
                                                 handler ProviderHandler) error
                                                 
// Register an handler for progress interaction
func (cctx *ClientContext) RegisterProgressHandler(area UShort, areaVersion UOctet, service UShort, operation UShort,
                                                   handler ProviderHandler) error
                                                   
// Register an handler for PubSub interaction
func (cctx *ClientContext) RegisterBrokerHandler(area UShort, areaVersion UOctet, service UShort, operation UShort,
                                                 handler ProviderHandler) error
```

The definition of transactions entities are:

```go
type SendTransaction interface {
	Transaction
}

// MAL Submit interaction
type SubmitTransaction interface {
	Transaction
	Ack(err error) error
}


func (tx *SubmitTransactionX) Ack(err error) error

// MAL Request interaction
type RequestTransaction interface {
	Reply([]byte, error) error
}


func (tx *RequestTransactionX) Reply(body []byte, err error) error

// MAL Invoke interaction
type InvokeTransaction interface {
	Ack(error) error
	Reply([]byte, error) error
}

func (tx *InvokeTransactionX) Ack(err error) error
func (tx *InvokeTransactionX) Reply(body []byte, err error) error

// MAL Progress interaction
type ProgressTransaction interface {
	Ack(error) error
	Update([]byte, error) error
	Reply([]byte, error) error
}

func (tx *ProgressTransactionX) Ack(err error) error
func (tx *ProgressTransactionX) Update(body []byte, err error) error
func (tx *ProgressTransactionX) Reply(body []byte, err error) error

// SubscriberTransaction
type SubscriberTransaction interface {
	Transaction
	AckRegister(error) error
	AckDeregister(error) error
	Notify([]byte, error) error
}

func (tx *SubscriberTransactionX) AckRegister(err error) error
func (tx *SubscriberTransactionX) Notify(body []byte, err error) error
func (tx *SubscriberTransactionX) AckDeregister(err error) error

// PublisherTransaction
type PublisherTransaction interface {
	Transaction
	AckRegister(error) error
	AckDeregister(error) error
}

func (tx *PublisherTransactionX) AckRegister(err error) error
func (tx *PublisherTransactionX) AckDeregister(err error) error

```

###Creating consumer's operation

The **ClientContext** entity defines a set of 7 methods allowing consumers to create entities to request providers. Each created entity encapsulates
the resources needed to enable a MAL consumer to initiate and manage interactions with a MAL provider or broker.

  - **NewSendOperation** creates an operation entity allowing a consumer to request a provider with a SEND. It returns a **SendOperation** entity that
  define a **Send** method allowing to request the provider with the SEND message. This method returns when the message is sent.
  - **NewSubmitOperation** creates an operation entity allowing a consumer to request a provider with a SUBMIT. It returns a **SubmitOperation** entity
  that define a **Submit** method to request the provider with the SUBMIT message. This method returns when the acknowledge from provider is received.
  - **NewRequestOperation** creates an operation entity allowing a consumer to request a provider with a REQUEST. It returns a **RequestOperation** entity
  that define a **Request** method allowing to request the provider with the REQUEST message. This method returns the RESPONSE message from provider.
  - **NewInvokeOperation** creates an operation entity allowing a consumer to request a provider with a INVOKE. It returns an **InvokeOperation** entity
  that define an **Invoke** method allowing to request the provider with the INVOKE message. This method returns when the acknowledge from provider is
  received. The operation also defines a GetResponse method to wait the RESPONSE message from the provider.
  - **NewProgressOperation** creates an operation entity allowing a consumer to request a provider with a PROGRESS. It returns an **InvokeOperation** entity
  that define a **Progress** method allowing to request the provider with the PROGRESS message.  This method returns when the acknowledge from provider is
  received. The operation define also 2 other methods:
    - the getUpdate method returns the UPDATE messages from the provider or nil if there is no other updates. 
    - the getResponse method returns the RESPONSE message from the provider. If there is remaining updates to deliver
    to the consumer they are deleted.
  - **NewSubscriberOperation** creates an operation entity allowing a subscriber to interact with a broker, it returns a **SubscriberOperation** entity
  that define 3 methods:
    - the register method allows to send a REGISTER message to the broker.
    - the getNotify method allows to retrieve the NOTIFY message sent by the broker.
    - the deregister method allows to send a DEREGISTER message to the broker.
  - **NewPublisherOperation** creates an operation entity allowing a publisher to interact with a broker, it returns a **PublisherOperation** entity
  that define 3 methods:
    - the register method allows to send a PUBLISH\_REGISTER message to the broker.
    - the getNotify method allows to send a PUBLISH message to the broker.
    - the deregister method allows to send a PUBLISH\_DEREGISTER message to the broker.

Additionally all **Operation** entities provide a **Close** method and a **Reset** method. The **Reset** method allows to reuse the operation after
its completion.

The interface of ClientContext is:

```go
func NewClientContext(ctx *Context, service string) (*ClientContext, error)
func (cctx *ClientContext) TransactionId() ULong

func (cctx *ClientContext) NewSendOperation(urito *URI,
	area UShort, areaVersion UOctet, service UShort, operation UShort) SendOperation 
	
func (cctx *ClientContext) NewSubmitOperation(urito *URI,
	area UShort, areaVersion UOctet, service UShort, operation UShort) SubmitOperation
	
func (cctx *ClientContext) NewRequestOperation(urito *URI,
	area UShort, areaVersion UOctet, service UShort, operation UShort) RequestOperation

func (cctx *ClientContext) NewInvokeOperation(urito *URI,
	area UShort, areaVersion UOctet, service UShort, operation UShort) InvokeOperation

func (cctx *ClientContext) NewProgressOperation(urito *URI,
	area UShort, areaVersion UOctet, service UShort, operation UShort) ProgressOperation

func (cctx *ClientContext) NewSubscriberOperation(urito *URI,
	area UShort, areaVersion UOctet, service UShort, operation UShort) SubscriberOperation

func (cctx *ClientContext) NewPublisherOperation(urito *URI,
	area UShort, areaVersion UOctet, service UShort, operation UShort) PublisherOperation
```

The definition of operation entities are:

```go
type Operation interface {
	GetTid() ULong
	Close() error
	Reset() error
}

type SendOperation interface {
	Operation
	Send(body []byte) error
}

type SubmitOperation interface {
	Operation
	Submit(body []byte) (*Message, error)
}

type RequestOperation interface {
	Operation
	Request(body []byte) (*Message, error)
}

type InvokeOperation interface {
	Operation
	Invoke(body []byte) (*Message, error)
	GetResponse() (*Message, error)
}

type ProgressOperation interface {
	Operation
	Progress(body []byte) (*Message, error)
	GetUpdate() (*Message, error)
	GetResponse() (*Message, error)
}

type SubscriberOperation interface {
	Operation
	Register(body []byte) (*Message, error)
	GetNotify() (*Message, error)
	Deregister(body []byte) (*Message, error)
}

type PublisherOperation interface {
	Operation
	Register(body []byte) (*Message, error)
	Publish(body []byte) error
	Deregister(body []byte) (*Message, error)
}
```

###Closing a client context

The **Close** method allows to close and unregister the context.s

```go
func (cctx *ClientContext) Close() error
```

User's Guide of the MAL/GO implementationn
==========================================
Creating a MAL program using the MAL/GO implementation needs first to create and initialize a MAL context, then you can either use the low level API with
**EndPoint** or use the high level API with **ClientContext**, Handlers and operations. These two modes are demonstrated through code samples below.

Creation of MAL context
----------------------

A MAL context is created and initialized from a call to **NewContext** function of the mal package.
This function takes in parameter an URL corresponding to the URI of this context:

  - The scheme part of this URL defines the underlying transport (This transport factory needs to be first registered, it is done by importing 
  for side effect the corresponding package).
  - The host and port determines the binding of listen socket.
  - The query part of this URL contains the optional parameters needed by the transport.
  
After use the MAL context shall be closed using the Close method, this call closes all registered listeners (end-point, etc) and finalizes the underlying transport.

```go
func NewContext(url string) (*Context, error)
func (ctx *Context) Close() error
```

Example: The code below creates a MAL context using MALTCP transport and listening on port 16000 of the local network interface.
The Close method of this context will be called at the end of the current function.

```go
consumer_ctx, err := NewContext("maltcp://127.0.0.1:16000")
if err != nil 
	t.Fatal("Error creating context, ", err
	return
defer consumer_ctx.Close()
```

Using a simple end-point
------------------------

The end-point is created using the **NewEndPoint** function, the parameters are:

  - the underlying MAL context.
  - the service name (the service URI will be built by concatenation of MAL context URI with this name).
  - an optional channel (if nil a channel is created during the end-point initialization).

A code sample is available in the TCP transport tests (mal/transport/tcp/transport_test.go).

```go
service, err := NewEndPoint(ctx, "service", nil)
if err != nil {
	logger.Fatalf("Error creating service end-point: ", err)
}
```

###Sending and Receiving a message

```go
msg := &Message{
	UriFrom:          service.Uri,
	UriTo:            providerUri,
	TransactionId:    service.TransactionId(),
	...
	Body:             payload,
}
service.Send(msg)

msg, err := service.Recv()
if msg != nil {
	logger.Errorf("Error receiving a message: ", err)
}
```

Using the High level API
------------------------

###Implementing a Provider

The example below shows a provider handling a unique progress interaction.
Code samples are available in the MAL API tests (mal/api/*_test.go).

```go
// Provider context containing its datas
type MyProvider struct {
	ctx   *Context
	cctx  *ClientContext
	nbmsg int
}

// Creation and initialization of the provider
func newMyProvider(url string, service string) (*MyProvider, error) {
	// Creates and initializes the MAL context
	ctx, err := NewContext(url)
	if err != nil {
		return nil, err
	}
	// Creates the client context for the provider
	cctx, err := NewClientContext(ctx, service)
	if err != nil {
		return nil, err
	}
	// Allocates and initializes the provider structure
	provider := &MyProvider{ctx, cctx, 0}

	// Handler of a Progress operation
	progressHandler := func(msg *Message, t Transaction) error {
		provider.nbmsg += 1
		if msg != nil {
			transaction := t.(ProgressTransaction)
			transaction.Ack(nil, false)
			for i := 0; i < 10; i++ {
				transaction.Update(update, false)
			}
			transaction.Reply(response, false)
		} else {
			logger.Warnf("receive: nil")
		}
		return nil
	}
	// Registers the handler above
	cctx.RegisterProgressHandler(200, 1, 1, 1, progressHandler1)

	// Declares and registers other handlers..

	return provider, nil
}

// Close the provider.
func (provider *MyProvider) close() {
	provider.ctx.Close()
}
```

The code below allows to create the provider.

```go
// Creates a provider registered with maltcp://127.0.0.1:16000/service URI
provider, err := newMyProvider("maltcp://127.0.0.1:16000", "service")
if err != nil {
	logger.Errorf("Error creating provider: ", err)
	return
}
defer provider.close()
```

###Implementing a Consumer

The code below shows how to request the progress interaction of the provider above.

```go
// Creates a new MAL context.
ctx, err := NewContext("maltcp://127.0.0.1:16001")
if err != nil {
	return err
}
defer ctx.Close()

// Creates a client context for consumer's operations.
consumer, err := NewClientContext(consumer_ctx, "consumer")
if err != nil {
	return err
}

// Creates a new ProgressOperation from the consumer's context.
op := consumer.NewProgressOperation("maltcp://127.0.0.1:16000/service", 200, 1, 1, 1)
// Initiates a Progress interaction with payload ([]byte encoded content).
op.Progress(payload)

// Gets Update messages from service
updt, err := op1.GetUpdate()
if err != nil {
	return err
}
for updt != nil {
	updt, err = op1.GetUpdate()
	if err != nil {
		return err
	}
}
// There is no more updates, gets the Response message from service.
rep, err := op1.GetResponse()
if err != nil {
	return err
}
```
  
