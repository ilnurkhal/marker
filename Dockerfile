FROM golang:latest AS build
WORKDIR /forge/

COPY . ./
WORKDIR /forge/cmd/marker
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o marker


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=build /forge/cmd/marker .
RUN adduser -S marker --uid 1000
USER marker
ENTRYPOINT ["/app/marker"]
