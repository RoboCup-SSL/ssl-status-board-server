#!/usr/bin/env bash

k6 run --out influxdb=http://influxdb-svc:8086/k6 config.js