FROM golang:1.8
  

WORKDIR /go/src/askGamblersApi

COPY . .

RUN ls -la
RUN go get github.com/mattn/go-sqlite3
RUN go build -o ask-gamblers-api server/main.go

EXPOSE 5001
CMD ["./ask-gamblers-api"]

