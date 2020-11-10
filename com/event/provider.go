package event

import (
  "errors"
  "github.com/CNES/ccsdsmo-malgo/mal"
  malapi "github.com/CNES/ccsdsmo-malgo/mal/api"
  "github.com/CNES/ccsdsmo-malgo/com"
)


// service provider internal interface
type ProviderInterface interface {
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
func MonitorEventDummy() error {
  return errors.New(string(com.AREA_NAME))
}
