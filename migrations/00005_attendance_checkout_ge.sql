-- +goose Up
-- check_out boleh sama dengan check_in (check-out dalam detik yang sama;
-- DATETIME tanpa fraksi membuat > gagal). Aturan bisnis: check-out tak boleh
-- lebih awal dari check-in.
ALTER TABLE attendance DROP CHECK chk_attendance_out;
ALTER TABLE attendance ADD CONSTRAINT chk_attendance_out
  CHECK (check_out IS NULL OR check_in IS NULL OR check_out >= check_in);

-- +goose Down
ALTER TABLE attendance DROP CHECK chk_attendance_out;
ALTER TABLE attendance ADD CONSTRAINT chk_attendance_out
  CHECK (check_out IS NULL OR check_in IS NULL OR check_out > check_in);
