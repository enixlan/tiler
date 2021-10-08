FROM ubuntu:20.04

RUN apt-get update -y && apt-get install --no-install-recommends -y \
    ca-certificates \
    curl \
    gnupg && \
    rm -rf /var/lib/apt/lists/* \

RUN echo 'deb http://ppa.launchpad.net/osmadmins/ppa/ubuntu bionic main\n\
deb-src http://ppa.launchpad.net/osmadmins/ppa/ubuntu bionic main' > \
    /etc/apt/sources.list.d/osmadmins-ppa.list

RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 \
    --recv A438A16C88C6BE41CB1616B8D57F48750AC4F2CB

RUN apt-get update && apt-get install --no-install-recommends -y \
    osm2pgsql \
    osmium-tool \
    gdal-bin \
    python3-psycopg2 \
    python3-yaml \
    python3-requests \
    postgresql-client && \
    rm -rf /var/lib/apt/lists/*

RUN ldconfig

RUN curl -L https://github.com/gravitystorm/openstreetmap-carto/archive/v5.4.0.tar.gz | tar -zxf - && \
    mv openstreetmap-carto-5.4.0 /openstreetmap-carto

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT /entrypoint.sh
