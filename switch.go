package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	texttable "github.com/syohex/go-texttable"
	"github.com/urfave/cli/v2"
)

const toggleOn = "on"
const toggleOff = "off"

var switchCommand = &cli.Command{
	Name:    "switch",
	Aliases: []string{"s"},
	Usage:   "Switch on/off lights",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "light",
			Aliases: []string{"l"},
			Value:   "all",
			Usage:   "ID, example E859, for the light to control. If not provided all lights will be updated",
		},
		&cli.BoolFlag{Name: toggleOn, Usage: "Toggle light on"},
		&cli.BoolFlag{Name: toggleOff, Usage: "Toggle light off"},
		&cli.IntFlag{
			Name:    "brightness",
			Aliases: []string{"b"},
			Value:   10,
			Usage:   "Set brightness of the lights (0 to 100)",
		},
		&cli.IntFlag{
			Name:    "temperature",
			Value:   3000, // minimum of 331 (~3000k)
			Aliases: []string{"t"},
			Usage:   "Set temperature of the lights in kelvin (3000 to 7000)",
		},
		&cli.IntFlag{
			Name:  "timeout",
			Value: 2, // 2 seconds
			Usage: "Timeout in seconds",
		},
	},
	Action: switchAction,
}

func switchAction(c *cli.Context) error {
	timeout := time.Duration(c.Int("timeout")) * time.Second

	toggleSwitch := toggleOn
	if c.Bool("off") {
		toggleSwitch = toggleOff
	}
	lightID := c.String("light")

	fmt.Printf("Command: Switch - lights: %s, toggle: %s \n", lightID, toggleSwitch)

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
			if lightID != "all" && !strings.HasSuffix(device.Name, lightID) {
				fmt.Printf("Light ID ending with %s not found\n", lightID)
				return nil
			}

			group, err := device.FetchLightGroup(context.Background())
			if err != nil {
				log.Println("failed to retrieve light group: ", err.Error())
				return err
			}

			newGroup := group.Copy()

			for _, light := range newGroup.Lights {
				light.On = getPowerStateInt(toggleSwitch)
				light.Brightness = c.Int("brightness")
				light.Temperature = int(math.Floor(987007 * math.Pow(float64(c.Int("temperature")), -0.999)))
			}

			_, err = device.UpdateLightGroup(context.Background(), newGroup)

			if err != nil {
				log.Println("failed to update light group: ", err.Error())
				return err
			}
			addTableRow(tbl, device, newGroup)
			count++
		case <-time.After(timeout):
			drawTable(tbl, count)
			return nil
		}
	}
}
