FROM alpine:3.18.3
RUN apk add --no-cache ca-certificates git
WORKDIR /home
ADD ./gitcache .
ENTRYPOINT [ "sh", "-c", "./gitcache"]
