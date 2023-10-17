package main

import (
	"fmt"
	"mooreMealyConversion/machine"
)

func main() {
	myMachine, err := machine.ReadMachineFromFile("./mealy-in-2.txt")
	if err != nil {
		fmt.Println("error reading a machine from file", err)
		return
	}

	err = myMachine.Print("original-machine.txt")
	if err != nil {
		fmt.Println("error printing the machine", err)
		return
	}
	myMachine.DrawGraph("original-machine")

	if myMachine.Type == machine.Mealy {
		err = myMachine.ConvertToMachine(machine.Moore)
	} else {
		err = myMachine.ConvertToMachine(machine.Mealy)
	}
	if err != nil {
		fmt.Println("error converting the machine to Moore", err)
		return
	}

	err = myMachine.Print("converted-machine.txt")
	if err != nil {
		fmt.Println("error printing the machine", err)
		return
	}
	myMachine.DrawGraph("converted-machine")
}
