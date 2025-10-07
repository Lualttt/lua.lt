FROM golang:1.22-alpine AS build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/lualt

FROM scratch

WORKDIR /app

COPY --from=build /go/bin/app .

EXPOSE 8080

ENTRYPOINT ["./app"]
