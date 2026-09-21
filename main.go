package main

import (
	"context"
	"log"

	"github.com/achachw/terraform-provider-servicenow/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version = "dev"

func main() {
	opts := providerserver.ServeOpts{
		Address:         "registry.terraform.io/achachw/servicenow",
		ProtocolVersion: 6,
	}

	if err := providerserver.Serve(context.Background(), provider.New(version), opts); err != nil {
		log.Fatal(err)
	}
}
