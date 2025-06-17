package ges

import (
	"context"
	"fmt"

	"github.com/pentops/o5-aws-tool/gen/o5/ges/v1/ges"
	"github.com/pentops/o5-aws-tool/libo5"
	"github.com/pentops/runner/commander"
)

func CommandSet() *commander.CommandSet {
	eventSet := commander.NewCommandSet()
	eventSet.Add("ls", commander.NewCommand(runEventLS))

	upsertSet := commander.NewCommandSet()

	set := commander.NewCommandSet()
	set.Add("e", eventSet)
	set.Add("u", upsertSet)
	return set
}

func runEventLS(ctx context.Context, cfg struct {
	libo5.APIConfig
},
) error {
	client := cfg.APIClient()
	queryClient := ges.NewQueryService(client)

	req := &ges.EventsListRequest{}
	if err := libo5.Paged(ctx,
		req,
		queryClient.EventsList,
		func(event *ges.Event) error {
			fmt.Println("=====")
			fmt.Printf("GrpcService: %s\n", event.GrpcService)
			return nil
		}); err != nil {
		return err
	}

	return nil
}
