package machine

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"mooreMealyConversion/graph"
	"os"
	"sort"
)

type Symbol string

type Transitions[T any] map[Symbol]T

type State struct {
	Name    string
	isFinal bool
}

type Transition map[State]State

type Machine struct {
	InputSymbols []Symbol
	States       map[State]bool
	Transitions  Transitions[Transition]
}

func NewMachine() *Machine {
	return &Machine{make([]Symbol, 0), make(map[State]bool), make(Transitions[Transition])}
}

func ReadMachineFromFile(filePath string) (*Machine, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	// Print each record
	for _, record := range records {
		fmt.Println(record)
	}

	return nil, nil
}

func (m *Machine) Draw(outputFileName string) {
	graphView := graph.NewGraph()

	nodes := make(map[State]int)
	for state := range m.States {
		nodes[state] = graphView.AddNode(state.Name)
	}

	type Edge struct {
		First  int
		Second int
	}

	edges := make(map[Edge]string)

	for inputSymbol, transition := range m.Transitions {
		for state, transitionOutput := range transition {
			first := nodes[state]
			second := nodes[State{Name: transitionOutput.Name}]

			edge := Edge{First: first, Second: second}
			label := string(inputSymbol)

			if _, ok := edges[edge]; ok {
				edges[edge] += ", " + label
			} else {
				edges[edge] = label
			}
		}
	}

	for edge, label := range edges {
		graphView.AddEdge(edge.First, edge.Second, label)
	}

	graphView.GenerateImage(outputFileName)
}

func (m *Machine) Print(outputFileName string) error {
	file, err := os.Create(outputFileName)
	defer file.Close()

	if err != nil {
		return fmt.Errorf("invalid output file")
	}

	writer := bufio.NewWriter(file)

	var sortedStates []State
	for state := range m.States {
		sortedStates = append(sortedStates, state)
	}
	sort.Slice(sortedStates, func(i, j int) bool {
		return sortedStates[i].Name < sortedStates[j].Name
	})

	for _, transition := range m.Transitions {
		for _, state := range sortedStates {
			fmt.Fprint(writer, transition[state].Name, " ")
		}

		fmt.Fprintln(writer)
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	return nil
}
