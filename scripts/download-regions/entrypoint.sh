#!/bin/bash

export PGPASSWORD="${DB_PASSWORD}"
REGIONS=${REGIONS}
REGIONS_ARRAY=($REGIONS)

mkdir -p /tile-data
mkdir -p /shapefiles

FILES=""

echo "== download regions ======================="
for i in "${!REGIONS_ARRAY[@]}"; do
  name=${REGIONS_ARRAY[$i]};
	url="http://download.geofabrik.de/$name-latest.osm.pbf";

  echo "== fetching $name";

  curl "$url" -o "/tile-data/$i.pbf";

  FILES="$FILES /tile-data/$i.pbf"
done

echo "== merge regions =========================="
osmium merge $FILES -o /tile-data/sum.pbf

echo "== upload regions to postgres ============="
osm2pgsql \
  --username="${DB_USER:-renderer}" \
  --host="${DB_HOST:-localhost}" \
  --database="${DB_NAME:-gis}" \
  --port="${DB_PORT:-5432}" \
  --cache="${CACHE_SIZE:-800}" \
  --create \
  --slim \
  -G \
  --hstore \
  --number-processes "${THREADS:-$(nproc)}" \
  --multi-geometry \
  --tag-transform-script /openstreetmap-carto/openstreetmap-carto.lua \
  --style /openstreetmap-carto/openstreetmap-carto.style \
   /tile-data/sum.pbf

echo "== downloading needed shapefiles =========="
/openstreetmap-carto/scripts/get-external-data.py \
  --data="/shapefiles" \
  --username="${DB_USER:-renderer}" \
  --host="${DB_HOST:-localhost}" \
  --database="${DB_NAME:-gis}" \
  --port="${DB_PORT:-5432}" \
  --password="${DB_PASSWORD}" \
  --config="/openstreetmap-carto/external-data.yml"

echo "== cleanup ================================"
rm -r /tile-data
rm -r /shapefiles
