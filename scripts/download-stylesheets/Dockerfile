FROM node:16.4.2-alpine AS stylesheet

ARG OPENSTREETMAP_CARTO_VERSION=5.4.0

RUN apk add --no-cache curl

RUN npm install -g carto

RUN curl -L https://github.com/gravitystorm/openstreetmap-carto/archive/v${OPENSTREETMAP_CARTO_VERSION}.tar.gz | tar -zxf - && \
    mv openstreetmap-carto-${OPENSTREETMAP_CARTO_VERSION} /openstreetmap-carto

RUN sed -i "/^    type: \"postgis\"/a\    host: \"{{DB_HOST}}\"" /openstreetmap-carto/project.mml
RUN sed -i "/^    type: \"postgis\"/a\    port: \"{{DB_PORT}}\"" /openstreetmap-carto/project.mml
RUN sed -i "/^    type: \"postgis\"/a\    password: \"{{DB_PASSWORD}}\"" /openstreetmap-carto/project.mml
RUN sed -i "/^    type: \"postgis\"/a\    user: \"{{DB_USER}}\"" /openstreetmap-carto/project.mml
RUN sed -i "s/dbname: \"gis\"/dbname: \"{{DB_NAME}}\"/" /openstreetmap-carto/project.mml

RUN carto /openstreetmap-carto/project.mml -f /openstreetmap-carto/mapnik.xml
