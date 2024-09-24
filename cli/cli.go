package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
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

type Command struct {
	Name    string
	Short   string
	Summary string
	Run     func() error
}

var ErrExitLoop = errors.New("exit command loop")

func RunOneCommand(ctx context.Context, commands []*Command) error {
	for {

		labels := make([]string, len(commands))
		length := 0
		for i, cmd := range commands {
			label := ""
			if cmd.Short == "" {
				label = cmd.Name
			} else if cmd.Name == "" {
				label = fmt.Sprintf("[%s]", cmd.Short)
			} else if cmd.Short == cmd.Name[:1] {
				label = fmt.Sprintf("[%s]%s", cmd.Short, cmd.Name[1:])
			} else {
				label = fmt.Sprintf("[%s] %s", cmd.Short, cmd.Name)
			}
			labels[i] = label
			if len(label) > length {
				length = len(label)
			}
		}
		for i, cmd := range commands {
			fmt.Printf("% *s : %s\n", length, labels[i], cmd.Summary)
		}

		option := Ask("Action")
		if option == "" {
			fmt.Printf("invalid option\n")
			continue
		}
		requested := strings.ToLower(option)
		var found *Command
		for _, cmd := range commands {
			if requested == cmd.Short || requested == cmd.Name {
				found = cmd
				break
			}
		}
		if found == nil {
			fmt.Printf("invalid option\n")
			continue
		}

		return found.Run()
	}

}

// Edit passes input into $EDITOR and waits for an exit
func Edit(ctx context.Context, input string) (string, error) {

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	tmpfile, err := os.CreateTemp("", "o5-edit-")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(input); err != nil {
		return "", fmt.Errorf("write temp file: %w", err)
	}

	if err := tmpfile.Close(); err != nil {
		return "", fmt.Errorf("close temp file: %w", err)
	}

	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("run editor: %w", err)
	}

	edited, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return "", fmt.Errorf("read edited file: %w", err)
	}

	return string(edited), nil
}
