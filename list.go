package main

import (
	"context"
	"fmt"
	"log"
	"time"

	texttable "github.com/syohex/go-texttable"
	"github.com/urfave/cli/v2"
)

var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "Discover and list available lights",
	Action:  listAction,
}

func listAction(c *cli.Context) error {
	timeout := time.Duration(c.Int("timeout")) * time.Second

	fmt.Println("Command: List - lights: all")

	devicesCh, err := discoverLights()
	if err != nil {
		return err
	}
	tbl := &texttable.TextTable{}
	tbl.SetHeader("Name", "Power State", "Brightness", "Temperature", "Address")

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
