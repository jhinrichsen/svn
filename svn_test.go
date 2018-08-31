package svn

import (
	"encoding/xml"
	"io/ioutil"
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

// TestSvnList requires network access to svn.apache.org
func TestOfflineList(t *testing.T) {
	filename := "testdata/list_apache_subversion_tags.xml"
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	var l ListElement
	if err := xml.Unmarshal(f, &l); err != nil {
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
	r := NewRepository("https://svn.apache.org/repos/asf/subversion")
	got, err := r.List("tags")
	if err != nil {
		t.Fatal(err)
	}
	if want != got[0] {
		t.Fatalf("want %q but got %q\n", want, got)
	}
}
