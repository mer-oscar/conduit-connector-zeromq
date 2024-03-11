package main

import (
	sdk "github.com/conduitio/conduit-connector-sdk"

	zeromq "github.com/mer-oscar/conduit-connector-zeromq"
)

func main() {
	sdk.Serve(zeromq.Connector)
}
