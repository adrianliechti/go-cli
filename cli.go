package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

// Commands, flags and arguments come from urfave/cli; the aliases below let a
// program build its commands with this package alone.

type Command = cli.Command
type ActionFunc = cli.ActionFunc
type Args = cli.Args

type Flag = cli.Flag
type Argument = cli.Argument

type BoolFlag = cli.BoolFlag
type DurationFlag = cli.DurationFlag
type FloatFlag = cli.FloatFlag
type FloatSliceFlag = cli.FloatSliceFlag
type IntFlag = cli.IntFlag
type IntSliceFlag = cli.IntSliceFlag
type Int64Flag = cli.Int64Flag
type UintFlag = cli.UintFlag
type UintSliceFlag = cli.UintSliceFlag
type StringFlag = cli.StringFlag
type StringSliceFlag = cli.StringSliceFlag
type StringMapFlag = cli.StringMapFlag
type TimestampFlag = cli.TimestampFlag
type GenericFlag = cli.GenericFlag

type FloatArg = cli.FloatArg
type IntArg = cli.IntArg
type IntArgs = cli.IntArgs
type StringArg = cli.StringArg
type StringArgs = cli.StringArgs
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

func Title(val string) {
	fmt.Println(bold(underline(themeHighlight(val))))
}
