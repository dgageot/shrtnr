FROM golang:1.16.4-alpine3.13 as build
WORKDIR /app
COPY . ./
RUN go build -o shrtnr

FROM alpine:3.13.5
EXPOSE 8888
WORKDIR /root
ENTRYPOINT ["/shrtnr"]
COPY --from=build /app/shrtnr /
