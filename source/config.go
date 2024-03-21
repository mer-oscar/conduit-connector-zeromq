package source

import "github.com/mer-oscar/conduit-connector-zeromq/common"

//go:generate paramgen -output=config_paramgen.go Config
type Config struct {
	common.Config
	// PortBindings is a comma separated list of ports that we wish to bind to
	PortBindings string `json:"portBindings"`
}
