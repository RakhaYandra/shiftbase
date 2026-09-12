# ADR-002: Migrasi dengan goose (SQL polos)

Status: diterima.

Konteks: butuh migrasi berversi untuk MySQL yang bisa jalan di CI, compose,
dan laptop tanpa toolchain Go.

Keputusan: file `migrations/NNNNN_nama.sql` format goose (`-- +goose Up/Down`),
dijalankan via CLI `goose` (bukan library — hindari dependency tree raksasa
driver yang tak dipakai; lihat go.mod).

Alasan:
- SQL polos: CHECK/UNIQUE/FK terlihat langsung, mudah diaudit.
- CLI single-binary: `go install .../cmd/goose@<versi>` lalu
  `goose -dir migrations mysql "$DB_DSN" up`.
- Down migration wajib tiap file agar rollback aman.

Konsekuensi: tim harus install CLI goose sekali; versi di-pin di README/CI.
