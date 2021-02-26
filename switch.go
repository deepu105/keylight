package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

const toggleOn = "on"

var timeoutFlag = cli.IntFlag{
	Name:  "timeout",
	Value: 2, // 2 seconds
	Usage: "Timeout for light discovery in seconds",
}

var switchCommand = &cli.Command{
	Name:    "switch",
	Aliases: []string{"s"},
	Usage:   "Toggle the light switch",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "light",
			Aliases: []string{"l"},
			Value:   "all",
			Usage:   "ID, example E859, for the light to control. If not provided all lights will be updated",
		},
		&cli.BoolFlag{
			Name:    toggleOn,
			Aliases: []string{"o"},
			Usage:   "Switch light on. If not provided the light power state will be toggled",
		},
		&cli.IntFlag{
			Name:    "brightness",
			Aliases: []string{"b"},
			Value:   -1,
			Usage:   "Set brightness of the lights (0 to 100)",
		},
		&cli.IntFlag{
			Name:    "temperature",
			Value:   -1,
			Aliases: []string{"t"},
			Usage:   "Set temperature of the lights in kelvin (3000 to 7000)",
		},
		&timeoutFlag,
	},
	Action: switchAction,
}

func switchAction(c *cli.Context) error {
	timeout := time.Duration(c.Int("timeout")) * time.Second
	brightness := c.Int("brightness")
	temperature := c.Int("temperature")
	toggleSwitchOn := c.Bool("on")
	lightID := c.String("light")

	fmt.Printf("Command: Switch - lights: %s, toggleOn: %v \n", lightID, toggleSwitchOn)

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
				light.On = togglePowerState(toggleSwitchOn, light.On)
				if brightness > -1 {
					light.Brightness = brightness
				}
				if temperature > -1 {
					light.Temperature = convertFromKelvin(temperature)
				}
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
