package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	keylight "github.com/endocrimes/keylight-go"
	"github.com/oleksandr/bonjour"
)

func discoverLights(ctx context.Context) (chan *Light, error) {
	resultsCh := make(chan *bonjour.ServiceEntry)
	lightsCh := make(chan *Light, 5)

	resolver, err := bonjour.NewResolver(nil)
	if err != nil {
		return nil, err
	}
	err = resolver.Browse("_elg._tcp", "local.", resultsCh)
	if err != nil {
		return nil, err
	}

	for {
		select {
		case <-ctx.Done():
			close(resultsCh)
			resolver.Exit <- true
			return lightsCh, nil
		case e := <-resultsCh:
			lightsCh <- &Light{
				Name:     e.Instance,
				HostName: e.HostName,
				Port:     e.Port,
			}
		}
	}
}

func getPowerState(state int) string {
	if state == 0 {
		return "off"
	}
	return "on"
}

func fetchLightGroup(d *keylight.Device) (*keylight.LightGroup, error) {
	o := &keylight.LightGroup{Lights: make([]*keylight.Light, 0)}
	url := fmt.Sprintf("http://%s:%d/%s", d.DNSAddr, d.Port, "elgato/lights")
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(o)
	return o, err
}
