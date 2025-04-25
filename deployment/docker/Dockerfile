FROM golang:1.19.2-alpine3.15 as build

WORKDIR /app

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:3.15 as runner

COPY --from=build /app/principal-credential-creation-service /opt/principal-credential-creation-service

WORKDIR /opt

CMD ["./principal-credential-creation-service"]
