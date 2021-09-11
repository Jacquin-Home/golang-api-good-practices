# build stage
FROM golang:alpine3.14 as build-stage

WORKDIR /opt/app/build

COPY . .

RUN go mod download \
    && go build -o app

# final
FROM golang:alpine3.14

WORKDIR /opt/app

COPY --from=build-stage /opt/app/build/app .

CMD ["./app"]
