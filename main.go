package main

import (
	"errors"
	"fmt"
	"machine"
	"time"
)

type Cd74hc4067 struct {
	Signal  machine.Pin
	Enabled machine.Pin
	S1      machine.Pin
	S2      machine.Pin
	S3      machine.Pin
	S4      machine.Pin
}

func (mux *Cd74hc4067) GetSignal(pin int) (bool, error) {

	if pin > 15 {
		return false, errors.New(fmt.Sprintf("Invalid pin %d", pin))
	}

	mux.S1.Set(pins[pin][0] == 1)
	mux.S2.Set(pins[pin][1] == 1)
	mux.S3.Set(pins[pin][2] == 1)
	mux.S4.Set(pins[pin][3] == 1)

	time.Sleep(10 * time.Microsecond)

	value := mux.Signal.Get()

	return value, nil
}

func (mux *Cd74hc4067) HandlePinsSignal() (int, bool, error) {

	for index := range pins {

		value, error := mux.GetSignal(index)

		if error != nil {
			return -1, false, error
		}

		if value {
			return index, value, nil
		}
	}

	return -1, false, nil
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

func main() {

	machine.D6.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	machine.D10.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.D9.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.D8.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.D7.Configure(machine.PinConfig{Mode: machine.PinOutput})

	mux := &Cd74hc4067{
		Signal:  machine.D6,
		Enabled: 0,
		S1:      machine.D10,
		S2:      machine.D9,
		S3:      machine.D8,
		S4:      machine.D7,
	}

	for {
		pin, value, error := mux.HandlePinsSignal()

		if error != nil {
			println("Deu bosta")
		}

		if value {
			println(value, pin)
		}
	}

}
