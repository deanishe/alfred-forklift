//
// Copyright (c) 2017 Dean Jackson <deanishe@deanishe.net>
//
// MIT Licence. See http://opensource.org/licenses/MIT
//
// Created on 2017-07-20
//

// alfred-forklift is a Script Filter for Alfred 3+ to open ForkLift favourites.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
)

var (
	// workflow configuration
	repo        = "deanishe/alfred-forklift"
	helpURL     = "https://github.com/deanishe/alfred-forklift/issues"
	iconDefault = &aw.Icon{Value: "icon.png"}
	iconUpdate  = &aw.Icon{Value: "update-available.png"}
	wf          *aw.Workflow

	// CLI options and workflow settings
	cli         *flag.FlagSet
	query       string
	demoMode    bool
	doUpdate    bool
	ignoreLocal bool
)

func init() {
	wf = aw.New(
		update.GitHub(repo),
		aw.HelpURL(helpURL),
	)

	cli = flag.NewFlagSet("forklift", flag.ExitOnError)
	cli.Usage = func() {
		fmt.Fprintf(os.Stderr, `forklift [options] [<query>]

Alfred 3+ workflow for searching ForkLift 3 favourites.

Usage:
    forklift [-demo] [<query>]
    forklift -update
    forklift -h

Options:
`)
		cli.PrintDefaults()
	}

	cli.BoolVar(&demoMode, "demo", false, "use demo data instead of real favourites")
	cli.BoolVar(&doUpdate, "update", false, "check whether an update is available")
}

// run starts the workflow
func run() {
	if err := cli.Parse(wf.Args()); err != nil {
		wf.FatalError(err)
	}
	query = cli.Arg(0)
	ignoreLocal = wf.Config.GetBool("IGNORE_LOCAL", false)

	// Alternate actions ----------------------------------
	if doUpdate {
		wf.Configure(aw.TextErrors(true))
		log.Printf("checking for update...")
		if err := wf.CheckForUpdate(); err != nil {
			wf.FatalError(err)
		}
		return
	}

	// Script Filter -------------------------------------
	log.Printf("query=%q", query)

	// Notify updates
	if wf.UpdateCheckDue() && !wf.IsRunning("update") {
		log.Printf("checking for update ...")
		cmd := exec.Command(os.Args[0], "-update")
		if err := wf.RunInBackground("update", cmd); err != nil {
			log.Printf("[ERROR] update check failed: %v", err)
		}
	}

	if query == "" && wf.UpdateAvailable() {
		log.Printf("update available")
		wf.NewItem("An update is available").
			Subtitle("↩ or ⇥ to install update").
			Valid(false).
			Autocomplete("workflow:update").
			Icon(iconUpdate)
		wf.Configure(aw.SuppressUIDs(true))
	}

	var (
		faves []Favourite
		err   error
	)

	if demoMode {
		faves = demoFavourites()
	} else {
		faves, err = loadFavourites(favesFile)
		if err != nil {
			wf.Fatalf("couldn't load favourites: %v", err)
		}
	}
	log.Printf("%d favourite(s)", len(faves))

	for _, f := range faves {
		log.Printf("%s (%s)", f.Name, f.Type)
		wf.NewItem(f.Name).
			Subtitle(f.Server).
			Arg(f.UUID).
			UID(f.UUID).
			Match(fmt.Sprintf("%s %s", f.Name, f.Server)).
			Icon(f.Icon()).
			Valid(true)
	}

	if query != "" {
		res := wf.Filter(query)
		log.Printf("%d favourite(s) match %q", len(res), query)
	}

	wf.WarnEmpty("No favourites found", "Try a different query?")

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
