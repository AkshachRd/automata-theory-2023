package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/AkshachRd/automata-theory-2023/NFAToDFA/machine"
	"github.com/AkshachRd/automata-theory-2023/NFAToDFA/moore"
)

type Args struct {
	SourceFilePath      string
	DestinationFilePath string
}

func NewArgs(sourceFilePath, destinationFilePath string) (*Args, error) {
	return &Args{
		SourceFilePath:      sourceFilePath,
		DestinationFilePath: destinationFilePath,
	}, nil
}

func ParseArgs(args []string) (*Args, error) {
	if len(args) != 2 {
		return nil, errors.New("incorrect arguments count")
	}

	return NewArgs(args[0], args[1])
}

func getInfoFromFile(filePath string) ([]string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    return lines, scanner.Err()
}

func printDataToFile(data, filePath string) error {
    return os.WriteFile(filePath, []byte(data), 0644)
}

func processData(infoFromFile []string) (machine.IMachineInfo, error) {
    machineInfo := moore.NewMooreMachineInfo(infoFromFile)
    machineInfo.Determine()

    return machineInfo, nil
}

func main() {
    parsedArgs, err := ParseArgs(os.Args[1:])
    if err != nil {
        fmt.Println(err)
        return
    }

    infoFromFile, err := getInfoFromFile(parsedArgs.SourceFilePath)
    if err != nil {
        fmt.Println(err)
        return
    }

    machineInfo, err := processData(infoFromFile)
    if err != nil {
        fmt.Println(err)
        return
    }

    csvData := machineInfo.GetCsvData()

    err = printDataToFile(csvData, parsedArgs.DestinationFilePath)
    if err != nil {
        fmt.Println(err)
        return
    }
}