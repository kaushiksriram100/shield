clean:
	rm -rf shieldapp
shield: clean
	go build -o shieldapp server/*.go
