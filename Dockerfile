FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY github-actions-exporter app
CMD ["./app"]