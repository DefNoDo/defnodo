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
	Daemon        dockerDaemonJson        `json:"daemon.json"`
	UserConfig    dockerUserConfigJson    `json:"user.config"`
	DockerVersion dockerVersionConfigJson `json:"docker.version"`
}

type dockerDaemonJson struct {
	Content string `json:"content"`
}

type dockerUserConfigJson struct {
	Content string `json:"content"`
}

type dockerVersionConfigJson struct {
	Content string `json:"content"`
}

func (defnodo *DefNoDo) generateMetadata(directory string) (filename string, err error) {
	// Read any specified daemon.json
	daemon_json := "{}"
	if defnodo.config.ContainerRuntime.DaemonJson != "" {
		daemon, err := os.ReadFile(defnodo.config.ContainerRuntime.DaemonJson)
		if err != nil {
			return filename, err
		}
		if !json.Valid(daemon) {
			msg := fmt.Sprintf("%s does not contain valid JSON: \n%s\n", defnodo.config.ContainerRuntime.DaemonJson, string(daemon))
			err = errors.New(msg)
			return filename, err
		}
		daemon_split := strings.Split(string(daemon), "\n")
		daemon_json = strings.Join(daemon_split, "\n")
	}

	runtime_version, err := getRuntimeVersionHash(defnodo.config.ContainerRuntime.VersionsFile, defnodo.config.ContainerRuntime.Version)
	if err != nil {
		return
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
				DockerVersion: dockerVersionConfigJson{
					Content: runtime_version,
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

func getRuntimeVersionHash(filename string, version string) (result string, err error) {
	result = ""
	if version == "latest" {
		result = "latest latest"
	}

	data, err := ioutil.ReadFile(filename)
	versions := strings.Split(string(data), "\n")
	for _, line := range versions {
		if strings.HasPrefix(line, version) {
			// Don't break out of the loop so it gets the last value for the highest matching version
			result = line
		}
	}

	if result == "" {
		msg := fmt.Sprintf("%s does not contain an entry for %s\n", filename, version)
		err = errors.New(msg)
	}
	return
}
