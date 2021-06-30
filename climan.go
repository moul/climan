package climan

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/peterbourgon/ff/v3"
)

type Command struct {
	Name           string
	Exec           func(context.Context, []string) error
	FlagSetBuilder func(fs *flag.FlagSet)
	Subcommands    []*Command
	ShortUsage     string
	ShortHelp      string
	LongHelp       string
	FFOptions      []ff.Option
	FlagSet        *flag.FlagSet
	UsageFunc      func(c *Command) string

	// internal
	selected *Command
	args     []string
}

func (c *Command) Parse(args []string) error {
	if c.selected != nil {
		return nil
	}

	if c.Name == "" && len(os.Args) > 0 {
		c.Name = filepath.Base(os.Args[0])
	}

	if c.FlagSet == nil {
		c.FlagSet = flag.NewFlagSet(c.Name, flag.ContinueOnError)
	}

	if c.FlagSetBuilder != nil {
		c.FlagSetBuilder(c.FlagSet)
	}

	if c.UsageFunc == nil {
		c.UsageFunc = DefaultUsageFunc
	}

	c.FlagSet.Usage = func() {
		fmt.Fprintln(c.FlagSet.Output(), c.UsageFunc(c))
	}

	if err := ff.Parse(c.FlagSet, args, c.FFOptions...); err != nil {
		return err
	}

	c.args = c.FlagSet.Args()
	if len(c.args) > 0 {
		for _, subcommand := range c.Subcommands {
			if strings.EqualFold(c.args[0], subcommand.Name) {
				c.selected = subcommand
				return subcommand.Parse(c.args[1:])
			}
		}
	}

	c.selected = c
	if c.Exec == nil {
		c.Exec = func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		}
	}
	return nil
}

func (c *Command) Run(ctx context.Context) error {
	if c.selected == nil {
		return errors.New("command is unparsed, cannot run")
	}

	if c.selected == c {
		err := c.selected.Exec(ctx, c.args)
		if err != nil {
			if errors.Is(err, flag.ErrHelp) {
				c.FlagSet.Usage()
			}
			return err
		}
		return nil
	}

	return c.selected.Run(ctx)
}

func DefaultUsageFunc(c *Command) string {
	var b strings.Builder

	fmt.Fprintf(&b, "USAGE\n")
	if c.ShortUsage != "" {
		fmt.Fprintf(&b, "  %s\n", c.ShortUsage)
	} else {
		fmt.Fprintf(&b, "  %s\n", c.Name)
	}
	fmt.Fprintf(&b, "\n")

	if c.LongHelp != "" {
		fmt.Fprintf(&b, "%s\n\n", c.LongHelp)
	}

	if len(c.Subcommands) > 0 {
		fmt.Fprintf(&b, "SUBCOMMANDS\n")
		tw := tabwriter.NewWriter(&b, 0, 2, 2, ' ', 0)
		for _, subcommand := range c.Subcommands {
			fmt.Fprintf(tw, "  %s\t%s\n", subcommand.Name, subcommand.ShortHelp)
		}
		tw.Flush()
		fmt.Fprintf(&b, "\n")
	}

	if countFlags(c.FlagSet) > 0 {
		fmt.Fprintf(&b, "FLAGS\n")
		tw := tabwriter.NewWriter(&b, 0, 2, 2, ' ', 0)

		c.FlagSet.VisitAll(func(f *flag.Flag) {
			def := f.DefValue
			if def == "" {
				def = "..."
			}
			fmt.Fprintf(tw, "  -%s %s\t%s\n", f.Name, def, f.Usage)
		})
		tw.Flush()
		fmt.Fprintf(&b, "\n")
	}

	return strings.TrimSpace(b.String())
}

func countFlags(fs *flag.FlagSet) int {
	var n int
	fs.VisitAll(func(*flag.Flag) { n++ })
	return n
}
