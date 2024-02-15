package moore

import (
	"errors"
	"slices"
	"strconv"
	"strings"
)

type MooreMachineInfo struct {
    States              []string
    TransitionFunctions [][]string
    OutputAlphabet      []string
    InputAlphabet       []string
}

func NewMooreMachineInfo(info []string) (*MooreMachineInfo, error) {
    m := &MooreMachineInfo{
        OutputAlphabet:      []string{},
        States:              []string{},
        InputAlphabet:       []string{},
        TransitionFunctions: [][]string{},
    }

    if info == nil || len(info) == 0 {
        return m, nil
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

    if len(m.States) <= 1 {
        return nil, errors.New("Incorrect input machine. Number of states can not be less than 2")
    }

    return m, nil
}

func (m *MooreMachineInfo) GetCsvData() string {
    csvData := ";"
    csvData += strings.Join(m.OutputAlphabet, ";")
    csvData += "\n"

    csvData += ";"
    csvData += strings.Join(m.States, ";")
    csvData += "\n"

    for i := 0; i < len(m.InputAlphabet); i++ {
        csvData += m.InputAlphabet[i] + ";"
        csvData += strings.Join(m.TransitionFunctions[i], ";")
        csvData += "\n"
    }

    return csvData
}

func (m *MooreMachineInfo) Minimize() {
    m.deleteUnreachableStates()
    previousMatchingMinimizedStatesToStates := make(map[string]map[string]struct{})
    matchingMinimizedStatesToStates := make(map[string]map[string]struct{})
    for outputLetter := range toSet(m.OutputAlphabet) {
        outputLetterStates := make(map[string]struct{})
        for i := range m.OutputAlphabet {
            if m.OutputAlphabet[i] == outputLetter {
                outputLetterStates[m.States[i]] = struct{}{}
            }
        }
        matchingMinimizedStatesToStates[outputLetter] = outputLetterStates
    }
    newTransitionFunctions := m.getNewTransitionFunctions(matchingMinimizedStatesToStates, m.TransitionFunctions)

    for len(matchingMinimizedStatesToStates) != len(previousMatchingMinimizedStatesToStates) && len(m.States) != len(matchingMinimizedStatesToStates) {
        previousMatchingMinimizedStatesToStates = matchingMinimizedStatesToStates
        matchingMinimizedStatesToStates = m.getMatchingMinimizedStatesToStates(previousMatchingMinimizedStatesToStates, newTransitionFunctions)
        newTransitionFunctions = m.getNewTransitionFunctions(matchingMinimizedStatesToStates, m.TransitionFunctions)
    }

    minimizedStates := make([]string, 0, len(matchingMinimizedStatesToStates))
	for key := range matchingMinimizedStatesToStates {
		minimizedStates = append(minimizedStates, key)
	}
    minimizedTransitionFunctions := m.getMinimizedTransitionFunctions(matchingMinimizedStatesToStates)
    minimizedOutputAlphabet := m.getMinimizedOutputAlphabet(matchingMinimizedStatesToStates)
    m.States = minimizedStates
    m.TransitionFunctions = minimizedTransitionFunctions
    m.OutputAlphabet = minimizedOutputAlphabet
}

func toSet(slice []string) map[string]struct{} {
    set := make(map[string]struct{})
    for _, item := range slice {
        set[item] = struct{}{}
    }
    return set
}

func (m *MooreMachineInfo) deleteUnreachableStates() {
    reachableStates := make(map[string]struct{})

    for _, innerStateTransitionFunctions := range m.TransitionFunctions {
        for i, transitionFunction := range innerStateTransitionFunctions {
            if transitionFunction != m.States[i] || i == 0 {
                reachableStates[transitionFunction] = struct{}{}
            }
        }
    }

    removedStates := []string{}
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
        m.OutputAlphabet = slices.Delete(m.OutputAlphabet, indexOfRemovedState, indexOfRemovedState + 1)
    }
}

func (m *MooreMachineInfo) getMatchingMinimizedStatesToStates(matchingEquivalenceClassesToStates map[string]map[string]struct{}, transitionFunctions [][]string) map[string]map[string]struct{} {
    matchingNewStatesToPreviousStates := make(map[string]map[string]struct{})
    matchingNewStatesToTransitionFunctions := make(map[string][]string)

    for i := 0; i < len(m.States); i++ {
        transitionsSequence := []string{}
        for _, innerStateTransitionFunctions := range transitionFunctions {
            transitionsSequence = append(transitionsSequence, innerStateTransitionFunctions[i])
        }
        isExistMinimizedState := false
        for matchingNewState, transitionFunctions := range matchingNewStatesToTransitionFunctions {
            if slices.Equal(transitionFunctions, transitionsSequence) {
                firstElementOfEquivalenceClass := ""
                for element := range matchingNewStatesToPreviousStates[matchingNewState] {
                    firstElementOfEquivalenceClass = element
                    break
                }
                firstEquivalenceClass := ""
                secondEquivalenceClass := ""
                for equivalenceClass, states := range matchingEquivalenceClassesToStates {
                    if _, ok := states[firstElementOfEquivalenceClass]; ok {
                        firstEquivalenceClass = equivalenceClass
                    }
                    if _, ok := states[m.States[i]]; ok {
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

func (m *MooreMachineInfo) getNewTransitionFunctions(matchingNewStatesToStates map[string]map[string]struct{}, oldTransitionFunctions [][]string) [][]string {
    newTransitionFunctions := make([][]string, 0)

    for _, innerStateTransitionFunctions := range oldTransitionFunctions {
        newInnerStateTransitionFunctions := make([]string, 0)
        for _, oldTransitionFunction := range innerStateTransitionFunctions {
            oldState := oldTransitionFunction
            if oldState == "" {
                newInnerStateTransitionFunctions = append(newInnerStateTransitionFunctions, "")
            } else {
                for matchingNewState := range matchingNewStatesToStates {
                    if _, ok := matchingNewStatesToStates[matchingNewState][oldState]; ok {
                        newInnerStateTransitionFunctions = append(newInnerStateTransitionFunctions, matchingNewState)
                    }
                }
            }
        }
        newTransitionFunctions = append(newTransitionFunctions, newInnerStateTransitionFunctions)
    }

    return newTransitionFunctions
}

func (m *MooreMachineInfo) getMinimizedTransitionFunctions(matchingMinimizedStatesToStates map[string]map[string]struct{}) [][]string {
    minimizedTransitionFunctions := make([][]string, 0)

    for _, innerStateTransitionFunctions := range m.TransitionFunctions {
        innerStateMinimizedTransitionFunctions := make([]string, 0)
        for _, innerStateTransitionFunction := range innerStateTransitionFunctions[:len(matchingMinimizedStatesToStates)] {
            if innerStateTransitionFunction == "" {
                innerStateMinimizedTransitionFunctions = append(innerStateMinimizedTransitionFunctions, "")
            } else {
                for matchingMinimizedState := range matchingMinimizedStatesToStates {
                    if _, ok := matchingMinimizedStatesToStates[matchingMinimizedState][innerStateTransitionFunction]; ok {
                        innerStateMinimizedTransitionFunctions = append(innerStateMinimizedTransitionFunctions, matchingMinimizedState)
                    }
                }
            }
        }
        minimizedTransitionFunctions = append(minimizedTransitionFunctions, innerStateMinimizedTransitionFunctions)
    }

    return minimizedTransitionFunctions
}

func (m *MooreMachineInfo) getMinimizedOutputAlphabet(matchingMinimizedStatesToStates map[string]map[string]struct{}) []string {
    minimizedOutputAlphabet := make([]string, 0)

    for _, matchingMinimizedStateToStates := range matchingMinimizedStatesToStates {
        for state := range matchingMinimizedStateToStates {
            for i, s := range m.States {
                if s == state {
                    minimizedOutputAlphabet = append(minimizedOutputAlphabet, m.OutputAlphabet[i])
                    break
                }
            }
            break
        }
    }

    return minimizedOutputAlphabet
}