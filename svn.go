package main

import (
	"encoding/xml"
	"time"
)

type SvnList struct {
	XMLName xml.Name       `xml:"lists"`
	Entries []SvnListEntry `xml:"list>entry"`
}

type SvnListEntry struct {
	Kind   string    `xml:"kind,attr"`
	Name   string    `xml:"name"`
	Commit SvnCommit `xml:"commit"`
}

type SvnCommit struct {
	Revision string    `xml:"revision,attr"`
	Author   string    `xml:"author"`
	Date     time.Time `xml:"date"`
}

func svnListBzd(segment string) (SvnList, error) {
	var l SvnList
	return l, nil
}
