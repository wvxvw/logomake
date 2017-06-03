package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mndrix/golog"
	"github.com/mndrix/golog/read"
	"github.com/mndrix/golog/term"

	"github.com/wvxvw/logomake/logging"
	"github.com/wvxvw/logomake/native"
)

type LogomakeOpions struct {
	Makefile string `default:"Makefile.logomake"`
	Goal     string `default:"all."`
	Home     string `default:"~/.logomake"`
}

func main() {
	opts := parseFlags()
	m := initMachine(opts.Home)

	f, err := os.OpenFile(opts.Makefile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("cannot open %s: %s", opts.Makefile, err))
	}
	m = m.Consult(f)
	logging.Info.Printf("Executing %s", opts.Goal)
	goal, err := read.Term(opts.Goal)
	if err != nil {
		panic(fmt.Sprintf("Cannot parse %s: %s", opts.Goal, err))
	}
	if !m.CanProve(goal) {
		panic(fmt.Sprintf("Cannot prove %s", opts.Goal))
	}
	variables := term.Variables(goal)
	answers := m.ProveAll(goal)
	if len(answers) == 0 {
		logging.Info.Printf("no")
		return
	}
	if variables.Size() == 0 {
		logging.Info.Printf("yes")
		return
	}
	lines := make([]string, 0)
	for i, answer := range answers {
		variables.ForEach(func(name string, variable interface{}) {
			v := variable.(*term.Variable)
			val := answer.Resolve_(v)
			line := fmt.Sprintf("%s = %s", name, val)
			lines = append(lines, line)
		})
		if i == len(answers)-1 {
			lines = append(lines, ".")
		} else {
			lines = append(lines, ";")
		}
	}
	logging.Info.Printf(strings.Join(lines, "\n"))
}

func parseFlags() *LogomakeOpions {
	fs := flag.NewFlagSet("logomake", flag.ContinueOnError)
	makefilePtr := fs.String(
		"makefile",
		"Makefile.logomake",
		"Read recipes from this file",
	)
	goalPtr := fs.String("goal", "all.", "Goal to build")
	homePtr := fs.String("home", "~/.logomake", "Logomake Prolog sources")
	logging.Dbug.Printf("Args: %s", os.Args)
	if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}
	if !fs.Parsed() {
		panic(fmt.Sprintf("Couldn't parse arguments: %s", os.Args))
	}
	logging.Info.Printf("Loading recipes from: %s", *makefilePtr)
	return &LogomakeOpions{
		Makefile: *makefilePtr,
		Goal:     *goalPtr,
		Home:     *homePtr,
	}
}

func initMachine(home string) golog.Machine {
	builtins := golog.NewMachine().RegisterForeign(
		map[string]golog.ForeignPredicate{
			"c/2":    native.C2,
			"glob/2": native.Glob2,
		},
	)
	prolog, err := filepath.Glob(home + "/*.pl")
	if err == nil {
		for _, f := range prolog {
			file, err := os.OpenFile(f, os.O_RDONLY, os.ModePerm)
			if err == nil {
				builtins = builtins.Consult(file)
			} else {
				logging.Warn.Printf("Couldn't read file: %s: %s", f, err)
			}
		}
	} else {
		logging.Warn.Printf(
			"Couldn't access Logomake home in: %s: %s",
			home, err,
		)
	}
	return builtins
}
