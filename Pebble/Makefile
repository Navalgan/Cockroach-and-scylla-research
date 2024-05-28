# Makefile

install: depend-install go-mod

depend-install:
	sudo apt-get install libgflags-dev libsnappy-dev zlib1g-dev libbz2-dev libzstd-dev

rocksdb-install:
	git clone https://github.com/facebook/rocksdb.git && cd rocksdb
	DEBUG_LEVEL=0 make shared_lib
	sudo make install-shared
	sudo ldconfig

go-mod:
	go mod tidy

start:
	rm -rf data
	CGO_CFLAGS="-I/path/to/rocksdb/include" CGO_LDFLAGS="-L/path/to/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd" go build
	./main
