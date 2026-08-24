BUCKET_NAME=toukibo-parser-samples
URL=https://pub-a26a7972d1ea437b983bf6696a7d847e.r2.dev
SNAPSHOT_MANIFEST=testdata-snapshot.json
DATA_DIR=testdata
export NUM_SAMPLE=1778

build:
	mkdir -p bin
	go build -o bin/toukibo-parser scripts/main.go

run: build
	./bin/toukibo-parser -mode=run -path=$(TARGET).pdf

run/sample: build
ifndef TARGET
	$(error TARGET is not set. Usage: make run/sample TARGET=sample1)
endif
	./bin/toukibo-parser -mode=run -path="$(DATA_DIR)/pdf/$(TARGET).pdf"

find/sample: build
ifndef TARGET
	$(error TARGET is not set. Usage: make find/sample TARGET=sample1 FIND="text")
endif
ifndef FIND
	$(error FIND is not set. Usage: make find/sample TARGET=sample1 FIND="text")
endif
	./bin/toukibo-parser -mode=find -path="$(DATA_DIR)/pdf/$(TARGET).pdf" -target="$(FIND)"

find/all: build
	FIND=$(FIND) ./scripts/find-samples.sh

rename:
	IDX=$(IDX) ./scripts/rename.sh

edit:
ifndef TARGET
	$(error TARGET is not set. Usage: make edit TARGET=sample1)
endif
	cat $(DATA_DIR)/yaml/$(TARGET).yaml

check:
	make open/sample TARGET=$(TARGET)
	make edit TARGET=$(TARGET)

annotate: build
ifndef TARGET
	$(error TARGET is not set. Usage: make annotate TARGET=sample1)
endif
	./bin/toukibo-parser -mode=run -path="$(DATA_DIR)/pdf/$(TARGET).pdf" > $(DATA_DIR)/yaml/$(TARGET).yaml
	make check TARGET=$(TARGET)

annotate/all: build
	./scripts/annotate-samples.sh

test: build
	go test -p 4 -coverprofile=coverage.out -shuffle=on ./...

bench: build
	go test -benchmem -run=^$$ -bench ^BenchmarkMain$$ -cpuprofile cpu.out -memprofile mem.out github.com/tychy/toukibo-parser
#	go tool pprof -http=":8887" cpu.out
#	go tool pprof -http=":8888" mem.out

zip/sample:
	rm -f testdata.zip
	zip -X -r testdata.zip testdata -x '*/.DS_Store' '*/bak_*'

put/sample:
	BUCKET_NAME=$(BUCKET_NAME) SNAPSHOT_MANIFEST=$(SNAPSHOT_MANIFEST) ./scripts/put-sample-snapshot.sh

get/sample: clean/data
	SAMPLE_URL=$(URL) SNAPSHOT_MANIFEST=$(SNAPSHOT_MANIFEST) ./scripts/get-sample-snapshot.sh

open/sample:
ifndef TARGET
	$(error TARGET is not set. Usage: make open/sample TARGET=sample1)
endif
ifeq ($(shell uname -s),Linux)
	xdg-open $(DATA_DIR)/pdf/$(TARGET).pdf
else
	open $(DATA_DIR)/pdf/$(TARGET).pdf
endif

clean: clean/bin clean/data

clean/bin:
	rm -rf bin

clean/data:
	rm -rf $(DATA_DIR)
	rm -rf testdata.zip
