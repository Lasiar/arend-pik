FROM amd64/golang as builder
WORKDIR /go/src/pik-arenda/
COPY ./*.go /go/src/pik-arenda/
COPY ./web/*.go /go/src/pik-arenda/web/
COPY ./base/*.go /go/src/pik-arenda/base/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main

FROM amd64/alpine
COPY --from=builder /go/src/pik-arenda/main /app/
COPY pik-arenda/config.prod.json /etc/pik-arenda/config.json
WORKDIR /app
EXPOSE 80/tcp
CMD ["/app/main"]