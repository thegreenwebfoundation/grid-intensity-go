FROM alpine:3.18

RUN apk add --no-cache ca-certificates

ADD ./grid-intensity /grid-intensity

EXPOSE 8000/tcp

ENTRYPOINT ["/grid-intensity"]
