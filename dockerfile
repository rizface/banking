## Build
FROM golang:1.22.1-alpine AS build

WORKDIR $GOPATH/src/banking

# manage dependencies
COPY . .
RUN go mod download

RUN go build -a -o /banking-server ./main.go


## Deploy
FROM alpine:latest
RUN apk add tzdata
COPY --from=build /banking-server /banking-server

EXPOSE 8080

ENTRYPOINT ["/banking-server"]