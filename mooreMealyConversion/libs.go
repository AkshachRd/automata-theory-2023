package main

import "fmt"

type Symbol string

type Transitions[T any] map[Symbol]T

type MachineType int

const (
	Moore MachineType = 2
	Mealy             = 1
)

type IMachine interface {
	Draw()
	ConvertToMoore() *MooreMachine
	ConvertToMealy() *MealyMachine
}

type Machine struct {
	Type           MachineType
	Implementation IMachine
}

func (m *Machine) ConvertToMachine(machineType MachineType) error {
	switch machineType {
	case Moore:
		m.Implementation = m.Implementation.ConvertToMoore()
	case Mealy:
		m.Implementation = m.Implementation.ConvertToMealy()
	default:
		return fmt.Errorf("unknown machine type")
	}

	return nil
}
