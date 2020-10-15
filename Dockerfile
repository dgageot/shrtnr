FROM golang:1.15.2-alpine3.12 as build
WORKDIR /app
COPY . ./
RUN go build -o shrtnr

FROM alpine:3.12
EXPOSE 8888
WORKDIR /root
ENTRYPOINT ["/shrtnr"]
COPY --from=build /app/shrtnr /
