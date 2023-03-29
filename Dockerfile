FROM alpine:3

EXPOSE 8080

RUN apk add ca-certificates

COPY ./bin/cmd /

CMD ["./cmd"]
