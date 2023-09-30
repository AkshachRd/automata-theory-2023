package machine

func ConvertMooreToMealy(moore *MooreMachine) MealyMachine {
	var mealy MealyMachine
	mealy.States = make(map[MealyState]bool)
	mealy.Transitions = make(Transitions[MealyTransition])

	mealy.InputSymbolsNum = moore.InputSymbolsNum
	for inputSymbol, mooreTransition := range moore.Transitions {
		for mooreState, mooreTransitionOutput := range mooreTransition {
			mealyState := MealyState{
				Name: mooreState.Name,
			}
			mealy.States[mealyState] = true

			if moore.CurrentState.Name == mealyState.Name {
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

	return mealy
}
