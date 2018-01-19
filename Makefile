build: bin bin/vms

bin:
	mkdir bin

bin/vms:
	go build -o bin/vms vms.go

clean:
	rm -rf bin/

container: build
	docker build -f Dockerfile -t 2pg/vms .
