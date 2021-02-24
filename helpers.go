package main

import (
	"context"
	"fmt"
	"log"
	"math"

	keylight "github.com/endocrimes/keylight-go"
	texttable "github.com/syohex/go-texttable"
)

func discoverLights() (<-chan *keylight.Device, error) {
	discovery, err := keylight.NewDiscovery()
	if err != nil {
		log.Println("failed to initialize keylight discovery: ", err.Error())
		return nil, err
	}

	go func() {
		err := discovery.Run(context.Background())
		if err != nil {
			log.Fatalln("Failed to discover lights: ", err.Error())
		}
	}()

	return discovery.ResultsCh(), nil
}

func getPowerState(state int) string {
	if state == 0 {
		return "off"
	}
	return "on"
}

func getPowerStateInt(state string) int {
	if state == "off" {
		return 0
	}
	return 1
}

func addTableRow(tbl *texttable.TextTable, d *keylight.Device, group *keylight.LightGroup) {
	tbl.AddRow(
		d.Name,
		getPowerState(group.Lights[0].On),
		fmt.Sprintf("%d", group.Lights[0].Brightness),
		fmt.Sprintf("%d (%d K)", group.Lights[0].Temperature, int(math.Round(1000000*math.Pow(float64(group.Lights[0].Temperature), -1)))),
		fmt.Sprintf("%s:%d", d.DNSAddr, d.Port),
	)
}

func drawTable(tbl *texttable.TextTable, count int) {
	if count > 0 {
		fmt.Println(tbl.Draw())
	} else {
		fmt.Println("Either timedout or no lights found! Try again")
	}
}
