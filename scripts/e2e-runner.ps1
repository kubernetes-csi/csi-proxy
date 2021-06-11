# hyperv
Install-WindowsFeature -Name Hyper-V -IncludeManagementTools -Restart

# installers
iwr https://chocolatey.org/install.ps1 -UseBasicParsing | iex
choco upgrade chocolatey
choco install -y golang
choco install -y vim
choco install -y git
$env:GOPATH="$HOME\go"
$env:GOOS="windows"
$env:GOARCH="amd64"
refreshenv

# check env
go version
git version

# go setup
New-Item -ItemType Directory -Force -Path $HOME\go\bin
New-Item -ItemType Directory -Force -Path $HOME\go\src
New-Item -ItemType Directory -Force -Path $HOME\go\src\github.com\kubernetes-csi\
cd $HOME/go/src/github.com/kubernetes-csi

# XXX: update the url to your repo & branch
Write-Output "building csi-proxy"
git clone https://github.com/mauriciopoppe/csi-proxy
# csi-proxy-setup (taken from the github action workflow)
cd $HOME/go/src/github.com/kubernetes-csi/csi-proxy
git checkout volume-api-changes

go build -v -a -o ./bin/csi-proxy.exe ./cmd/csi-proxy
go build -v -a -o ./build/csi-proxy-api-gen ./cmd/csi-proxy-api-gen

# start the CSI Proxy before running tests on windows
Write-Output "starting csi-proxy"
$csiproxy = Start-Job -Name CSIProxy -ScriptBlock {
    cd $HOME/go/src/github.com/kubernetes-csi/csi-proxy
    .\bin\csi-proxy.exe --v=5
};
Start-Sleep -Seconds 10;

Write-Output "checking that csi-proxy is still up"
Get-Job

Write-Output "getting named pipes"
[System.IO.Directory]::GetFiles("\\.\\pipe\\")

Write-Output "running e2e tests"
cd $HOME/go/src/github.com/kubernetes-csi/csi-proxy
go test -v -race ./integrationtests/
