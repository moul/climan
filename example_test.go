package climan_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"moul.io/climan"
)

func Example() {
	var opts struct {
		Debug bool
	}

	root := &climan.Command{
		Name:       "example",
		ShortUsage: "example [global flags] <subcommand> [flags] [args...]",
		ShortHelp:  "example's short help",
		LongHelp:   "example's longer help.\nwith more details.",
		FlagSetBuilder: func(fs *flag.FlagSet) {
			fs.BoolVar(&opts.Debug, "debug", opts.Debug, "debug mode")
		},
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println("args", args)
			return nil
		},
		Subcommands: []*climan.Command{
			&climan.Command{
				Name: "sub",
			},
		},
		// Options: []climan.Option{climan.WithEnvVarPrefix("EXAMPLE")},
	}
	if err := root.Parse(os.Args[1:]); err != nil {
		log.Fatal(fmt.Errorf("parse error: %w", err))
	}

	if err := root.Run(context.Background()); err != nil {
		log.Fatal(fmt.Errorf("run error: %w", err))
	}
}
