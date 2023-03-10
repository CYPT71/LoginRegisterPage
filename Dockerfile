ARG  DISTROLESS_IMAGE=gcr.io/distroless/static:nonroot
ARG AppListen=80

# Use the official Go image as the base image
FROM golang:alpine3.17 AS build

# Set the working directory to /app
WORKDIR /app

# Install deps
RUN apk update --no-cache && apk add pkgconf openssl-dev gcc libc-dev

# Copy the source code to the container
COPY . .

ENV GO111MODULE=on
RUN go mod download

# Build the binary
RUN GOOS=linux GOARCH=amd64 go build -tags static -o /go/bin/userApp .

# Use a lightweight image for the final stage
FROM ${DISTROLESS_IMAGE}

USER 65532:65532

# Copy the binary from the previous stage
COPY --from=build /go/bin/userApp /go/bin/userApp
COPY --from=builder /lib/libssl.so.3     /lib/libssl.so.3
COPY --from=builder /lib/libcrypto.so.3     /lib/libcrypto.so.3
COPY --from=builder /lib/ld-musl-x86_64.so.1  /lib/ld-musl-x86_64.so.1


ENV PostgresHost=local_pgdb
ENV PostgresUser=postgres
ENV PostgresPassword=admin
ENV PostgresDatabase=postgres
ENV PostgresPort=5432
ENV Origins=http://localhost:8080
ENV APIS=http://127.0.0.1:5000
ENV RPDisplayName=LocalTest
ENV RPID=localhost
ENV RPOrigin=http://localhost:80
ENV RPIcon=https://duo.com/logo.png
ENV AppListen=":80"

EXPOSE 80

CMD ["/go/bin/userApp"]



