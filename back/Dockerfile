FROM golang as builder

WORKDIR /src

# cache dependencies
COPY go.sum .
COPY go.mod .
RUN go list -e $(go list -f '{{.Path}}' -m all 2>/dev/null)

# build app
ADD . /src
RUN CGO_ENABLED=0 go build -o main .

# "DOCKER_BUILDKIT=1" would skip this stage for the release target
# https://github.com/docker/cli/issues/1134#issuecomment-406449342
FROM builder as test
WORKDIR /src
RUN go test -v ./...

FROM alpine:latest as release
WORKDIR /root
RUN ls
COPY --from=builder /src/main .
RUN ls

CMD ["/root/main"]
