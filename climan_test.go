package climan_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"moul.io/climan"
	"moul.io/u"
)

func TestCliman(t *testing.T) {

	cases := []struct {
		name                 string
		root                 *climan.Command
		args                 []string
		output               string
		parseShouldErrFlag   bool
		parseShouldOtherFail bool
		runShouldErrFlag     bool
		runShouldOtherFail   bool
	}{
		{
			name: "empty",
			root: &climan.Command{},
			args: []string{},
			output: `
USAGE
  climan.test
`[1:],
			runShouldErrFlag: true,
		}, {
			name:   "example-no-args",
			root:   exampleCommand(),
			args:   []string{},
			output: "debug false\nargs []\n",
		}, {
			name:   "example-flags",
			root:   exampleCommand(),
			args:   []string{"-debug"},
			output: "debug true\nargs []\n",
		}, {
			name:   "example-args",
			root:   exampleCommand(),
			args:   []string{"foo", "bar"},
			output: "debug false\nargs [foo bar]\n",
		}, {
			name:   "example-flags-and-args",
			root:   exampleCommand(),
			args:   []string{"-debug", "foo", "bar"},
			output: "debug true\nargs [foo bar]\n",
		}, {
			name: "example-help",
			root: exampleCommand(),
			args: []string{"-h"},
			output: `
USAGE
  example

FLAGS
  -debug false  debug mode
`[1:],
			parseShouldErrFlag: true,
		}, {
			name:   "example-sub-no-args",
			root:   exampleCommandWithSubcommands(),
			args:   []string{},
			output: "debug false\nfoo \nbar 0\nargs []\nfoo called false\nbar called false\n",
		}, {
			name:   "example-sub-flags",
			root:   exampleCommandWithSubcommands(),
			args:   []string{"-debug"},
			output: "debug true\nfoo \nbar 0\nargs []\nfoo called false\nbar called false\n",
		}, {
			name:   "example-sub-args",
			root:   exampleCommandWithSubcommands(),
			args:   []string{"blah", "blih"},
			output: "debug false\nfoo \nbar 0\nargs [blah blih]\nfoo called false\nbar called false\n",
		}, {
			name:   "example-sub-flags-and-args",
			root:   exampleCommandWithSubcommands(),
			args:   []string{"-debug", "blah", "blih"},
			output: "debug true\nfoo \nbar 0\nargs [blah blih]\nfoo called false\nbar called false\n",
		}, {
			name:   "example-sub-foo-no-args",
			root:   exampleCommandWithSubcommands(),
			args:   []string{"foo"},
			output: "debug false\nfoo \nbar 0\nargs []\nfoo called true\nbar called false\n",
		}, {
			name:   "example-sub-foo-pre-flags",
			root:   exampleCommandWithSubcommands(),
			args:   []string{"-debug", "foo", "-flag=baz"},
			output: "debug true\nfoo baz\nbar 0\nargs []\nfoo called true\nbar called false\n",
		}, {
			name:   "example-sub-foo-post-flags",
			root:   exampleCommandWithSubcommands(),
			args:   []string{"foo", "-debug", "-flag=baz"},
			output: "debug true\nfoo baz\nbar 0\nargs []\nfoo called true\nbar called false\n",
		}, {
			name:   "example-sub-foo-cancel-flags",
			root:   exampleCommandWithSubcommands(),
			args:   []string{"-debug", "foo", "-debug=false"},
			output: "debug false\nfoo \nbar 0\nargs []\nfoo called true\nbar called false\n",
		}, {
			name:   "example-sub-foo-args",
			root:   exampleCommandWithSubcommands(),
			args:   []string{"foo", "bar", "baz"},
			output: "debug false\nfoo \nbar 0\nargs [bar baz]\nfoo called true\nbar called false\n",
		}, {
			name:   "example-sub-foo-args-and-flags",
			root:   exampleCommandWithSubcommands(),
			args:   []string{"-debug", "foo", "-debug", "-flag=baz", "bar", "baz"},
			output: "debug true\nfoo baz\nbar 0\nargs [bar baz]\nfoo called true\nbar called false\n",
		}, {
			name: "example-sub-help",
			root: exampleCommandWithSubcommands(),
			args: []string{"-h"},
			output: `
USAGE
  example

SUBCOMMANDS
  foo  dolor sit amet
  bar  dolor sit amet2

FLAGS
  -debug false  debug mode
`[1:],
			parseShouldErrFlag: true,
		}, {
			name: "example-foo-help",
			root: exampleCommandWithSubcommands(),
			args: []string{"-debug", "foo", "-h"},
			output: `
USAGE
  lorem ipsum

FLAGS
  -debug true  debug mode
  -flag ...    foo's flag
`[1:],
			parseShouldErrFlag: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := tc.root

			// parse
			{
				done, err := u.CaptureStdoutAndStderr()
				require.NoError(t, err)

				err = root.Parse(tc.args)
				out := done()

				if tc.parseShouldOtherFail {
					require.Error(t, err)
					return
				}

				if tc.parseShouldErrFlag {
					require.True(t, errors.Is(err, flag.ErrHelp))
					require.Equal(t, out, tc.output)
					return
				}

				require.NoError(t, err)
			}

			// run
			{
				done, err := u.CaptureStdoutAndStderr()
				require.NoError(t, err)

				err = root.Run(context.Background())
				out := done()

				if tc.runShouldOtherFail {
					require.Error(t, err)
					return
				}

				if tc.runShouldErrFlag {
					require.True(t, errors.Is(err, flag.ErrHelp))
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, out, tc.output)
			}
		})
	}
}

func exampleCommand() *climan.Command {
	var opts struct {
		Debug bool
	}
	return &climan.Command{
		Name: "example",
		FlagSetBuilder: func(fs *flag.FlagSet) {
			fs.BoolVar(&opts.Debug, "debug", opts.Debug, "debug mode")
		},
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println("debug", opts.Debug)
			fmt.Println("args", args)
			return nil
		},
	}
}

func exampleCommandWithSubcommands() *climan.Command {
	var opts struct {
		Debug     bool
		FooFlag   string
		BarFlag   int
		FooCalled bool
		BarCalled bool
	}
	exec := func(ctx context.Context, args []string) error {
		fmt.Println("debug", opts.Debug)
		fmt.Println("foo", opts.FooFlag)
		fmt.Println("bar", opts.BarFlag)
		fmt.Println("args", args)
		fmt.Println("foo called", opts.FooCalled)
		fmt.Println("bar called", opts.BarCalled)
		return nil
	}
	commonFlagsSetup := func(fs *flag.FlagSet) {
		fs.BoolVar(&opts.Debug, "debug", opts.Debug, "debug mode")
	}
	return &climan.Command{
		Name:           "example",
		FlagSetBuilder: commonFlagsSetup,
		Exec:           exec,
		Subcommands: []*climan.Command{
			{
				Name:       "foo",
				ShortUsage: "lorem ipsum",
				ShortHelp:  "dolor sit amet",
				FlagSetBuilder: func(fs *flag.FlagSet) {
					commonFlagsSetup(fs)
					opts.FooCalled = true
					fs.StringVar(&opts.FooFlag, "flag", opts.FooFlag, "foo's flag")
				},
				Exec: exec,
			}, {
				Name:       "bar",
				ShortUsage: "lorem ipsum2",
				ShortHelp:  "dolor sit amet2",
				FlagSetBuilder: func(fs *flag.FlagSet) {
					opts.BarCalled = true
					fs.IntVar(&opts.BarFlag, "flag", opts.BarFlag, "bar's flag")
				},
				Exec: exec,
			},
		},
	}
}
