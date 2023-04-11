# Releasing Cuttlefish

The release process is not automated through GitHub Actions yet, mainly due to the complications related to Wails.

Right now the build is done through scripts, with the archiving and publishing process done through goreleaser. All executed locally.

The MacOS and Windows builds can be executed on MacOS. The Linux build must be executed on Linux (basically, Wails for MacOS and Linux doesn't support cross-compilation).

## MacOS and Windows build
Run on MacOS. Run the `scripts/build.sh` script.

## Linux build
In a golang-based docker container on the amd64 platform, run `scripts/build-linux.sh`. Then, move the binaries (from `build/bin`) to the host machine, in the same directories.

## Publish
Back on the MacOS machine, run `scripts/publish.sh`. This will create a set of zip-files in `build/release` that you can upload to GitHub.

Then, create a release on GitHub and attach the zip-files.

Finally, open https://github.com/cube2222/homebrew-cube2222/blob/main/cuttlefish.rb and exchange the old download paths with the new release's download paths. Update the sha256 hashes with those printed by the `publish.sh` script. 
