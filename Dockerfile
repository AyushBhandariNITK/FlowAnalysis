FROM golang:1.22.2 AS build
ENV GO1111ODUE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY pkg ./pkg
COPY main.go .

COPY go.mod .
RUN go mod tidy
RUN go mod download

RUN go build -o flowanalysis .

WORKDIR /dist
RUN cp /build/flowanalysis .

FROM ubuntu:22.04

RUN apt-get update -y && \
    apt-get install -y --no-install-recommends curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /
COPY --from=build /dist/flowanalysis .
EXPOSE 5010

ENTRYPOINT ["./flowanalysis"]
