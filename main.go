package main

import (
	"image/color"
	"machine"
	"machine/usb/midi"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

type ADCChannel struct {
	Value uint8
}

var adc0 = make(chan ADCChannel)

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
	btnState := make(map[int]bool)
	for {
		for pin := range pins {
			mux.S1.Set(pins[pin][0] == 1)
			mux.S2.Set(pins[pin][1] == 1)
			mux.S3.Set(pins[pin][2] == 1)
			mux.S4.Set(pins[pin][3] == 1)

			time.Sleep(20 * time.Microsecond)

			pinValue := mux.Signal.Get()

			if btnState[pin] == pinValue {
				continue
			}

			btnState[pin] = pinValue

			mux.ActiveChannels <- ActiveChannel{
				Pin:   pin,
				Value: pinValue,
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

func setup() {
	machine.InitADC()

	machine.D10.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.D9.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.D8.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.D7.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.NEOPIXEL.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.NEOPIXEL_POWER.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.NEOPIXEL_POWER.High()
}

func SendCC() {
	sensor := machine.ADC{Pin: machine.ADC0}
	sensor.Configure(machine.ADCConfig{})

	var adc0State uint8

	for {
		normal := float32(sensor.Get()) / float32(65535)

		value := uint8((normal * 255) / 2)

		if adc0State == value {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		adc0State = value

		adc0 <- ADCChannel{
			Value: value,
		}

		time.Sleep(10 * time.Microsecond)
	}
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

	go SendCC()
	go mux.Select()

	for {
		select {
		case pin := <-mux.ActiveChannels:

			if pin.Value {
				Midi.NoteOn(0, 0, notes[pin.Pin], 127)
				ws.WriteColors([]color.RGBA{{R: 255, G: 0, B: 0, A: 100}})
			} else {
				Midi.NoteOff(0, 0, notes[pin.Pin], 127)
				ws.WriteColors([]color.RGBA{{R: 0, G: 0, B: 0, A: 0}})
			}
		case cc := <-adc0:
			Midi.SendCC(0, 0, 1, cc.Value)
		}
	}
}
