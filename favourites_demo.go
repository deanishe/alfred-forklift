//
// Copyright (c) 2017 Dean Jackson <deanishe@deanishe.net>
//
// MIT Licence. See http://opensource.org/licenses/MIT
//
// Created on 2017-07-28
//

package main

// demoFavourites returns sample data
func demoFavourites() []Favourite {
	return []Favourite{
		{
			UUID:   "BA892C03-0300-4A56-AB60-A116FBC4B84C",
			Name:   "Warez",
			Group:  "Data",
			Server: "www.warez.ru",
			Type:   "FTP",
		},
		{
			UUID:   "87ECFC95-FE2F-4327-839B-CE9C2E43C96C",
			Name:   "Ubuntu ISOs",
			Group:  "Data",
			Server: "ftp.ubuntu.com",
			Type:   "FTP",
		},
		{
			UUID:   "16E522BD-1213-4894-8FB1-F86C5C017A0A",
			Name:   "OpenStreetMap",
			Group:  "Data",
			Server: "ftp5.gwdg.de",
			Type:   "FTP",
		},
		{
			UUID:   "3DD969C4-773B-4600-8F8F-13416CC68A57",
			Name:   "Homepage",
			Group:  "Data",
			Server: "webdav.example.com",
			Type:   "WebDAV",
		},
		{
			UUID:   "711B838F-E565-470E-8086-53ADFFDAA32E",
			Name:   "NAS",
			Group:  "Data",
			Server: "192.168.0.5",
			Type:   "Workspace",
		},
		{
			UUID:   "8935EF9D-BBB6-4B0B-87E4-80CCA3E84DE6",
			Name:   "S3 Bucket",
			Group:  "Data",
			Server: "mybucket.amazon.com",
			Type:   "S3",
		},
		{
			UUID:   "FA1B24E2-32ED-42BE-B28E-00534223EBFD",
			Name:   "www.example.com",
			Group:  "Data",
			Server: "server.example.com",
			Type:   "SFTP",
		},
		{
			UUID:   "751DA871-C5D5-44C6-BA29-B8CF6CBFA0FF",
			Name:   "demo.example.com",
			Group:  "Data",
			Server: "server.example.com",
			Type:   "SFTP",
		},
		{
			UUID:   "56AD2BFC-A7CD-4085-BED0-E50E2E5FB6D1",
			Name:   "reynolds.com",
			Group:  "Data",
			Server: "reynolds.com",
			Type:   "Workspace",
		},
		{
			UUID:   "1857118C-08C6-4AB1-B2BA-3A6B02318F59",
			Name:   "ullmann.org",
			Group:  "Data",
			Server: "ullmann.org",
			Type:   "SFTP",
		},
		{
			UUID:   "B21FBB02-048F-44B5-93D8-0E08E8A8D4D6",
			Name:   "Server Logs",
			Group:  "Data",
			Server: "fiebig.net",
			Type:   "Sync",
		},
		{
			UUID:   "45125B47-75BC-4CE3-A864-941B4CE24422",
			Name:   "iPhone (8080)",
			Group:  "Data",
			Server: "192.168.0.2",
			Type:   "WebDAV",
		},
	}
}
