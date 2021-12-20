build:
	mkdir -p bin
	rm -rf bin/*
	go build -o bin/azkaban_exporter cmd/main.go
	cp conf/azkaban.yml ./bin

run:
	$(MAKE) build
	./bin/azkaban_exporter