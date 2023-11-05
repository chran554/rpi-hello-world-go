
.PHONY: clean
clean:
	rm -rf bin/

.PHONY: build
build: clean
	GOOS=linux GOARCH=arm GOARM=5 go build -o bin/hello-world ./cmd

.PHONY: deploy
deploy:
	scp ./bin/hello-world pi@192.168.86.50:tmp
