FROM golang as build

WORKDIR /app

COPY . .
RUN go build

FROM alpine:latest as product

WORKDIR /app
COPY --from=build /app/userApi /bin/userApi

RUN chmod 777 /bin/userApi

ARG PostgresHost=localhost
ARG PostgresUser=postgres
ARG PostgresPassword=admin
ARG PostgresDatabase=postgres
ARG postgresPort=5432
ARG RPDisplayName=LocalTest
ARG RPID=localhost
ARG RPOrigin=http://localhost:8080
ARG RPIcon=https://duo.com/logo.png
ARG AppListen=":80"

EXPOSE 80
