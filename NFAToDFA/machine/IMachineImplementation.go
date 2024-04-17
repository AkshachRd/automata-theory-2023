package machine

type IMachineImplementation interface {
	Draw(outputFileName string)
	ReadFromFile(filePath string) error
	Print(outputFileName string) error
	Determine() error
	Minimize() error
}
