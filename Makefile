URL=https://pub-a26a7972d1ea437b983bf6696a7d847e.r2.dev
OUTPUT_DIR=sample
NUM_SAMPLE=778

build:
	mkdir -p bin
	go build -o bin/toukibo-parser main.go

run: build
	./bin/toukibo-parser -path=$(TARGET)

run-all-sample: build
	./run-all-sample.sh
# ちょっとおかしいもの
# 133 住所が途中で切れている
# 770 住所のパースがおかしい

get-sample: 
	mkdir -p $(OUTPUT_DIR)
	@for i in {1..500}; do \
		curl -s -o "$(OUTPUT_DIR)/sample$$i.pdf" "$(URL)/remote-sample$$i.pdf" & \
	done
	wait
	sleep 5
	@for i in {500..$(NUM_SAMPLE)}; do \
		curl -s -o "$(OUTPUT_DIR)/sample$$i.pdf" "$(URL)/remote-sample$$i.pdf" & \
	done
	wait


clean:
	rm -rf bin
	rm -rf $(OUTPUT_DIR)