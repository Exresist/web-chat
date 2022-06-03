FROM golang:latest as builder

WORKDIR $GOPATH/src/web-chat/
COPY . .

RUN export CGO_ENABLED=0 make build


FROM alpine:latest

COPY --from=builder /go/src/web-chat/config.json config.json

COPY --from=builder /go/src/web-chat/bin/web-chat go/bin/web-chat


ENTRYPOINT ["go/bin/web-chat"]