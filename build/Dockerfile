FROM golang:1.15.7-alpine3.13 AS builder

ENV GO111MODULE=on
RUN apk --no-cache add git make build-base
RUN go get github.com/rakyll/statik
WORKDIR /src

COPY . .

RUN mkdir -p /app
RUN make statik
# -a 強制リビルド
RUN go build -a -ldflags="-s -w" -o /app /src/cmd/server/

FROM alpine:3.10.2
WORKDIR /app
# Timezone = Tokyo
RUN apk --no-cache add tzdata ca-certificates && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

COPY --from=builder /app /app
RUN chmod u+x /app

ARG PRIVATE_KEY
ENV PRIVATE_KEY ${PRIVATE_KEY}

ENTRYPOINT ["/app/server"]
