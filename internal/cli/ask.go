package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	s "github.com/mmfallacy/flakeup/internal/style"
	"github.com/mmfallacy/flakeup/internal/utils"
)

var scanner = bufio.NewScanner(os.Stdin)

var ErrAskInvalidAnswer = errors.New("answer is outside choices")

// TODO: should this be named pick better? or askFromChoices?
func ask[T ~string](question string, choices []T) (T, error) {
	prettyChoices := ""

	for i, choice := range choices {
		if i > 0 {
			prettyChoices += " | "
		}
		prettyChoices += "[" + string(choice[:1]) + "]" + string(choice[1:])
	}

	fmt.Print(s.Warnf("%s\n\t", question), s.Info(prettyChoices), fmt.Sprintf(" %s ", s.Icons.Ask))

	if !scanner.Scan() {
		return "", scanner.Err()
	}

	answer := strings.ToLower(strings.TrimSpace(scanner.Text()))

	if mapped, ok := utils.LooseMapStringToType(answer, choices); ok {
		return mapped, nil
	} else {
		return "", ErrAskInvalidAnswer
	}
}

// open ended ask, return string of result instead
func prompt(question, d string) (string, error) {
	fmt.Print(s.Infof(" %s", question))
	if d != "" {
		fmt.Print(s.Infof(" [%s]", d))
	}
	fmt.Print(s.Info(": "))
	if !scanner.Scan() {
		return "", scanner.Err()
	}

	return scanner.Text(), nil
}
