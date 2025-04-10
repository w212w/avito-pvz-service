FROM golang:1.23

WORKDIR /app

COPY ../go.mod ../go.sum ./

RUN go mod download

COPY . .

RUN go build -o avito-pvz-service ./cmd/pvz-service \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/app/avito-pvz-service"]