FROM alpine:latest as release
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY github-actions-exporter .
CMD ["./github-actions-exporter"]