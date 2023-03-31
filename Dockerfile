FROM alpine:3
ARG JOB_NAME

EXPOSE 8080

RUN apk add ca-certificates

COPY ./cmd/$JOB_NAME/bin/$JOB_NAME /application
COPY ./cmd/$JOB_NAME/config/config.json /config/config.json

CMD ["./application"]
