#! /bin/bash

mkdir build/release
zip --junk-paths build/release/cuttlefish_darwin_amd64.zip build/bin/darwin/amd64/cuttlefish
zip --junk-paths build/release/cuttlefish_darwin_arm64.zip build/bin/darwin/arm64/cuttlefish
zip --junk-paths build/release/cuttlefish_windows_amd64.zip build/bin/windows/amd64/cuttlefish
zip --junk-paths build/release/cuttlefish_windows_arm64.zip build/bin/windows/arm64/cuttlefish
zip --junk-paths build/release/cuttlefish_linux_amd64.zip build/bin/linux/amd64/cuttlefish

openssl dgst -sha256 build/release/cuttlefish_darwin_amd64.zip
openssl dgst -sha256 build/release/cuttlefish_darwin_arm64.zip
