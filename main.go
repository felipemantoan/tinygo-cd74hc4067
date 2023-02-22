package main

import (
	"errors"
	"fmt"
	"image/color"
	"machine"
	"machine/usb/midi"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

type Cd74hc4067 struct {
	Signal  machine.Pin
	Enabled machine.Pin
	S1      machine.Pin
	S2      machine.Pin
	S3      machine.Pin
	S4      machine.Pin
}

func (mux *Cd74hc4067) GetSignal(pin int) (int, bool, error) {

	if pin > 15 {
		return -2, false, errors.New(fmt.Sprintf("Invalid pin %d", pin))
	}

	mux.S1.Set(pins[pin][0] == 1)
	mux.S2.Set(pins[pin][1] == 1)
	mux.S3.Set(pins[pin][2] == 1)
	mux.S4.Set(pins[pin][3] == 1)

	time.Sleep(10 * time.Microsecond)

	return pin, mux.Signal.Get(), nil
}

var pins = [][]uint8{
	{0, 0, 0, 0},
	{0, 0, 0, 1},
	{0, 0, 1, 0},
	{0, 0, 1, 1},
	{0, 1, 0, 0},
	{0, 1, 0, 1},
	{0, 1, 1, 0},
	{0, 1, 1, 1},
	{1, 0, 0, 0},
	{1, 0, 0, 1},
	{1, 0, 1, 0},
	{1, 0, 1, 1},
	{1, 1, 0, 0},
	{1, 1, 0, 1},
	{1, 1, 1, 0},
	{1, 1, 1, 1},
}

var notes = []midi.Note{
	midi.C3,
	midi.D3,
	midi.E3,
	midi.F3,
	midi.G3,
	midi.A3,
	midi.B3,
	midi.C4,
	midi.D4,
	midi.E4,
	midi.F4,
	midi.G4,
	midi.A4,
	midi.B4,
	midi.C5,
	midi.D5,
}

var Midi = midi.New()

func setup() {
	machine.D6.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	machine.D10.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.D9.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.D8.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.D7.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.NEOPIXEL.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.NEOPIXEL_POWER.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.NEOPIXEL_POWER.High()

}

func main() {

	go setup()

	ws := ws2812.New(machine.NEOPIXEL)
	color := []color.RGBA{{R: 0, G: 255, B: 255, A: 1}}

	mux := &Cd74hc4067{
		Signal:  machine.D6,
		Enabled: 0,
		S1:      machine.D10,
		S2:      machine.D9,
		S3:      machine.D8,
		S4:      machine.D7,
	}

	active := make([]bool, 16)

	for {

		go ws.WriteColors(color)

		for index := range pins {

			pin, value, error := mux.GetSignal(index)

			if error != nil {
				println("Deu bosta")
			}

			if value == active[pin] {
				continue
			}

			active[pin] = value

			if active[pin] {
				Midi.NoteOn(0, 0, notes[pin], 255)
			} else {
				Midi.NoteOff(0, 0, notes[pin], 255)
			}
		}
	}
}
