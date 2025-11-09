package style

import "github.com/jedib0t/go-pretty/v6/text"

var (
	Info  = text.FgBlue.Sprint
	Infof = text.FgBlue.Sprintf

	Warn  = text.FgYellow.Sprint
	Warnf = text.FgYellow.Sprintf

	Err  = text.FgRed.Sprint
	Errf = text.FgRed.Sprintf

	Mkdir    = func(a ...interface{}) string { return text.FgCyan.Sprint(prepend(Icons.Dir+" ", a...)...) }
	Ignore   = func(a ...interface{}) string { return text.Faint.Sprint(prepend(Icons.Ignore+" ", a...)...) }
	Conflict = func(a ...interface{}) string { return text.FgHiRed.Sprint(prepend(Icons.Conflict+" ", a...)...) }
	Clean    = func(a ...interface{}) string { return text.FgHiGreen.Sprint(prepend(Icons.Clean+" ", a...)...) }
)

func prepend(prepend string, args ...interface{}) []interface{} {
	return append([]any{prepend}, args...)
}

type icons struct {
	Dir      string
	Warn     string
	Err      string
	Clean    string
	Conflict string
	Ignore   string
	Ask      string
}

var Icons = icons{
	Dir:      "",
	Warn:     "",
	Err:      "",
	Clean:    "",
	Conflict: "󰬳",
	Ignore:   "⊘",
	Ask:      "?",
}
