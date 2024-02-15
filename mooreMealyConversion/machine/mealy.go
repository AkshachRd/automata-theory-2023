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

			state := MealyState{Name: fmt.Sprintf("s%d", j)}
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

type MealyPartition []MealyState

func (m *MealyMachine) Minimize() error {
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

func (m *MealyMachine) partitionsToMachine(partitions []MealyPartition) {
	newStates := make(map[MealyState]bool)
	newTransitions := make(Transitions[MealyTransition])
	for inputSymbol := range m.Transitions {
		newTransitions[inputSymbol] = make(MealyTransition)
	}
	statesToNewStates := make(map[MealyState]MealyState)

	for _, partition := range partitions {
		name := ""
		for i, state := range partition {
			name += state.Name
			if i != len(partition) - 1 {
				name += ","
			}
		}
		state := MealyState{Name: name}
		newStates[state] = true

		for _, oldState := range partition {
			statesToNewStates[oldState] = state
		}
	}

	for _, partition := range partitions {
		name := ""
		for i, state := range partition {
			name += state.Name
			if i != len(partition) - 1 {
				name += ","
			}
		}
		state := MealyState{Name: name}

		for inputSymbol, transition := range m.Transitions {
			newTransition := MealyTransitionOutput{
				State:        statesToNewStates[transition[partition[0]].State],
				OutputSymbol: transition[partition[0]].OutputSymbol,
			}

			newTransitions[inputSymbol][state] = newTransition
		}
	}

	m.States = newStates
	m.Transitions = newTransitions
	m.CurrentState = MealyState{Name: "q0"}
}

func (m *MealyMachine) getInitialPartitions() []MealyPartition {
	var partitions []MealyPartition

	nextPartitionIndex := 0
	outputSymbolsVariants := make(map[string]int)
	var outputSymbolsVariant string

	for state := range m.States {
		outputSymbolsVariant = ""
		for _, transition := range m.Transitions {
			outputSymbolsVariant += string(transition[state].OutputSymbol)
		}

		if partitionIndex, ok := outputSymbolsVariants[outputSymbolsVariant]; !ok {
			partitions = append(partitions, MealyPartition{state})
			outputSymbolsVariants[outputSymbolsVariant] = nextPartitionIndex
			nextPartitionIndex++
		} else {
			partitions[partitionIndex] = append(partitions[partitionIndex], state)
		}
	}

	return partitions
}

func (m *MealyMachine) calculatePartitions(partitions []MealyPartition) []MealyPartition {
	var newPartitions []MealyPartition

	for i := 0; i < len(partitions); i++ {
		partition := partitions[i]
		var newPartition MealyPartition
		var restPartition MealyPartition

		for _, state := range partition {
			if len(newPartition) == 0 {
				newPartition = append(newPartition, state)
				continue
			}

			if m.checkIfInSamePartition(partitions, newPartition[0], state) {
				newPartition = append(newPartition, state)
			} else {
				restPartition = append(restPartition, state)
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


func (m *MealyMachine) checkIfInSamePartition(partitions []MealyPartition, first, second MealyState) bool {
	inSamePartition := true

	for _, transition := range m.Transitions {
		transitionOutputFirst := transition[first]
		transitionOutputSecond := transition[second]

		transitionInSamePartition := false
		for _, partition := range partitions {
			if slices.Contains(partition, transitionOutputFirst.State) && slices.Contains(partition, transitionOutputSecond.State) {
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
