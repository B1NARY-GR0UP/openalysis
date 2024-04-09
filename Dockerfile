FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o oa .


FROM alpine

WORKDIR /src
COPY --from=build /app/oa /app/oa
COPY default.yaml /src/default.yaml

ENTRYPOINT ["/app/oa"]