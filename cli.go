package cli

import (
	"errors"
	"fmt"
	"os"

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
type StringMapArgs = cli.StringMapArgs
type TimestampArg = cli.TimestampArg
type UintArg = cli.UintArg

var ErrUserAborted = errors.New("user aborted")

func IsTerminal() bool {
	return isTerminalCheck()
}

func Info(v ...interface{}) {
	os.Stdout.WriteString(fmt.Sprintln(v...))
}

func Infof(format string, a ...interface{}) {
	v := fmt.Sprintf(format, a...)
	Info(v)
}

func Warn(v ...any) {
	s := themeWarning(fmt.Sprint(v...))
	os.Stderr.WriteString(s + "\n")
}

func Warnf(format string, a ...interface{}) {
	v := fmt.Sprintf(format, a...)
	Warn(v)
}

func Error(v ...any) {
	s := themeError(fmt.Sprint(v...))
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

func Debug(v ...any) {
	s := themeMuted(fmt.Sprint(v...))
	os.Stderr.WriteString(s + "\n")
}

func Debugf(format string, a ...interface{}) {
	v := fmt.Sprintf(format, a...)
	Debug(v)
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

func Title(val string) {
	fmt.Println(bold(underline(themeHighlight(val))))
}
