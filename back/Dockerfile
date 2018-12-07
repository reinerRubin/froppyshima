FROM golang

WORKDIR /src

# cache dependencies
COPY go.sum .
COPY go.mod .
RUN go list -e $(go list -f '{{.Path}}' -m all 2>/dev/null)

# build app
ADD . /src
RUN go build -o main .

CMD ["/src/main"]
