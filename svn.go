package svn

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
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

// List will execute an `svn list` subcommand.
// Any non-nil xmlWriter will receive the XML content
func (a *Repository) List(relpath string, xmlWriter io.Writer) ([]Entry, error) {
	log.Printf("listing %s\n", relpath)
	fp := a.FullPath(relpath)
	cmd := exec.Command("svn", "list", "--xml", fp)
	log.Printf("executing %+v\n", cmd)
	buf, err := cmd.CombinedOutput()
	if xmlWriter != nil {
		io.Copy(xmlWriter, bytes.NewReader(buf))
	}
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

// Since returns all entries created after t
func Since(entries []Entry, t time.Time) []Entry {
	var es []Entry
	for _, e := range entries {
		if e.Commit.Date.After(t) {
			es = append(es, e)
		}
	}
	return es
}

func (a *Repository) Export(relpath string, into string, w io.Writer) error {
	log.Printf("exporting %s\n", relpath)
	fp := a.FullPath(relpath)
	// TODO force?
	cmd := exec.Command("svn", "export", fp, into)
	log.Printf("executing %+v\n", cmd)
	buf, err := cmd.CombinedOutput()
	if w != nil {
		io.Copy(w, bytes.NewReader(buf))
	}
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s", buf)
		return fmt.Errorf("Cannot export %s into %s: %s", fp, into, err)
	}
	return nil
}
