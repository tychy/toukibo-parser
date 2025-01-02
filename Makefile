BUCKET_NAME=toukibo-parser-samples
URL=https://pub-a26a7972d1ea437b983bf6696a7d847e.r2.dev
DATA_DIR=testdata
export NUM_SAMPLE=1452

build:
	mkdir -p bin
	go build -o bin/toukibo-parser main.go toukibo_parser.go

run: build
	./bin/toukibo-parser -mode=run -path=$(TARGET).pdf

run/sample: build
	./bin/toukibo-parser -mode=run -path="$(DATA_DIR)/pdf/$(TARGET).pdf"

find/sample: build
	./bin/toukibo-parser -mode=find -path="$(DATA_DIR)/pdf/$(TARGET).pdf" -target="$(FIND)"

find/all: build
	FIND=$(FIND) ./scripts/find-samples.sh

rename:
	IDX=$(IDX) ./scripts/rename.sh

edit:
	cat $(DATA_DIR)/yaml/$(TARGET).yaml

check:
	make open/sample TARGET=$(TARGET)
	make edit TARGET=$(TARGET)

annotate: build
	./bin/toukibo-parser -mode=run -path="$(DATA_DIR)/pdf/$(TARGET).pdf" > $(DATA_DIR)/yaml/$(TARGET).yaml
	make check TARGET=$(TARGET)

annotate/all: build
	./scripts/annotate-samples.sh

test: build
	go test -coverprofile=coverage.out -shuffle=on ./...

bench: build
	go test -benchmem -run=^$$ -bench ^BenchmarkMain$$ -cpuprofile cpu.out -memprofile mem.out github.com/tychy/toukibo-parser
#	go tool pprof -http=":8887" cpu.out
#	go tool pprof -http=":8888" mem.out

zip/sample:
	zip -r testdata.zip testdata

put/sample: zip/sample
	wrangler r2 object delete $(BUCKET_NAME)/testdata.zip
	wrangler r2 object put $(BUCKET_NAME)/testdata.zip --file testdata.zip
	
get/sample: clean/data
	mkdir -p $(DATA_DIR)
	curl -o testdata.zip $(URL)/testdata.zip
	unzip testdata.zip

open/sample:
	open $(DATA_DIR)/pdf/$(TARGET).pdf

clean: clean/bin clean/data

clean/bin:
	rm -rf bin

clean/data:
	rm -rf $(DATA_DIR)
	rm -rf testdata.zip
