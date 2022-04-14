FROM golang:1.18.1 as builder

RUN mkdir /build
COPY . /build/
WORKDIR /build

RUN go test ./...

RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN apt-get update
RUN apt install -y protobuf-compiler
RUN ./protoc.sh

RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -o service-binary .


FROM alpine:latest
COPY --from=builder /build/service-binary .
RUN mkdir config
COPY ./config/config.json ./config/config.json

# executable
ENTRYPOINT [ "./service-binary" ]