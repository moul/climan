package climan_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/peterbourgon/ff/v3"
	"moul.io/climan"
)

var opts struct {
	debug   bool
	fooFlag string
}

func Example() {
	root := &climan.Command{
		Name:       "example",
		ShortUsage: "example [global flags] <subcommand> [flags] [args...]",
		ShortHelp:  "example's short help",
		LongHelp:   "example's longer help.\nwith more details.",
		FlagSetBuilder: func(fs *flag.FlagSet) {
			fs.BoolVar(&opts.debug, "debug", opts.debug, "debug mode")
		},
		Exec: doRoot,
		Subcommands: []*climan.Command{
			&climan.Command{
				Name: "foo",
				FlagSetBuilder: func(fs *flag.FlagSet) {
					fs.BoolVar(&opts.debug, "debug", opts.debug, "debug mode")
					fs.StringVar(&opts.fooFlag, "flag", opts.fooFlag, "foo's flag")
				},
				ShortUsage: "foo [flags]",
				ShortHelp:  "foo things",
				Exec:       doFoo,
			},
		},
		FFOptions: []ff.Option{ff.WithEnvVarPrefix("EXAMPLE")},
	}

	if err := root.Parse(os.Args[1:]); err != nil {
		log.Fatal(fmt.Errorf("parse error: %w", err))
	}

	if err := root.Run(context.Background()); err != nil {
		log.Fatal(fmt.Errorf("run error: %w", err))
	}
}

func doRoot(ctx context.Context, args []string) error {
	fmt.Println("args", args)
	return nil
}

func doFoo(ctx context.Context, args []string) error {
	fmt.Println("flag", opts.fooFlag)
	return nil
}
