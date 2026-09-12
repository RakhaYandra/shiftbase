-- Seed fiktif untuk portofolio (bukan data asli).
-- Password: Admin123! / Manager123! / Staff123! (bcrypt cost 10)
INSERT INTO users (email, password_hash, role) VALUES
('admin@shiftbase.local',   '$2a$10$zAd94cVAZNzD9MgmLXmn4ekeLuHtX7BMU.UW153nVJYzUH7.rUh.u', 'admin'),
('manager@shiftbase.local', '$2a$10$LppHuVQKOLM2TAqBRPrr8egY.Rq0HrTnD3F4PJtNXevtHqnS6N4fG', 'manager'),
('staff@shiftbase.local',   '$2a$10$KCVxjVlcGmGM1fifTK0uZeEiLier4SB2LVPvjWZnBdeH2ZNxbMdF6', 'staff');

INSERT INTO employees (name, email, phone, position, hire_date) VALUES
('Ayu Lestari',  'ayu.lestari@example.com',  '0812000001', 'kasir',   '2024-01-15'),
('Budi Santoso', 'budi.santoso@example.com', '0812000002', 'barista', '2024-03-01'),
('Citra Dewi',   'citra.dewi@example.com',   '0812000003', 'kasir',   '2024-06-10'),
('Dedi Kurnia',  'dedi.kurnia@example.com',  '0812000004', 'koki',    '2025-01-05'),
('Eka Putri',    'eka.putri@example.com',    '0812000005', 'barista', '2025-07-20');
