package archive

import (
  "errors"
  "github.com/CNES/ccsdsmo-malgo/mal"
  malapi "github.com/CNES/ccsdsmo-malgo/mal/api"
  "github.com/CNES/ccsdsmo-malgo/com"
)


// service provider internal interface
type ProviderInterface interface {
  Retrieve(opHelper *RetrieveHelper, objType *com.ObjectType, domain *mal.IdentifierList, objInstIds *mal.LongList) error
  Query(opHelper *QueryHelper, returnBody *mal.Boolean, objType *com.ObjectType, archiveQuery *ArchiveQueryList, queryFilter QueryFilterList) error
  Count(opHelper *CountHelper, objType *com.ObjectType, archiveQuery *ArchiveQueryList, queryFilter QueryFilterList) error
  Store(opHelper *StoreHelper, returnObjInstIds *mal.Boolean, objType *com.ObjectType, domain *mal.IdentifierList, objDetails *ArchiveDetailsList, objBodies mal.ElementList) error
  Update(opHelper *UpdateHelper, objType *com.ObjectType, domain *mal.IdentifierList, objDetails *ArchiveDetailsList, objBodies mal.ElementList) error
  Delete(opHelper *DeleteHelper, objType *com.ObjectType, domain *mal.IdentifierList, objInstIds *mal.LongList) error
}


// service provider structure
type Provider struct {
  Cctx *malapi.ClientContext
  provider ProviderInterface
}

// create a service provider
func NewProvider(ctx *mal.Context, uri string, providerImpl ProviderInterface) (*Provider, error) {
  cctx, err := malapi.NewClientContext(ctx, uri)
  if err != nil {
    return nil, err
  }
  // define the handler for operation Retrieve
  RetrieveHandler := func(msg *mal.Message, t malapi.Transaction) error {
  opHelper, err := NewRetrieveHelper(t)
  if err != nil {
    return err
  }
    if msg == nil {
      err := errors.New("missing Message")
      return opHelper.ReturnError(err)
    }
    // decode in parameters
    inElem_objType, err := msg.DecodeParameter(com.NullObjectType)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_objType, ok := inElem_objType.(*com.ObjectType)
    if !ok {
      err = errors.New("unexpected type for parameter objType")
      return opHelper.ReturnError(err)
    }
    inElem_domain, err := msg.DecodeParameter(mal.NullIdentifierList)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_domain, ok := inElem_domain.(*mal.IdentifierList)
    if !ok {
      err = errors.New("unexpected type for parameter domain")
      return opHelper.ReturnError(err)
    }
    inElem_objInstIds, err := msg.DecodeLastParameter(mal.NullLongList, false)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_objInstIds, ok := inElem_objInstIds.(*mal.LongList)
    if !ok {
      err = errors.New("unexpected type for parameter objInstIds")
      return opHelper.ReturnError(err)
    }
    // call the provider implementation
    err = providerImpl.Retrieve(opHelper, inParam_objType, inParam_domain, inParam_objInstIds)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    return nil
  }
  // register the handler
  err = cctx.RegisterInvokeHandler(com.AREA_NUMBER, com.AREA_VERSION, SERVICE_NUMBER, RETRIEVE_OPERATION_NUMBER, RetrieveHandler)
  if err != nil {
    return nil, err
  }
  // define the handler for operation Query
  QueryHandler := func(msg *mal.Message, t malapi.Transaction) error {
  opHelper, err := NewQueryHelper(t)
  if err != nil {
    return err
  }
    if msg == nil {
      err := errors.New("missing Message")
      return opHelper.ReturnError(err)
    }
    // decode in parameters
    inElem_returnBody, err := msg.DecodeParameter(mal.NullBoolean)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_returnBody, ok := inElem_returnBody.(*mal.Boolean)
    if !ok {
      err = errors.New("unexpected type for parameter returnBody")
      return opHelper.ReturnError(err)
    }
    inElem_objType, err := msg.DecodeParameter(com.NullObjectType)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_objType, ok := inElem_objType.(*com.ObjectType)
    if !ok {
      err = errors.New("unexpected type for parameter objType")
      return opHelper.ReturnError(err)
    }
    inElem_archiveQuery, err := msg.DecodeParameter(NullArchiveQueryList)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_archiveQuery, ok := inElem_archiveQuery.(*ArchiveQueryList)
    if !ok {
      err = errors.New("unexpected type for parameter archiveQuery")
      return opHelper.ReturnError(err)
    }
    inElem_queryFilter, err := msg.DecodeLastParameter(nil, true)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_queryFilter, ok := inElem_queryFilter.(QueryFilterList)
    if !ok {
      if inElem_queryFilter == mal.NullElement {
        inParam_queryFilter = NullQueryFilterList
      } else {
        err = errors.New("unexpected type for parameter queryFilter")
        return opHelper.ReturnError(err)
      }
    }
    // call the provider implementation
    err = providerImpl.Query(opHelper, inParam_returnBody, inParam_objType, inParam_archiveQuery, inParam_queryFilter)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    return nil
  }
  // register the handler
  err = cctx.RegisterProgressHandler(com.AREA_NUMBER, com.AREA_VERSION, SERVICE_NUMBER, QUERY_OPERATION_NUMBER, QueryHandler)
  if err != nil {
    return nil, err
  }
  // define the handler for operation Count
  CountHandler := func(msg *mal.Message, t malapi.Transaction) error {
  opHelper, err := NewCountHelper(t)
  if err != nil {
    return err
  }
    if msg == nil {
      err := errors.New("missing Message")
      return opHelper.ReturnError(err)
    }
    // decode in parameters
    inElem_objType, err := msg.DecodeParameter(com.NullObjectType)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_objType, ok := inElem_objType.(*com.ObjectType)
    if !ok {
      err = errors.New("unexpected type for parameter objType")
      return opHelper.ReturnError(err)
    }
    inElem_archiveQuery, err := msg.DecodeParameter(NullArchiveQueryList)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_archiveQuery, ok := inElem_archiveQuery.(*ArchiveQueryList)
    if !ok {
      err = errors.New("unexpected type for parameter archiveQuery")
      return opHelper.ReturnError(err)
    }
    inElem_queryFilter, err := msg.DecodeLastParameter(nil, true)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_queryFilter, ok := inElem_queryFilter.(QueryFilterList)
    if !ok {
      if inElem_queryFilter == mal.NullElement {
        inParam_queryFilter = NullQueryFilterList
      } else {
        err = errors.New("unexpected type for parameter queryFilter")
        return opHelper.ReturnError(err)
      }
    }
    // call the provider implementation
    err = providerImpl.Count(opHelper, inParam_objType, inParam_archiveQuery, inParam_queryFilter)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    return nil
  }
  // register the handler
  err = cctx.RegisterInvokeHandler(com.AREA_NUMBER, com.AREA_VERSION, SERVICE_NUMBER, COUNT_OPERATION_NUMBER, CountHandler)
  if err != nil {
    return nil, err
  }
  // define the handler for operation Store
  StoreHandler := func(msg *mal.Message, t malapi.Transaction) error {
  opHelper, err := NewStoreHelper(t)
  if err != nil {
    return err
  }
    if msg == nil {
      err := errors.New("missing Message")
      return opHelper.ReturnError(err)
    }
    // decode in parameters
    inElem_returnObjInstIds, err := msg.DecodeParameter(mal.NullBoolean)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_returnObjInstIds, ok := inElem_returnObjInstIds.(*mal.Boolean)
    if !ok {
      err = errors.New("unexpected type for parameter returnObjInstIds")
      return opHelper.ReturnError(err)
    }
    inElem_objType, err := msg.DecodeParameter(com.NullObjectType)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_objType, ok := inElem_objType.(*com.ObjectType)
    if !ok {
      err = errors.New("unexpected type for parameter objType")
      return opHelper.ReturnError(err)
    }
    inElem_domain, err := msg.DecodeParameter(mal.NullIdentifierList)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_domain, ok := inElem_domain.(*mal.IdentifierList)
    if !ok {
      err = errors.New("unexpected type for parameter domain")
      return opHelper.ReturnError(err)
    }
    inElem_objDetails, err := msg.DecodeParameter(NullArchiveDetailsList)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_objDetails, ok := inElem_objDetails.(*ArchiveDetailsList)
    if !ok {
      err = errors.New("unexpected type for parameter objDetails")
      return opHelper.ReturnError(err)
    }
    inElem_objBodies, err := msg.DecodeLastParameter(nil, true)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_objBodies, ok := inElem_objBodies.(mal.ElementList)
    if !ok {
      if inElem_objBodies == mal.NullElement {
        inParam_objBodies = mal.NullElementList
      } else {
        err = errors.New("unexpected type for parameter objBodies")
        return opHelper.ReturnError(err)
      }
    }
    // call the provider implementation
    err = providerImpl.Store(opHelper, inParam_returnObjInstIds, inParam_objType, inParam_domain, inParam_objDetails, inParam_objBodies)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    return nil
  }
  // register the handler
  err = cctx.RegisterRequestHandler(com.AREA_NUMBER, com.AREA_VERSION, SERVICE_NUMBER, STORE_OPERATION_NUMBER, StoreHandler)
  if err != nil {
    return nil, err
  }
  // define the handler for operation Update
  UpdateHandler := func(msg *mal.Message, t malapi.Transaction) error {
  opHelper, err := NewUpdateHelper(t)
  if err != nil {
    return err
  }
    if msg == nil {
      err := errors.New("missing Message")
      return opHelper.ReturnError(err)
    }
    // decode in parameters
    inElem_objType, err := msg.DecodeParameter(com.NullObjectType)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_objType, ok := inElem_objType.(*com.ObjectType)
    if !ok {
      err = errors.New("unexpected type for parameter objType")
      return opHelper.ReturnError(err)
    }
    inElem_domain, err := msg.DecodeParameter(mal.NullIdentifierList)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_domain, ok := inElem_domain.(*mal.IdentifierList)
    if !ok {
      err = errors.New("unexpected type for parameter domain")
      return opHelper.ReturnError(err)
    }
    inElem_objDetails, err := msg.DecodeParameter(NullArchiveDetailsList)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_objDetails, ok := inElem_objDetails.(*ArchiveDetailsList)
    if !ok {
      err = errors.New("unexpected type for parameter objDetails")
      return opHelper.ReturnError(err)
    }
    inElem_objBodies, err := msg.DecodeLastParameter(nil, true)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_objBodies, ok := inElem_objBodies.(mal.ElementList)
    if !ok {
      if inElem_objBodies == mal.NullElement {
        inParam_objBodies = mal.NullElementList
      } else {
        err = errors.New("unexpected type for parameter objBodies")
        return opHelper.ReturnError(err)
      }
    }
    // call the provider implementation
    err = providerImpl.Update(opHelper, inParam_objType, inParam_domain, inParam_objDetails, inParam_objBodies)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    return nil
  }
  // register the handler
  err = cctx.RegisterSubmitHandler(com.AREA_NUMBER, com.AREA_VERSION, SERVICE_NUMBER, UPDATE_OPERATION_NUMBER, UpdateHandler)
  if err != nil {
    return nil, err
  }
  // define the handler for operation Delete
  DeleteHandler := func(msg *mal.Message, t malapi.Transaction) error {
  opHelper, err := NewDeleteHelper(t)
  if err != nil {
    return err
  }
    if msg == nil {
      err := errors.New("missing Message")
      return opHelper.ReturnError(err)
    }
    // decode in parameters
    inElem_objType, err := msg.DecodeParameter(com.NullObjectType)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_objType, ok := inElem_objType.(*com.ObjectType)
    if !ok {
      err = errors.New("unexpected type for parameter objType")
      return opHelper.ReturnError(err)
    }
    inElem_domain, err := msg.DecodeParameter(mal.NullIdentifierList)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_domain, ok := inElem_domain.(*mal.IdentifierList)
    if !ok {
      err = errors.New("unexpected type for parameter domain")
      return opHelper.ReturnError(err)
    }
    inElem_objInstIds, err := msg.DecodeLastParameter(mal.NullLongList, false)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    inParam_objInstIds, ok := inElem_objInstIds.(*mal.LongList)
    if !ok {
      err = errors.New("unexpected type for parameter objInstIds")
      return opHelper.ReturnError(err)
    }
    // call the provider implementation
    err = providerImpl.Delete(opHelper, inParam_objType, inParam_domain, inParam_objInstIds)
    if err != nil {
      return opHelper.ReturnError(err)
    }
    return nil
  }
  // register the handler
  err = cctx.RegisterRequestHandler(com.AREA_NUMBER, com.AREA_VERSION, SERVICE_NUMBER, DELETE_OPERATION_NUMBER, DeleteHandler)
  if err != nil {
    return nil, err
  }
  provider := &Provider { cctx, providerImpl }
  return provider, nil
}

func (receiver *Provider) Close() error {
  if receiver.Cctx != nil {
    err := receiver.Cctx.Close()
    if err != nil {
      return err
    }
  }
  return nil
}

// generated code for operation retrieve
type RetrieveHelper struct {
  acked bool
  transaction malapi.InvokeTransaction
}
func NewRetrieveHelper(transaction malapi.Transaction) (*RetrieveHelper, error) {
  iptransaction, ok := transaction.(malapi.InvokeTransaction)
  if !ok {
    return nil, errors.New("Unexpected transaction type")
  }
helper := &RetrieveHelper{false, iptransaction}
  return helper, nil
}
func (receiver *RetrieveHelper) Ack() error {
  transaction := receiver.transaction

  // interaction call
  err := transaction.Ack(nil, false)
  if err != nil {
    return err
  }

  receiver.acked = true
  return nil
}
func (receiver *RetrieveHelper) Reply(objDetails *ArchiveDetailsList, objBodies mal.ElementList) error {
  transaction := receiver.transaction
  // create a body for the interaction call
  body := transaction.NewBody()
  // encode parameters
  err := body.EncodeParameter(objDetails)
  if err != nil {
    return nil
  }
  err = body.EncodeLastParameter(objBodies, true)
  if err != nil {
    return nil
  }

  // interaction call
  err = transaction.Reply(body, false)
  if err != nil {
    return err
  }

  receiver.acked = true
  return nil
}

func (receiver *RetrieveHelper) ReturnError(e error) error {
  transaction := receiver.transaction
  body := transaction.NewBody()
  var errCode *mal.UInteger
  var errExtraInfo mal.Element
  var errIsAbstract bool
  malErr, ok := e.(*malapi.MalError)
  if ok {
    errCode = &malErr.Code
    errExtraInfo = malErr.ExtraInfo
    errIsAbstract = false
  } else {
    // return an UNKNOWN error with a String information
    errCode = mal.NewUInteger(uint32(mal.ERROR_UNKNOWN))
    errExtraInfo = mal.NewString(e.Error())
    errIsAbstract = false
  }
  // encode parameters
  if body.EncodeParameter(errCode) != nil {
    return errors.New("Fatal error in error handling code")
  } else if body.EncodeLastParameter(errExtraInfo, errIsAbstract) != nil {
    return errors.New("Fatal error in error handling code")
  }

  // interaction call
  var err error
  if !receiver.acked {
    err = transaction.Ack(body, true)
  } else {
    err = transaction.Reply(body, true)
  }
  if err != nil {
    return err
  }
  return nil
}

// generated code for operation query
type QueryHelper struct {
  acked bool
  transaction malapi.ProgressTransaction
}
func NewQueryHelper(transaction malapi.Transaction) (*QueryHelper, error) {
  iptransaction, ok := transaction.(malapi.ProgressTransaction)
  if !ok {
    return nil, errors.New("Unexpected transaction type")
  }
helper := &QueryHelper{false, iptransaction}
  return helper, nil
}
func (receiver *QueryHelper) Ack() error {
  transaction := receiver.transaction

  // interaction call
  err := transaction.Ack(nil, false)
  if err != nil {
    return err
  }

  receiver.acked = true
  return nil
}
func (receiver *QueryHelper) Update(objType *com.ObjectType, domain *mal.IdentifierList, objDetails *ArchiveDetailsList, objBodies mal.ElementList) error {
  transaction := receiver.transaction
  // create a body for the interaction call
  body := transaction.NewBody()
  // encode parameters
  err := body.EncodeParameter(objType)
  if err != nil {
    return nil
  }
  err = body.EncodeParameter(domain)
  if err != nil {
    return nil
  }
  err = body.EncodeParameter(objDetails)
  if err != nil {
    return nil
  }
  err = body.EncodeLastParameter(objBodies, true)
  if err != nil {
    return nil
  }

  // interaction call
  err = transaction.Update(body, false)
  if err != nil {
    return err
  }

  receiver.acked = true
  return nil
}
func (receiver *QueryHelper) Reply(objType *com.ObjectType, domain *mal.IdentifierList, objDetails *ArchiveDetailsList, objBodies mal.ElementList) error {
  transaction := receiver.transaction
  // create a body for the interaction call
  body := transaction.NewBody()
  // encode parameters
  err := body.EncodeParameter(objType)
  if err != nil {
    return nil
  }
  err = body.EncodeParameter(domain)
  if err != nil {
    return nil
  }
  err = body.EncodeParameter(objDetails)
  if err != nil {
    return nil
  }
  err = body.EncodeLastParameter(objBodies, true)
  if err != nil {
    return nil
  }

  // interaction call
  err = transaction.Reply(body, false)
  if err != nil {
    return err
  }

  receiver.acked = true
  return nil
}

func (receiver *QueryHelper) ReturnError(e error) error {
  transaction := receiver.transaction
  body := transaction.NewBody()
  var errCode *mal.UInteger
  var errExtraInfo mal.Element
  var errIsAbstract bool
  malErr, ok := e.(*malapi.MalError)
  if ok {
    errCode = &malErr.Code
    errExtraInfo = malErr.ExtraInfo
    errIsAbstract = false
  } else {
    // return an UNKNOWN error with a String information
    errCode = mal.NewUInteger(uint32(mal.ERROR_UNKNOWN))
    errExtraInfo = mal.NewString(e.Error())
    errIsAbstract = false
  }
  // encode parameters
  if body.EncodeParameter(errCode) != nil {
    return errors.New("Fatal error in error handling code")
  } else if body.EncodeLastParameter(errExtraInfo, errIsAbstract) != nil {
    return errors.New("Fatal error in error handling code")
  }

  // interaction call
  var err error
  if !receiver.acked {
    err = transaction.Ack(body, true)
  } else {
    err = transaction.Update(body, true)
  }
  if err != nil {
    return err
  }
  return nil
}

// generated code for operation count
type CountHelper struct {
  acked bool
  transaction malapi.InvokeTransaction
}
func NewCountHelper(transaction malapi.Transaction) (*CountHelper, error) {
  iptransaction, ok := transaction.(malapi.InvokeTransaction)
  if !ok {
    return nil, errors.New("Unexpected transaction type")
  }
helper := &CountHelper{false, iptransaction}
  return helper, nil
}
func (receiver *CountHelper) Ack() error {
  transaction := receiver.transaction

  // interaction call
  err := transaction.Ack(nil, false)
  if err != nil {
    return err
  }

  receiver.acked = true
  return nil
}
func (receiver *CountHelper) Reply(counts *mal.LongList) error {
  transaction := receiver.transaction
  // create a body for the interaction call
  body := transaction.NewBody()
  // encode parameters
  err := body.EncodeLastParameter(counts, false)
  if err != nil {
    return nil
  }

  // interaction call
  err = transaction.Reply(body, false)
  if err != nil {
    return err
  }

  receiver.acked = true
  return nil
}

func (receiver *CountHelper) ReturnError(e error) error {
  transaction := receiver.transaction
  body := transaction.NewBody()
  var errCode *mal.UInteger
  var errExtraInfo mal.Element
  var errIsAbstract bool
  malErr, ok := e.(*malapi.MalError)
  if ok {
    errCode = &malErr.Code
    errExtraInfo = malErr.ExtraInfo
    errIsAbstract = false
  } else {
    // return an UNKNOWN error with a String information
    errCode = mal.NewUInteger(uint32(mal.ERROR_UNKNOWN))
    errExtraInfo = mal.NewString(e.Error())
    errIsAbstract = false
  }
  // encode parameters
  if body.EncodeParameter(errCode) != nil {
    return errors.New("Fatal error in error handling code")
  } else if body.EncodeLastParameter(errExtraInfo, errIsAbstract) != nil {
    return errors.New("Fatal error in error handling code")
  }

  // interaction call
  var err error
  if !receiver.acked {
    err = transaction.Ack(body, true)
  } else {
    err = transaction.Reply(body, true)
  }
  if err != nil {
    return err
  }
  return nil
}

// generated code for operation store
type StoreHelper struct {
  transaction malapi.RequestTransaction
}
func NewStoreHelper(transaction malapi.Transaction) (*StoreHelper, error) {
  iptransaction, ok := transaction.(malapi.RequestTransaction)
  if !ok {
    return nil, errors.New("Unexpected transaction type")
  }
helper := &StoreHelper{iptransaction}
  return helper, nil
}
func (receiver *StoreHelper) Reply(objInstIds *mal.LongList) error {
  transaction := receiver.transaction
  // create a body for the interaction call
  body := transaction.NewBody()
  // encode parameters
  err := body.EncodeLastParameter(objInstIds, false)
  if err != nil {
    return nil
  }

  // interaction call
  err = transaction.Reply(body, false)
  if err != nil {
    return err
  }

  return nil
}

func (receiver *StoreHelper) ReturnError(e error) error {
  transaction := receiver.transaction
  body := transaction.NewBody()
  var errCode *mal.UInteger
  var errExtraInfo mal.Element
  var errIsAbstract bool
  malErr, ok := e.(*malapi.MalError)
  if ok {
    errCode = &malErr.Code
    errExtraInfo = malErr.ExtraInfo
    errIsAbstract = false
  } else {
    // return an UNKNOWN error with a String information
    errCode = mal.NewUInteger(uint32(mal.ERROR_UNKNOWN))
    errExtraInfo = mal.NewString(e.Error())
    errIsAbstract = false
  }
  // encode parameters
  if body.EncodeParameter(errCode) != nil {
    return errors.New("Fatal error in error handling code")
  } else if body.EncodeLastParameter(errExtraInfo, errIsAbstract) != nil {
    return errors.New("Fatal error in error handling code")
  }

  // interaction call
  var err error
  err = transaction.Reply(body, true)
  if err != nil {
    return err
  }
  return nil
}

// generated code for operation update
type UpdateHelper struct {
  transaction malapi.SubmitTransaction
}
func NewUpdateHelper(transaction malapi.Transaction) (*UpdateHelper, error) {
  iptransaction, ok := transaction.(malapi.SubmitTransaction)
  if !ok {
    return nil, errors.New("Unexpected transaction type")
  }
helper := &UpdateHelper{iptransaction}
  return helper, nil
}
func (receiver *UpdateHelper) Ack() error {
  transaction := receiver.transaction

  // interaction call
  err := transaction.Ack(nil, false)
  if err != nil {
    return err
  }

  return nil
}

func (receiver *UpdateHelper) ReturnError(e error) error {
  transaction := receiver.transaction
  body := transaction.NewBody()
  var errCode *mal.UInteger
  var errExtraInfo mal.Element
  var errIsAbstract bool
  malErr, ok := e.(*malapi.MalError)
  if ok {
    errCode = &malErr.Code
    errExtraInfo = malErr.ExtraInfo
    errIsAbstract = false
  } else {
    // return an UNKNOWN error with a String information
    errCode = mal.NewUInteger(uint32(mal.ERROR_UNKNOWN))
    errExtraInfo = mal.NewString(e.Error())
    errIsAbstract = false
  }
  // encode parameters
  if body.EncodeParameter(errCode) != nil {
    return errors.New("Fatal error in error handling code")
  } else if body.EncodeLastParameter(errExtraInfo, errIsAbstract) != nil {
    return errors.New("Fatal error in error handling code")
  }

  // interaction call
  var err error
  err = transaction.Ack(body, true)
  if err != nil {
    return err
  }
  return nil
}

// generated code for operation delete
type DeleteHelper struct {
  transaction malapi.RequestTransaction
}
func NewDeleteHelper(transaction malapi.Transaction) (*DeleteHelper, error) {
  iptransaction, ok := transaction.(malapi.RequestTransaction)
  if !ok {
    return nil, errors.New("Unexpected transaction type")
  }
helper := &DeleteHelper{iptransaction}
  return helper, nil
}
func (receiver *DeleteHelper) Reply(deletedObjInstIds *mal.LongList) error {
  transaction := receiver.transaction
  // create a body for the interaction call
  body := transaction.NewBody()
  // encode parameters
  err := body.EncodeLastParameter(deletedObjInstIds, false)
  if err != nil {
    return nil
  }

  // interaction call
  err = transaction.Reply(body, false)
  if err != nil {
    return err
  }

  return nil
}

func (receiver *DeleteHelper) ReturnError(e error) error {
  transaction := receiver.transaction
  body := transaction.NewBody()
  var errCode *mal.UInteger
  var errExtraInfo mal.Element
  var errIsAbstract bool
  malErr, ok := e.(*malapi.MalError)
  if ok {
    errCode = &malErr.Code
    errExtraInfo = malErr.ExtraInfo
    errIsAbstract = false
  } else {
    // return an UNKNOWN error with a String information
    errCode = mal.NewUInteger(uint32(mal.ERROR_UNKNOWN))
    errExtraInfo = mal.NewString(e.Error())
    errIsAbstract = false
  }
  // encode parameters
  if body.EncodeParameter(errCode) != nil {
    return errors.New("Fatal error in error handling code")
  } else if body.EncodeLastParameter(errExtraInfo, errIsAbstract) != nil {
    return errors.New("Fatal error in error handling code")
  }

  // interaction call
  var err error
  err = transaction.Reply(body, true)
  if err != nil {
    return err
  }
  return nil
}
