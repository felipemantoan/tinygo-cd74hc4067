package main

import (
	"machine"
	"time"
)

var pins = [][]int{
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

	for {
		for index, pin := range pins {

			machine.D10.Set(pin[0] == 1)
			machine.D9.Set(pin[1] == 1)
			machine.D8.Set(pin[2] == 1)
			machine.D7.Set(pin[3] == 1)
			time.Sleep(10 * time.Microsecond)

			if machine.D6.Get() {
				println(index)
			}
		}
	}
}
