package machine

import (
	"bufio"
	"mooreMealyConversion/graph"
	"os"
)

type IMachineImplementation interface {
	Draw(graph graph.IGraph)
	ReadFromFile(scanner *bufio.Scanner, statesNum, inputSymbolsNum uint64) error
	Print(file *os.File) error
	Minimize() error
}
