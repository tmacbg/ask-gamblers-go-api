dev:
	go run server/main.go
live:
	go build -o askGamblersApi server/main.go  && ./askGamblersApi