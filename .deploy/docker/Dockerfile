# BUILD BACKEND
FROM golang:1.17-alpine as builder

# Versions
ARG MAPNIK_VERSION="v3.1.0"
ARG PROJ4_VERSION="4.9.3"

# Download Locations
ARG MAPNIK_SOURCE="https://github.com/mapnik/mapnik/releases/download/$MAPNIK_VERSION/mapnik-$MAPNIK_VERSION.tar.bz2"
ARG PROJ4_SOURCE="https://download.osgeo.org/proj/proj-$PROJ4_VERSION.tar.gz"

# Install deps.
RUN apk upgrade --update

RUN apk add --virtual .build-deps \
    build-base \
    wget \
    bash \
    git \
    python3-dev \
    boost-dev \
    cairo-dev \
    libwebp-dev \
    harfbuzz-dev \
    icu-dev \
    jpeg-dev \
    tiff-dev \
    sqlite-dev

RUN apk add --virtual .run-deps \
    python3 \
    boost \
    boost-python3 \
    py3-cairo \
    libwebp \
    harfbuzz \
    icu \
    gdal-dev \
    openssl-dev \
    postgresql-dev \
    font-noto-all \
    unifont \
    ttf-dejavu

# Install proj4.
RUN wget $PROJ4_SOURCE -O /tmp/proj-$PROJ4_VERSION.tar.gz \
    && mkdir -p /usr/src \
    && tar zxf /tmp/proj-$PROJ4_VERSION.tar.gz -C /usr/src \
    && cd /usr/src/proj-$PROJ4_VERSION \
    && ./configure \
    && make \
    && make install

# Install mapnik.
RUN wget https://github.com/mapnik/mapnik/releases/download/v3.1.0/mapnik-v3.1.0.tar.bz2 -O /tmp/mapnik-v3.1.0.tar.bz2 \
    && tar xjf /tmp/mapnik-v3.1.0.tar.bz2 -C /usr/src \
    && cd /usr/src/mapnik-v3.1.0 \
    && ln -sf python3 /usr/bin/python \
    && ./configure \
    && make \
    && make install

WORKDIR /src

COPY go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY . .

ARG SERVICE_NAME=tiler
ARG APP_NAME=github.com/enixlan/tiler

RUN GOOS=linux GOARCH=amd64 \
    CGO_LDFLAGS="$(mapnik-config --libs)" \
    CGO_ENABLED=1 \
    CGO_CXXFLAGS="$(mapnik-config --cxxflags --includes --dep-includes | tr '\n' ' ')" \
    go build \
    -ldflags="-w -s \
    -X github.com/loghole/tron/internal/app.ServiceName=$SERVICE_NAME \
    -X github.com/loghole/tron/internal/app.AppName=$APP_NAME \
    -X github.com/loghole/tron/internal/app.GitHash=$(git rev-parse HEAD) \
    -X github.com/loghole/tron/internal/app.Version=$(git describe --tags --always) \
    -X github.com/loghole/tron/internal/app.BuildAt=$(date --utc +%FT%TZ) \
    " -o /app cmd/$SERVICE_NAME/*.go

# Clean up files
RUN apk del .build-deps \
    && rm -rf /tmp/* /var/cache/apk/* /usr/src/* /src

RUN apk add --no-cache bash

COPY stylesheets/patterns/ /etc/mapnik/patterns
COPY stylesheets/symbols/ /etc/mapnik/symbols
COPY stylesheets/mapnik.xml /etc/mapnik/mapnik.xml

ENV MAPNIK_CONFIG_PATH=/etc/mapnik/mapnik.xml
