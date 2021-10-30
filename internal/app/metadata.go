package app

import (
	"encoding/json"
	"io/ioutil"
)

// {
//   "docker": {
//     "entries": {
//       "daemon.json": {
//         "content": "{\n  \"debug\" : true,\n  \"experimental\" : true\n}\n"
//       },
//       "user.config": {
//         "content": "/Users/mbrace"
//       }
//     }
//   }
// }

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
	// Read and build the different parts of the metadata
	data := &metadataJson{
		Docker: dockerMetadata{
			Entries: dockerEntriesMetadata{
				Daemon: dockerDaemonJson{
					Content: "{\n  \"debug\" : true,\n  \"experimental\" : true\n}\n",
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
	js, err := json.Marshal(data)
	_, err = file.Write(js)
	return
}
