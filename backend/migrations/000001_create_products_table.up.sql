CREATE TYPE product_status AS ENUM (
    'Active', 
    'Inactive'
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    
    sku VARCHAR(100) NOT NULL UNIQUE,
    
    name TEXT NOT NULL,
    description TEXT NULL,
    
    price NUMERIC(10, 2) NOT NULL,
    
    stock INT NOT NULL DEFAULT 0,
    
    category VARCHAR(100) NOT NULL DEFAULT 'Uncategorized',
    
    status product_status NOT NULL DEFAULT 'Active',
    
    image_url TEXT NULL,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Memasukkan 50 data dummy ke dalam tabel 'products'

INSERT INTO products (
    sku, name, description, price, stock, category, status, image_url
) VALUES 
-- Kategori: Electronics (1-5)
('SKU-001', '4K Smart TV 65"', 'Ultra HD Smart TV with HDR and built-in streaming apps.', 12500000.00, 15, 'Electronics', 'Active', 'https://placehold.co/600x400/blue/white?text=SmartTV'),
('SKU-002', 'Wireless Noise-Cancelling Headphones', 'Over-ear headphones with 30-hour battery life and Bluetooth 5.0.', 3200000.00, 45, 'Electronics', 'Active', 'https://placehold.co/600x400/333/white?text=Headphones'),
('SKU-003', 'Gaming Mouse RGB', NULL, 750000.00, 120, 'Electronics', 'Active', NULL),
('SKU-004', 'Ultra-thin Laptop 14"', 'Lightweight laptop with 16GB RAM and 512GB SSD.', 18000000.00, 10, 'Electronics', 'Active', 'https://placehold.co/600x400/grey/white?text=Laptop'),
('SKU-005', 'Portable Power Bank 20000mAh', 'Fast charging power bank with 2 USB-C ports.', 450000.00, 200, 'Electronics', 'Active', NULL),

-- Kategori: Furniture (6-10)
('SKU-006', 'Ergonomic Office Chair', 'Mesh back office chair with lumbar support.', 2800000.00, 25, 'Furniture', 'Active', 'https://placehold.co/600x400/brown/white?text=Chair'),
('SKU-007', 'Modern Oak Coffee Table', 'Solid oak wood coffee table, minimalist design.', 1900000.00, 12, 'Furniture', 'Active', 'https://placehold.co/600x400/tan/white?text=Table'),
('SKU-008', 'King Size Bed Frame', NULL, 4500000.00, 5, 'Furniture', 'Active', NULL),
('SKU-009', '3-Seater Sofa (Grey)', 'Comfortable fabric sofa for living room.', 5200000.00, 8, 'Furniture', 'Active', 'https://placehold.co/600x400/6c757d/white?text=Sofa'),
('SKU-010', 'Bookshelf (5-tier)', 'Tall wooden bookshelf for storage.', 1300000.00, 18, 'Furniture', 'Active', NULL),

-- Kategori: Apparel (11-15)
('SKU-011', 'Men''s Denim Jacket', 'Classic blue denim jacket.', 550000.00, 75, 'Apparel', 'Active', 'https://placehold.co/600x400/0d6efd/white?text=Jacket'),
('SKU-012', 'Women''s Running Shoes', 'Lightweight and breathable.', 890000.00, 110, 'Apparel', 'Active', NULL),
('SKU-013', 'Cotton T-Shirt (Black)', NULL, 120000.00, 300, 'Apparel', 'Active', 'https://placehold.co/600x400/212529/white?text=TShirt'),
('SKU-014', 'Leather Wallet', 'Genuine leather wallet with RFID blocking.', 350000.00, 90, 'Apparel', 'Active', NULL),
('SKU-015', 'Silk Scarf (Old model)', 'Discontinued pattern.', 250000.00, 0, 'Apparel', 'Inactive', NULL),

-- Kategori: Groceries (16-20)
('SKU-016', 'Organic Arabica Coffee Beans 1kg', 'Whole bean, medium roast.', 220000.00, 150, 'Groceries', 'Active', NULL),
('SKU-017', 'Italian Olive Oil 500ml', 'Extra virgin olive oil.', 180000.00, 80, 'Groceries', 'Active', 'https://placehold.co/600x400/84a98c/white?text=Oil'),
('SKU-018', 'Almond Milk (Unsweetened)', NULL, 45000.00, 60, 'Groceries', 'Active', NULL),
('SKU-019', 'Premium Dark Chocolate 100g', '70% Cacao.', 35000.00, 250, 'Groceries', 'Active', 'https://placehold.co/600x400/583101/white?text=Choc'),
('SKU-020', 'Imported Truffle Oil (Expired)', 'Past expiration date.', 300000.00, 5, 'Groceries', 'Inactive', NULL),

-- Kategori: Books (21-25)
('SKU-021', 'The Go Programming Language', 'By Donovan and Kernighan.', 450000.00, 30, 'Books', 'Active', 'https://placehold.co/600x400/007bff/white?text=GoBook'),
('SKU-022', 'Designing Data-Intensive Applications', 'By Martin Kleppmann.', 550000.00, 22, 'Books', 'Active', NULL),
('SKU-023', 'Sapiens: A Brief History of Humankind', NULL, 210000.00, 60, 'Books', 'Active', 'https://placehold.co/600x400/fca311/white?text=Sapiens'),
('SKU-024', '1984 by George Orwell', 'Classic dystopian novel.', 130000.00, 0, 'Books', 'Active', NULL),
('SKU-025', 'Old Programming Manual (1995)', 'Outdated content.', 50000.00, 3, 'Books', 'Inactive', NULL),

-- Kategori: Tools (26-30)
('SKU-026', 'Cordless Drill Set 18V', 'Includes 2 batteries and charger.', 1400000.00, 40, 'Tools', 'Active', 'https://placehold.co/600x400/ffc107/black?text=Drill'),
('SKU-027', 'Wrench Set (24-piece)', 'Chrome vanadium steel.', 600000.00, 70, 'Tools', 'Active', NULL),
('SKU-028', 'Digital Multimeter', 'For electrical testing.', 350000.00, 90, 'Tools', 'Active', 'https://placehold.co/600x400/dc3545/white?text=Meter'),
('SKU-029', 'Hammer (Claw)', NULL, 80000.00, 150, 'Tools', 'Active', NULL),
('SKU-030', 'Hand Saw (Rusted)', 'Damaged stock.', 120000.00, 10, 'Tools', 'Inactive', NULL),

-- Kategori: Toys (31-35)
('SKU-031', 'LEGO City Space Port', '600-piece building set.', 900000.00, 35, 'Toys', 'Active', 'https://placehold.co/600x400/fd7e14/white?text=LEGO'),
('SKU-032', 'Plush Teddy Bear (Large)', 'Soft and cuddly, 1m tall.', 450000.00, 60, 'Toys', 'Active', NULL),
('SKU-033', 'Remote Control Car', '1:16 scale, 2.4GHz.', 300000.00, 0, 'Toys', 'Inactive', 'https://placehold.co/600x400/198754/white?text=RCCar'),
('SKU-034', 'Jigsaw Puzzle (1000-piece)', 'Landscape scene.', 180000.00, 110, 'Toys', 'Active', NULL),
('SKU-035', 'Action Figure (Vintage)', NULL, 750000.00, 5, 'Toys', 'Active', NULL),

-- Kategori: Sports (36-40)
('SKU-036', 'Yoga Mat (Eco-friendly)', 'Non-slip TPE material.', 250000.00, 130, 'Sports', 'Active', 'https://placehold.co/600x400/20c997/white?text=Yoga'),
('SKU-037', 'Dumbbell Set (20kg)', 'Adjustable cast iron dumbbells.', 800000.00, 50, 'Sports', 'Active', NULL),
('SKU-038', 'Basketball (Size 7)', 'Official NBA size and weight.', 350000.00, 80, 'Sports', 'Active', NULL),
('SKU-039', 'Running Treadmill (Foldable)', 'Foldable home treadmill, up to 12km/h.', 6500000.00, 7, 'Sports', 'Active', 'https://placehold.co/600x400/6f42c1/white?text=Treadmill'),
('SKU-040', 'Bicycle Helmet (Old Design)', 'Discontinued model.', 300000.00, 15, 'Sports', 'Inactive', NULL),

-- Kategori: Home & Garden (41-45)
('SKU-041', 'Robot Vacuum Cleaner', 'Smart mapping and auto-charging.', 4200000.00, 18, 'Home & Garden', 'Active', 'https://placehold.co/600x400/e83e8c/white?text=Vacuum'),
('SKU-042', 'Gardening Tool Set (3-piece)', 'Trowel, fork, and cultivator.', 150000.00, 120, 'Home & Garden', 'Active', NULL),
('SKU-043', 'Air Fryer (5.5L)', NULL, 1100000.00, 65, 'Home & Garden', 'Active', NULL),
('SKU-044', 'LED Desk Lamp', 'Adjustable brightness and color temp.', 280000.00, 90, 'Home & Garden', 'Active', 'https://placehold.co/600x400/f8f9fa/black?text=Lamp'),
('SKU-045', 'Electric Kettle 1.7L', 'Stainless steel, fast boil.', 220000.00, 0, 'Home & Garden', 'Inactive', NULL),

-- Kategori: Automotive (46-50)
('SKU-046', 'Car Tire Inflator (Portable)', '12V DC portable air compressor.', 400000.00, 55, 'Automotive', 'Active', NULL),
('SKU-047', 'Dash Cam 1080p', 'Full HD dash cam with night vision.', 750000.00, 30, 'Automotive', 'Active', 'https://placehold.co/600x400/343a40/white?text=DashCam'),
('SKU-048', 'Car Wax (Carnauba)', 'Premium carnauba wax, 200g.', 180000.00, 100, 'Automotive', 'Active', NULL),
('SKU-049', 'Wiper Blades (Set of 2)', NULL, 120000.00, 140, 'Automotive', 'Active', NULL),
('SKU-050', 'Engine Oil 5W-30 (Old Stock)', 'Old packaging, clearance.', 200000.00, 20, 'Automotive', 'Inactive', NULL);