up:
	docker build -t modern_art .
	docker run -d -p 5000:5000 --name modern_art modern_art

down:
	docker rm modern_art
	docker rmi modern_art

local:
	go run main.go

compile:
	go build -o modern_art .