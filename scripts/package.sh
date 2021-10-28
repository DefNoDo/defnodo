#!/usr/bin/env bash
set -e

tar -cvf defnodo.tar --exclude='defnodo-state/*' -s '#^#defnodo/~#' defnodo-data
tar -uvf defnodo.tar -s '#^#defnodo/~#' metadata.json
tar -uvf defnodo.tar -s '#^output#defnodo/bin#' output/*
gzip --best defnodo.tar