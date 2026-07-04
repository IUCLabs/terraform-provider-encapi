// Command terraform-provider-encapi is the Terraform/OpenTofu provider for
// encapi: it manages Puppet ENC roles, statuses (environments), and per-node
// role assignments.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"git.unkin.net/unkin/terraform-provider-encapi/internal/provider"
)

var version = "0.0.1"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "enable debug mode")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "git.unkin.net/unkin/encapi",
		Debug:   debug,
	}

	if err := providerserver.Serve(context.Background(), provider.New(version), opts); err != nil {
		log.Fatal(err)
	}
}
