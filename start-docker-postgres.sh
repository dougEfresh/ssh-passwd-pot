#!/bin/bash
docker run  -p 5432:5432 -v "$PWD/schema/psql:/docker-entrypoint-initdb.d/" -e POSTGRES_USER=postgres  -d postgres
