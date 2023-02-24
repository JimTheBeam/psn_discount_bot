#!/bin/bash

# make migrations
goose -dir ./migrations postgres "postgres://postgres:postgres@db_postgresql:5432/psn_discount?sslmode=disable" up

# run the service
./psn_discounter -c config.yaml
