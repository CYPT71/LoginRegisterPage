FROM golang

WORKDIR /app

COPY . .
RUN go build

