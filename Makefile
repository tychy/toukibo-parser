URL=https://toukibo-sample-files.s3.ap-northeast-1.amazonaws.com
OUTPUT_DIR=testdata/pdf
NUM_SAMPLE=778

build:
	mkdir -p bin
	go build -o bin/toukibo-parser main.go

run: build
	./bin/toukibo-parser -path=$(TARGET)

run-sample: build
	./bin/toukibo-parser -path="testdata/pdf/$(TARGET).pdf"
	
run-all-sample: build
	./run-all-sample.sh
# ちょっとおかしいもの
# 133 住所が途中で切れている
# 770 住所のパースがおかしい

get-sample: 
	mkdir -p $(OUTPUT_DIR)
	@for i in {1..$(NUM_SAMPLE)}; do \
		curl -s -o $(OUTPUT_DIR)/sample$$i.pdf $(URL)/sample$$i.pdf & \
	done


clean:
	rm -rf bin
	rm -rf $(OUTPUT_DIR)
