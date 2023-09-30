package main

import "mooreMealyConversion/machine"

func main() {
	myMachine, err := machine.ReadMachineFromFile("./mealy2.txt")
	if err != nil {
		return
	}

	myMachine.DrawGraph("test1")
	err = myMachine.ConvertToMachine(machine.Moore)
	if err != nil {
		return
	}

	myMachine.DrawGraph("test2")
}
