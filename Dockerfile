FROM golang:1.17.5-alpine3.15 as build
WORKDIR /app
COPY . ./
RUN go build -o shrtnr

FROM alpine:3.15.0
EXPOSE 8888
WORKDIR /root
ENTRYPOINT ["/shrtnr"]
COPY --from=build /app/shrtnr /
