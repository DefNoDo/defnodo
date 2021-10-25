#!/usr/bin/env bash
set -e

echo "Building defnodo..."

mkdir -p output

go build -o output/defnodo cmd/defnodo/main.go

echo "Built output/defnodo"