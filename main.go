package main

import (
	"image/color"
	"machine"
	"machine/usb/midi"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

type ActiveChannel struct {
	Value bool
	Pin   int
}

type Cd74hc4067 struct {
	Signal         machine.Pin
	Enabled        machine.Pin
	S1             machine.Pin
	S2             machine.Pin
	S3             machine.Pin
	S4             machine.Pin
	ActiveChannels chan ActiveChannel
}

func (mux *Cd74hc4067) Select() {
	for {
		for pin := range pins {
			mux.S1.Set(pins[pin][0] == 1)
			mux.S2.Set(pins[pin][1] == 1)
			mux.S3.Set(pins[pin][2] == 1)
			mux.S4.Set(pins[pin][3] == 1)

			time.Sleep(12 * time.Microsecond)

			mux.ActiveChannels <- ActiveChannel{
				Pin:   pin,
				Value: mux.Signal.Get(),
			}
		}
	}
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

var notesActive = make(map[int]bool)

func setup() {

	machine.D0.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	machine.D1.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	machine.D2.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	machine.D3.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	machine.D4.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	machine.D5.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

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

	setup()

	mux := &Cd74hc4067{
		Signal:         machine.D6,
		Enabled:        0,
		S1:             machine.D10,
		S2:             machine.D9,
		S3:             machine.D8,
		S4:             machine.D7,
		ActiveChannels: make(chan ActiveChannel, 16),
	}

	ws := ws2812.New(machine.NEOPIXEL)

	go mux.Select()

	for pin := range mux.ActiveChannels {

		if notesActive[pin.Pin] == pin.Value {
			continue
		}

		notesActive[pin.Pin] = pin.Value

		if notesActive[pin.Pin] {
			Midi.NoteOn(0, 0, notes[pin.Pin], 127)
			ws.WriteColors([]color.RGBA{{R: 255, G: 0, B: 0, A: 100}})
		} else {
			Midi.NoteOff(0, 0, notes[pin.Pin], 127)
			ws.WriteColors([]color.RGBA{{R: 0, G: 0, B: 0, A: 0}})
		}
	}
}
