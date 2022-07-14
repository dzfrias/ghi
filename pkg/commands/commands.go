// Package commands holds the underlying functions called when the user uses a
// command.
// Ex:
//    `ghi list`
//    Calls commands.List
package commands

import "fmt"

type errNotEnoughArgs struct {
	command string
}

func (e errNotEnoughArgs) Error() string {
	return fmt.Sprintf("not enough args to command `%s`", e.command)
}
