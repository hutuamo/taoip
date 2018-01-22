.PHONEY: clean

BIN = taoip

$(BIN): src/*.go
	GOPATH=`pwd` go build -o $@ $^

clean:
	rm -f $(BIN)

