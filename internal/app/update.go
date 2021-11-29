package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/codeclysm/extract"
	"github.com/google/go-github/v40/github"
)

func UpdateVM(vmDirectory string, force bool) (err error) {
	log.Printf("Checking for updates against: %s....\n", vmDirectory)
	client := github.NewClient(nil)
	// Retrieve the latest non-draft, non-prerelease release
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "defnodo", "vm-image")
	if err != nil {
		return err
	}

	latest := release.GetTagName()

	log.Printf("Found Github version: %s\n", latest)
	// Assume an error is a non-existent file / not a match
	current, _ := os.ReadFile(filepath.Join(vmDirectory, "VERSION"))

	if err != nil {
		log.Fatal(err)
	}

	if string(current) != latest || force {
		log.Printf("Updating to VM version %x\n", latest)
		// Iterate over the assets to find the one for the VM Image
		for _, asset := range release.Assets {
			if asset.GetName() == "vm-image.tar.gz" {
				// Retrieves a reader to fetch the associated asset
				reader, _, err := client.Repositories.DownloadReleaseAsset(context.Background(), "defnodo", "vm-image", asset.GetID(), http.DefaultClient)
				if err != nil {
					return err
				}
				// Extract the content from the reader directly to disk.  This saves disk space/copying/extracting and
				// different archive format issues.
				// TODO: Clean this up to remove previous?  This assumes it's the same files being overwritten
				err = extract.Archive(context.Background(), reader, vmDirectory, nil)
				if err != nil {
					// Error here is bad, local contents might be in a bad state, so don't write VERSION file
					return err
				}

				// Write the new version number
				err = os.WriteFile(filepath.Join(vmDirectory, "VERSION"), []byte(latest), 0644)
				if err != nil {
					return err
				}
			}
		}
	} else {
		log.Printf("Latest VM version installed")
	}
	return
}
