# Makefile

install: go-install

go-install:
	wget https://dl.google.com/go/go1.22.0.linux-amd64.tar.gz
	sudo tar -C /opt -xzf go1.22.0.linux-amd64.tar.gz
	bash export-go-path.sh
