FROM golang:1.24.1 as build
WORKDIR /go/src/
COPY go.* ./
RUN go mod download

COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o config-puller main.go

# Final image in the build process
FROM alpine:3.20.1
COPY --from=build /go/src/config-puller /config-puller

ENTRYPOINT ["/config-puller"]
