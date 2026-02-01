-- 002_seed_data.up.sql
-- Initial seed data for development

-- Insert default super admin (password: admin123)
-- Hash generated with bcrypt cost 10
INSERT INTO stores (id, name, code, address, phone, is_active)
VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'Kaori Cafe Main',
    'MAIN-01',
    'Jl. Contoh No. 123, Jakarta',
    '021-12345678',
    true
) ON CONFLICT (code) DO NOTHING;

-- Super admin user (password: admin123)
-- Use bcrypt hash: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92lL1tVqE1Y.L0KVUlGWy
INSERT INTO users (id, email, password_hash, name, role, store_id, is_active)
VALUES (
    'b0000000-0000-0000-0000-000000000001',
    'admin@kaori.pos',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92lL1tVqE1Y.L0KVUlGWy',
    'Super Admin',
    'super_admin',
    NULL,
    true
) ON CONFLICT (email) DO NOTHING;

-- Store admin (password: store123)
INSERT INTO users (id, email, password_hash, name, role, store_id, pin, is_active)
VALUES (
    'b0000000-0000-0000-0000-000000000002',
    'store@kaori.pos',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92lL1tVqE1Y.L0KVUlGWy',
    'Store Manager',
    'store_admin',
    'a0000000-0000-0000-0000-000000000001',
    '1234',
    true
) ON CONFLICT (email) DO NOTHING;

-- Cashier (password: cashier123)
INSERT INTO users (id, email, password_hash, name, role, store_id, pin, is_active)
VALUES (
    'b0000000-0000-0000-0000-000000000003',
    'cashier@kaori.pos',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92lL1tVqE1Y.L0KVUlGWy',
    'Kasir 1',
    'cashier',
    'a0000000-0000-0000-0000-000000000001',
    '0000',
    true
) ON CONFLICT (email) DO NOTHING;

-- Kitchen staff (password: kitchen123)
INSERT INTO users (id, email, password_hash, name, role, store_id, is_active)
VALUES (
    'b0000000-0000-0000-0000-000000000004',
    'kitchen@kaori.pos',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92lL1tVqE1Y.L0KVUlGWy',
    'Chef 1',
    'kitchen',
    'a0000000-0000-0000-0000-000000000001',
    true
) ON CONFLICT (email) DO NOTHING;

-- Sample tables
INSERT INTO tables (id, store_id, table_number, capacity, location, is_active)
VALUES 
    ('c0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'A1', 4, 'Indoor', true),
    ('c0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'A2', 4, 'Indoor', true),
    ('c0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 'A3', 6, 'Indoor', true),
    ('c0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001', 'B1', 2, 'Outdoor', true),
    ('c0000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000001', 'B2', 2, 'Outdoor', true)
ON CONFLICT DO NOTHING;

-- Sample categories
INSERT INTO categories (id, store_id, name, icon, sort_order, is_active)
VALUES 
    ('d0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'Coffee', 'coffee', 1, true),
    ('d0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'Non-Coffee', 'cup', 2, true),
    ('d0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 'Food', 'utensils', 3, true),
    ('d0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001', 'Snacks', 'cookie', 4, true)
ON CONFLICT DO NOTHING;

-- Sample products
INSERT INTO products (id, store_id, category_id, name, description, base_price, is_available, has_variants, has_modifiers)
VALUES 
    -- Coffee
    ('e0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'd0000000-0000-0000-0000-000000000001', 
     'Americano', 'Espresso with hot water', 25000, true, true, true),
    ('e0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'd0000000-0000-0000-0000-000000000001', 
     'Cafe Latte', 'Espresso with steamed milk', 30000, true, true, true),
    ('e0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 'd0000000-0000-0000-0000-000000000001', 
     'Cappuccino', 'Espresso with foamed milk', 30000, true, true, true),
    -- Non-Coffee
    ('e0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001', 'd0000000-0000-0000-0000-000000000002', 
     'Matcha Latte', 'Japanese green tea with milk', 32000, true, true, false),
    ('e0000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000001', 'd0000000-0000-0000-0000-000000000002', 
     'Chocolate', 'Rich chocolate drink', 28000, true, true, false),
    -- Food
    ('e0000000-0000-0000-0000-000000000006', 'a0000000-0000-0000-0000-000000000001', 'd0000000-0000-0000-0000-000000000003', 
     'Nasi Goreng', 'Indonesian fried rice', 35000, true, false, true),
    ('e0000000-0000-0000-0000-000000000007', 'a0000000-0000-0000-0000-000000000001', 'd0000000-0000-0000-0000-000000000003', 
     'Mie Goreng', 'Indonesian fried noodles', 32000, true, false, true),
    -- Snacks
    ('e0000000-0000-0000-0000-000000000008', 'a0000000-0000-0000-0000-000000000001', 'd0000000-0000-0000-0000-000000000004', 
     'French Fries', 'Crispy potato fries', 20000, true, false, false),
    ('e0000000-0000-0000-0000-000000000009', 'a0000000-0000-0000-0000-000000000001', 'd0000000-0000-0000-0000-000000000004', 
     'Croissant', 'Buttery pastry', 18000, true, false, false)
ON CONFLICT DO NOTHING;

-- Product variants (sizes)
INSERT INTO product_variants (id, product_id, name, price_adjustment, is_default, sort_order)
VALUES 
    -- Americano sizes
    ('f0000000-0000-0000-0000-000000000001', 'e0000000-0000-0000-0000-000000000001', 'Small', 0, true, 1),
    ('f0000000-0000-0000-0000-000000000002', 'e0000000-0000-0000-0000-000000000001', 'Medium', 5000, false, 2),
    ('f0000000-0000-0000-0000-000000000003', 'e0000000-0000-0000-0000-000000000001', 'Large', 10000, false, 3),
    -- Latte sizes
    ('f0000000-0000-0000-0000-000000000004', 'e0000000-0000-0000-0000-000000000002', 'Small', 0, true, 1),
    ('f0000000-0000-0000-0000-000000000005', 'e0000000-0000-0000-0000-000000000002', 'Medium', 5000, false, 2),
    ('f0000000-0000-0000-0000-000000000006', 'e0000000-0000-0000-0000-000000000002', 'Large', 10000, false, 3),
    -- Cappuccino sizes
    ('f0000000-0000-0000-0000-000000000007', 'e0000000-0000-0000-0000-000000000003', 'Small', 0, true, 1),
    ('f0000000-0000-0000-0000-000000000008', 'e0000000-0000-0000-0000-000000000003', 'Medium', 5000, false, 2),
    ('f0000000-0000-0000-0000-000000000009', 'e0000000-0000-0000-0000-000000000003', 'Large', 10000, false, 3),
    -- Matcha sizes
    ('f0000000-0000-0000-0000-000000000010', 'e0000000-0000-0000-0000-000000000004', 'Regular', 0, true, 1),
    ('f0000000-0000-0000-0000-000000000011', 'e0000000-0000-0000-0000-000000000004', 'Large', 8000, false, 2),
    -- Chocolate sizes
    ('f0000000-0000-0000-0000-000000000012', 'e0000000-0000-0000-0000-000000000005', 'Regular', 0, true, 1),
    ('f0000000-0000-0000-0000-000000000013', 'e0000000-0000-0000-0000-000000000005', 'Large', 8000, false, 2)
ON CONFLICT DO NOTHING;

-- Product modifiers (add-ons)
INSERT INTO product_modifiers (id, product_id, name, price, sort_order)
VALUES 
    -- Coffee modifiers
    ('10000000-0000-0000-0000-000000000001', 'e0000000-0000-0000-0000-000000000001', 'Extra Shot', 5000, 1),
    ('10000000-0000-0000-0000-000000000002', 'e0000000-0000-0000-0000-000000000001', 'Oat Milk', 8000, 2),
    ('10000000-0000-0000-0000-000000000003', 'e0000000-0000-0000-0000-000000000002', 'Extra Shot', 5000, 1),
    ('10000000-0000-0000-0000-000000000004', 'e0000000-0000-0000-0000-000000000002', 'Oat Milk', 8000, 2),
    ('10000000-0000-0000-0000-000000000005', 'e0000000-0000-0000-0000-000000000002', 'Vanilla Syrup', 5000, 3),
    ('10000000-0000-0000-0000-000000000006', 'e0000000-0000-0000-0000-000000000003', 'Extra Shot', 5000, 1),
    ('10000000-0000-0000-0000-000000000007', 'e0000000-0000-0000-0000-000000000003', 'Oat Milk', 8000, 2),
    -- Food modifiers
    ('10000000-0000-0000-0000-000000000008', 'e0000000-0000-0000-0000-000000000006', 'Extra Egg', 5000, 1),
    ('10000000-0000-0000-0000-000000000009', 'e0000000-0000-0000-0000-000000000006', 'Extra Spicy', 0, 2),
    ('10000000-0000-0000-0000-000000000010', 'e0000000-0000-0000-0000-000000000007', 'Extra Egg', 5000, 1),
    ('10000000-0000-0000-0000-000000000011', 'e0000000-0000-0000-0000-000000000007', 'Extra Spicy', 0, 2)
ON CONFLICT DO NOTHING;
