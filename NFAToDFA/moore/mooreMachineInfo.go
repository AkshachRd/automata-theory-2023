package moore

import (
	"errors"
	"maps"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const (
    FINISH_OUTPUT_SYMBOL = "F"
    EMPTY_SYMBOL        = "e"
)

type MooreMachineInfo struct {
    States              []string
    TransitionFunctions [][]string
    OutputAlphabet      []string
    InputAlphabet       []string
}

func NewMooreMachineInfo(info []string) *MooreMachineInfo {
    m := &MooreMachineInfo{
        OutputAlphabet:      make([]string, 0),
        States:              make([]string, 0),
        InputAlphabet:       make([]string, 0),
        TransitionFunctions: make([][]string, 0),
    }

    if len(info) == 0 {
        return m
    }

    dirtyOutputAlphabet := strings.Split(info[0], ";")
    m.OutputAlphabet = dirtyOutputAlphabet[1:]
    dirtyStates := strings.Split(info[1], ";")
    m.States = dirtyStates[1:]
    for i := 2; i < len(info); i++ {
        values := strings.Split(info[i], ";")
        m.InputAlphabet = append(m.InputAlphabet, values[0])
        m.TransitionFunctions = append(m.TransitionFunctions, values[1:])
    }
    for _, inputSymbolTransitionFunctions := range m.TransitionFunctions {
        for i := range inputSymbolTransitionFunctions {
            if inputSymbolTransitionFunctions[i] == "-" {
                inputSymbolTransitionFunctions[i] = ""
            }
        }
    }

    if len(m.States) <= 1 {
        panic(errors.New("incorrect input machine. Number of states can not be less than 2"))
    }

    return m
}

func (m *MooreMachineInfo) GetCsvData() string {
    csvData := ";"

    if len(m.OutputAlphabet) != 0 {
        csvData += strings.Join(m.OutputAlphabet, ";")
    }
    csvData += "\n;"

    csvData += strings.Join(m.States, ";")
    csvData += "\n"

    for i := range m.InputAlphabet {
        csvData += m.InputAlphabet[i] + ";"
        for _, part := range m.TransitionFunctions[i] {
            if part == "" {
                part = "-"
            }
            csvData += part + ";"
        }
        csvData += "\n"
    }

    return csvData
}

func (m *MooreMachineInfo) Determine() {
	// Инициализируем новое состояние автомата
	determinedStates := make([]string, len(m.States))
	copy(determinedStates, m.States)

	determinedOutputAlphabet := make([]string, len(m.OutputAlphabet))
	copy(determinedOutputAlphabet, m.OutputAlphabet)

	determinedInputAlphabet := make([]string, len(m.InputAlphabet))
	copy(determinedInputAlphabet, m.InputAlphabet)

	eclosures := make(map[string]string)

	finishStates := make([]string, 0)
	for i, v := range m.OutputAlphabet {
		if v == FINISH_OUTPUT_SYMBOL {
			finishStates = append(finishStates, determinedStates[i])
		}
	}

	newStates := make([]string, 0)

	determinedTransitionFunctions := make([][]string, len(m.TransitionFunctions))
	for i := range m.TransitionFunctions {
		copy(determinedTransitionFunctions[i], m.TransitionFunctions[i])
	}

	for {
		newStates = make([]string, 0)
		// Определение новых состояний в ДКА
		indexOfEmptySymbol := slices.Index(determinedInputAlphabet, EMPTY_SYMBOL)
		if indexOfEmptySymbol != -1 {
			for i := 0; i < len(m.States); i++ {
				if determinedTransitionFunctions[indexOfEmptySymbol][i] != "" {
					stateIndex := i
					statesSet := make(map[string]struct{})
					statesIndexQueue := make([]int, 0)

					statesIndexQueue = append(statesIndexQueue, stateIndex)
					for stateIndex != -1 {
						statesSet[determinedStates[stateIndex]] = struct{}{}
						if determinedTransitionFunctions[indexOfEmptySymbol][stateIndex] != "" {
							if strings.Contains(determinedTransitionFunctions[indexOfEmptySymbol][stateIndex], ",") {
								states := strings.Split(determinedTransitionFunctions[indexOfEmptySymbol][stateIndex], ",")
								for _, state := range states {
									indexOfState := slices.Index(determinedStates, state)
									statesIndexQueue = append(statesIndexQueue, indexOfState)
								}
							} else {
								indexOfState := slices.Index(determinedStates, determinedTransitionFunctions[indexOfEmptySymbol][stateIndex])
								statesIndexQueue = append(statesIndexQueue, indexOfState)
							}
						}
						newStateIndex := stateIndex
						for newStateIndex == stateIndex && len(statesIndexQueue) > 0 {
							_, contains := statesSet[determinedTransitionFunctions[indexOfEmptySymbol][newStateIndex]]
							if !contains {
								newStateIndex, statesIndexQueue = statesIndexQueue[0], statesIndexQueue[1:]
							}
						}
						if newStateIndex == stateIndex {
							stateIndex = -1
						} else {
							stateIndex = newStateIndex
						}
					}
					keys := make([]string, 0, len(statesSet))
					for k := range statesSet {
						keys = append(keys, k)
					}
					eclosures[m.States[i]] = strings.Join(keys, ",")
				} else {
					eclosures[m.States[i]] = m.States[i]
				}
			}
			determinedInputAlphabet = append(determinedInputAlphabet[:indexOfEmptySymbol], determinedInputAlphabet[indexOfEmptySymbol+1:]...)
			determinedTransitionFunctions = append(determinedTransitionFunctions[:indexOfEmptySymbol], determinedTransitionFunctions[indexOfEmptySymbol+1:]...)
			newStates = []string{eclosures[m.States[0]]}
		} else {
			newSet := make(map[string]struct{})
			for _, v := range determinedStates[len(eclosures):] {
				newSet[v] = struct{}{}
			}

			for _, inputSymbolTransitionFunction := range determinedTransitionFunctions {
				if len(eclosures) != 0 {
					for i := len(eclosures); i < len(inputSymbolTransitionFunction); i++ {
						_, exists := newSet[sortString(inputSymbolTransitionFunction[i])]
						if !exists && inputSymbolTransitionFunction[i] != "" {
							newStates = append(newStates, inputSymbolTransitionFunction[i])
						}
					}
				} else {
					for _, transitionFunction := range inputSymbolTransitionFunction {
						if strings.Contains(transitionFunction, ",") {
							newStates = append(newStates, transitionFunction)
						}
					}
				}
			}
		}

		for _, newState := range newStates {
			determinedState := strings.Replace(newState, ",", "", -1)
			if len(eclosures) != 0 {
				state := sortString(determinedState)
				if slices.Contains(determinedStates, state) && len(eclosures) != 0 && slices.Index(determinedStates, state) >= len(m.States) {
					for _, inputSymbolTransitionFunctions := range determinedTransitionFunctions {
						determinedStateCharHashSet := toSet(strings.Split(determinedState, ""))
						for indexOfTransitionFunction, transitionFunction := range inputSymbolTransitionFunctions {
							if maps.Equal(toSet(strings.Split(strings.Replace(transitionFunction, ",", "", -1), "")), determinedStateCharHashSet) {
								inputSymbolTransitionFunctions[indexOfTransitionFunction] = determinedState
							}
						}
					}
					continue
				}
			}
			if slices.Contains(determinedStates, determinedState) && len(eclosures) == 0 {
				for _, inputSymbolTransitionFunctions := range determinedTransitionFunctions {
					determinedStateCharHashSet := toSet(strings.Split(determinedState, ""))
					for indexOfTransitionFunction, transitionFunction := range inputSymbolTransitionFunctions {
						if maps.Equal(toSet(strings.Split(strings.Replace(transitionFunction, ",", "", -1), "")), determinedStateCharHashSet) {
							inputSymbolTransitionFunctions[indexOfTransitionFunction] = determinedState
						}
					}
				}
				continue
			}
			determinedState = sortString(determinedState)
			determinedOutputAlphabet = append(determinedOutputAlphabet, "")
			determinedStatesCharHashSet := toSet(strings.Split(newState, ","))
			for i := 0; i < len(finishStates); i++ {
				if _, contains := determinedStatesCharHashSet[finishStates[i]]; contains {
					determinedOutputAlphabet[len(determinedOutputAlphabet)-1] = FINISH_OUTPUT_SYMBOL
				}
			}
			states := strings.Split(newState, ",")
			for _, inputSymbolTransitionFunctions := range determinedTransitionFunctions {
				inputSymbolTransitionFunctions = append(inputSymbolTransitionFunctions, "")
			}
			for _, state := range states {
				indexOfState := slices.Index(determinedStates, state)
				for _, inputSymbolTransitionFunctions := range determinedTransitionFunctions {
					transitionFunction := inputSymbolTransitionFunctions[indexOfState]
					if transitionFunction == "" {
						continue
					}
					if val, ok := eclosures[transitionFunction]; ok {
						transitionFunction = val
					} else {
						if len(eclosures) != 0 {
							transitionFunctionStates := strings.Split(transitionFunction, ",")
							transitionFunction = ""
							for i := 0; i < len(transitionFunctionStates); i++ {
								if transitionFunction == "" {
									transitionFunction += eclosures[transitionFunctionStates[i]]
								} else if !strings.Contains(transitionFunction, transitionFunctionStates[i]) {
									transitionFunction += "," + eclosures[transitionFunctionStates[i]]
								}
							}
						}
					}
					if inputSymbolTransitionFunctions[len(inputSymbolTransitionFunctions)-1] == "" {
						inputSymbolTransitionFunctions[len(inputSymbolTransitionFunctions)-1] += transitionFunction
					} else {
						if !strings.Contains(inputSymbolTransitionFunctions[len(inputSymbolTransitionFunctions)-1], transitionFunction) {
							transitions := strings.Split(transitionFunction, ",")
							for _, trans := range transitions {
								if !strings.Contains(inputSymbolTransitionFunctions[len(inputSymbolTransitionFunctions)-1], trans) {
									inputSymbolTransitionFunctions[len(inputSymbolTransitionFunctions)-1] += "," + trans
								}
							}
						}
					}
				}
			}
			determinedStates = append(determinedStates, determinedState)
		}

		if len(newStates) == 0 {
			break
		}
	}

	newStatesToDeterminedStates := make(map[string]string)
	for i := len(eclosures); i < len(determinedStates); i++ {
		newStatesToDeterminedStates[determinedStates[i]] = "S" + strconv.Itoa(i-len(eclosures))
	}

	newTransitionFunctions := make([][]string, 0)
	for _, inputSymbolTransitionFunctions := range determinedTransitionFunctions {
		newTransitionFunctions = append(newTransitionFunctions, inputSymbolTransitionFunctions[len(eclosures):])
	}

	determinedOutputAlphabet = determinedOutputAlphabet[len(eclosures):]

	for _, inputSymbolTransitionFunctions := range newTransitionFunctions {
		for i := 0; i < len(inputSymbolTransitionFunctions); i++ {
			if inputSymbolTransitionFunctions[i] != "" {
				inputSymbolTransitionFunctions[i] = newStatesToDeterminedStates[sortString(inputSymbolTransitionFunctions[i])]
			}
		}
	}

	m.InputAlphabet = determinedInputAlphabet

	m.States = make([]string, 0, len(newStatesToDeterminedStates))
	for _, value := range newStatesToDeterminedStates {
		m.States = append(m.States, value)
	}

	m.TransitionFunctions = newTransitionFunctions
	m.OutputAlphabet = determinedOutputAlphabet
}

func sortString(s string) string {
    r := []rune(s)
    sort.Slice(r, func(i, j int) bool {
        return r[i] < r[j]
    })
    return string(r)
}

func toSet(s []string) map[string]struct{} {
    res := make(map[string]struct{})
    for _, item := range s {
        res[item] = struct{}{}
    }
    return res
}