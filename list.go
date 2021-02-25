package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/urfave/cli/v2"
)

var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "Discover and list available lights",
	Flags: []cli.Flag{
		&timeoutFlag,
	},
	Action: listAction,
}

func listAction(c *cli.Context) error {
	timeout := time.Duration(c.Int("timeout")) * time.Second

	fmt.Println("Command: List - lights: all")

	devicesCh, err := discoverLights()
	if err != nil {
		return err
	}
	tbl := createTable()

	count := 0
	for {
		select {
		case device := <-devicesCh:
			if device == nil {
				return nil
			}
			group, err := device.FetchLightGroup(context.Background())
			if err != nil {
				log.Println("failed to retrieve light group: ", err.Error())
				return err
			}
			addTableRow(tbl, device, group)
			count++
		case <-time.After(timeout):
			drawTable(tbl, count)
			return nil
		}
	}
}
