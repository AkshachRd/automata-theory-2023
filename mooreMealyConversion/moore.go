package main

import (
	"bufio"
	"fmt"
	"github.com/mzohreva/GoGraphviz/graphviz"
	"log"
	"strings"
)

type MooreState struct {
	Name         string
	OutputSymbol Symbol
}

type MooreTransition map[MooreState]string

type MooreMachine struct {
	InputSymbolsNum uint64
	States          map[MooreState]bool
	Transitions     Transitions[MooreTransition]
	CurrentState    MooreState
}

func ReadMooreFromFile(scanner *bufio.Scanner, statesNum, inputSymbolsNum uint64) (*MooreMachine, error) {
	var moore MooreMachine

	moore.States = make(map[MooreState]bool, statesNum)
	moore.Transitions = make(Transitions[MooreTransition], inputSymbolsNum)
	moore.InputSymbolsNum = inputSymbolsNum

	var outputSymbols []Symbol

	scanner.Scan()
	line := scanner.Text()
	outputSymbolStrings := strings.Fields(line)

	if len(outputSymbolStrings) != int(statesNum) {
		return nil, fmt.Errorf(
			"error reading output symbols of moore machine: it must have %d symbols instead of %d",
			statesNum,
			len(outputSymbolStrings),
		)
	}

	for _, outputSymbolString := range outputSymbolStrings {
		outputSymbols = append(outputSymbols, Symbol(outputSymbolString))
	}

	for i := uint64(0); i < inputSymbolsNum; i++ {
		inputSymbol := Symbol(fmt.Sprintf("x%d", i))
		moore.Transitions[inputSymbol] = make(MooreTransition, statesNum)

		scanner.Scan()
		line = scanner.Text()
		transitionOutputStrings := strings.Fields(line)

		if len(transitionOutputStrings) != int(statesNum) {
			return nil, fmt.Errorf(
				"error reading transition name of moore machine: it must have %d names instead of %d",
				statesNum,
				len(transitionOutputStrings),
			)
		}

		for j, transitionOutputString := range transitionOutputStrings {
			state := MooreState{OutputSymbol: outputSymbols[j], Name: fmt.Sprintf("s%d", j+1)}
			moore.States[state] = true

			moore.Transitions[inputSymbol][state] = transitionOutputString
		}
	}

	return &moore, nil
}

func (m *MooreMachine) ConvertToMealy() *MealyMachine {
	var mealy MealyMachine
	mealy.States = make(map[MealyState]bool)
	mealy.Transitions = make(Transitions[MealyTransition])

	mealy.InputSymbolsNum = m.InputSymbolsNum
	for inputSymbol, mooreTransition := range m.Transitions {
		for mooreState, mooreTransitionOutput := range mooreTransition {
			mealyState := MealyState{
				Name: mooreState.Name,
			}
			mealy.States[mealyState] = true

			if m.CurrentState.Name == mealyState.Name {
				mealy.CurrentState = mealyState
			}
			if _, ok := mealy.Transitions[inputSymbol]; !ok {
				mealy.Transitions[inputSymbol] = make(MealyTransition)
			}
			mealy.Transitions[inputSymbol][mealyState] = MealyTransitionOutput{
				State:        MealyState{Name: mooreTransitionOutput},
				OutputSymbol: mooreState.OutputSymbol,
			}
		}
	}

	return &mealy
}

func (m *MooreMachine) ConvertToMoore() *MooreMachine {
	return m
}

func (m *MooreMachine) Draw() {
	graph := graphviz.Graph{}
	graph.MakeDirected()

	nodes := make(map[MooreState]int)
	for state := range m.States {
		nodes[state] = graph.AddNode(fmt.Sprintf("%s/%s", state.Name, state.OutputSymbol))
	}

	type Edge struct {
		First  int
		Second int
	}

	edges := make(map[Edge]string)

	for inputSymbol, transition := range m.Transitions {
		for state, transitionOutput := range transition {
			first := nodes[state]

			var secondNodeSymbol Symbol
			found := false
			for mooreState := range m.States {
				if mooreState.Name == transitionOutput {
					found = true
					secondNodeSymbol = mooreState.OutputSymbol
				}
			}
			if !found {
				panic("invalid state")
			}

			second := nodes[MooreState{Name: transitionOutput, OutputSymbol: secondNodeSymbol}]

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
		graph.AddEdge(edge.First, edge.Second, label)
	}

	err := graph.GenerateImage("dot", "moore.png", "png")
	if err != nil {
		log.Fatal(err)
	}
}
