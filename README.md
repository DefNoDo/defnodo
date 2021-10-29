# DefNoDo
DefNoDo is a docker service runtime environment that doesn't require any root or elevated permissions, does not require user interaction to install, and has no required upgrade path.

The goal is to provide a docker service that can be installed, configured and run in a headless, server environment as well as used for desktop MacOS machines.  Prototype support for running podman instead of docker has been successful, and may be included in the future.

Some high-level goals:
* Install and runnable without needing any UI interactions/fully automatable
* Fully usable without root or elevated permissions
* Docker version independent of DefNoDo version
* Ability to run podman instead of or alongside docker
* Interactions/usage that match running docker on Linux systems (`/etc/docker/daemon.json`, `/var/run/docker.sock`, etc)
* No account necessary to use

## Usage
```
NAME:
   DefNoDo - A new cli application

USAGE:
   defnodo [global options] command [command options] [arguments...]

VERSION:
   v0.1

COMMANDS:
   run, r      Run defnodo and the underlying docker service
   service, s  Control defnodo service
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE  Load defnodo configuration from FILE (default: "~/.defnodorc") [$DEFNODORC]
   --help, -h              show help (default: false)
   --version, -v           print the version (default: false)
```

## Current Limitations
Currently, nearly everything is hardcoded, but all tested docker functionality works.

* `metadata.json` must be updated with the user's home directory, no other values will result in successful mounts
* `/var/run/docker.sock` is not symlinked to the socket yet
* dockerd is version `18.06` due to underlying containerd conflicts.
* Running is done via the `linuxkit` run command rather than building the actual `hyperkit` and `vpnkit` commands
* Zero tests
* Locations are all relative to the `defnodo` binary and must be maintained