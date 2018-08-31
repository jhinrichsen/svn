package svn

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// First entry from Subversion's `tags`
func firstEntry() Entry {
	return Entry{
		Kind: "dir",
		Name: "0.10.0",
		Commit: Commit{
			Revision: "849186",
			Author:   "cmpilato",
			// year, month, day, hour, minute, seconds, nanoseconds
			Date: time.Date(2004, 3, 18, 17, 35, 35, 355023*1000, time.UTC),
		},
	}
}

func listSample() (*ListElement, error) {
	filename := "testdata/list_apache_subversion_tags.xml"
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var l ListElement
	if err := xml.Unmarshal(f, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

func svnRepo() *Repository {
	return NewRepository("https://svn.apache.org/repos/asf/subversion")
}

// TestSvnList requires network access to svn.apache.org
func TestOfflineList(t *testing.T) {
	l, err := listSample()
	if err != nil {
		t.Fatal(err)
	}

	// Check number of entries
	wantN := 235
	gotN := len(l.Entries)
	if wantN != gotN {
		t.Fatalf("want %d but got %d\n", wantN, gotN)
	}

	// Check first entry
	want := firstEntry()
	got := l.Entries[0]
	if want != got {
		t.Fatalf("want %q but got %q\n", want, got)
	}
}

func TestOnlineList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping online test in short mode.")
	}
	want := firstEntry()
	r := svnRepo()
	got, err := r.List("tags", os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
	if want != got[0] {
		t.Fatalf("want %q but got %q\n", want, got)
	}
}

func TestSince(t *testing.T) {
	want := 99
	l, err := listSample()
	if err != nil {
		t.Fatal(err)
	}
	got := len(Since(l.Entries, time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)))
	if want != got {
		t.Fatalf("want %d but got %d\n", want, got)
	}
}

func TestExportOnline(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping online test in short mode.")
	}
	r := svnRepo()
	d, err := ioutil.TempDir("", "svntest-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(d); err != nil {
			t.Fatalf("error cleaning up %s: %s", d, err)
		}
	}()
	into := filepath.Join(d, "notes.txt")
	from := "tags/1.9.9/CHANGES"
	if err := r.Export(from, into, os.Stdout, make(chan string)); err != nil {
		t.Fatalf("error exporting %s into %s: %s", from, into, err)
	}

	// Check first line
	const want = "Version 1.9.9"
	buf, err := ioutil.ReadFile(into)
	if err != nil {
		t.Fatal(err)
	}
	got := string(buf)
	if !strings.HasPrefix(got, want) {
		t.Fatalf("want %s but got %s\n", want, got)
	}
}

func TestExportChannelOnline(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping online test in short mode.")
	}
	r := svnRepo()
	d, err := ioutil.TempDir("", "svntest-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(d); err != nil {
			t.Fatalf("error cleaning up %s: %s", d, err)
		}
	}()

	from := "tags/1.9.9"
	into := filepath.Join(d, "1.9.9")
	c := make(chan string)
	go func(filenames chan string) {
		for filename := range filenames {
			log.Printf("exported %q\n", filename)
		}
	}(c)
	if err := r.Export(from, into, os.Stdout, c); err != nil {
		t.Fatalf("error exporting %s into %s: %s", from, into, err)
	}
}
