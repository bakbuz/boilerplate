-- Seed Brands
INSERT INTO catalog.brands (id, name, slug, created_by) VALUES
(1, 'Nike', 'nike', '00000000-0000-0000-0000-000000000000'),
(2, 'Adidas', 'adidas', '00000000-0000-0000-0000-000000000000'),
(3, 'Apple', 'apple', '00000000-0000-0000-0000-000000000000'),
(4, 'Samsung', 'samsung', '00000000-0000-0000-0000-000000000000'),
(5, 'Sony', 'sony', '00000000-0000-0000-0000-000000000000'),
(6, 'Microsoft', 'microsoft', '00000000-0000-0000-0000-000000000000'),
(7, 'LG', 'lg', '00000000-0000-0000-0000-000000000000'),
(8, 'Dell', 'dell', '00000000-0000-0000-0000-000000000000'),
(9, 'HP', 'hp', '00000000-0000-0000-0000-000000000000'),
(10, 'Lenovo', 'lenovo', '00000000-0000-0000-0000-000000000000'),
(11, 'Asus', 'asus', '00000000-0000-0000-0000-000000000000'),
(12, 'Acer', 'acer', '00000000-0000-0000-0000-000000000000'),
(13, 'Canon', 'canon', '00000000-0000-0000-0000-000000000000'),
(14, 'Nikon', 'nikon', '00000000-0000-0000-0000-000000000000'),
(15, 'Nintendo', 'nintendo', '00000000-0000-0000-0000-000000000000'),
(16, 'PlayStation', 'playstation', '00000000-0000-0000-0000-000000000000'),
(17, 'Xbox', 'xbox', '00000000-0000-0000-0000-000000000000'),
(18, 'Coca-Cola', 'coca-cola', '00000000-0000-0000-0000-000000000000'),
(19, 'Pepsi', 'pepsi', '00000000-0000-0000-0000-000000000000'),
(20, 'Nestle', 'nestle', '00000000-0000-0000-0000-000000000000'),
(21, 'Procter & Gamble', 'procter-gamble', '00000000-0000-0000-0000-000000000000'),
(22, 'Unilever', 'unilever', '00000000-0000-0000-0000-000000000000'),
(23, 'L''Oreal', 'loreal', '00000000-0000-0000-0000-000000000000'),
(24, 'Toyota', 'toyota', '00000000-0000-0000-0000-000000000000'),
(25, 'Honda', 'honda', '00000000-0000-0000-0000-000000000000'),
(26, 'Ford', 'ford', '00000000-0000-0000-0000-000000000000'),
(27, 'BMW', 'bmw', '00000000-0000-0000-0000-000000000000'),
(28, 'Mercedes-Benz', 'mercedes-benz', '00000000-0000-0000-0000-000000000000'),
(29, 'Audi', 'audi', '00000000-0000-0000-0000-000000000000'),
(30, 'Volkswagen', 'volkswagen', '00000000-0000-0000-0000-000000000000'),
(31, 'Tesla', 'tesla', '00000000-0000-0000-0000-000000000000'),
(32, 'Amazon', 'amazon', '00000000-0000-0000-0000-000000000000'),
(33, 'Google', 'google', '00000000-0000-0000-0000-000000000000'),
(34, 'Facebook', 'facebook', '00000000-0000-0000-0000-000000000000'),
(35, 'Disney', 'disney', '00000000-0000-0000-0000-000000000000'),
(36, 'Netflix', 'netflix', '00000000-0000-0000-0000-000000000000'),
(37, 'Starbucks', 'starbucks', '00000000-0000-0000-0000-000000000000'),
(38, 'McDonald''s', 'mcdonalds', '00000000-0000-0000-0000-000000000000'),
(39, 'KFC', 'kfc', '00000000-0000-0000-0000-000000000000'),
(40, 'Burger King', 'burger-king', '00000000-0000-0000-0000-000000000000'),
(41, 'Pizza Hut', 'pizza-hut', '00000000-0000-0000-0000-000000000000'),
(42, 'Domino''s', 'dominos', '00000000-0000-0000-0000-000000000000'),
(43, 'Subway', 'subway', '00000000-0000-0000-0000-000000000000'),
(44, 'Zara', 'zara', '00000000-0000-0000-0000-000000000000'),
(45, 'H&M', 'hm', '00000000-0000-0000-0000-000000000000'),
(46, 'Gucci', 'gucci', '00000000-0000-0000-0000-000000000000'),
(47, 'Prada', 'prada', '00000000-0000-0000-0000-000000000000'),
(48, 'Louis Vuitton', 'louis-vuitton', '00000000-0000-0000-0000-000000000000'),
(49, 'Hermes', 'hermes', '00000000-0000-0000-0000-000000000000'),
(50, 'Chanel', 'chanel', '00000000-0000-0000-0000-000000000000');

-- Restart ID sequence
ALTER SEQUENCE catalog.brands_id_seq RESTART WITH 51;

-- Seed Products
INSERT INTO catalog.products (brand_id, name, sku, summary, stock_quantity, price, created_by) VALUES
-- Nike (1)
(1, 'Air Force 1', 'NIKE-AF1-001', 'Classic basketball shoe', 1000, 100.00, '00000000-0000-0000-0000-000000000000'),
(2, 'Air Max 90', 'NIKE-AM90-001', 'Comfort and style combined', 850, 120.00, '00000000-0000-0000-0000-000000000000'),
-- Adidas (2)
(3, 'Ultraboost', 'ADI-UB-001', 'High performance running shoe', 500, 180.00, '00000000-0000-0000-0000-000000000000'),
(4, 'Stan Smith', 'ADI-SS-001', 'Timeless tennis shoe', 1200, 85.00, '00000000-0000-0000-0000-000000000000'),
-- Apple (3)
(5, 'iPhone 15 Pro', 'APPLE-IP15P-001', 'The ultimate iPhone', 300, 999.00, '00000000-0000-0000-0000-000000000000'),
(6, 'MacBook Air M2', 'APPLE-MBA-M2', 'Supercharged by M2', 250, 1199.00, '00000000-0000-0000-0000-000000000000'),
-- Samsung (4)
(7, 'Galaxy S24 Ultra', 'SAM-S24U-001', 'Galaxy AI is here', 400, 1299.00, '00000000-0000-0000-0000-000000000000'),
(8, 'Galaxy Watch 6', 'SAM-GW6-001', 'Advanced sleep coaching', 600, 299.00, '00000000-0000-0000-0000-000000000000'),
-- Sony (5)
(9, 'WH-1000XM5', 'SONY-XM5-001', 'Noise canceling headphones', 450, 399.00, '00000000-0000-0000-0000-000000000000'),
(10, 'PlayStation 5', 'SONY-PS5-001', 'Play Has No Limits', 200, 499.00, '00000000-0000-0000-0000-000000000000'),
-- Microsoft (6)
(11, 'Surface Pro 9', 'MS-SP9-001', 'Tablet flexibility, laptop performance', 150, 999.00, '00000000-0000-0000-0000-000000000000'),
(12, 'Xbox Series X', 'MS-XSX-001', 'Power your dreams', 300, 499.00, '00000000-0000-0000-0000-000000000000'),
-- LG (7)
(13, 'OLED C3 TV', 'LG-OLED-C3', 'Advanced OLED picture quality', 100, 1499.00, '00000000-0000-0000-0000-000000000000'),
(14, 'UltraGear Monitor', 'LG-UG-MON', 'Gaming monitor 144Hz', 200, 399.00, '00000000-0000-0000-0000-000000000000'),
-- Dell (8)
(15, 'XPS 13', 'DELL-XPS13-001', 'Premium laptop', 200, 1099.00, '00000000-0000-0000-0000-000000000000'),
(16, 'Alienware m16', 'DELL-AW-M16', 'High-performance gaming laptop', 100, 1899.00, '00000000-0000-0000-0000-000000000000'),
-- HP (9)
(17, 'Spectre x360', 'HP-SPEC-360', 'Convertible laptop', 180, 1299.00, '00000000-0000-0000-0000-000000000000'),
(18, 'LaserJet Pro', 'HP-LJP-001', 'Efficient office printer', 300, 249.00, '00000000-0000-0000-0000-000000000000'),
-- Lenovo (10)
(19, 'ThinkPad X1 Carbon', 'LEN-TP-X1', 'Business ultrabook', 220, 1499.00, '00000000-0000-0000-0000-000000000000'),
(20, 'Legion Pro 7i', 'LEN-LEG-7I', 'Gaming powerhouse', 120, 2199.00, '00000000-0000-0000-0000-000000000000'),
-- Asus (11)
(21, 'ROG Zephyrus G14', 'ASUS-ROG-G14', 'Compact gaming laptop', 150, 1499.00, '00000000-0000-0000-0000-000000000000'),
(22, 'Zenbook S 13', 'ASUS-ZEN-S13', 'Ultra-thin OLED laptop', 180, 1099.00, '00000000-0000-0000-0000-000000000000'),
-- Acer (12)
(23, 'Predator Helios', 'ACER-PRED-HEL', 'Gaming laptop', 140, 1399.00, '00000000-0000-0000-0000-000000000000'),
(24, 'Swift Go 14', 'ACER-SWIFT-14', 'Lightweight laptop', 200, 799.00, '00000000-0000-0000-0000-000000000000'),
-- Canon (13)
(25, 'EOS R6 Mark II', 'CANON-EOS-R6', 'Full-frame mirrorless camera', 80, 2499.00, '00000000-0000-0000-0000-000000000000'),
(26, 'PowerShot G7 X', 'CANON-PS-G7X', 'Premium compact camera', 120, 749.00, '00000000-0000-0000-0000-000000000000'),
-- Nikon (14)
(27, 'Z8', 'NIKON-Z8', 'Professional mirrorless camera', 50, 3999.00, '00000000-0000-0000-0000-000000000000'),
(28, 'D850', 'NIKON-D850', 'High-res DSLR', 70, 2999.00, '00000000-0000-0000-0000-000000000000'),
-- Nintendo (15)
(29, 'Switch OLED', 'NIN-SW-OLED', 'Handheld console with OLED screen', 800, 349.00, '00000000-0000-0000-0000-000000000000'),
(30, 'Pro Controller', 'NIN-PRO-CON', 'Premium controller for Switch', 500, 69.00, '00000000-0000-0000-0000-000000000000'),
-- PlayStation (16)
(31, 'DualSense Edge', 'PS-DS-EDGE', 'Pro wireless controller', 300, 199.00, '00000000-0000-0000-0000-000000000000'),
(32, 'Pulse 3D Headset', 'PS-Pulse-3D', '3D Audio headset', 400, 99.00, '00000000-0000-0000-0000-000000000000'),
-- Xbox (17)
(33, 'Elite Controller 2', 'XB-ELITE-2', 'Adjustable tension thumbsticks', 250, 179.00, '00000000-0000-0000-0000-000000000000'),
(34, 'Game Pass Ultimate', 'XB-GP-ULT', '3 Month subscription', 1000, 44.99, '00000000-0000-0000-0000-000000000000'),
-- Coca-Cola (18)
(35, 'Coke Original', 'COKE-ORG-12PK', 'Original taste 12 pack', 5000, 8.99, '00000000-0000-0000-0000-000000000000'),
(36, 'Coke Zero Sugar', 'COKE-ZERO-12PK', 'Zero sugar 12 pack', 4000, 8.99, '00000000-0000-0000-0000-000000000000'),
-- Pepsi (19)
(37, 'Pepsi Cola', 'PEPSI-REG-12PK', 'Refreshing cola 12 pack', 4500, 7.99, '00000000-0000-0000-0000-000000000000'),
(38, 'Mountain Dew', 'MTN-DEW-12PK', 'Citrus soda 12 pack', 3000, 7.99, '00000000-0000-0000-0000-000000000000'),
-- Nestle (20)
(39, 'KitKat', 'NESTLE-KITKAT', 'Wafer bar', 10000, 1.50, '00000000-0000-0000-0000-000000000000'),
(40, 'Nescafe Gold', 'NESTLE-NESCAFE', 'Instant coffee jar', 2000, 9.99, '00000000-0000-0000-0000-000000000000'),
-- P&G (21)
(41, 'Tide Pods', 'PG-TIDE-PODS', 'Laundry detergent pacs', 1500, 19.99, '00000000-0000-0000-0000-000000000000'),
(42, 'Gillette Fusion5', 'PG-GILL-F5', 'Men''s razor', 2000, 12.99, '00000000-0000-0000-0000-000000000000'),
-- Unilever (22)
(43, 'Dove Soap', 'UNI-DOVE-SOAP', 'Beauty bar 4 pack', 3000, 6.99, '00000000-0000-0000-0000-000000000000'),
(44, 'Hellmann''s Mayo', 'UNI-HELL-MAYO', 'Real mayonnaise', 1500, 5.99, '00000000-0000-0000-0000-000000000000'),
-- L'Oreal (23)
(45, 'Revitalift Serum', 'LOR-REV-SERUM', 'Anti-aging serum', 1000, 25.99, '00000000-0000-0000-0000-000000000000'),
(46, 'Voluminous Mascara', 'LOR-VOL-MASC', 'Volume building mascara', 2000, 9.99, '00000000-0000-0000-0000-000000000000'),
-- Toyota (24)
(47, 'Camry Floor Mats', 'TOY-CAM-MATS', 'All-weather floor liners', 200, 150.00, '00000000-0000-0000-0000-000000000000'),
(48, 'TRD Oil Filter', 'TOY-TRD-FILT', 'High performance oil filter', 500, 15.00, '00000000-0000-0000-0000-000000000000'),
-- Honda (25)
(49, 'Civic Type R Wing', 'HON-CTR-WING', 'Rear spoiler', 50, 450.00, '00000000-0000-0000-0000-000000000000'),
(50, 'H-Badge Emblem', 'HON-BADGE-RED', 'Red JDM emblem', 1000, 45.00, '00000000-0000-0000-0000-000000000000'),
-- Ford (26)
(51, 'Mustang Keychain', 'FORD-MUST-KEY', 'Pony logo keychain', 2000, 12.00, '00000000-0000-0000-0000-000000000000'),
(52, 'F-150 Bed Cover', 'FORD-F150-CVR', 'Tonneau cover', 100, 899.00, '00000000-0000-0000-0000-000000000000'),
-- BMW (27)
(53, 'M Performance Exhaust', 'BMW-MP-EXH', 'Sport exhaust system', 30, 2500.00, '00000000-0000-0000-0000-000000000000'),
(54, 'BMW Cap', 'BMW-CAP-BLK', 'Logo baseball cap', 800, 35.00, '00000000-0000-0000-0000-000000000000'),
-- Mercedes (28)
(55, 'AMG Wheels', 'MB-AMG-19', '19-inch alloy wheels', 40, 3000.00, '00000000-0000-0000-0000-000000000000'),
(56, 'Star Emblem', 'MB-HOOD-STAR', 'Hood ornament', 300, 50.00, '00000000-0000-0000-0000-000000000000'),
-- Audi (29)
(57, 'Quattro Decal', 'AUDI-QUAT-DEC', 'Side door decal', 500, 25.00, '00000000-0000-0000-0000-000000000000'),
(58, 'Sport Mug', 'AUDI-SPORT-MUG', 'Carbon fiber look mug', 400, 20.00, '00000000-0000-0000-0000-000000000000'),
-- VW (30)
(59, 'GTI Floor Mats', 'VW-GTI-MATS', 'Red stitched mats', 300, 120.00, '00000000-0000-0000-0000-000000000000'),
(60, 'Vintage Bus Toy', 'VW-BUS-TOY', 'Diecast T1 bus', 1000, 30.00, '00000000-0000-0000-0000-000000000000'),
-- Tesla (31)
(61, 'Wall Connector', 'TESLA-WC-GEN3', 'Home charging station', 500, 475.00, '00000000-0000-0000-0000-000000000000'),
(62, 'Model Y Floor Liners', 'TESLA-MY-LINERS', 'All-weather liners', 600, 225.00, '00000000-0000-0000-0000-000000000000'),
-- Amazon (32)
(63, 'Echo Dot', 'AMZN-ECHO-DOT', 'Smart speaker with Alexa', 5000, 49.99, '00000000-0000-0000-0000-000000000000'),
(64, 'Kindle Paperwhite', 'AMZN-KINDLE-PW', 'E-reader', 2000, 139.99, '00000000-0000-0000-0000-000000000000'),
-- Google (33)
(65, 'Pixel 8', 'GOOG-PIXEL-8', 'Android phone', 800, 699.00, '00000000-0000-0000-0000-000000000000'),
(66, 'Nest Hub', 'GOOG-NEST-HUB', 'Smart display', 1000, 99.99, '00000000-0000-0000-0000-000000000000'),
-- Facebook (Meta) (34)
(67, 'Quest 3', 'META-QUEST-3', 'VR Headset', 600, 499.00, '00000000-0000-0000-0000-000000000000'),
(68, 'Ray-Ban Meta', 'META-RAY-BAN', 'Smart glasses', 400, 299.00, '00000000-0000-0000-0000-000000000000'),
-- Disney (35)
(69, 'Mickey Ears', 'DIS-MIC-EARS', 'Classic headband', 5000, 29.99, '00000000-0000-0000-0000-000000000000'),
(70, 'Elsa Doll', 'DIS-ELSA-DOLL', 'Frozen character doll', 3000, 19.99, '00000000-0000-0000-0000-000000000000'),
-- Netflix (36) - Selling merch?
(71, 'Stranger Things Tee', 'NET-ST-TEE', 'Hellfire Club t-shirt', 2000, 24.99, '00000000-0000-0000-0000-000000000000'),
(72, 'Squid Game Mask', 'NET-SG-MASK', 'Front Man mask', 1000, 15.99, '00000000-0000-0000-0000-000000000000'),
-- Starbucks (37)
(73, 'Pike Place Roast', 'SB-PIKE-1LB', 'Whole bean coffee', 3000, 13.99, '00000000-0000-0000-0000-000000000000'),
(74, 'Cold Cup Tumbler', 'SB-TUMB-24', 'Reusable straw cup', 2000, 19.99, '00000000-0000-0000-0000-000000000000'),
-- McDonald's (38)
(75, 'Big Mac Sauce', 'MCD-MAC-SAUCE', 'Signature sauce bottle', 500, 9.99, '00000000-0000-0000-0000-000000000000'),
(76, 'Golden Arches Tee', 'MCD-TEE-YEL', 'Yellow graphic tee', 1000, 14.99, '00000000-0000-0000-0000-000000000000'),
-- KFC (39)
(77, '11 Herbs Spice Mix', 'KFC-SPICE-MIX', 'Seasoning blend', 200, 5.99, '00000000-0000-0000-0000-000000000000'),
(78, 'Bucket Hat', 'KFC-BUCKET-HAT', 'Logo hat', 500, 19.99, '00000000-0000-0000-0000-000000000000'),
-- Burger King (40)
(79, 'Crown', 'BK-CROWN-CARD', 'Cardboard crown', 10000, 0.99, '00000000-0000-0000-0000-000000000000'),
(80, 'Whopper Tee', 'BK-WHOP-TEE', 'Flame grilled tee', 500, 19.99, '00000000-0000-0000-0000-000000000000'),
-- Pizza Hut (41)
(81, 'Red Roof Hat', 'PH-ROOF-HAT', 'Vintage style hat', 300, 15.99, '00000000-0000-0000-0000-000000000000'),
(82, 'Pizza Cutter', 'PH-CUTTER', 'Branded pizza cutter', 1000, 8.99, '00000000-0000-0000-0000-000000000000'),
-- Domino's (42)
(83, 'Heatwave Bag', 'DOM-HEAT-BAG', 'Insulated delivery bag', 100, 29.99, '00000000-0000-0000-0000-000000000000'),
(84, 'Domino Set', 'DOM-GAME-SET', 'Double six dominoes', 500, 14.99, '00000000-0000-0000-0000-000000000000'),
-- Subway (43)
(85, 'Cookie Pack', 'SUB-COOKIE-12', 'Fresh baked cookies', 5000, 5.99, '00000000-0000-0000-0000-000000000000'),
(86, 'Sandwich Artist Apron', 'SUB-APRON-GRN', 'Green apron', 200, 9.99, '00000000-0000-0000-0000-000000000000'),
-- Zara (44)
(87, 'Basic Tee', 'ZARA-TEE-WHT', 'Cotton t-shirt', 5000, 9.90, '00000000-0000-0000-0000-000000000000'),
(88, 'Slim Fit Jeans', 'ZARA-JEANS-BLU', 'Denim pant', 2000, 49.90, '00000000-0000-0000-0000-000000000000'),
-- H&M (45)
(89, 'Hoodie', 'HM-HOODIE-BLK', 'Cotton blend hoodie', 4000, 24.99, '00000000-0000-0000-0000-000000000000'),
(90, 'Chino Pants', 'HM-CHINO-BEIGE', 'Slim fit chinos', 3000, 29.99, '00000000-0000-0000-0000-000000000000'),
-- Gucci (46)
(91, 'GG Belt', 'GUCCI-BELT-BLK', 'Leather belt with double G', 200, 450.00, '00000000-0000-0000-0000-000000000000'),
(92, 'Ace Sneaker', 'GUCCI-ACE-WHT', 'Embroidered leather sneaker', 150, 790.00, '00000000-0000-0000-0000-000000000000'),
-- Prada (47)
(93, 'Galleria Bag', 'PRADA-GAL-BLK', 'Saffiano leather bag', 50, 2900.00, '00000000-0000-0000-0000-000000000000'),
(94, 'Bucket Hat', 'PRADA-NYL-HAT', 'Re-Nylon hat', 200, 595.00, '00000000-0000-0000-0000-000000000000'),
-- Louis Vuitton (48)
(95, 'Neverfull MM', 'LV-NEVER-MONO', 'Monogram canvas tote', 100, 2030.00, '00000000-0000-0000-0000-000000000000'),
(96, 'Keepall 55', 'LV-KEEP-55', 'Travel bag', 80, 2350.00, '00000000-0000-0000-0000-000000000000'),
-- Hermes (49)
(97, 'Birkin 30', 'HER-BIRK-30', 'Black togo leather', 10, 12000.00, '00000000-0000-0000-0000-000000000000'),
(98, 'Silk Scarf', 'HER-SCARF-90', '90cm silk twill', 300, 495.00, '00000000-0000-0000-0000-000000000000'),
-- Chanel (50)
(99, 'No. 5 Perfume', 'CHAN-NO5-100', 'Eau de Parfum 100ml', 1000, 160.00, '00000000-0000-0000-0000-000000000000'),
(100, 'Classic Flap Bag', 'CHAN-FLAP-BLK', 'Quilted leather bag', 20, 10000.00, '00000000-0000-0000-0000-000000000000');
