#! /bin/bash

mkdir build/release
cd build/bin/darwin/amd64
zip ../../../release/cuttlefish_darwin_amd64.zip cuttlefish
cd ../../../..
cd build/bin/darwin/arm64
zip ../../../release/cuttlefish_darwin_arm64.zip cuttlefish
cd ../../../..
cd build/bin/windows/amd64
zip ../../../release/cuttlefish_windows_amd64.zip cuttlefish
cd ../../../..
cd build/bin/windows/arm64
zip ../../../release/cuttlefish_windows_arm64.zip cuttlefish
cd ../../../..
cd build/bin/linux/amd64
zip ../../../release/cuttlefish_linux_amd64.zip cuttlefish
cd ../../../..

openssl dgst -sha256 build/release/cuttlefish_darwin_amd64.zip
openssl dgst -sha256 build/release/cuttlefish_darwin_arm64.zip
