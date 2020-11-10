package archive

import (
  "errors"
  "github.com/CNES/ccsdsmo-malgo/mal"
  malapi "github.com/CNES/ccsdsmo-malgo/mal/api"
  "github.com/CNES/ccsdsmo-malgo/com"
)

var Cctx *malapi.ClientContext
func Init(cctxin *malapi.ClientContext) error {
  if cctxin == nil {
    return errors.New("Illegal nil client context in Init")
  }
  Cctx = cctxin
  return nil
}

// generated code for operation retrieve
type RetrieveOperation struct {
  op malapi.InvokeOperation
}
func NewRetrieveOperation(providerURI *mal.URI) (*RetrieveOperation, error) {
  op := Cctx.NewInvokeOperation(providerURI, com.AREA_NUMBER, com.AREA_VERSION, SERVICE_NUMBER, RETRIEVE_OPERATION_NUMBER)
  consumer := &RetrieveOperation{op}
  return consumer, nil
}
func (receiver *RetrieveOperation) Invoke(objType *com.ObjectType, domain *mal.IdentifierList, objInstIds *mal.LongList) error {
  // create a body for the operation call
  body := receiver.op.NewBody()
  // encode in parameters
  err := body.EncodeParameter(objType)
  if err != nil {
    return err
  }
  err = body.EncodeParameter(domain)
  if err != nil {
    return err
  }
  err = body.EncodeLastParameter(objInstIds, false)
  if err != nil {
    return err
  }

  // operation call
  resp, err := receiver.op.Invoke(body)
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return err
    }
    err = receiver.decodeError(resp, err)
    return err
  }

  // decode out parameters
  return nil
}
func (receiver *RetrieveOperation) GetResponse() (*ArchiveDetailsList, mal.ElementList, error) {

  // operation call
  resp, err := receiver.op.GetResponse()
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return nil, nil, err
    }
    err = receiver.decodeError(resp, err)
    return nil, nil, err
  }

  // decode out parameters
  outElem_objDetails, err := resp.DecodeParameter(NullArchiveDetailsList)
  if err != nil {
    return nil, nil, err
  }
  outParam_objDetails, ok := outElem_objDetails.(*ArchiveDetailsList)
  if !ok {
    err = errors.New("unexpected type for parameter objDetails")
    return nil, nil, err
  }

  outElem_objBodies, err := resp.DecodeLastParameter(nil, true)
  if err != nil {
    return nil, nil, err
  }
  outParam_objBodies, ok := outElem_objBodies.(mal.ElementList)
  if !ok {
    if outElem_objBodies == mal.NullElement {
      outParam_objBodies = mal.NullElementList
    } else {
      err = errors.New("unexpected type for parameter objBodies")
      return nil, nil, err
    }
  }

  return outParam_objDetails, outParam_objBodies, nil
}

func (receiver *RetrieveOperation) decodeError(resp *mal.Message, e error) error {
  // decode err parameters
  outElem_code, err := resp.DecodeParameter(mal.NullUInteger)
  if err != nil {
    return err
  }
  outParam_code, ok := outElem_code.(*mal.UInteger)
  if ! ok {
    err = errors.New("unexpected type for parameter code")
    return err
  }
  nullValue := mal.NullElement
  errIsAbstract := false
  switch *outParam_code {
  case com.ERROR_INVALID:
    nullValue = mal.NullString
  case mal.ERROR_UNKNOWN:
    nullValue = mal.NullUIntegerList
  default:
    nullValue = mal.NullString
  }
  outElem_extraInfo, err := resp.DecodeLastParameter(nullValue, errIsAbstract)
  if err != nil {
    return err
  }
  return malapi.NewMalError(*outParam_code, outElem_extraInfo)
}

// generated code for operation query
type QueryOperation struct {
  op malapi.ProgressOperation
}
func NewQueryOperation(providerURI *mal.URI) (*QueryOperation, error) {
  op := Cctx.NewProgressOperation(providerURI, com.AREA_NUMBER, com.AREA_VERSION, SERVICE_NUMBER, QUERY_OPERATION_NUMBER)
  consumer := &QueryOperation{op}
  return consumer, nil
}
func (receiver *QueryOperation) Progress(returnBody *mal.Boolean, objType *com.ObjectType, archiveQuery *ArchiveQueryList, queryFilter QueryFilterList) error {
  // create a body for the operation call
  body := receiver.op.NewBody()
  // encode in parameters
  err := body.EncodeParameter(returnBody)
  if err != nil {
    return err
  }
  err = body.EncodeParameter(objType)
  if err != nil {
    return err
  }
  err = body.EncodeParameter(archiveQuery)
  if err != nil {
    return err
  }
  err = body.EncodeLastParameter(queryFilter, true)
  if err != nil {
    return err
  }

  // operation call
  resp, err := receiver.op.Progress(body)
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return err
    }
    err = receiver.decodeError(resp, err)
    return err
  }

  // decode out parameters
  return nil
}
func (receiver *QueryOperation) GetUpdate() (*com.ObjectType, *mal.IdentifierList, *ArchiveDetailsList, mal.ElementList, error) {

  // operation call
  resp, err := receiver.op.GetUpdate()
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return nil, nil, nil, nil, err
    }
    err = receiver.decodeError(resp, err)
    return nil, nil, nil, nil, err
  }
  if resp == nil {
    return nil, nil, nil, nil, err
  }

  // decode out parameters
  outElem_objType, err := resp.DecodeParameter(com.NullObjectType)
  if err != nil {
    return nil, nil, nil, nil, err
  }
  outParam_objType, ok := outElem_objType.(*com.ObjectType)
  if !ok {
    err = errors.New("unexpected type for parameter objType")
    return nil, nil, nil, nil, err
  }

  outElem_domain, err := resp.DecodeParameter(mal.NullIdentifierList)
  if err != nil {
    return nil, nil, nil, nil, err
  }
  outParam_domain, ok := outElem_domain.(*mal.IdentifierList)
  if !ok {
    err = errors.New("unexpected type for parameter domain")
    return nil, nil, nil, nil, err
  }

  outElem_objDetails, err := resp.DecodeParameter(NullArchiveDetailsList)
  if err != nil {
    return nil, nil, nil, nil, err
  }
  outParam_objDetails, ok := outElem_objDetails.(*ArchiveDetailsList)
  if !ok {
    err = errors.New("unexpected type for parameter objDetails")
    return nil, nil, nil, nil, err
  }

  outElem_objBodies, err := resp.DecodeLastParameter(nil, true)
  if err != nil {
    return nil, nil, nil, nil, err
  }
  outParam_objBodies, ok := outElem_objBodies.(mal.ElementList)
  if !ok {
    if outElem_objBodies == mal.NullElement {
      outParam_objBodies = mal.NullElementList
    } else {
      err = errors.New("unexpected type for parameter objBodies")
      return nil, nil, nil, nil, err
    }
  }

  return outParam_objType, outParam_domain, outParam_objDetails, outParam_objBodies, nil
}
func (receiver *QueryOperation) GetResponse() (*com.ObjectType, *mal.IdentifierList, *ArchiveDetailsList, mal.ElementList, error) {

  // operation call
  resp, err := receiver.op.GetResponse()
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return nil, nil, nil, nil, err
    }
    err = receiver.decodeError(resp, err)
    return nil, nil, nil, nil, err
  }

  // decode out parameters
  outElem_objType, err := resp.DecodeParameter(com.NullObjectType)
  if err != nil {
    return nil, nil, nil, nil, err
  }
  outParam_objType, ok := outElem_objType.(*com.ObjectType)
  if !ok {
    err = errors.New("unexpected type for parameter objType")
    return nil, nil, nil, nil, err
  }

  outElem_domain, err := resp.DecodeParameter(mal.NullIdentifierList)
  if err != nil {
    return nil, nil, nil, nil, err
  }
  outParam_domain, ok := outElem_domain.(*mal.IdentifierList)
  if !ok {
    err = errors.New("unexpected type for parameter domain")
    return nil, nil, nil, nil, err
  }

  outElem_objDetails, err := resp.DecodeParameter(NullArchiveDetailsList)
  if err != nil {
    return nil, nil, nil, nil, err
  }
  outParam_objDetails, ok := outElem_objDetails.(*ArchiveDetailsList)
  if !ok {
    err = errors.New("unexpected type for parameter objDetails")
    return nil, nil, nil, nil, err
  }

  outElem_objBodies, err := resp.DecodeLastParameter(nil, true)
  if err != nil {
    return nil, nil, nil, nil, err
  }
  outParam_objBodies, ok := outElem_objBodies.(mal.ElementList)
  if !ok {
    if outElem_objBodies == mal.NullElement {
      outParam_objBodies = mal.NullElementList
    } else {
      err = errors.New("unexpected type for parameter objBodies")
      return nil, nil, nil, nil, err
    }
  }

  return outParam_objType, outParam_domain, outParam_objDetails, outParam_objBodies, nil
}

func (receiver *QueryOperation) decodeError(resp *mal.Message, e error) error {
  // decode err parameters
  outElem_code, err := resp.DecodeParameter(mal.NullUInteger)
  if err != nil {
    return err
  }
  outParam_code, ok := outElem_code.(*mal.UInteger)
  if ! ok {
    err = errors.New("unexpected type for parameter code")
    return err
  }
  nullValue := mal.NullElement
  errIsAbstract := false
  switch *outParam_code {
  case com.ERROR_INVALID:
    nullValue = mal.NullUIntegerList
  default:
    nullValue = mal.NullString
  }
  outElem_extraInfo, err := resp.DecodeLastParameter(nullValue, errIsAbstract)
  if err != nil {
    return err
  }
  return malapi.NewMalError(*outParam_code, outElem_extraInfo)
}

// generated code for operation count
type CountOperation struct {
  op malapi.InvokeOperation
}
func NewCountOperation(providerURI *mal.URI) (*CountOperation, error) {
  op := Cctx.NewInvokeOperation(providerURI, com.AREA_NUMBER, com.AREA_VERSION, SERVICE_NUMBER, COUNT_OPERATION_NUMBER)
  consumer := &CountOperation{op}
  return consumer, nil
}
func (receiver *CountOperation) Invoke(objType *com.ObjectType, archiveQuery *ArchiveQueryList, queryFilter QueryFilterList) error {
  // create a body for the operation call
  body := receiver.op.NewBody()
  // encode in parameters
  err := body.EncodeParameter(objType)
  if err != nil {
    return err
  }
  err = body.EncodeParameter(archiveQuery)
  if err != nil {
    return err
  }
  err = body.EncodeLastParameter(queryFilter, true)
  if err != nil {
    return err
  }

  // operation call
  resp, err := receiver.op.Invoke(body)
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return err
    }
    err = receiver.decodeError(resp, err)
    return err
  }

  // decode out parameters
  return nil
}
func (receiver *CountOperation) GetResponse() (*mal.LongList, error) {

  // operation call
  resp, err := receiver.op.GetResponse()
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return nil, err
    }
    err = receiver.decodeError(resp, err)
    return nil, err
  }

  // decode out parameters
  outElem_counts, err := resp.DecodeLastParameter(mal.NullLongList, false)
  if err != nil {
    return nil, err
  }
  outParam_counts, ok := outElem_counts.(*mal.LongList)
  if !ok {
    err = errors.New("unexpected type for parameter counts")
    return nil, err
  }

  return outParam_counts, nil
}

func (receiver *CountOperation) decodeError(resp *mal.Message, e error) error {
  // decode err parameters
  outElem_code, err := resp.DecodeParameter(mal.NullUInteger)
  if err != nil {
    return err
  }
  outParam_code, ok := outElem_code.(*mal.UInteger)
  if ! ok {
    err = errors.New("unexpected type for parameter code")
    return err
  }
  nullValue := mal.NullElement
  errIsAbstract := false
  switch *outParam_code {
  case com.ERROR_INVALID:
    nullValue = mal.NullUIntegerList
  default:
    nullValue = mal.NullString
  }
  outElem_extraInfo, err := resp.DecodeLastParameter(nullValue, errIsAbstract)
  if err != nil {
    return err
  }
  return malapi.NewMalError(*outParam_code, outElem_extraInfo)
}

// generated code for operation store
type StoreOperation struct {
  op malapi.RequestOperation
}
func NewStoreOperation(providerURI *mal.URI) (*StoreOperation, error) {
  op := Cctx.NewRequestOperation(providerURI, com.AREA_NUMBER, com.AREA_VERSION, SERVICE_NUMBER, STORE_OPERATION_NUMBER)
  consumer := &StoreOperation{op}
  return consumer, nil
}
func (receiver *StoreOperation) Request(returnObjInstIds *mal.Boolean, objType *com.ObjectType, domain *mal.IdentifierList, objDetails *ArchiveDetailsList, objBodies mal.ElementList) (*mal.LongList, error) {
  // create a body for the operation call
  body := receiver.op.NewBody()
  // encode in parameters
  err := body.EncodeParameter(returnObjInstIds)
  if err != nil {
    return nil, err
  }
  err = body.EncodeParameter(objType)
  if err != nil {
    return nil, err
  }
  err = body.EncodeParameter(domain)
  if err != nil {
    return nil, err
  }
  err = body.EncodeParameter(objDetails)
  if err != nil {
    return nil, err
  }
  err = body.EncodeLastParameter(objBodies, true)
  if err != nil {
    return nil, err
  }

  // operation call
  resp, err := receiver.op.Request(body)
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return nil, err
    }
    err = receiver.decodeError(resp, err)
    return nil, err
  }

  // decode out parameters
  outElem_objInstIds, err := resp.DecodeLastParameter(mal.NullLongList, false)
  if err != nil {
    return nil, err
  }
  outParam_objInstIds, ok := outElem_objInstIds.(*mal.LongList)
  if !ok {
    err = errors.New("unexpected type for parameter objInstIds")
    return nil, err
  }

  return outParam_objInstIds, nil
}

func (receiver *StoreOperation) decodeError(resp *mal.Message, e error) error {
  // decode err parameters
  outElem_code, err := resp.DecodeParameter(mal.NullUInteger)
  if err != nil {
    return err
  }
  outParam_code, ok := outElem_code.(*mal.UInteger)
  if ! ok {
    err = errors.New("unexpected type for parameter code")
    return err
  }
  nullValue := mal.NullElement
  errIsAbstract := false
  switch *outParam_code {
  case com.ERROR_DUPLICATE:
    nullValue = mal.NullUIntegerList
  case com.ERROR_INVALID:
    nullValue = mal.NullUIntegerList
  default:
    nullValue = mal.NullString
  }
  outElem_extraInfo, err := resp.DecodeLastParameter(nullValue, errIsAbstract)
  if err != nil {
    return err
  }
  return malapi.NewMalError(*outParam_code, outElem_extraInfo)
}

// generated code for operation update
type UpdateOperation struct {
  op malapi.SubmitOperation
}
func NewUpdateOperation(providerURI *mal.URI) (*UpdateOperation, error) {
  op := Cctx.NewSubmitOperation(providerURI, com.AREA_NUMBER, com.AREA_VERSION, SERVICE_NUMBER, UPDATE_OPERATION_NUMBER)
  consumer := &UpdateOperation{op}
  return consumer, nil
}
func (receiver *UpdateOperation) Submit(objType *com.ObjectType, domain *mal.IdentifierList, objDetails *ArchiveDetailsList, objBodies mal.ElementList) error {
  // create a body for the operation call
  body := receiver.op.NewBody()
  // encode in parameters
  err := body.EncodeParameter(objType)
  if err != nil {
    return err
  }
  err = body.EncodeParameter(domain)
  if err != nil {
    return err
  }
  err = body.EncodeParameter(objDetails)
  if err != nil {
    return err
  }
  err = body.EncodeLastParameter(objBodies, true)
  if err != nil {
    return err
  }

  // operation call
  resp, err := receiver.op.Submit(body)
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return err
    }
    err = receiver.decodeError(resp, err)
    return err
  }

  // decode out parameters
  return nil
}

func (receiver *UpdateOperation) decodeError(resp *mal.Message, e error) error {
  // decode err parameters
  outElem_code, err := resp.DecodeParameter(mal.NullUInteger)
  if err != nil {
    return err
  }
  outParam_code, ok := outElem_code.(*mal.UInteger)
  if ! ok {
    err = errors.New("unexpected type for parameter code")
    return err
  }
  nullValue := mal.NullElement
  errIsAbstract := false
  switch *outParam_code {
  case mal.ERROR_UNKNOWN:
    nullValue = mal.NullUIntegerList
  case com.ERROR_INVALID:
    nullValue = mal.NullUIntegerList
  default:
    nullValue = mal.NullString
  }
  outElem_extraInfo, err := resp.DecodeLastParameter(nullValue, errIsAbstract)
  if err != nil {
    return err
  }
  return malapi.NewMalError(*outParam_code, outElem_extraInfo)
}

// generated code for operation delete
type DeleteOperation struct {
  op malapi.RequestOperation
}
func NewDeleteOperation(providerURI *mal.URI) (*DeleteOperation, error) {
  op := Cctx.NewRequestOperation(providerURI, com.AREA_NUMBER, com.AREA_VERSION, SERVICE_NUMBER, DELETE_OPERATION_NUMBER)
  consumer := &DeleteOperation{op}
  return consumer, nil
}
func (receiver *DeleteOperation) Request(objType *com.ObjectType, domain *mal.IdentifierList, objInstIds *mal.LongList) (*mal.LongList, error) {
  // create a body for the operation call
  body := receiver.op.NewBody()
  // encode in parameters
  err := body.EncodeParameter(objType)
  if err != nil {
    return nil, err
  }
  err = body.EncodeParameter(domain)
  if err != nil {
    return nil, err
  }
  err = body.EncodeLastParameter(objInstIds, false)
  if err != nil {
    return nil, err
  }

  // operation call
  resp, err := receiver.op.Request(body)
  if err != nil {
    // Verify if an error occurs during the operation
    if !resp.IsErrorMessage {
      return nil, err
    }
    err = receiver.decodeError(resp, err)
    return nil, err
  }

  // decode out parameters
  outElem_deletedObjInstIds, err := resp.DecodeLastParameter(mal.NullLongList, false)
  if err != nil {
    return nil, err
  }
  outParam_deletedObjInstIds, ok := outElem_deletedObjInstIds.(*mal.LongList)
  if !ok {
    err = errors.New("unexpected type for parameter deletedObjInstIds")
    return nil, err
  }

  return outParam_deletedObjInstIds, nil
}

func (receiver *DeleteOperation) decodeError(resp *mal.Message, e error) error {
  // decode err parameters
  outElem_code, err := resp.DecodeParameter(mal.NullUInteger)
  if err != nil {
    return err
  }
  outParam_code, ok := outElem_code.(*mal.UInteger)
  if ! ok {
    err = errors.New("unexpected type for parameter code")
    return err
  }
  nullValue := mal.NullElement
  errIsAbstract := false
  switch *outParam_code {
  case mal.ERROR_UNKNOWN:
    nullValue = mal.NullUIntegerList
  case com.ERROR_INVALID:
    nullValue = mal.NullString
  default:
    nullValue = mal.NullString
  }
  outElem_extraInfo, err := resp.DecodeLastParameter(nullValue, errIsAbstract)
  if err != nil {
    return err
  }
  return malapi.NewMalError(*outParam_code, outElem_extraInfo)
}
