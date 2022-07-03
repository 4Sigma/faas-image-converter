FROM --platform=${TARGETPLATFORM:-linux/amd64} ghcr.io/openfaas/of-watchdog:0.9.6 as watchdog
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.16-bullseye as build
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ENV GO111MODULE=on

RUN apt update -y && apt install -y \
    build-essential \
    curl \
    git \
    sudo \
    libcurl4-openssl-dev \
    libssl-dev \
    libxml2-dev \
    libxslt1-dev \
    libyaml-dev \
    zlib1g-dev

WORKDIR /app
COPY --from=watchdog /fwatchdog /usr/bin/fwatchdog
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./faas-image-processor

FROM --platform=${TARGETPLATFORM:-linux/amd64} debian:bullseye
#  #Add non root user and certs
RUN apt update -y && apt install -y \
       ca-certificates \
    && adduser --disabled-password --gecos '' app \
    && adduser app sudo \
    && echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

#Split instructions so that buildkit can run & cache
#the previous command ahead of time.
RUN mkdir -p /home/app \
    && chown app /home/app

WORKDIR /home/app

COPY --from=build --chown=app /app/faas-image-processor         .
COPY --from=build --chown=app /usr/bin/fwatchdog                .

USER app

ENV fprocess="./faas-image-processor"
ENV mode="http"
ENV upstream_url="http://127.0.0.1:8082"
ENV prefix_logs="false"
#
CMD ["./fwatchdog"]
