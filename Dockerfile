FROM golang:1.19 as build
WORKDIR /go/src/
COPY go.* ./
RUN go mod download

COPY main.go .
# --mount flag allows for speedier builds of the binary
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o config-puller main.go

# Final image in the build process
FROM scratch
COPY --from=build /go/src/config-puller /config-puller

ENTRYPOINT ["/config-puller"]
