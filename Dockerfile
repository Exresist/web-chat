
# Stage 1 - build executable in go container
FROM golang:latest as builder

WORKDIR $GOPATH/src/web-chat/
COPY . .

RUN export CGO_ENABLED=0 && make build

# Stage 2 - build final image
FROM alpine:latest

# Copy our static executable and config
COPY --from=builder /go/src/web-chat/config.json config.json

COPY --from=builder /go/src/web-chat/bin/web-chat go/bin/web-chat

RUN chmod +x go/bin/web-chat

# Run the binary.
RUN ["go/bin/web-chat"]