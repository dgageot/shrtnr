FROM golang:1.15.2-alpine3.12 as build
WORKDIR /app
COPY . ./
RUN go build -o shrtnr

FROM alpine:3.12
EXPOSE 8080
ENTRYPOINT ["/shrtnr"]
COPY --from=build /app/shrtnr /
