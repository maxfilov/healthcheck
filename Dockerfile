# syntax=docker/dockerfile:1
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN make

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/bin/healthcheck /usr/bin/healthcheck
COPY config.yaml /etc/healthcheck/config.yaml

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/usr/bin/healthcheck", "--config=/etc/healthcheck/config.yaml"]
