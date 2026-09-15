# shiftbase

Shift scheduling & attendance REST API — **Go/Gin + MySQL + JWT + Docker**.
`database/sql` tanpa ORM, migrasi goose. Data contoh fiktif (portofolio).

## Purpose, Output & Expectations

**Purpose.** Small businesses (cafes, retail, kitchens) schedule staff on
spreadsheets or chat: shifts overlap, no one tracks hours, overtime pay is
guesswork. Shiftbase replaces that with one API: roster, attendance, and
payroll inputs in a single system.

**Output.** A running API (`:8080`) with JWT auth, role-based access
(admin/manager/staff), conflict-free shift planning, check-in/out attendance,
CSV employee import, and overtime/coverage reports — plus Swagger + Postman
contracts and a green CI pipeline (Newman included).

**Expectations.** After deploying: no double-booked shifts (409 on overlap),
every workday has check-in/out records, overtime is computed (> 8 h/day,
Asia/Jakarta) instead of estimated, and headcount gaps are visible per date.

## Features

| Feature | Description |
|---|---|
| Auth | - Register (default role staff), login returning a 24 h JWT, `/me` profile. - Purpose: single sign-in for all roles. Output: bearer token used by every endpoint. |
| Employees | - Create, read, update, delete employees (name, email, phone, position, hire date). - Purpose: one master roster instead of scattered lists. Output: employee IDs referenced by shifts and attendance. |
| Shifts | - Plan shifts per employee/date with start/end times; overlapping shifts rejected (409). - Purpose: conflict-free rostering. Output: a roster where no employee works two places at once. |
| Attendance | - Check-in/out (staff locked to their own record); history filterable by date range. - Purpose: proof of presence for payroll. Output: daily hours feeding overtime math. |
| CSV Import | - Bulk import employees (`name,email,phone,position,hire_date`, max 2 MB) with per-row error report. - Purpose: onboard dozens of staff at once. Output: imported/failed counts + row-level errors. |
| Reports | - Overtime per employee and headcount coverage per date. - Purpose: payroll inputs and gap detection. Output: overtime hours and daily headcount. |

## How It Works

```mermaid
flowchart TD
    C[Client] --> A[POST /v1/auth/login]
    A --> T[JWT 24h]
    T --> R[RBAC middleware: admin/manager/staff]
    R --> M[Modules: employees, shifts, attendance, reports]
    M --> DB[(MySQL)]
    M --> O[Conflict check: overlap = 409]
    M --> OT[Overtime = GREATEST(hours-8, 0)/day, Asia/Jakarta]
```

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
# Swagger UI (tanpa tambah kode):
docker run --rm -d --name swagger-ui -p 8081:8080 \
  -v $PWD/api/swagger.yaml:/usr/share/nginx/html/swagger.yaml:ro \
  -e "URLS=[{url: '/swagger.yaml', name: 'shiftbase'}]" swaggerapi/swagger-ui
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
