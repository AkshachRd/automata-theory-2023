package mealy

import (
	"errors"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type MealyMachineInfo struct {
	States              []string
	InnerStates         []string
	TransitionFunctions [][]string
}

func NewMealyMachineInfo(info []string) (*MealyMachineInfo, error) {
	m := &MealyMachineInfo{
		States:              make([]string, 0),
		InnerStates:         make([]string, 0),
		TransitionFunctions: make([][]string, 0),
	}

	if len(info) == 0 {
		return m, nil
	}

	dirtyStates := strings.Split(info[0], ";")
	m.States = append(m.States, dirtyStates[1:]...)

	for i := 1; i < len(info); i++ {
		values := strings.Split(info[i], ";")
		m.InnerStates = append(m.InnerStates, values[0])
		m.TransitionFunctions = append(m.TransitionFunctions, values[1:])
	}

	if len(m.States) <= 1 {
		return nil, errors.New("incorrect input machine. Number of states can not be less than 2")
	}

	return m, nil
}

func (m *MealyMachineInfo) GetCsvData() string {
	sort.Strings(m.States)

	csvData := ";"
	csvData += strings.Join(m.States, ";")
	csvData += "\n"

	for i := 0; i < len(m.InnerStates); i++ {
		csvData += m.InnerStates[i] + ";"
		for j := 0; j < len(m.TransitionFunctions[i]); j++ {
            csvData += m.TransitionFunctions[i][j]
            if j != len(m.TransitionFunctions[i])-1 {
                csvData += ";"
            }
        }
		csvData += "\n"
	}

	return csvData
}

func (m *MealyMachineInfo) Minimize() {
    m.deleteUnreachableStates()
    previousMatchingMinimizedStatesToStates := make(map[string]map[string]struct{})
    matchingMinimizedStatesToStates := make(map[string]map[string]struct{})
    for _, state := range m.States {
        matchingMinimizedStatesToStates[state] = make(map[string]struct{})
        for _, fillingState := range m.States {
            matchingMinimizedStatesToStates[state][fillingState] = struct{}{}
        }
    }
    newTransitionFunctions := m.TransitionFunctions

	previousMatchingMinimizedStatesToStates = matchingMinimizedStatesToStates
	matchingMinimizedStatesToStates = m.getMatchingMinimizedStatesToStates(previousMatchingMinimizedStatesToStates, newTransitionFunctions)
	newTransitionFunctions = m.getNewTransitionFunctions(matchingMinimizedStatesToStates, m.TransitionFunctions)
    for len(matchingMinimizedStatesToStates) != len(previousMatchingMinimizedStatesToStates) && len(m.States) != len(matchingMinimizedStatesToStates) {
        previousMatchingMinimizedStatesToStates = matchingMinimizedStatesToStates
        matchingMinimizedStatesToStates = m.getMatchingMinimizedStatesToStates(previousMatchingMinimizedStatesToStates, newTransitionFunctions)
        newTransitionFunctions = m.getNewTransitionFunctions(matchingMinimizedStatesToStates, m.TransitionFunctions)
    }

    minimizedStates := make([]string, 0)
    for key := range matchingMinimizedStatesToStates {
        minimizedStates = append(minimizedStates, key)
    }
    minimizedTransitionFunctions := m.getMinimizedTransitionFunctions(matchingMinimizedStatesToStates)
    m.States = minimizedStates
    m.TransitionFunctions = minimizedTransitionFunctions
}

func (m *MealyMachineInfo) deleteUnreachableStates() {
	reachableStates := make(map[string]struct{})

	for _, innerStateTransitionFunctions := range m.TransitionFunctions {
		for i, transitionFunction := range innerStateTransitionFunctions {
			state := strings.Split(transitionFunction, "/")[0]
			if state != m.States[i] || i == 0 {
				reachableStates[state] = struct{}{}
			}
		}
	}

	removedStates := make([]string, 0)
	if len(reachableStates) != len(m.States) {
		for _, state := range m.States {
			if _, ok := reachableStates[state]; !ok {
				removedStates = append(removedStates, state)
			}
		}
	}

	for _, state := range removedStates {
		indexOfRemovedState := slices.Index(m.States, state)
		m.States = slices.Delete(m.States, indexOfRemovedState, indexOfRemovedState + 1)
		for i := range m.TransitionFunctions {
			m.TransitionFunctions[i] = slices.Delete(m.TransitionFunctions[i], indexOfRemovedState, indexOfRemovedState + 1)
		}
	}
}

func (m *MealyMachineInfo) getMatchingMinimizedStatesToStates(matchingEquivalenceClassesToStates map[string]map[string]struct{}, transitionFunctions [][]string) map[string]map[string]struct{} {
	matchingNewStatesToPreviousStates := make(map[string]map[string]struct{})
	matchingNewStatesToTransitionFunctions := make(map[string][]string)

	for i := range m.States {
		transitionsSequence := make([]string, 0)
		for _, innerStateTransitionFunctions := range transitionFunctions {
			if strings.Contains(innerStateTransitionFunctions[i], "/") {
				transitionsSequence = append(transitionsSequence, strings.Split(innerStateTransitionFunctions[i], "/")[1])
			} else {
				transitionsSequence = append(transitionsSequence, innerStateTransitionFunctions[i])
			}
		}
		isExistMinimizedState := false
		for matchingNewState, transitionFunctions := range matchingNewStatesToTransitionFunctions {
			if slices.Equal(transitionFunctions, transitionsSequence) {
				firstElementOfEquivalenceClass := getFirstElement(matchingNewStatesToPreviousStates[matchingNewState])
				firstEquivalenceClass := ""
				secondEquivalenceClass := ""
				for equivalenceClass, states := range matchingEquivalenceClassesToStates {
					if _, contains := states[firstElementOfEquivalenceClass]; contains {
						firstEquivalenceClass = equivalenceClass
					}
					if _, contains := states[m.States[i]]; contains {
						secondEquivalenceClass = equivalenceClass
					}
				}
				if firstEquivalenceClass == secondEquivalenceClass {
					matchingNewStatesToPreviousStates[matchingNewState][m.States[i]] = struct{}{}
					isExistMinimizedState = true
					break
				}
			}
		}
		if !isExistMinimizedState {
			newState := "q" + strconv.Itoa(len(matchingNewStatesToPreviousStates))
			matchingNewStatesToPreviousStates[newState] = map[string]struct{}{m.States[i]: {}}
			matchingNewStatesToTransitionFunctions[newState] = transitionsSequence
		}
	}

	return matchingNewStatesToPreviousStates
}

func getFirstElement(set map[string]struct{}) string {
    for k := range set {
        return k
    }
    return ""
}

func (m *MealyMachineInfo) getNewTransitionFunctions(matchingNewStatesToStates map[string]map[string]struct{}, oldTransitionFunctions [][]string) [][]string {
    newTransitionFunctions := make([][]string, 0)

    for _, innerStateTransitionFunctions := range oldTransitionFunctions {
        newInnerStateTransitionFunctions := make([]string, 0)
        for _, oldTransitionFunction := range innerStateTransitionFunctions {
            oldState := oldTransitionFunction
            if strings.Contains(oldTransitionFunction, "/") {
                oldState = strings.Split(oldTransitionFunction, "/")[0]
            }
            for matchingNewState := range matchingNewStatesToStates {
                if _, ok := matchingNewStatesToStates[matchingNewState][oldState]; ok {
                    newInnerStateTransitionFunctions = append(newInnerStateTransitionFunctions, matchingNewState)
                }
            }
        }
        newTransitionFunctions = append(newTransitionFunctions, newInnerStateTransitionFunctions)
    }

    return newTransitionFunctions
}

func (m *MealyMachineInfo) getMinimizedTransitionFunctions(
	matchingMinimizedStatesToStates map[string]map[string]struct{},
) [][]string {
	minimizedTransitionFunctions := make([][]string, 0)

	for _, states := range matchingMinimizedStatesToStates {
		innerStateMinimizedTransitionFunctions := make([]string, 0)
		for _, innerStateTransitionFunctions := range m.TransitionFunctions {
			for localMatchingMinimizedState, localStates := range matchingMinimizedStatesToStates {
				for state := range states {
					stateIndex := slices.Index(m.States, state)
					if _, contains := localStates[strings.Split(innerStateTransitionFunctions[stateIndex], "/")[0]]; stateIndex != -1 && contains {
						innerStateMinimizedTransitionFunctions = append(innerStateMinimizedTransitionFunctions, localMatchingMinimizedState+"/"+strings.Split(innerStateTransitionFunctions[stateIndex], "/")[1])
					}
				}
			}
		}
		if len(minimizedTransitionFunctions) != len(innerStateMinimizedTransitionFunctions) {
			for _, transitionFunction := range innerStateMinimizedTransitionFunctions {
				minimizedTransitionFunctions = append(minimizedTransitionFunctions, []string{transitionFunction})
			}
		} else {
			for i := range minimizedTransitionFunctions {
				minimizedTransitionFunctions[i] = append(minimizedTransitionFunctions[i], innerStateMinimizedTransitionFunctions[i])
			}
		}
	}

	return minimizedTransitionFunctions
}