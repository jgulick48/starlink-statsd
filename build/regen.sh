#!/bin/sh
cd "$(dirname "$0")/.." && go run ./build/gen "$@"
