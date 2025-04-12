package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/table"
	"github.com/charmbracelet/x/term"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v3"
)

type Command = cli.Command

type Flag = cli.Flag
type Argument = cli.Argument

type IntFlag = cli.IntFlag
type IntSliceFlag = cli.IntSliceFlag
type StringFlag = cli.StringFlag
type StringSliceFlag = cli.StringSliceFlag
type BoolFlag = cli.BoolFlag

type FloatArg = cli.FloatArg
type IntArg = cli.IntArg
type StringArg = cli.StringArg
type StringMapArg = cli.StringMapArg
type TimestampArg = cli.TimestampArg
type UintArg = cli.UintArg

func IsTerminal() {
	term.IsTerminal(os.Stdout.Fd())
}

func Info(v ...interface{}) {
	os.Stdout.WriteString(fmt.Sprintln(v...))
}

func Infof(format string, a ...interface{}) {
	v := fmt.Sprintf(format, a...)
	Info(v)
}

func Warn(v ...any) {
	color := lipgloss.Yellow

	var style = lipgloss.NewStyle().
		Foreground(color)

	s := style.Render(fmt.Sprintln(v...))
	os.Stderr.WriteString(s + "\n")
}

func Warnf(format string, a ...interface{}) {
	v := fmt.Sprintf(format, a...)
	Warn(v)
}

func Error(v ...any) {
	color := lipgloss.Red

	var style = lipgloss.NewStyle().
		Foreground(color)

	s := style.Render(fmt.Sprintln(v...))
	os.Stderr.WriteString(s + "\n")
}

func Errorf(format string, a ...interface{}) {
	v := fmt.Sprintf(format, a...)
	Error(v)
}

func Fatal(v ...any) {
	Error(v...)
	os.Exit(1)
}

func Fatalf(format string, a ...interface{}) {
	v := fmt.Sprintf(format, a...)
	Fatal(v)
}

func ShowAppHelp(cmd *Command) error {
	return cli.ShowAppHelp(cmd)
}

func ShowCommandHelp(cmd *Command) error {
	return cli.ShowSubcommandHelp(cmd)
}

func OpenFile(name string) error {
	err := browser.OpenFile(name)

	if err != nil {
		Error("Unable to open file. try manually")
		Error(name)
	}

	return nil
}

func OpenURL(url string) error {
	err := browser.OpenURL(url)

	if err != nil {
		Error("Unable to start your browser. try manually.")
		Error(url)
	}

	return nil
}

func Select(label string, items []string) (int, string, error) {
	s := huh.NewSelect[int]()

	if label != "" {
		s.Title(label)
	}

	options := make([]huh.Option[int], 0)

	for i, item := range items {
		options = append(options, huh.NewOption(item, i))
	}

	var index int

	s.Value(&index)
	s.Options(options...)

	if err := s.Run(); err != nil {
		return 0, "", err
	}

	result := items[index]

	if result != "" {
		fmt.Println("> " + result)
	}

	return index, result, nil
}

func MustSelect(label string, items []string) (int, string) {
	index, value, err := Select(label, items)

	if err != nil {
		Fatal(err)
	}

	return index, value
}

func Prompt(label, placeholder string) (string, error) {
	i := huh.NewInput()

	if label != "" {
		i.Title(label)
	}

	if placeholder != "" {
		i.Placeholder(placeholder)
	}

	var result string
	i.Value(&result)

	if err := i.Run(); err != nil {
		return "", err
	}

	if result != "" {
		fmt.Println("> " + result)
	}

	return result, nil
}

func MustPrompt(label, placeholder string) string {
	value, err := Prompt(label, placeholder)

	if err != nil {
		Fatal(err)
	}

	return value
}

func Confirm(label string, placeholder bool) (bool, error) {
	c := huh.NewConfirm()

	if label != "" {
		c.Title(label)
	}

	var result bool
	c.Value(&result)

	return result, c.Run()
}

func MustConfirm(label string, placeholder bool) bool {
	value, err := Confirm(label, placeholder)

	if err != nil {
		Fatal(err)
	}

	return value
}

func Title(val string) {
	color := lipgloss.Green

	var style = lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Underline(true)

	fmt.Println(style.Render(val))
}

func Table(headers []string, rows [][]string) {
	table := table.New().
		Headers(headers...).
		Rows(rows...).
		Width(80)

	fmt.Println(table)
}
