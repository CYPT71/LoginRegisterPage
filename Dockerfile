FROM golang as build

WORKDIR /app

COPY . .
ENV CGO_ENABLED=0
RUN go build -o userApp

FROM alpine:latest as product

WORKDIR /app
COPY --from=build /app/userApp /bin/userApp

RUN chmod 777 /bin/userApp

ENV PostgresHost=localhost
ENV PostgresUser=postgres
ENV PostgresPassword=admin
ENV PostgresDatabase=postgres
ENV postgresPort=5432
ENV RPDisplayName=LocalTest
ENV RPID=localhost
ENV RPOrigin=http://localhost:8080
ENV RPIcon=https://duo.com/logo.png
ENV AppListen=":80"

EXPOSE 80

CMD ["userApp"]
