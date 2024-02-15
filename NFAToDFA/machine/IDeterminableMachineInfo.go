package machine

type IDeterminableMachineInfo interface {
	IMachineInfo
	Determine()
}
