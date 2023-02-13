FROM golang:1.18.0-alpine3.15 AS build
RUN apk add --update make
WORKDIR /etc

ARG GOBIN=/go/bin

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
#Build the code
RUN make shield
#RUN go build -o shieldapp server/*

FROM alpine:3.15
WORKDIR /bin

COPY --from=build /etc/shieldapp ./
COPY --from=build /etc/server/blacklist.json ./

ENTRYPOINT ["/bin/shieldapp"]
