sudo apt-get update
sudo apt-get install go git
mkdir ~/go
mkdir ~/go/src ~/go/pkg ~/go/bin
echo "export GOPATH=~/go" > ~/.bashrc
echo "export GOBIN=$GOPATH/bin" > ~/.bashrc


