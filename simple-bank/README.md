# Simple Bank Example

## Create Database Schema

https://dbdiagram.io/d

## Golang migrate

https://github.com/golang-migrate/migrate

migrate create -ext sql -dir db/migration -seq init_schema

migrate -path db/migration -database "postgresql://dbadmin:admin123@localhost:5432/simple-bank?sslmode=disable"  -verbose up