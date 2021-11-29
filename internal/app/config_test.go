package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

// Entrypoint to run the suite
func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

// Perform suite initialization
func (suite *ConfigTestSuite) SetupTest() {

}

// Represents an os.File to be loaded
type MockConfigFile struct {
	mock.Mock
	name     string
	contents string
	index    int
	reader   *strings.Reader
}

// Return filename
func (cf *MockConfigFile) Name() string {
	return cf.name
}

// Perform read operations on the contents as if the file contents
func (cf *MockConfigFile) Read(p []byte) (n int, err error) {
	if cf.reader == nil {
		cf.reader = strings.NewReader(cf.contents)
	}
	n, err = cf.reader.Read(p)
	return n, err
}

func (suite *ConfigTestSuite) TestDefaultsNoFile() {
	homedir, err := os.UserHomeDir()
	assert.Nil(suite.T(), err)

	datadir := filepath.Join(homedir, ".defnodo")

	config, err := LoadConfig(nil)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), datadir, config.DataDirectory)
	assert.Len(suite.T(), config.VolumeMounts, 1)
	assert.Equal(suite.T(), homedir, config.VolumeMounts[0])

	assert.Equal(suite.T(), "", config.ContainerRuntime.DaemonJson)
	assert.Equal(suite.T(), "docker", config.ContainerRuntime.Runtime)
	assert.Equal(suite.T(), "latest", config.ContainerRuntime.Version)
	assert.Equal(suite.T(), filepath.Join(datadir, "docker.versions"), config.ContainerRuntime.VersionsFile)

	assert.Equal(suite.T(), 2048, config.VM.Memory)
	assert.Equal(suite.T(), 1, config.VM.Cpus)
	assert.Equal(suite.T(), "10G", config.VM.DiskSize)

	assert.False(suite.T(), config.IsService)
	assert.False(suite.T(), config.Interactive)
	assert.Equal(suite.T(), homedir, config.ConfigBaseDirectory)
}

func (suite *ConfigTestSuite) TestFileAbsolutePaths() {
	file := &MockConfigFile{
		name: "/my/homedir/.defnodorc",
		contents: `---
data-directory: /another/dir/defnodo/defnodo-data
volume-mounts:
  - /a/third/dir
container-runtime:
  docker-daemon.json: /path/to/daemon.json
  runtime: notdefault
  version: 82.5
  versions-file: /custom/path/to/custom.versions
vm:
  memory: 8192
  cpus: 4
  disk-size: 5G
`,
	}

	config, err := LoadConfig(file)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), "/another/dir/defnodo/defnodo-data", config.DataDirectory)
	assert.Len(suite.T(), config.VolumeMounts, 1)
	assert.Equal(suite.T(), "/a/third/dir", config.VolumeMounts[0])

	assert.Equal(suite.T(), "/path/to/daemon.json", config.ContainerRuntime.DaemonJson)
	assert.Equal(suite.T(), "notdefault", config.ContainerRuntime.Runtime)
	assert.Equal(suite.T(), "82.5", config.ContainerRuntime.Version)
	assert.Equal(suite.T(), "/custom/path/to/custom.versions", config.ContainerRuntime.VersionsFile)

	assert.Equal(suite.T(), 8192, config.VM.Memory)
	assert.Equal(suite.T(), 4, config.VM.Cpus)
	assert.Equal(suite.T(), "5G", config.VM.DiskSize)

	assert.False(suite.T(), config.IsService)
	assert.False(suite.T(), config.Interactive)
	assert.Equal(suite.T(), "/my/homedir", config.ConfigBaseDirectory)
}

func (suite *ConfigTestSuite) TestFileRelativePaths() {
	file := &MockConfigFile{
		name: "/my/homedir/.defnodorc",
		contents: `---
data-directory: defnodo-data
volume-mounts:
  - dir
container-runtime:
  docker-daemon.json: daemon.json
  runtime: notdefault
  version: 82.5
  versions-file: custom.versions
vm:
  memory: 8192
  cpus: 4
  disk-size: 5G
`,
	}

	config, err := LoadConfig(file)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), "/my/homedir/defnodo-data", config.DataDirectory)
	assert.Len(suite.T(), config.VolumeMounts, 1)
	assert.Equal(suite.T(), "/my/homedir/dir", config.VolumeMounts[0])

	assert.Equal(suite.T(), "/my/homedir/defnodo-data/daemon.json", config.ContainerRuntime.DaemonJson)
	assert.Equal(suite.T(), "notdefault", config.ContainerRuntime.Runtime)
	assert.Equal(suite.T(), "82.5", config.ContainerRuntime.Version)
	assert.Equal(suite.T(), "/my/homedir/defnodo-data/custom.versions", config.ContainerRuntime.VersionsFile)

	assert.Equal(suite.T(), 8192, config.VM.Memory)
	assert.Equal(suite.T(), 4, config.VM.Cpus)
	assert.Equal(suite.T(), "5G", config.VM.DiskSize)

	assert.False(suite.T(), config.IsService)
	assert.False(suite.T(), config.Interactive)
	assert.Equal(suite.T(), "/my/homedir", config.ConfigBaseDirectory)
}

func (suite *ConfigTestSuite) TestFilePartialOverrides() {
	homedir, err := os.UserHomeDir()
	assert.Nil(suite.T(), err)

	file := &MockConfigFile{
		name: "/my/homedir/.defnodorc",
		contents: `---
data-directory: defnodo-data
container-runtime:
  docker-daemon.json: daemon.json
  versions-file: custom.versions
`,
	}

	config, err := LoadConfig(file)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), "/my/homedir/defnodo-data", config.DataDirectory)
	assert.Len(suite.T(), config.VolumeMounts, 1)
	assert.Equal(suite.T(), homedir, config.VolumeMounts[0])

	assert.Equal(suite.T(), "/my/homedir/defnodo-data/daemon.json", config.ContainerRuntime.DaemonJson)
	assert.Equal(suite.T(), "docker", config.ContainerRuntime.Runtime)
	assert.Equal(suite.T(), "latest", config.ContainerRuntime.Version)
	assert.Equal(suite.T(), "/my/homedir/defnodo-data/custom.versions", config.ContainerRuntime.VersionsFile)

	assert.Equal(suite.T(), 2048, config.VM.Memory)
	assert.Equal(suite.T(), 1, config.VM.Cpus)
	assert.Equal(suite.T(), "10G", config.VM.DiskSize)

	assert.False(suite.T(), config.IsService)
	assert.False(suite.T(), config.Interactive)
	assert.Equal(suite.T(), "/my/homedir", config.ConfigBaseDirectory)
}

func (suite *ConfigTestSuite) TestBadYaml() {
	file := &MockConfigFile{
		name: "/my/homedir/.defnodorc",
		contents: `---
data-directory: defnodo-data
# Below here is invalid tab/indentation
		foobar: baz
`,
	}

	_, err := LoadConfig(file)
	assert.NotNil(suite.T(), err)
}
