package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/AkshachRd/automata-theory-2023/minimization/machine"
	"github.com/AkshachRd/automata-theory-2023/minimization/mealy"
	"github.com/AkshachRd/automata-theory-2023/minimization/moore"
)

const (
    MEALY_MINIMIZATION_TYPE = "mealy"
    MOORE_MINIMIZATION_TYPE = "moore"
)

type Args struct {
    ConversionType      string
    SourceFilePath      string
    DestinationFilePath string
}

var AvailableConversionTypes = map[string]struct{}{
	MEALY_MINIMIZATION_TYPE: {},
	MOORE_MINIMIZATION_TYPE: {},
}

func NewArgs(conversionType, sourceFilePath, destinationFilePath string) (*Args, error) {
    if _, ok := AvailableConversionTypes[strings.ToLower(conversionType)]; !ok {
        return nil, errors.New("incorrect conversion type")
    }

    return &Args{
        ConversionType:      conversionType,
        SourceFilePath:      sourceFilePath,
        DestinationFilePath: destinationFilePath,
    }, nil
}

func ParseArgs(args []string) (*Args, error) {
    if len(args) != 3 {
        return nil, errors.New("incorrect arguments count")
    }

    return NewArgs(args[0], args[1], args[2])
}

func GetInfoFromFile(filePath string) ([]string, error) {
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

func PrintDataToFile(data, filePath string) error {
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = io.WriteString(file, data)
    if err != nil {
        return err
    }

    err = file.Sync()
    if err != nil {
        return err
    }

    return nil
}

func ProcessData(infoFromFile []string, conversionType string) (machine.IMachineInfo, error) {
    var machineInfo machine.IMachineInfo

    switch strings.ToLower(conversionType) {
    case MEALY_MINIMIZATION_TYPE:
        mealyMachineInfo, err := mealy.NewMealyMachineInfo(infoFromFile)
		if err != nil {
			return nil, err
		}
        mealyMachineInfo.Minimize()
        machineInfo = mealyMachineInfo
    case MOORE_MINIMIZATION_TYPE:
        mooreMachineInfo, err := moore.NewMooreMachineInfo(infoFromFile)
		if err != nil {
			return nil, err
		}
        mooreMachineInfo.Minimize()
        machineInfo = mooreMachineInfo
    default:
        return nil, errors.New("unavailable conversion type")
    }

    return machineInfo, nil
}

func main() {
    fmt.Println(os.Args[1:])
    parsedArgs, err := ParseArgs(os.Args[1:])
    if err != nil {
        fmt.Println(err)
        return
    }

    infoFromFile, err := GetInfoFromFile(parsedArgs.SourceFilePath)
    if err != nil {
        fmt.Println(err)
        return
    }

    machineInfo, err := ProcessData(infoFromFile, parsedArgs.ConversionType)
    if err != nil {
        fmt.Println(err)
        return
    }

    err = PrintDataToFile(machineInfo.GetCsvData(), parsedArgs.DestinationFilePath)
    if err != nil {
        fmt.Println(err)
        return
    }
}