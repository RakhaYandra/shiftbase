# shiftbase

Shift scheduling & attendance REST API — Go/Gin + MySQL + JWT + Docker.

> Status: M1 — scaffold + auth. Employees/shifts (M2), attendance/CSV import (M3),
> reports + Swagger/Postman (M4), tests + CI + docs poles (M5) menyusul.

## Quickstart (5 menit)

```bash
cp .env.example .env            # sesuaikan JWT_SECRET di prod
docker compose up --build       # api:8080, mysql:3306, adminer:8082
# migrasi + seed (goose CLI, tanpa lib ORM):
go install github.com/pressly/goose/v3/cmd/goose@v3.28.0
goose -dir migrations mysql "shift:shiftpass@tcp(localhost:3306)/shiftbase?parseTime=true" up
mysql -h127.0.0.1 -ushift -pshiftpass shiftbase < seed/seed.sql
```

## Contoh curl

```bash
curl -s localhost:8080/healthz
curl -s -X POST localhost:8080/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"budi@example.com","password":"Rahasia123"}'
TOKEN=$(curl -s -X POST localhost:8080/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@shiftbase.local","password":"Admin123!"}' | cut -d'"' -f4)
curl -s localhost:8080/v1/me -H "Authorization: Bearer $TOKEN"
```

Seed login: `admin@shiftbase.local / Admin123!`, `manager@shiftbase.local / Manager123!`,
`staff@shiftbase.local / Staff123!` (data fiktif untuk portofolio).
