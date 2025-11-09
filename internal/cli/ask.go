package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	s "github.com/mmfallacy/flakeup/internal/style"
)

var scanner = bufio.NewScanner(os.Stdin)

var ErrAskInvalidAnswer = errors.New("answer is outside choices")

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

	for _, choice := range choices {
		if len(answer) < 1 {
			return "", ErrAskInvalidAnswer
		}
		if answer == string(choice) {
			return choice, nil
		}
		// match prefix
		if len(choice) >= len(answer) && string(choice)[:len(answer)] == answer {
			return choice, nil
		}
	}

	return "", ErrAskInvalidAnswer
}
