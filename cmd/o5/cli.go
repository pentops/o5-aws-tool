package main

import (
	"github.com/pentops/o5-aws-tool/cli/aws"
	"github.com/pentops/o5-aws-tool/cli/dante"
	"github.com/pentops/o5-aws-tool/cli/deployer"
	"github.com/pentops/o5-aws-tool/cli/registry"
	"github.com/pentops/runner/commander"
)

var Version string

func main() {

	cmdGroup := commander.NewCommandSet()

	cmdGroup.Add("aws", aws.CommandSet())
	cmdGroup.Add("o5", deployer.O5CommandSet())
	cmdGroup.Add("registry", registry.RegistryCommandSet())
	cmdGroup.Add("dante", dante.DanteCommandSet())

	cmdGroup.RunMain("o5-aws-tool", Version)

}
