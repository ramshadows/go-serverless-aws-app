.PHONY: build clean deploy gomodgen

build: gomodgen
	env GOARCH=amd64 GOOS=linux go build -o bin/main cmd/main.go
	

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
