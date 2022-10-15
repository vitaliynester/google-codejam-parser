FROM golang:1.18-alpine AS build
RUN apk add --no-cache --update alpine-sdk
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /build/app cmd/app/main.go


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /build/app ./
CMD ["./app"]