package main

import (
	"github.com/pentops/o5-aws-tool/cli/aws"
	"github.com/pentops/o5-aws-tool/cli/builds"
	"github.com/pentops/o5-aws-tool/cli/dante"
	"github.com/pentops/o5-aws-tool/cli/deployer"
	"github.com/pentops/o5-aws-tool/cli/logs"
	"github.com/pentops/o5-aws-tool/cli/rds"
	"github.com/pentops/runner/commander"
)

var Version string

func main() {

	cmdGroup := commander.NewCommandSet()

	cmdGroup.Add("aws", aws.CommandSet())
	cmdGroup.Add("logs", logs.CommandSet())
	cmdGroup.Add("o5", deployer.O5CommandSet())
	cmdGroup.Add("builds", builds.BuildsCommandSet())
	cmdGroup.Add("dante", dante.DanteCommandSet())
	cmdGroup.Add("rds", rds.CommandSet())

	cmdGroup.RunMain("o5-aws-tool", Version)

}
