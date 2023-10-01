package machine

import (
	"bufio"
	"errors"
	"fmt"
	"mooreMealyConversion/graph"
	"os"
	"strconv"
	"strings"
)

type Symbol string

type Transitions[T any] map[Symbol]T

type MachineType int

const (
	Moore MachineType = 2
	Mealy             = 1
)

type Machine struct {
	Type           MachineType
	Implementation IMachineImplementation
}

func NewMachine(machineType MachineType, machineImplementation IMachineImplementation) *Machine {
	return &Machine{Type: machineType, Implementation: machineImplementation}
}

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
		moore := MooreMachine{}
		err := moore.ReadFromFile(scanner, statesNum, inputSymbolsNum)
		if err != nil {
			return nil, err
		}

		return NewMachine(Moore, &moore), nil
	case Mealy:
		mealy := MealyMachine{}
		err := mealy.ReadFromFile(scanner, statesNum, inputSymbolsNum)
		if err != nil {
			return nil, err
		}

		return NewMachine(Mealy, &mealy), nil
	}

	return nil, errors.New("error unknown type of machine")
}

func (m *Machine) DrawGraph(outputFileName string) {
	graphView := graph.NewGraph()
	m.Implementation.Draw(graphView)
	graphView.GenerateImage(outputFileName)
}

func (m *Machine) ConvertToMachine(machineType MachineType) error {
	switch machineType {
	case Moore:
		mealy, ok := m.Implementation.(*MealyMachine)
		if !ok {
			return fmt.Errorf("cannot convert to MealyMachine")
		}
		moore := ConvertMealyToMoore(mealy)
		m.Implementation = &moore
	case Mealy:
		moore, ok := m.Implementation.(*MooreMachine)
		if !ok {
			return fmt.Errorf("cannot convert to MooreMachine")
		}
		mealy := ConvertMooreToMealy(moore)
		m.Implementation = &mealy
	default:
		return fmt.Errorf("unknown machine type")
	}

	return nil
}

func (m *Machine) Print(outputFileName string) error {
	file, err := os.Create(outputFileName)
	defer file.Close()

	if err != nil {
		return fmt.Errorf("invalid output file")
	}

	err = m.Implementation.Print(file)
	if err != nil {
		return fmt.Errorf("can't print file: %+v\n", err)
	}

	return nil
}
