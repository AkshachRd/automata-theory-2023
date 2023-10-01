package main

import (
	"fmt"
	"mooreMealyConversion/machine"
)

func main() {
	myMachine, err := machine.ReadMachineFromFile("./mealy.txt")
	if err != nil {
		fmt.Println("error reading a machine from file", err)
		return
	}

	err = myMachine.Print("test1")
	if err != nil {
		fmt.Println("error printing the machine", err)
		return
	}
	myMachine.DrawGraph("test1")

	err = myMachine.ConvertToMachine(machine.Moore)
	if err != nil {
		fmt.Println("error converting the machine to Moore", err)
		return
	}

	err = myMachine.Print("test2")
	if err != nil {
		fmt.Println("error printing the machine", err)
		return
	}
	myMachine.DrawGraph("test2")
}
