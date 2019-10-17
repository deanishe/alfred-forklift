// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"
)

// ForkLift's icon directory
const resources = "/Applications/ForkLift.app/Contents/Resources/"

var (
	// path to favourites file
	favesFile = os.ExpandEnv("$HOME/Library/Application Support/ForkLift/Favorites/Favorites.json")
)

// Favourite is a ForkLift favourite
type Favourite struct {
	UUID   string // Favourite UUID
	Name   string // Name of favourite
	Group  string // Name of favourite's group
	Path   string // Favourite's path
	Server string // Hostname of favourite
	Type   string // Type of favourite (SFTP, Local etc.)
}

// ByName sorts a slice of Favourites by name.
type ByName []Favourite

func (s ByName) Len() int           { return len(s) }
func (s ByName) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s ByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Icon returns the appropriate icon for the favourite's type
func (f Favourite) Icon() *aw.Icon {
	switch f.Type {
	case "Local":
		return &aw.Icon{Value: f.Path, Type: aw.IconTypeFileIcon}
	case "Backblaze":
		return &aw.Icon{Value: resources + "ConnectionBackblaze.icns"}
	case "FTP":
		return &aw.Icon{Value: resources + "ConnectionFTP.icns"}
	case "GoogleDrive":
		return &aw.Icon{Value: resources + "ConnectionGoogleDrive.icns"}
	case "NFS":
		return &aw.Icon{Value: resources + "ConnectionNFS.icns"}
	case "Rackspace":
		return &aw.Icon{Value: resources + "ConnectionRackspace.icns"}
	case "S3":
		return &aw.Icon{Value: resources + "ConnectionS3.icns"}
	case "SFTP":
		return &aw.Icon{Value: resources + "ConnectionSFTP.icns"}
	case "Search":
		return &aw.Icon{Value: resources + "ConnectionSearch.icns"}
	case "Sync":
		return &aw.Icon{Value: resources + "ConnectionSync.icns"}
	case "VNC":
		return &aw.Icon{Value: resources + "ConnectionVNC.icns"}
	case "WebDAV", "WebDAVHTTPS":
		return &aw.Icon{Value: resources + "ConnectionWebDAV.icns"}
	case "Workspace":
		return &aw.Icon{Value: resources + "ConnectionWorkspace.icns"}
	default:
		log.Printf("[WARN] unknown type: %s", f.Type)
		return iconDefault
	}
}

// loadFavourites reads favourites from JSON file at path.
func loadFavourites(path string) ([]Favourite, error) {
	var faves = []Favourite{}

	if !util.PathExists(path) {
		return nil, fmt.Errorf("favourites file does not exist: %s", path)
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	raw := struct {
		Groups []struct {
			UUID string `json:"UUID"`
			Attr struct {
				Name   string
				Path   string
				Server string
			} `json:"attributes"`
			Children []struct {
				UUID string `json:"UUID"`
				Attr struct {
					Name   string
					Path   string
					Server string
				} `json:"attributes"`
				Type string `json:"type"`
			} `json:"childItems"`
			Type string `json:"type"`
		} `json:"favorites"`
	}{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("error unmarshalling %s: %s", path, err)
	}
	log.Printf("%2d group(s)", len(raw.Groups))
	for _, fg := range raw.Groups {
		log.Printf("%2d favourite(s) in group %q", len(fg.Children), fg.Attr.Name)
		for _, f := range fg.Children {

			fav := Favourite{
				UUID:   f.UUID,
				Name:   f.Attr.Name,
				Group:  fg.Attr.Name,
				Path:   f.Attr.Path,
				Server: f.Attr.Server,
				Type:   f.Type,
			}

			if ignoreLocal && fav.Type == "Local" {
				continue
			}

			// Ignore local favourites whose path doesn't exist
			if fav.Type == "Local" && !util.PathExists(fav.Path) {
				continue
			}

			if fav.Name == "" && fav.Path != "" {
				fav.Name = filepath.Base(fav.Path)
			}

			faves = append(faves, fav)
		}
	}

	sort.Sort(ByName(faves))

	return faves, nil
}
