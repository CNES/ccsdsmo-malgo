package activitytracking

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
