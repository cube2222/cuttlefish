#! /bin/bash

wails build --noPackage --platform darwin/arm64 -o darwin/arm64/cuttlefish
wails build --noPackage --platform darwin/amd64 -o darwin/amd64/cuttlefish
wails build --noPackage --platform windows/arm64 -o windows/arm64/cuttlefish.exe
wails build --noPackage --platform windows/amd64 -o windows/amd64/cuttlefish.exe
