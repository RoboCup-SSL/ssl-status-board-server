#!/usr/bin/env bash

users=100

k6 run config.js -u $users -i $users