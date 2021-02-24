package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	keylight "github.com/endocrimes/keylight-go"
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
					&cli.StringFlag{
						Name:    "brightness",
						Aliases: []string{"b"},
						Usage:   "Set brightness of the lights",
					},
					&cli.StringFlag{
						Name:    "temperature",
						Aliases: []string{"t"},
						Usage:   "Set temperature of the lights in kelvin",
					},
				},
				Action: func(c *cli.Context) error {
					toggleSwitch := toggleOn
					if c.Bool("off") {
						toggleSwitch = toggleOff
					}
					lightID := c.String("light")

					fmt.Printf("Command: Switch - lights: %s, toggle: %s \n", lightID, toggleSwitch)
					return nil
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "Discover and list available lights",
				Action: func(c *cli.Context) error {
					fmt.Println("Command: List - lights: all")

					discovery, err := keylight.NewDiscovery()
					if err != nil {
						log.Println("failed to initialize keylight discovery: ", err.Error())
						return err
					}
					discoveryCtx, cancelFn := context.WithCancel(context.Background())
					defer cancelFn()

					go func() {
						err := discovery.Run(discoveryCtx)
						if err != nil {
							log.Fatalln("Failed to discover lights: ", err.Error())
						}
					}()

					devicesCh := discovery.ResultsCh()

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
								tbl.AddRow(
									d.Name,
									getPowerState(group.Lights[0].On),
									fmt.Sprintf("%d", group.Lights[0].Brightness),
									fmt.Sprintf("%d (%d K)", group.Lights[0].Temperature, int(math.Round(1000000*math.Pow(float64(group.Lights[0].Temperature), -1)))),
									fmt.Sprintf("%s:%d", d.DNSAddr, d.Port),
								)
								count++
							}
						case <-time.After(timeout):
							if count > 0 {
								fmt.Println(tbl.Draw())
							} else {
								fmt.Println("Either timedout or no lights found! Try again")
							}
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
