FROM golang:alpine AS builder

# Add all the source code (except what's ignored
# under `.dockerignore`) to the build context.
ADD ./ /go/src/

WORKDIR /go/src

RUN apk add --no-cache git
RUN go env -w GO111MODULE=auto
RUN go get github.com/satori/go.uuid
RUN go get github.com/labstack/echo
RUN go get github.com/golang-jwt/jwt
RUN go get golang.org/x/time/rate

RUN set -ex && \
  CGO_ENABLED=0 GOOS=linux go build \
        -tags netgo \
        -installsuffix cgo \
        -v -a \
        -ldflags '-extldflags "-static"' \
        -o mock-sdc .

RUN ls -la

FROM scratch

# Retrieve the binary from the previous stage
COPY --from=builder /go/src/mock-sdc /app/mock-sdc
WORKDIR /app

# Set the binary as the entrypoint of the container
ENTRYPOINT [ "./mock-sdc" ]
