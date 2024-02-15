package main

import (
	"fmt"
	"mooreMealyConversion/machine"
)

func main() {
	myMachine, err := machine.ReadMachineFromFile("./moore-in-5.txt")
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

	err = myMachine.Implementation.Minimize()
	if err != nil {
		fmt.Println("error minimizing the machine", err)
		return
	}

	err = myMachine.Print("converted-machine.txt")
	if err != nil {
		fmt.Println("error printing the machine", err)
		return
	}
	myMachine.DrawGraph("converted-machine")

	// if myMachine.Type == machine.Mealy {
	// 	err = myMachine.ConvertToMachine(machine.Moore)
	// } else {
	// 	err = myMachine.ConvertToMachine(machine.Mealy)
	// }
	// if err != nil {
	// 	fmt.Println("error converting the machine to Moore", err)
	// 	return
	// }

	// err = myMachine.Print("converted-machine.txt")
	// if err != nil {
	// 	fmt.Println("error printing the machine", err)
	// 	return
	// }
	// myMachine.DrawGraph("converted-machine")
}
