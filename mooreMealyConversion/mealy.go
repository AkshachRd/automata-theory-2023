package main

import (
	"bufio"
	"fmt"
	"github.com/mzohreva/GoGraphviz/graphviz"
	"log"
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

func ReadMealyFromFile(scanner *bufio.Scanner, statesNum, inputSymbolsNum uint64) (*MealyMachine, error) {
	var mealy MealyMachine

	mealy.States = make(map[MealyState]bool, statesNum)
	mealy.Transitions = make(Transitions[MealyTransition], inputSymbolsNum)
	mealy.InputSymbolsNum = inputSymbolsNum

	for i := uint64(0); i < inputSymbolsNum; i++ {
		inputSymbol := Symbol(fmt.Sprintf("x%d", i))
		mealy.Transitions[inputSymbol] = make(MealyTransition, statesNum)

		scanner.Scan()
		line := scanner.Text()
		transitionOutputStrings := strings.Fields(line)

		if len(transitionOutputStrings) != int(statesNum) {
			return nil, fmt.Errorf(
				"error reading transition name of mealy machine: it must have %d names instead of %d",
				statesNum,
				len(transitionOutputStrings),
			)
		}

		for j, transitionOutputString := range transitionOutputStrings {
			transitionOutput := strings.Split(transitionOutputString, "/")

			if len(transitionOutput) != 2 {
				return nil, fmt.Errorf(
					"error reading transition output of mealy machine: it must have both state and output symbol")
			}

			state := MealyState{Name: fmt.Sprintf("s%d", j+1)}
			mealy.States[state] = true

			mealy.Transitions[inputSymbol][state] = MealyTransitionOutput{
				State: MealyState{Name: transitionOutput[0]}, OutputSymbol: Symbol(transitionOutput[1]),
			}
		}
	}

	return &mealy, nil
}

func (m *MealyMachine) ConvertToMealy() *MealyMachine {
	return m
}

func (m *MealyMachine) ConvertToMoore() *MooreMachine {
	var moore MooreMachine
	moore.States = make(map[MooreState]bool)
	moore.Transitions = make(Transitions[MooreTransition])

	moore.InputSymbolsNum = m.InputSymbolsNum
	for inputSymbol, mealyTransition := range m.Transitions {
		for mealyState, mealyTransitionOutput := range mealyTransition {
			mooreState := MooreState{
				Name:         mealyState.Name,
				OutputSymbol: mealyTransitionOutput.OutputSymbol,
			}
			moore.States[mooreState] = true

			if m.CurrentState.Name == mooreState.Name {
				moore.CurrentState = mooreState
			}
			if _, ok := moore.Transitions[inputSymbol]; !ok {
				moore.Transitions[inputSymbol] = make(MooreTransition)
			}

			moore.Transitions[inputSymbol][mooreState] = mealyTransitionOutput.State.Name
		}
	}

	return &moore
}

func (m *MealyMachine) Draw() {
	graph := graphviz.Graph{}
	graph.MakeDirected()

	nodes := make(map[MealyState]int)
	for state := range m.States {
		nodes[state] = graph.AddNode(fmt.Sprintf("%s", state.Name))
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

	err := graph.GenerateImage("dot", "mealy.png", "png")
	if err != nil {
		log.Fatal(err)
	}
}
