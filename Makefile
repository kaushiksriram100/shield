clean:
	rm -rf shield
shield: clean
	go build -o shield server/*.go
