package serve9p

// Heavily cribbed from https://github.com/knusbaum/go9p/blob/master/cmd/export9p/main.go
// License at https://github.com/knusbaum/go9p/blob/master/LICENSE

import (
	"errors"
	"log"
	"net/url"
	"path/filepath"

	"github.com/knusbaum/go9p"
	"github.com/knusbaum/go9p/fs"
	"github.com/knusbaum/go9p/fs/real"
)

type Serve9P struct {
	address     *url.URL
	directories []string
	verbose     bool
}

func NewServe9P(listenAddress string, directories []string, verbose bool) (serve9p *Serve9P, err error) {
	address, err := url.Parse(listenAddress)
	if err != nil {
		return
	}

	if address.Scheme != "tcp" {
		err = errors.New("go9p fileshares only support TCP listeners currently")
		return
	}

	if len(directories) == 0 {
		err = errors.New("No directories specified.")
		return
	}

	for i, directory := range directories {
		dir, err := filepath.Abs(directory)
		if err != nil {
			return nil, err
		}
		directories[i] = dir
	}
	serve9p = &Serve9P{
		address:     address,
		directories: directories,
		verbose:     verbose,
	}

	return
}

func (server *Serve9P) Run() {
	go9p.Verbose = server.verbose
	sharedFS := &fs.FS{}
	sharedFS.Root = &real.Dir{
		Path: server.directories[0],
	}

	fs.WithCreateFile(real.CreateFile)(sharedFS)
	fs.WithCreateDir(real.CreateDir)(sharedFS)
	fs.WithRemoveFile(real.Remove)(sharedFS)

	fs.IgnorePermissions()(sharedFS)

	if server.verbose {
		log.Printf("Sharing %s on %s", server.directories[0], server.address.Host)
	}
	err := go9p.Serve(server.address.Host, sharedFS.Server())
	if err != nil {
		log.Fatal(err)
	}
	return
}
