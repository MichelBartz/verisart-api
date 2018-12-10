FROM golang:1.11.2-stretch
# Installing dep to fetch our dependencies
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/verisart-api

## Only copy the bare minimum
COPY main.go .
COPY internal/ ./internal
COPY Gopkg.toml .
COPY Gopkg.lock .

# Fetch the dependencies
RUN dep ensure

# Build & Install
RUN go install -v ./...

CMD ["verisart-api"]