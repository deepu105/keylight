package main

import (
	"context"
	"fmt"
	"log"
	"math"

	keylight "github.com/endocrimes/keylight-go"
	texttable "github.com/syohex/go-texttable"
)

// discover all Keylights in your network and return a channel with results
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

// convert powerstate from int to string
func getPowerState(state int) string {
	if state == 0 {
		return "off"
	}
	return "on"
}

// convert power state from string to int
func getPowerStateInt(state string) int {
	if state == "off" {
		return 0
	}
	return 1
}

// create a pretty table
func createTable() *texttable.TextTable {
	tbl := &texttable.TextTable{}
	tbl.SetHeader("Name", "Power State", "Brightness", "Temperature", "Address")
	return tbl
}

// add a row to Table to pretty print
func addTableRow(tbl *texttable.TextTable, d *keylight.Device, group *keylight.LightGroup) {
	tbl.AddRow(
		d.Name,
		getPowerState(group.Lights[0].On),
		fmt.Sprintf("%d", group.Lights[0].Brightness),
		fmt.Sprintf("%d (%d K)", group.Lights[0].Temperature, convertToKelvin(group.Lights[0].Temperature)),
		fmt.Sprintf("%s:%d", d.DNSAddr, d.Port),
	)
}

// draw pretty table
func drawTable(tbl *texttable.TextTable, count int) {
	if count > 0 {
		fmt.Println(tbl.Draw())
	} else {
		fmt.Println("Either timedout or no lights found! Try again")
	}
}

// converts the Elgato API temperatures to Kelvin.
func convertToKelvin(temp int) int {
	return int(math.Round(1000000 * math.Pow(float64(temp), -1)))
}

// converts Kelvin temperatures to those of the Elgato API.
func convertFromKelvin(kelvin int) int {
	return int(math.Floor(987007 * math.Pow(float64(kelvin), -0.999)))
}
