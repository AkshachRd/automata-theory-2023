package machine

import (
	"bufio"
	"fmt"
	"mooreMealyConversion/graph"
	"os"
	"sort"
	"strings"
)

type MealyState struct {
	Name string
}

type MealyTransitionOutput struct {
	State        MealyState
	OutputSymbol Symbol
}

type MealyTransition map[MealyState]MealyTransitionOutput

type MealyMachine struct {
	InputSymbolsNum uint64
	States          map[MealyState]bool
	Transitions     Transitions[MealyTransition]
	CurrentState    MealyState
}

func (m *MealyMachine) ReadFromFile(scanner *bufio.Scanner, statesNum, inputSymbolsNum uint64) error {
	m.States = make(map[MealyState]bool, statesNum)
	m.Transitions = make(Transitions[MealyTransition], inputSymbolsNum)
	m.InputSymbolsNum = inputSymbolsNum

	for i := uint64(0); i < inputSymbolsNum; i++ {
		inputSymbol := Symbol(fmt.Sprintf("x%d", i))
		m.Transitions[inputSymbol] = make(MealyTransition, statesNum)

		scanner.Scan()
		line := scanner.Text()
		transitionOutputStrings := strings.Fields(line)

		if len(transitionOutputStrings) != int(statesNum) {
			return fmt.Errorf(
				"error reading transition name of mealy machine: it must have %d names instead of %d",
				statesNum,
				len(transitionOutputStrings),
			)
		}

		for j, transitionOutputString := range transitionOutputStrings {
			transitionOutput := strings.Split(transitionOutputString, "/")

			if len(transitionOutput) != 2 {
				return fmt.Errorf(
					"error reading transition output of mealy machine: it must have both state and output symbol")
			}

			state := MealyState{Name: fmt.Sprintf("s%d", j+1)}
			m.States[state] = true

			m.Transitions[inputSymbol][state] = MealyTransitionOutput{
				State: MealyState{Name: transitionOutput[0]}, OutputSymbol: Symbol(transitionOutput[1]),
			}
		}
	}

	return nil
}

func (m *MealyMachine) Draw(graph graph.IGraph) {
	nodes := make(map[MealyState]int)
	for state := range m.States {
		nodes[state] = graph.AddNode(state.Name)
	}

	type Edge struct {
		First  int
		Second int
	}

	edges := make(map[Edge]string)

	for inputSymbol, transition := range m.Transitions {
		for state, transitionOutput := range transition {
			first := nodes[state]
			second := nodes[MealyState{Name: transitionOutput.State.Name}]

			edge := Edge{First: first, Second: second}
			label := string(inputSymbol) + "/" + string(transitionOutput.OutputSymbol)

			if _, ok := edges[edge]; ok {
				edges[edge] += ", " + label
			} else {
				edges[edge] = label
			}
		}
	}

	for edge, label := range edges {
		graph.AddEdge(edge.First, edge.Second, label)
	}
}

func (m *MealyMachine) Print(file *os.File) error {
	writer := bufio.NewWriter(file)

	var sortedStates []MealyState
	for state := range m.States {
		sortedStates = append(sortedStates, state)
	}
	sort.Slice(sortedStates, func(i, j int) bool {
		return sortedStates[i].Name < sortedStates[j].Name
	})

	for _, transition := range m.Transitions {
		for _, state := range sortedStates {
			fmt.Fprint(writer, transition[state].State.Name, "/", transition[state].OutputSymbol, " ")
		}

		fmt.Fprintln(writer)
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}
