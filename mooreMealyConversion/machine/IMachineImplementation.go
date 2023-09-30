package machine

import (
	"bufio"
	"mooreMealyConversion/graph"
)

type IMachineImplementation interface {
	Draw(graph graph.IGraph)
	ReadFromFile(scanner *bufio.Scanner, statesNum, inputSymbolsNum uint64) error
}
