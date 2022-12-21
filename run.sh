#!/bin/zsh

go build -o bookings cmd/web/*.go && ./bookings -dbname=bookings -dbuser=postgres -dbport=55000 -dbpass=postgrespw -cache=false -production=false


