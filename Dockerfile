FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o openalysis .


FROM alpine

WORKDIR /src
COPY --from=build /app/openalysis /app/openalysis
COPY default.yaml /src/default.yaml

ENTRYPOINT ["/app/openalysis"]