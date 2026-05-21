-- FinanceKu Seed Data

-- Default admin user (password: admin123)
INSERT INTO users (id, name, email, password, role, salary, meal_allowance)
VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'Admin',
    'admin@financeku.com',
    '$2b$10$7otIdbN6PsItOq2GMxvIaOYoBa5e3YtAvy2Szg4Nb9sAFu9H6N5zW',
    'admin',
    5000000,
    30000
) ON CONFLICT (email) DO NOTHING;

-- Default income categories
INSERT INTO categories (name, type, icon, color, is_default) VALUES
    ('Gaji', 'income', 'wallet', '#4CAF50', true),
    ('Overtime', 'income', 'clock', '#FF9800', true),
    ('Freelance', 'income', 'laptop', '#2196F3', true),
    ('Investasi', 'income', 'trending-up', '#9C27B0', true),
    ('Lainnya', 'income', 'plus-circle', '#607D8B', true)
ON CONFLICT DO NOTHING;

-- Default expense categories
INSERT INTO categories (name, type, icon, color, is_default) VALUES
    ('Makanan', 'expense', 'coffee', '#F44336', true),
    ('Transportasi', 'expense', 'truck', '#E91E63', true),
    ('Belanja', 'expense', 'shopping-bag', '#9C27B0', true),
    ('Tagihan', 'expense', 'file-text', '#673AB7', true),
    ('Hiburan', 'expense', 'film', '#3F51B5', true),
    ('Kesehatan', 'expense', 'heart', '#009688', true),
    ('Pendidikan', 'expense', 'book', '#00BCD4', true),
    ('Lainnya', 'expense', 'more-horizontal', '#795548', true)
ON CONFLICT DO NOTHING;

-- Default site settings
INSERT INTO site_settings (key, value, description) VALUES
    ('app_name', 'FinanceKu', 'Application name'),
    ('overtime_period_day_start', '4', 'Day of week overtime period starts (4=Thursday)'),
    ('overtime_period_day_end', '3', 'Day of week overtime period ends (3=Wednesday)'),
    ('overtime_payment_offset_days', '9', 'Days after period end for payment'),
    ('default_currency', 'IDR', 'Default currency'),
    ('max_login_attempts', '5', 'Maximum login attempts before lockout'),
    ('lockout_duration_minutes', '30', 'Account lockout duration in minutes')
ON CONFLICT (key) DO NOTHING;
