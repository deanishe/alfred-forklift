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

	aw "git.deanishe.net/deanishe/awgo"
	"github.com/docopt/docopt-go"
)

var (
	usage = `forklift [options] [<query>]

Filter ForkLift favourites in Alfred 3.

Usage:
    forklift [<query>]
    forklift --help | --version
    forklift --datadir | --cachedir | --distname | --logfile

Options:
    --datadir     Print path to workflow's data directory and exit.
    --cachedir    Print path to workflow's cache directory and exit.
    --logfile     Print path to workflow's log file and exit.
    --distname    Print filename of distributable .alfredworkflow file
	              (for the build script).
    -h, --help    Show this message and exit.

`
	favesFile = os.ExpandEnv("$HOME/Library/Application Support/ForkLift/Favorites/Favorites.json")
)

var (
	connectionTypes = []string{"Local", "SFTP", "NFS", "Rackspace", "S3", "Search", "FTP", "Sync", "VNC", "WebDAV", "Workspace"}
	iconDefault     = &aw.Icon{Value: "icon.png", Type: aw.IconTypeImageFile}
	connectionIcons map[string]*aw.Icon
)

func init() {
	connectionIcons = map[string]*aw.Icon{}
	for _, s := range connectionTypes {
		path := fmt.Sprintf("/Applications/ForkLift.app/Contents/Resources/Connection%s.icns", s)
		connectionIcons[s] = &aw.Icon{Value: path, Type: aw.IconTypeImageFile}
	}
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

	if !aw.PathExists(path) {
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
			if fav.Type == "Local" && !aw.PathExists(fav.Path) {
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

// run starts the workflow
func run() {
	var query string

	vstr := fmt.Sprintf("%s/%v (awgo/%v)", aw.Name(), aw.Version(), aw.AwGoVersion)
	args, err := docopt.Parse(usage, nil, true, vstr, false)
	if err != nil {
		log.Fatalf("error parsing CLI options: %s", err)
	}

	log.Printf("args=%+v", args)

	// ================ Alternate actions ====================

	if args["--datadir"] == true {
		fmt.Println(aw.DataDir())
		return
	}

	if args["--cachedir"] == true {
		fmt.Println(aw.CacheDir())
		return
	}

	if args["--logfile"] == true {
		fmt.Println(aw.LogFile())
		return
	}

	if args["--distname"] == true {
		name := strings.Replace(
			fmt.Sprintf("%s-%s.alfredworkflow", aw.Name(), aw.Version()),
			" ", "-", -1)
		fmt.Println(name)
		return
	}

	// ================ Script Filter ====================

	if args["<query>"] != nil {
		query = fmt.Sprintf("%v", args["<query>"])
	}
	log.Printf("query=%s", query)

	// Load favourites
	faves, err := loadFavourites(favesFile)
	if err != nil {
		panic(fmt.Sprintf("couldn't load favourites: %s", err))
	}
	log.Printf("%d favourites", len(faves))

	for _, f := range faves {
		log.Printf("%s (%s)", f.Name, f.Type)
		it := aw.NewItem(f.Name).
			Subtitle(f.Server).
			Arg(f.UUID).
			SortKey(fmt.Sprintf("%s %s", f.Name, f.Server)).
			Icon(f.Icon()).
			Valid(true)

		it.Var("UUID", f.UUID)
	}

	if query != "" {
		res := aw.Filter(query)
		log.Printf("%d favourites match '%s'", len(res), query)
	}

	aw.WarnEmpty("No matching favourites", "Try a different query?")

	aw.SendFeedback()
}

// main calls run via aw.Run()
func main() {
	aw.Run(run)
}
