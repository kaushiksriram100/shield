clean:
	rm -rf shield
shield: clean
	go build -o bin/shield server/*.go
