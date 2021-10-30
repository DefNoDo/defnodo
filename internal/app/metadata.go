package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// See https://github.com/linuxkit/linuxkit/blob/master/docs/metadata.md for info about
// formatting of the metadata.

type metadataJson struct {
	Docker dockerMetadata `json:"docker"`
}

type dockerMetadata struct {
	Entries dockerEntriesMetadata `json:"entries"`
}

type dockerEntriesMetadata struct {
	Daemon     dockerDaemonJson     `json:"daemon.json"`
	UserConfig dockerUserConfigJson `json:"user.config"`
}

type dockerDaemonJson struct {
	Content string `json:"content"`
}

type dockerUserConfigJson struct {
	Content string `json:"content"`
}

func (defnodo *DefNoDo) generateMetadata(directory string) (filename string, err error) {
	// Read any specified daemon.json
	daemon_json := "{}"
	if defnodo.config.DockerDaemonJson != "" {
		daemon, err := os.ReadFile(defnodo.config.DockerDaemonJson)
		if err != nil {
			return filename, err
		}
		if !json.Valid(daemon) {
			msg := fmt.Sprintf("%s does not contain valid JSON: \n%s\n", defnodo.config.DockerDaemonJson, string(daemon))
			err = errors.New(msg)
			return filename, err
		}
		daemon_split := strings.Split(string(daemon), "\n")
		daemon_json = strings.Join(daemon_split, "\n")
	}
	// Build the different parts of the metadata
	data := &metadataJson{
		Docker: dockerMetadata{
			Entries: dockerEntriesMetadata{
				Daemon: dockerDaemonJson{
					Content: daemon_json,
				},
				UserConfig: dockerUserConfigJson{
					Content: defnodo.config.VolumeMounts[0],
				},
			},
		},
	}
	// Create the temp file
	file, err := ioutil.TempFile(directory, "metadata.*.json")
	if err != nil {
		return
	}
	filename = file.Name()
	defer file.Close()

	// render the json to the temp file
	js, err := json.MarshalIndent(data, "", "  ")
	_, err = file.Write(js)
	return
}
