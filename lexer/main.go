package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Lexer > ")

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		lexer := NewLexer(response)
		tokens, err := lexer.MakeTokens()
		if err != nil {
			log.Fatal(err)
		}

		for _, token := range tokens {
			PrintToken(token)
		}
		fmt.Printf("\n")
	}
}

func PrintToken(token *Token) {
	fmt.Printf("%s", token.Type)
	if token.Value != nil {
		if f, ok := token.Value.(float32); ok {
			fmt.Printf(":%f", f)
		} else if i, ok := token.Value.(int); ok {
			fmt.Printf(":%d", i)
		}
	}
	fmt.Printf(", ")
}
