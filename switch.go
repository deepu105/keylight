package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	keylight "github.com/endocrimes/keylight-go"
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
			Usage:   "Switch light on. If not provided the light power state will be toggled based on last state",
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
		&cli.StringFlag{
			Name:    "preset",
			Aliases: []string{"p"},
			Usage:   fmt.Sprintf("Switch on and set a preset temperature and brigtness. Values: %v", presetKeys()),
		},
		&timeoutFlag,
	},
	Action: switchAction,
}

var presets = map[string]keylight.Light{
	"warm":       {Temperature: 3000, Brightness: 10},
	"warm-50":    {Temperature: 3000, Brightness: 50},
	"warm-100":   {Temperature: 3000, Brightness: 100},
	"cool":       {Temperature: 7000, Brightness: 10},
	"cool-50":    {Temperature: 7000, Brightness: 50},
	"cool-100":   {Temperature: 7000, Brightness: 100},
	"normal":     {Temperature: 5000, Brightness: 10},
	"normal-50":  {Temperature: 5000, Brightness: 50},
	"normal-100": {Temperature: 5000, Brightness: 100},
}

func presetKeys() []string {
	keys := make([]string, 0, len(presets))
	for k := range presets {
		keys = append(keys, k)
	}
	return keys
}

func switchAction(c *cli.Context) error {
	lightID := c.String("light")
	toggleSwitchOn := c.Bool("on")
	brightness := c.Int("brightness")
	temperature := c.Int("temperature")
	preset := c.String("preset")
	timeout := time.Duration(c.Int("timeout")) * time.Second

	if preset != "" {
		toggleSwitchOn = true
		brightness = presets[preset].Brightness
		temperature = presets[preset].Temperature
	}

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
