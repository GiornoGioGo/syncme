# SyncMe

Welcome to SyncMe! A free and open source Peer-to-Peer file synchronization tool written in Go.

Whether you want to syncronise game saves between systems or need to backup your files to another device, SyncMe has you covered.

## For Developers
For those who wish to contribute or build SyncMe from source, there are a few things you will need.

You will need:
- Go 1.26.6
- Git

## Cloning
```text
HTTPS:
git clone https://github.com/GiornoGioGo/syncme.git
SSH:
git clone git@github.com:GiornoGioGo/syncme.git
cd syncme
```

## Building
```text
Normal compilation:
go build .

Windows:
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o syncme.exe .

macOS (Apple Silicon - M1/M2/M3/M4):
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o syncme .

macOS (Intel):
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o syncme-intel .

Linux:
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o syncme.exe .
```

## Running
Now that you have built an executable for two of your devices, you can now run the program. Be sure to give both binaries execution permissions on your terminal.
```text
chmod += syncme
```
Server device (Device that will recieve files)
```text
./syncme.exe -mode=server -path="path/to/folder/"
```

Client device (Device sending out requested files)
```text
./syncme.exe -mode=client -target="Server device ip:9090" -path="path/to/folder/"
```

If all is successful, you will see all files transferred from the client folder to the server folder!