FROM golang:1.16

RUN mkdir /build
COPY tests/. /build/
RUN mkdir /proto
COPY proto/. /proto
RUN mkdir /build/config
COPY config/. /build/config
WORKDIR /build

RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN apt-get update
RUN apt install -y protobuf-compiler
RUN ./protoc.sh

RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -o test-binary .

ENTRYPOINT [ "./test-binary" ]
