package machine

import (
	"bufio"
	"fmt"
	"mooreMealyConversion/graph"
	"os"
	"reflect"
	"slices"
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
	for inputSymbol := range m.Transitions {
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

type MoorePartition []string

func (m *MooreMachine) Minimize() error {
	partitions := m.getInitialPartitions()

	for {
		newPartitions := m.calculatePartitions(partitions)

		if reflect.DeepEqual(newPartitions, partitions) {
			break
		}

		partitions = newPartitions
	}

	m.partitionsToMachine(partitions);

	return nil
}

func (m *MooreMachine) partitionsToMachine(partitions []MoorePartition) {
	newStates := make(map[MooreState]bool)
	newTransitions := make(Transitions[MooreTransition])
	for inputSymbol := range m.Transitions {
		newTransitions[inputSymbol] = make(MooreTransition)
	}
	oldStatesToNewStates := make(map[MooreState]MooreState)

	for _, partition := range partitions {
		state := MooreState{Name: strings.Join(partition, ","), OutputSymbol: m.findStateByName(partition[0]).OutputSymbol}
		newStates[state] = true

		for _, oldStateName := range partition {
			oldStatesToNewStates[m.findStateByName(oldStateName)] = state
		}
	}

	for _, partition := range partitions {
		state := MooreState{Name: strings.Join(partition, ","), OutputSymbol: m.findStateByName(partition[0]).OutputSymbol}

		for inputSymbol, transition := range m.Transitions {
			oldState := m.findStateByName(transition[m.findStateByName(partition[0])])

			newTransitions[inputSymbol][state] = oldStatesToNewStates[oldState].Name
		}
	}

	m.States = newStates
	m.Transitions = newTransitions
	m.CurrentState = MooreState{Name: "q0", OutputSymbol: m.CurrentState.OutputSymbol}
}

func (m *MooreMachine) getInitialPartitions() []MoorePartition {
	var partitions []MoorePartition

	partitionsMap := make(map[Symbol]MoorePartition)

	for state := range m.States {
		partitionsMap[state.OutputSymbol] = append(partitionsMap[state.OutputSymbol], state.Name)
	}

	for _, partition := range partitionsMap {
		partitions = append(partitions, partition)
	}

	return partitions
}

func (m *MooreMachine) calculatePartitions(partitions []MoorePartition) []MoorePartition {
	var newPartitions []MoorePartition

	for i := 0; i < len(partitions); i++ {
		partition := partitions[i]
		var newPartition MoorePartition
		var restPartition MoorePartition

		for _, stateName := range partition {
			if len(newPartition) == 0 {
				newPartition = append(newPartition, stateName)
				continue
			}

			if m.checkIfInSamePartition(partitions, newPartition[0], stateName) {
				newPartition = append(newPartition, stateName)
			} else {
				restPartition = append(restPartition, stateName)
			}
		}

		if len(newPartition) > 0 {
			newPartitions = append(newPartitions, newPartition)
		}

		if len(restPartition) > 0 {
			partitions = append(partitions, restPartition)
		}
	}

	return newPartitions
}

func (m *MooreMachine) checkIfInSamePartition(partitions []MoorePartition, firstName, secondName string) bool {
	inSamePartition := true

	for _, transition := range m.Transitions {
		mooreStateFirst := MooreState{Name: transition[m.findStateByName(firstName)]}
		mooreStateSecond := MooreState{Name: transition[m.findStateByName(secondName)]}

		for state := range m.States {
			if state.Name == mooreStateFirst.Name {
				mooreStateFirst = state
			}
			if state.Name == mooreStateSecond.Name {
				mooreStateSecond = state
			}
		}

		transitionInSamePartition := false
		for _, partition := range partitions {
			if slices.Contains(partition, mooreStateFirst.Name) && slices.Contains(partition, mooreStateSecond.Name) {
				transitionInSamePartition = true
			}
		}

		if !transitionInSamePartition {
			inSamePartition = false
			break
		}
	}

	return inSamePartition
}

func (m *MooreMachine) findStateByName(name string) MooreState {
	for state := range m.States {
		if state.Name == name {
			return state
		}
	}

	return MooreState{}
}