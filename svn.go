package svn

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

// List represents XML output of an `svn list` subcommand
type ListElement struct {
	XMLName xml.Name `xml:"lists"`
	Entries []Entry  `xml:"list>entry"`
}

// ListEntry represents XML output of an `svn list` subcommand
type Entry struct {
	Kind   string `xml:"kind,attr"`
	Name   string `xml:"name"`
	Commit Commit `xml:"commit"`
}

// Commit represents XML output of an `svn list` subcommand
type Commit struct {
	Revision string    `xml:"revision,attr"`
	Author   string    `xml:"author"`
	Date     time.Time `xml:"date"`
}

// Repository holds information about a (possibly remote) repository
type Repository struct {
	Location string
}

// NewRepository will initialize the internal structure of a possible remote
// repository, usually pointing to the parent of the default trunk/ tags/ branches
// structure.
func NewRepository(l string) *Repository {
	return &Repository{
		Location: l,
	}
}

// FullPath returns the full path into a repository
func (a *Repository) FullPath(relpath string) string {
	return fmt.Sprintf("%s/%s", a.Location, relpath)
}

// List will execute an `svn list` subcommand
func (a *Repository) List(relpath string) ([]Entry, error) {
	log.Printf("listing %s\n", relpath)
	fp := a.FullPath(relpath)
	cmd := exec.Command("svn", "list", "--xml", fp)
	buf, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s", buf)
		return nil, fmt.Errorf("Cannot list %s: %s", fp, err)
	}
	var l ListElement
	if err := xml.Unmarshal(buf, &l); err != nil {
		return nil, fmt.Errorf("cannot parse XML: %s: %s", buf, err)
	}
	return l.Entries, nil
}
