version = 1.0.0
temp_dir_name = azkaban_exporter-${version}
tar_name = azkaban_exporter-${version}.tar.gz

build:
	$(MAKE) clean
	mkdir -p bin
	go build -o bin/azkaban_exporter cmd/main.go
	cp conf/azkaban.yml bin

package:
	$(MAKE) build
	mkdir -p bin/${temp_dir_name}
	cp bin/azkaban.yml bin/${temp_dir_name}
	cp bin/azkaban_exporter bin/${temp_dir_name}
	tar zcvf bin/${tar_name} -C bin ${temp_dir_name}
	rm -rf bin/${temp_dir_name}

run:
	$(MAKE) build
	bin/azkaban_exporter --azkaban.conf=bin/azkaban.yml

clean:
	rm -rf bin
