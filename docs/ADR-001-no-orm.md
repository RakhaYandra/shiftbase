# ADR-001: Tanpa ORM, database/sql + query eksplisit

Status: diterima.

Konteks: shiftbase butuh skema kecil (4 tabel), query eksplisit (overlap shift,
agregasi lembur), dan citra minimal untuk portofolio backend.

Keputusan: `database/sql` + driver MySQL, SQL ditulis tangan di
`internal/repository`. Tanpa GORM/sqlc/ent.

Alasan:
- Query kompleks (overlap interval, `GREATEST` lembur) lebih jelas sebagai SQL
  daripada query-builder.
- Nol magic migrasi/skema; migrasi tetap eksplisit via goose (ADR-002).
- Biner kecil, dependensi sedikit, mudah direview perekrut.

Konsekuensi: mapping row↔struct manual (`Scan`); tambah tabel = tambah
repository. Risiko diterima untuk ukuran proyek ini.
