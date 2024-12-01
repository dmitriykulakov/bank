FROM golang:alpine

WORKDIR ./

COPY ./ ./

CMD ["go", "run", "cmd/main.go"]