repository=cloud
user=hatlonely
version=$(shell git describe --tags)

export GOPATH=$(shell pwd)/../../../../
export PATH:=${PATH}:${GOPATH}/bin:$(shell pwd)/third/go/bin:$(shell pwd)/third/protobuf/bin:$(shell pwd)/third/cloc-1.76:$(shell pwd)/third/redis-3.2.8/src

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	sedi=sed -i ""
else
	sedi=sed -i
endif

.PHONY: all
all: third vendor output test stat

deploy:
	mkdir -p /var/docker/${repository}/log
	docker stack deploy -c stack.yml ${repository}

remove:
	docker stack rm ${repository}

push:
	docker push ${user}/${repository}:${version}

.PHONY: buildenv
buildenv:
	docker run --rm --name go-build-env -d golang:1.12.5 tail -f /dev/null

.PHONY: image
image:
	docker exec -it go-build-env rm -rf /data/src/hpifu/${repository}
	docker exec -it go-build-env mkdir -p /data/src/hpifu/${repository}
	docker cp . go-build-env:/data/src/hpifu/${repository}
	docker exec -it go-build-env bash -c "cd /data/src/hpifu/${repository} && make output"
	mkdir -p docker/
	docker cp go-build-env:/data/src/hpifu/${repository}/output/${repository} docker/
	docker build --tag=hatlonely/${repository}:`git describe --tags` .
	${sedi} 's/image: ${user}\/${repository}:.*$$/image: ${user}\/${repository}:${version}/g' stack.yml

output: cmd/*/*.go internal/*/*.go scripts/version.sh Makefile vendor
	@echo "compile"
	@go build -ldflags "-X 'main.AppVersion=`sh scripts/version.sh`'" cmd/${repository}/main.go && \
	mkdir -p output/${repository}/bin && mv main output/${repository}/bin/${repository} && \
	mkdir -p output/${repository}/configs && cp configs/${repository}/* output/${repository}/configs && \
	mkdir -p output/${repository}/log

vendor: go.mod go.sum
	@echo "install golang dependency"
	go mod vendor

%_easyjson.go: %.go
	easyjson $<

%.pb.go: %.proto
	protoc --gofast_out=plugins=grpc:. $<

.PHONY: test
test: vendor
	@echo "Run unit tests"
	- cd internal && go test -cover ./...
	- cd pkg && go test -cover ./...

.PHONY: behave
behave: output
	@behave features

.PHONY: stat
stat: cloc gocyclo
	@echo "code statistics"
	@cloc pkg Makefile --by-file
	@echo "circle complexity statistics"
	@gocyclo pkg
	@gocyclo pkg | awk '{sum+=$$1}END{printf("complexity: %s", sum)}'

.PHONY: clean
clean:
	rm -rf output

.PHONY: deep_clean
deep_clean:
	rm -rf output vendor third

third: protoc golang cloc gocyclo easyjson

.PHONY: protoc
protoc: golang
	@hash protoc 2>/dev/null || { \
		echo "install protobuf codegen tool protoc" && \
		mkdir -p third && cd third && \
		wget https://github.com/google/protobuf/releases/download/v3.2.0/protobuf-cpp-3.2.0.tar.gz && \
		tar -xzvf protobuf-cpp-3.2.0.tar.gz && \
		cd protobuf-3.2.0 && \
		./configure --prefix=`pwd`/../protobuf && \
		make -j8 && \
		make install && \
		cd ../.. && \
		protoc --version; \
	}
	@hash protoc-gen-go 2>/dev/null || { \
		echo "install protobuf golang plugin protoc-gen-go" && \
		go get -u github.com/golang/protobuf/{proto,protoc-gen-go}; \
	}

.PHONY: golang
golang:
	@hash go 2>/dev/null || { \
		echo "install go1.9" && \
		mkdir -p third && cd third && \
		wget https://dl.google.com/go/go1.9.linux-amd64.tar.gz && \
    	tar -xzvf go1.9.linux-amd64.tar.gz && \
		cd .. && \
		go version; \
	}

.PHONY: cloc
cloc:
	@hash cloc 2>/dev/null || { \
		echo "install cloc" && \
		mkdir -p third && cd third && \
		wget https://github.com/AlDanial/cloc/archive/v1.76.zip && \
		unzip v1.76.zip; \
	}

.PHONY: gocyclo
gocyclo: golang
	@hash gocyclo 2>/dev/null || { \
		echo "install gocyclo" && \
		go get -u github.com/fzipp/gocyclo; \
	}

.PHONY: easyjson
easyjson: golang
	@hash easyjson 2>/dev/null || { \
		echo "install easyjson" && \
		go get -u github.com/mailru/easyjson/...; \
	}

.PHONY: redis
redis:
	@hash redis 2>/dev/null || { \
		echo "install redis" && \
		mkdir -p third && cd third && \
		wget http://download.redis.io/releases/redis-3.2.8.tar.gz && \
		tar -xzvf redis-3.2.8.tar.gz && \
		cd redis-3.2.8 && \
		make -j8; \
	}
