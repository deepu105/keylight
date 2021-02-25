package main

import (
	"testing"
)

func Test_convertFromKelvin(t *testing.T) {
	tests := []struct {
		name   string
		kelvin int
		want   int
	}{
		{
			"Convert max kelvin to API val",
			7000,
			142,
		},
		{
			"Convert min kelvin to API val",
			3000,
			331,
		},
		{
			"Convert a kelvin to API val",
			5000,
			199,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertFromKelvin(tt.kelvin); got != tt.want {
				t.Errorf("convertFromKelvin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToKelvin(t *testing.T) {
	tests := []struct {
		name string
		temp int
		want int
	}{
		{
			"Convert min API val to kelvin",
			142,
			7007,
		},
		{
			"Convert max API val to kelvin",
			331,
			3006,
		},
		{
			"Convert a API val to kelvin",
			199,
			5000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertToKelvin(tt.temp); got != tt.want {
				t.Errorf("convertToKelvin() = %v, want %v", got, tt.want)
			}
		})
	}
}
