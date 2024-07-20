FROM golang:1.22.2-alpine AS build-stage

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/main ./main

USER nonroot:nonroot

CMD ["./main"]