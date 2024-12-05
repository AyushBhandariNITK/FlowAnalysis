FROM golang:1.22.2-alpine AS build

RUN apk add --no-cache git
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./ 
COPY pkg ./pkg

RUN go build -o flowanalysis .

FROM alpine:3.18

RUN apk add --no-cache curl

WORKDIR /app
COPY --from=build /app/flowanalysis .

EXPOSE 5010

ENTRYPOINT ["./flowanalysis"]

