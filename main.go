package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	texttable "github.com/syohex/go-texttable"
	"github.com/urfave/cli/v2"
)

const toggleOn = "on"
const toggleOff = "off"

// Light holds Keylight structure
type Light struct {
	Name        string
	HostName    string
	Port        int
	Brightness  int
	Temperature int
	Power       int
}

const timeout = 2 * time.Second

func main() {
	app := &cli.App{
		Name:  "keylight",
		Usage: "A simple CLI to control your Elgato Keylights",
		Commands: []*cli.Command{
			{
				Name:    "switch",
				Aliases: []string{"s"},
				Usage:   "Switch on/off lights",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "light",
						Aliases: []string{"l"},
						Value:   "all",
						Usage:   "ID or address for the light to control. If not provided all lights will be updated",
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
				},
				Action: func(c *cli.Context) error {
					toggleSwitch := toggleOn
					if c.Bool("off") {
						toggleSwitch = toggleOff
					}
					lightID := c.String("light")

					fmt.Printf("Command: Switch - lights: %s, toggle: %s \n", lightID, toggleSwitch)

					devicesCh, cancelFn, err := discoverLights()
					if err != nil {
						return err
					}
					defer cancelFn()

					tbl := &texttable.TextTable{}
					tbl.SetHeader("Name", "Power State", "Brightness", "Temperature", "Address")

					count := 0
					for {
						select {
						case d := <-devicesCh:
							if d != nil {
								group, err := d.FetchLightGroup(context.Background())
								if err != nil {
									log.Println("failed to retrieve light group: ", err.Error())
									return err
								}

								newOpts := group.Copy()

								for _, light := range newOpts.Lights {
									light.On = getPowerStateInt(toggleSwitch)
									light.Brightness = c.Int("brightness")
									light.Temperature = int(math.Floor(987007 * math.Pow(float64(c.Int("temperature")), -0.999)))
								}

								_, err = d.UpdateLightGroup(context.Background(), newOpts)

								if err != nil {
									log.Println("failed to update light group: ", err.Error())
									return err
								}
								addTableRow(tbl, d, newOpts)
								count++
							}
						case <-time.After(timeout):
							drawTable(tbl, count)
							return nil
						}
					}
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "Discover and list available lights",
				Action: func(c *cli.Context) error {
					fmt.Println("Command: List - lights: all")

					devicesCh, cancelFn, err := discoverLights()
					if err != nil {
						return err
					}
					defer cancelFn()
					tbl := &texttable.TextTable{}
					tbl.SetHeader("Name", "Power State", "Brightness", "Temperature", "Address")

					count := 0
					for {
						select {
						case d := <-devicesCh:
							if d != nil {
								group, err := d.FetchLightGroup(context.Background())
								if err != nil {
									log.Println("failed to retrieve light group: ", err.Error())
									return err
								}
								addTableRow(tbl, d, group)
								count++
							}
						case <-time.After(timeout):
							drawTable(tbl, count)
							return nil
						}
					}
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
