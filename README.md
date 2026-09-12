# shiftbase

Shift scheduling & attendance REST API — **Go/Gin + MySQL + JWT + Docker**.
`database/sql` tanpa ORM, migrasi goose. Data contoh fiktif (portofolio).

## Quickstart 5 menit

```bash
git clone https://github.com/RakhaYandra/shiftbase.git && cd shiftbase
cp .env.example .env                     # ganti JWT_SECRET di prod
docker compose up --build                # api:8080, mysql:3306, adminer:8082
go install github.com/pressly/goose/v3/cmd/goose@v3.28.0
goose -dir migrations mysql "shift:shiftpass@tcp(localhost:3306)/shiftbase?parseTime=true" up
mysql -h127.0.0.1 -ushift -pshiftpass shiftbase < seed/seed.sql
curl -s localhost:8080/healthz          # {"status":"ok"}
```

Login seed: `admin@shiftbase.local / Admin123!`,
`manager@shiftbase.local / Manager123!`, `staff@shiftbase.local / Staff123!`.

## Contoh curl

```bash
B=localhost:8080
curl -s -X POST $B/v1/auth/register -H 'Content-Type: application/json' \
  -d '{"email":"budi@example.com","password":"Rahasia123"}'
TOKEN=$(curl -s -X POST $B/v1/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"admin@shiftbase.local","password":"Admin123!"}' | cut -d'"' -f4)
A="Authorization: Bearer $TOKEN"
curl -s $B/v1/me -H "$A"

# employees (admin)
E=$(curl -s -X POST $B/v1/employees -H "$A" -H 'Content-Type: application/json' \
  -d '{"name":"Fajar","email":"fajar@example.com","position":"kasir","hire_date":"2025-02-01"}')
echo $E; EID=$(echo $E | grep -o '"id":[0-9]*' | cut -d: -f2)
curl -s $B/v1/employees -H "$A"

# shifts (admin/manager; bentrok → 409)
curl -s -X POST $B/v1/shifts -H "$A" -H 'Content-Type: application/json' \
  -d "{\"employee_id\":$EID,\"date\":\"2030-01-05\",\"start_time\":\"08:00\",\"end_time\":\"16:00\"}"
curl -s $B/v1/shifts?date=2030-01-05 -H "$A"

# attendance
curl -s -X POST $B/v1/attendance/check-in -H "$A" -H 'Content-Type: application/json' \
  -d "{\"employee_id\":$EID}"
curl -s -X POST $B/v1/attendance/check-out -H "$A" -H 'Content-Type: application/json' \
  -d "{\"employee_id\":$EID}"
curl -s "$B/v1/attendance?from=2030-01-01&to=2030-01-31" -H "$A"

# import CSV (admin/manager; maks 2MB, header name,email,phone,position,hire_date)
printf 'name,email,phone,position,hire_date\nGita,gita@example.com,0812,kasir,2025-03-01\n' > /tmp/emp.csv
curl -s -X POST $B/v1/employees/import -H "$A" -F file=@/tmp/emp.csv

# reports (admin/manager)
curl -s "$B/v1/reports/overtime?from=2030-01-01&to=2030-01-31" -H "$A"
curl -s "$B/v1/reports/coverage?date=2030-01-05" -H "$A"
```

## RBAC

| Aksi | admin | manager | staff |
|---|---|---|---|
| auth register/login, `/me` | ✅ | ✅ | ✅ |
| employees CRUD | ✅ tulis | 👁 baca | ❌ |
| employees import CSV | ✅ | ✅ | ❌ |
| shifts CRUD | ✅ | ✅ | 👁 miliknya |
| attendance check-in/out, riwayat | ✅ | ✅ | 👁 miliknya |
| reports overtime/coverage | ✅ | ✅ | ❌ |

Overtime = jam kerja > 8 jam/hari; coverage = headcount/ tanggal.
Zona bisnis Asia/Jakarta (UTC+7, fixed — aman tanpa tzdata di distroless).

## Coba via Swagger / Postman

```bash
# Swagger UI ( Tanpa tambah kode):
docker run --rm -p 8081:8080 -e SWAGGER_JSON=/swagger.yaml -v $PWD/api:/swagger swaggerapi/swagger-ui
# → http://localhost:8081, Authorize dengan token dari /v1/auth/login

# Newman (butuh API + DB jalan):
npx newman run api/postman_collection.json --env-var baseUrl=http://localhost:8080
```

## Dev & CI

```bash
go build ./... && go vet ./... && gofmt -l .
go test ./... -cover            # service ≥ 60% (saat ini ~80%)
```

CI (`.github/workflows/ci.yml`): MySQL service → vet → golangci-lint →
test+coverage gate → goose migrate → seed → API boot → **Newman** →
swagger validate.

Struktur: `cmd/api/main.go`, `internal/{config,domain,repository,service,handler,middleware}`,
`migrations/` (goose), `seed/seed.sql`, `api/{swagger.yaml,postman_collection.json}`.

Keputusan arsitektur: [ADR-001 tanpa ORM](docs/ADR-001-no-orm.md),
[ADR-002 goose](docs/ADR-002-goose.md).
