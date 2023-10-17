package machine

import (
	"bufio"
	"fmt"
	"mooreMealyConversion/graph"
	"os"
	"sort"
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

func (m *MooreMachine) ReadFromFile(scanner *bufio.Scanner, statesNum, inputSymbolsNum uint64) error {
	m.States = make(map[MooreState]bool, statesNum)
	m.Transitions = make(Transitions[MooreTransition], inputSymbolsNum)
	m.InputSymbolsNum = inputSymbolsNum

	var outputSymbols []Symbol

	scanner.Scan()
	line := scanner.Text()
	outputSymbolStrings := strings.Fields(line)

	if len(outputSymbolStrings) != int(statesNum) {
		return fmt.Errorf(
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
		m.Transitions[inputSymbol] = make(MooreTransition, statesNum)

		scanner.Scan()
		line = scanner.Text()
		transitionOutputStrings := strings.Fields(line)

		if len(transitionOutputStrings) != int(statesNum) {
			return fmt.Errorf(
				"error reading transition name of moore machine: it must have %d names instead of %d",
				statesNum,
				len(transitionOutputStrings),
			)
		}

		for j, transitionOutputString := range transitionOutputStrings {
			state := MooreState{OutputSymbol: outputSymbols[j], Name: fmt.Sprintf("s%d", j)}
			m.States[state] = true

			m.Transitions[inputSymbol][state] = transitionOutputString
		}
	}

	return nil
}

func (m *MooreMachine) Draw(graph graph.IGraph) {
	nodes := make(map[MooreState]int)
	for state := range m.States {
		nodes[state] = graph.AddNode(state.Name + "/" + string(state.OutputSymbol))
	}

	type Edge struct {
		From int
		To   int
	}

	edges := make(map[Edge]string)

	for inputSymbol, transition := range m.Transitions {
		for state, transitionOutput := range transition {
			from := nodes[state]

			var toNodeSymbol Symbol
			found := false
			for mooreState := range m.States {
				if mooreState.Name == transitionOutput {
					found = true
					toNodeSymbol = mooreState.OutputSymbol
				}
			}
			if !found {
				panic("invalid state")
			}

			to := nodes[MooreState{Name: transitionOutput, OutputSymbol: toNodeSymbol}]

			edge := Edge{From: from, To: to}
			label := string(inputSymbol)

			if _, ok := edges[edge]; ok {
				edges[edge] += ", " + label
			} else {
				edges[edge] = label
			}
		}
	}

	for edge, label := range edges {
		graph.AddEdge(edge.From, edge.To, label)
	}
}

func (m *MooreMachine) Print(file *os.File) error {
	writer := bufio.NewWriter(file)

	var sortedStates []MooreState
	for state := range m.States {
		sortedStates = append(sortedStates, state)
	}
	sort.Slice(sortedStates, func(i, j int) bool {
		var n, q int
		fmt.Sscanf(sortedStates[i].Name, "s%d", &n)
		fmt.Sscanf(sortedStates[j].Name, "s%d", &q)

		return n < q
	})

	for _, state := range sortedStates {
		fmt.Fprint(writer, state.OutputSymbol, " ")
	}
	fmt.Fprintln(writer)

	var transitionInputSymbols []Symbol
	for inputSymbol, _ := range m.Transitions {
		transitionInputSymbols = append(transitionInputSymbols, inputSymbol)
	}
	sort.Slice(transitionInputSymbols, func(i, j int) bool {
		var n, q int
		fmt.Sscanf(string(transitionInputSymbols[i]), "x%d", &n)
		fmt.Sscanf(string(transitionInputSymbols[j]), "x%d", &q)

		return n < q
	})

	for _, inputSymbol := range transitionInputSymbols {
		for _, state := range sortedStates {
			fmt.Fprint(writer, m.Transitions[inputSymbol][state], " ")
		}

		fmt.Fprintln(writer)
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}
