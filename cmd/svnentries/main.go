// Error codes:
// 1: general error
// 2: bad commandline invocation (follow usage to resolve)
//

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/jhinrichsen/fio"
	"github.com/jhinrichsen/svn"
)

func main() {
	log.Println("starting svnentries")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--since|--sincefile] uri [uri]*\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "    Default timestamp format is RFC 3339, e.g. 'date --rfc-3339=seconds'")
		flag.PrintDefaults()
	}
	repository := flag.String("repository", "https://svn.apache.org/repos/asf/subversion",
		"Subversion repository to check")
	since := flag.String("since", DefaultSince().Format(time.RFC3339), "Use since timestamp")
	sincefile := flag.String("sincefile", "", "Use timestamp from file to check for new entries, takes precedence over since")
	sinceformat := flag.String("sinceformat", time.RFC3339, "Default timestamp format (RFC 3339)")
	flag.Parse()

	r := svn.NewRepository(*repository)

	// Prefer sincefile over since
	var t time.Time
	var err error
	if *sincefile != "" {
		t, err = fio.ModTime(*sincefile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot determine timestamp of file %q: %s\n", *sincefile, err)
		}
	} else {
		t, err = time.Parse(*sinceformat, *since)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing timestamp %s in format %q: %s", *since, *sinceformat, err)
			os.Exit(2)
		}
	}
	log.Printf("looking for entries newer than %s\n", t)
	for _, url := range flag.Args() {
		es, err := r.List(url, ioutil.Discard)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error checking %q: %s\n", url, err)
		}
		for _, newEntry := range svn.Since(es, t) {
			fmt.Fprintf(os.Stdout, fmt.Sprintf("%s/%s/%s\n", r.Location, url, newEntry.Name))

		}
	}
}

// DefaultSince returns the timestamp 24 hours ago
func DefaultSince() time.Time {
	t := time.Now()
	// minus one day
	return t.AddDate(0, 0, -1)
}
