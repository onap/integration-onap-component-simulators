FROM golang:1.17-alpine AS builder

# Add all the source code (except what's ignored
# under `.dockerignore`) to the build context.
ADD ./ /go/src/

WORKDIR /go/src

RUN apk add --no-cache git
RUN go mod init onap.com/mock-clamp
RUN go get github.com/satori/go.uuid
RUN go get github.com/labstack/echo
RUN go get github.com/dgrijalva/jwt-go

RUN set -ex && \
  CGO_ENABLED=0 GOOS=linux go build \
        -tags netgo \
        -installsuffix cgo \
        -v -a \
        -ldflags '-extldflags "-static"' \
        -o mock-clamp .

RUN ls -la

FROM scratch

# Retrieve the binary from the previous stage
COPY --from=builder /go/src/mock-clamp /app/mock-clamp
WORKDIR /app

# Set the binary as the entrypoint of the container
ENTRYPOINT [ "./mock-clamp" ]
