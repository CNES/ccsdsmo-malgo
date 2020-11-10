package archive

import (
  "github.com/CNES/ccsdsmo-malgo/mal"
)

const (
  // standard service identifiers
  SERVICE_NUMBER mal.UShort = 2
  SERVICE_NAME = mal.Identifier("Archive")

  // standard operation identifiers
  RETRIEVE_OPERATION_NUMBER mal.UShort = 1
  QUERY_OPERATION_NUMBER mal.UShort = 2
  COUNT_OPERATION_NUMBER mal.UShort = 3
  STORE_OPERATION_NUMBER mal.UShort = 4
  UPDATE_OPERATION_NUMBER mal.UShort = 5
  DELETE_OPERATION_NUMBER mal.UShort = 6
)

