#! /bin/bash
# Run on golang docker image.

go install github.com/wailsapp/wails/v2/cmd/wails@latest
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
nvm install node
apt update
apt install -y libgtk-3-dev libwebkit2gtk-4.0-dev

npm install

wails build --noPackage --platform linux/amd64 -o linux/amd64/cuttlefish
