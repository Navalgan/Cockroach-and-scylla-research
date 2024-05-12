# Makefile

install: go-mod

go-mod:
	go mod tidy

start:
	rm -rf data
	CGO_CFLAGS="-I/opt/homebrew/Cellar/rocksdb/9.1.1_1/include" CGO_LDFLAGS="-L/opt/homebrew/Cellar/rocksdb/9.1.1_1/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd -L/opt/homebrew/Cellar/snappy/1.2.1/lib -L/opt/homebrew/Cellar/lz4/1.9.4/lib/ -L /opt/homebrew/Cellar/zstd/1.5.6/lib/" go build
	./main