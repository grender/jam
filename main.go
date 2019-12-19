// Copyright (c) 2016, 2017 Evgeny Badin

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/mjibson/gmusic"

	"github.com/grender/jam/auth"
	"github.com/grender/jam/lastfm"
	"github.com/grender/jam/storage"
	"github.com/grender/jam/ui"
	"github.com/grender/jam/version"
)

const (
	// BANNER is what is printed for help/info output.
	BANNER = "Jam - %s\n"
)

var (
	vers   bool
	debug  bool
	lastFM bool
)

func init() {
	// parse flags
	flag.BoolVar(&vers, "version", false, "print version and exit")
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.BoolVar(&lastFM, "lastfm", false, "Enable LastFM scrobbler")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, version.Version))
		flag.PrintDefaults()
	}

	flag.Parse()

	if vers {
		fmt.Printf("%s\n", version.Version)
		os.Exit(0)
	}
}

func main() {
	db, err := storage.Open()
	if err != nil {
		log.Fatalf("Can't open database: %s", err)
	}
	gmusic, lmclient, lastfm, err := auth.CheckCreds(db, &lastFM)
	if err != nil {
		log.Fatalf("Can't connect to Google Music: %s", err)
	}
	defer db.Close()

	if err = doUI(gmusic, lmclient, lastfm, db); err != nil {
		log.Fatalf("Can't start UI: %s", err)
	}

}

func doUI(gmusic *gmusic.GMusic, lmclient *lastfm.Client, lastfm string, db *bolt.DB) error {
	app, err := ui.New(gmusic, lmclient, lastfm, db)
	if err != nil {
		return err
	}
	app.Run()
	return nil
}
