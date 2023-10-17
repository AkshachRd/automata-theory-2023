package machine

import (
	"fmt"
	"sort"
)

func ConvertMealyToMoore(mealy *MealyMachine) MooreMachine {
	var moore MooreMachine
	moore.States = make(map[MooreState]bool)
	moore.Transitions = make(Transitions[MooreTransition])

	moore.InputSymbolsNum = mealy.InputSymbolsNum
	for _, mealyTransition := range mealy.Transitions {
		for _, mealyTransitionOutput := range mealyTransition {
			mooreState := MooreState{
				Name:         mealyTransitionOutput.State.Name,
				OutputSymbol: mealyTransitionOutput.OutputSymbol,
			}
			moore.States[mooreState] = true
		}
	}

	var sortedStates []MooreState
	for state := range moore.States {
		sortedStates = append(sortedStates, state)
	}
	sort.Slice(sortedStates, func(i, j int) bool {
		var n, m, q, k int
		fmt.Sscanf(sortedStates[i].Name, "s%d", &n)
		fmt.Sscanf(sortedStates[j].Name, "s%d", &m)
		fmt.Sscanf(string(sortedStates[i].OutputSymbol), "y%d", &q)
		fmt.Sscanf(string(sortedStates[j].OutputSymbol), "y%d", &k)

		if n == m {
			return q < k
		}
		return n < m
	})

	moore.States = make(map[MooreState]bool)
	for i, mooreState := range sortedStates {
		currentMooreState := MooreState{
			Name:         fmt.Sprintf("s%d", i),
			OutputSymbol: mooreState.OutputSymbol,
		}
		moore.States[currentMooreState] = true
		for inputSymbol, transition := range mealy.Transitions {
			for mealyState, mealyTransitionOutput := range transition {
				if _, ok := moore.Transitions[inputSymbol]; !ok {
					moore.Transitions[inputSymbol] = make(MooreTransition)
				}
				if mooreState.Name == mealyState.Name {
					state := MooreState{
						Name:         mealyTransitionOutput.State.Name,
						OutputSymbol: mealyTransitionOutput.OutputSymbol,
					}
					for j, q := range sortedStates {
						if q == state {
							moore.Transitions[inputSymbol][currentMooreState] = fmt.Sprintf("s%d", j)
						}
					}
				}
			}
		}
	}

	sortedStates = make([]MooreState, 0)
	for state := range moore.States {
		sortedStates = append(sortedStates, state)
	}
	sort.Slice(sortedStates, func(i, j int) bool {
		var n, m int
		fmt.Sscanf(sortedStates[i].Name, "s%d", &n)
		fmt.Sscanf(sortedStates[j].Name, "s%d", &m)

		return n < m
	})

	moore.CurrentState = sortedStates[0]

	return moore
}
