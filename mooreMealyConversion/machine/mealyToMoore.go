package machine

func ConvertMealyToMoore(mealy *MealyMachine) MooreMachine {
	var moore MooreMachine
	moore.States = make(map[MooreState]bool)
	moore.Transitions = make(Transitions[MooreTransition])

	moore.InputSymbolsNum = mealy.InputSymbolsNum
	for inputSymbol, mealyTransition := range mealy.Transitions {
		for mealyState, mealyTransitionOutput := range mealyTransition {
			mooreState := MooreState{
				Name:         mealyState.Name,
				OutputSymbol: mealyTransitionOutput.OutputSymbol,
			}
			moore.States[mooreState] = true

			if mealy.CurrentState.Name == mooreState.Name {
				moore.CurrentState = mooreState
			}
			if _, ok := moore.Transitions[inputSymbol]; !ok {
				moore.Transitions[inputSymbol] = make(MooreTransition)
			}

			moore.Transitions[inputSymbol][mooreState] = mealyTransitionOutput.State.Name
		}
	}

	return moore
}
