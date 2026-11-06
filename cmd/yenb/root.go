/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package yenb

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/srtkkou/yenb/lib"
)

// Root command.
type RootCmd struct {
	fs     *pflag.FlagSet // Flag set.
	isHelp bool           // Help flag.
}

// Sub command interface.
type subCmd interface {
	usage() string
	exec(args []string) error
}

func NewRootCmd() *RootCmd {
	c := &RootCmd{
		fs: pflag.NewFlagSet("yenb", pflag.ExitOnError),
	}
	// Added flags.
	c.fs.BoolVarP(&c.isHelp, "help", "h", false,
		`Show help message.`)
	return c
}

func (c *RootCmd) Exec(args []string) error {
	// Check args.
	if len(args) < 2 {
		fmt.Println(c.usage())
		return nil
	}
	// Show usage if nothing is given after help subcommand.
	if args[1] == "help" && len(args) == 2 {
		fmt.Println(c.usage())
		return nil
	}
	// Initialize subcommand.
	var sub subCmd
	switch args[1] {
	case "help":
		fmt.Println("HELP")
		return nil
	case "install":
		sub = newInstallCmd(c)
	case "uninstall":
		sub = newUninstallCmd(c)
	default:
		fmt.Println(c.usage())
		return nil
	}
	// Initialize log.
	config := lib.Config()
	configBlob, _ := config.MarshalJSON()
	lib.Log().InfoContext(config.Context(), "yenb:start",
		slog.String("subcommand", args[1]),
		slog.Any("config", json.RawMessage(configBlob)))
	// Execute command.
	return sub.exec(os.Args[1:])
}

func (c *RootCmd) usage() string {
	var sb strings.Builder
	fmt.Fprintln(&sb, "usage: yenb <subcommand> <options>")
	fmt.Fprintln(&sb)
	c.writeVersion(&sb)
	fmt.Fprintln(&sb)
	fmt.Fprintln(&sb, "subcommands:")
	fmt.Fprintln(&sb, "  help      : Show help for sumcommand.")
	fmt.Fprintln(&sb, "  install   : Install yenb.")
	fmt.Fprintln(&sb, "  uninstall : Uninstall yenb.")
	fmt.Fprintln(&sb)
	fmt.Fprintln(&sb, "options:")
	fmt.Fprintln(&sb, c.fs.FlagUsages())
	return sb.String()
}

func (c *RootCmd) writeVersion(w io.Writer) {
	config := lib.Config()
	fmt.Fprintf(w, "version: %s",
		config.Version())
	fmt.Fprintln(w)
	fmt.Fprintf(w, "commit:  %s",
		config.CommitHash())
	fmt.Fprintln(w)
	fmt.Fprintf(w, "         %s",
		config.CommitAt().Format(time.RFC3339))
	fmt.Fprintln(w)
}
