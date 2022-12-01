FROM golang:1.18 as builder

WORKDIR /app
COPY . .
RUN bash ./build.sh

FROM alpine:3.17.0 as release
RUN apk update \
    && apk add --no-cache ca-certificates=20220614-r2 \
    && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=builder /app/bin/app .
CMD ["./app"]