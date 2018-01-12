//
// Copyright (c) 2017 Dean Jackson <deanishe@deanishe.net>
//
// MIT Licence. See http://opensource.org/licenses/MIT
//
// Created on 2017-07-20
//

/*
forklift.go
===========

A Script Filter for Alfred 3 to open ForkLift favourites.
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	"github.com/deanishe/awgo/util"
	"github.com/docopt/docopt-go"
)

var (
	usage = `forklift [options] [<query>]

Filter ForkLift favourites in Alfred 3.

Usage:
    forklift [<query>]
    forklift --help | --version
    forklift --distname
	forklift --logfile
	forklift --update

Options:
    --distname    Print filename of distributable .alfredworkflow file
	              (for the build script).
    -h, --help    Show this message and exit.
    --logfile     Print path to workflow's log file and exit.
	-u, --update  Check if an update is available.

`
	connectionTypes = []string{"Local", "SFTP", "NFS", "Rackspace", "S3", "Search", "FTP", "Sync", "VNC", "WebDAV", "Workspace"}
	connectionIcons map[string]*aw.Icon
	iconDefault     = &aw.Icon{Value: "icon.png", Type: aw.IconTypeImageFile}
	iconUpdate      = &aw.Icon{Value: "update-available.png", Type: aw.IconTypeImageFile}
	favesFile       = os.ExpandEnv("$HOME/Library/Application Support/ForkLift/Favorites/Favorites.json")
	// workflow configuration
	repo = "deanishe/alfred-forklift"
	wf   *aw.Workflow
	// CLI options
	query     string
	doLogfile bool
	doDist    bool
	doUpdate  bool
)

func init() {
	connectionIcons = map[string]*aw.Icon{}
	for _, s := range connectionTypes {
		path := fmt.Sprintf("/Applications/ForkLift.app/Contents/Resources/Connection%s.icns", s)
		connectionIcons[s] = &aw.Icon{Value: path, Type: aw.IconTypeImageFile}
	}
	wf = aw.New(update.GitHub(repo))
}

type faves struct {
	Groups []*faveGroup `json:"favorites"`
}

type fave struct {
	UUID string    `json:"UUID"`
	Attr *faveAttr `json:"attributes"`
	Type string    `json:"type"`
}

type faveGroup struct {
	UUID     string    `json:"UUID"`
	Attr     *faveAttr `json:"attributes"`
	Children []fave    `json:"childItems"`
	Type     string    `json:"type"`
}

type faveAttr struct {
	Name   string
	Path   string
	Server string
}

// Favourite is a ForkLift favourite
type Favourite struct {
	UUID   string // Favourite UUID
	Name   string // Name of favourite
	Group  string // Name of favourite's group
	Path   string // Favourite's path
	Server string // Hostname of favourite
	Type   string // Type of favourite (SFTP, Local etc.)
}

func (f *Favourite) Icon() *aw.Icon {
	if f.Type == "Local" {
		return &aw.Icon{Value: f.Path, Type: aw.IconTypeFileIcon}
	}
	for t, i := range connectionIcons {
		if strings.Index(f.Type, t) > -1 {
			return i
		}
	}
	return iconDefault
}

// newIcon creates a new workflow icon from an icon file in ForkLift's
// Resources directory.
func newIcon(filename string) *aw.Icon {
	path := fmt.Sprintf("/Applications/ForkLift.app/Contents/Resources/Connection%s.icns", filename)
	return &aw.Icon{Value: path, Type: aw.IconTypeImageFile}
}

// loadFavourites reads favourites from JSON file at path.
func loadFavourites(path string) ([]*Favourite, error) {
	var favourites = []*Favourite{}

	if !util.PathExists(path) {
		return nil, fmt.Errorf("file does not exist: %s", path)
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	rawFaves := &faves{}
	if err := json.Unmarshal(data, rawFaves); err != nil {
		return nil, fmt.Errorf("error unmarshalling %s: %s", path, err)
	}
	log.Printf("%d groups", len(rawFaves.Groups))
	for _, fg := range rawFaves.Groups {
		log.Printf("group '%s' contains %d favourites", fg.Attr.Name, len(fg.Children))
		for _, f := range fg.Children {

			fav := &Favourite{
				UUID:   f.UUID,
				Name:   f.Attr.Name,
				Group:  fg.Attr.Name,
				Path:   f.Attr.Path,
				Server: f.Attr.Server,
				Type:   f.Type,
			}

			// Ignore local favourites whose path doesn't exist
			if fav.Type == "Local" && !util.PathExists(fav.Path) {
				continue
			}

			if fav.Name == "" && fav.Path != "" {
				fav.Name = filepath.Base(fav.Path)
			}

			favourites = append(favourites, fav)
		}
	}
	return favourites, nil
}

func parseArgs() error {
	vstr := fmt.Sprintf("%s/%v (awgo/%v)", wf.Name(), wf.Version(), aw.AwGoVersion)
	args, err := docopt.Parse(usage, wf.Args(), true, vstr, false)
	if err != nil {
		return err
	}

	log.Printf("args=%+v", args)

	if args["--logfile"] == true {
		doLogfile = true
	}

	if args["--distname"] == true {
		doDist = true
	}

	if args["--update"] == true {
		doUpdate = true
	}

	if args["<query>"] != nil {
		query = args["<query>"].(string)
	}
	return nil
}

// run starts the workflow
func run() {
	// ================ Alternate actions ====================

	if err := parseArgs(); err != nil {
		panic(fmt.Sprintf("error parsing arguments: %s", err))
	}

	if doLogfile == true {
		fmt.Println(wf.LogFile())
		return
	}

	if doDist == true {
		name := strings.Replace(
			fmt.Sprintf("%s-%s.alfredworkflow", wf.Name(), wf.Version()),
			" ", "-", -1)
		fmt.Println(name)
		return
	}

	if doUpdate == true {
		wf.TextErrors = true
		log.Printf("checking for update...")
		if err := wf.CheckForUpdate(); err != nil {
			wf.FatalError(err)
		}
		return
	}

	// ================ Script Filter ====================

	var noUID bool
	log.Printf("query=%s", query)

	// Notify updates
	if wf.UpdateCheckDue() == true {
		log.Printf("update check due")
		wf.Var("check_update", "1")
	}

	if wf.UpdateAvailable() == true {
		log.Printf("update available")
		wf.NewItem("An update is available").
			Subtitle("↩ or ⇥ to install update").
			Valid(false).
			Autocomplete("workflow:update").
			Icon(iconUpdate)
		noUID = true
	}

	// Load favourites
	faves, err := loadFavourites(favesFile)
	if err != nil {
		panic(fmt.Sprintf("couldn't load favourites: %s", err))
	}
	log.Printf("%d favourite(s)", len(faves))

	for _, f := range faves {
		var uid string
		if noUID == false {
			uid = f.UUID
		}
		log.Printf("%s (%s)", f.Name, f.Type)
		it := wf.NewItem(f.Name).
			Subtitle(f.Server).
			Arg(f.UUID).
			UID(uid).
			Match(fmt.Sprintf("%s %s", f.Name, f.Server)).
			Icon(f.Icon()).
			Valid(true)

		it.Var("UUID", f.UUID)
	}

	if query != "" {
		res := wf.Filter(query)
		log.Printf("%d favourite(s) match '%s'", len(res), query)
	}

	wf.WarnEmpty("No favourite found", "Try a different query?")

	wf.SendFeedback()
}

// main calls run via aw.Run()
func main() {
	wf.Run(run)
}
