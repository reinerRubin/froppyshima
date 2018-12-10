FROM golang as builder

WORKDIR /src

# cache dependencies
COPY go.sum .
COPY go.mod .
RUN go list -e $(go list -f '{{.Path}}' -m all 2>/dev/null)

# build the app
ADD . /src
RUN CGO_ENABLED=0 go build -o main .

# run tests
FROM builder as test
WORKDIR /src
RUN go test -v ./...

# do a release container
FROM alpine:latest as release
WORKDIR /root
COPY --from=builder /src/main .

CMD ["/root/main"]
