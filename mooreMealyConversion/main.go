package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadMachineFromFile(filePath string) (*Machine, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	args := strings.Fields(line)
	if len(args) != 3 {
		fmt.Println("Invalid arguments count. Format: <states num> <input symbols num> <machine type>", err)
		return nil, err
	}

	statesNum, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Println("Error reading number of states:", err)
		return nil, err
	}

	inputSymbolsNum, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		fmt.Println("Error reading number of input symbols:", err)
		return nil, err
	}

	machineType, err := strconv.ParseInt(args[2], 10, 0)
	if err != nil {
		fmt.Println("Error reading type of machine:", err)
		return nil, err
	}

	switch MachineType(machineType) {
	case Moore:
		moore, err := ReadMooreFromFile(scanner, statesNum, inputSymbolsNum)
		if err != nil {
			return nil, err
		}

		return &Machine{Moore, moore}, nil
	case Mealy:
		mealy, err := ReadMealyFromFile(scanner, statesNum, inputSymbolsNum)
		if err != nil {
			return nil, err
		}

		return &Machine{Mealy, mealy}, nil
	}

	return nil, errors.New("error unknown type of machine")
}

func DrawMachine(machine *Machine) {
	machine.Implementation.Draw()
}

func main() {
	machine, err := ReadMachineFromFile("./mealy2.txt")
	if err != nil {
		return
	}

	DrawMachine(machine)
	err = machine.ConvertToMachine(Moore)
	if err != nil {
		return
	}

	DrawMachine(machine)
}
