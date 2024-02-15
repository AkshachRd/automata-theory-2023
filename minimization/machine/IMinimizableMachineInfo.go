package machine

type IMinimizableMachineInfo interface {
	IMachineInfo
	Minimize()
}
