package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var inputReader *bufio.Reader

func AskBool(prompt string) bool {
	answer := Ask(prompt + " [y/n]")
	return strings.HasPrefix(strings.ToLower(answer), "y")
}

func Ask(prompt string) string {
	if inputReader == nil {
		inputReader = bufio.NewReader(os.Stdin)
	}
	fmt.Printf("%s: ", prompt)
	answer, err := inputReader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(answer)
}

func Print(title string, msg interface{}) {
	prettyOut, _ := json.MarshalIndent(msg, " |", "  ")
	fmt.Printf("%s:\n |%s\n", title, string(prettyOut))
}
