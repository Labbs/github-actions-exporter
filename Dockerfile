FROM golang:1.18 as builder

WORKDIR /app
COPY . .
RUN bash ./build.sh

FROM alpine:latest as release
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /app/bin/app .
CMD ["./app"]