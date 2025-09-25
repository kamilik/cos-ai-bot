-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS user_states (
    user_id BIGINT PRIMARY KEY,
    step INTEGER DEFAULT 0,
    skin_type VARCHAR(50),
    age VARCHAR(20),
    gender VARCHAR(20),
    pregnancy VARCHAR(50),
    concerns TEXT,
    goal VARCHAR(50),
    budget VARCHAR(30),
    current_routine TEXT,
    allergies VARCHAR(100),
    preferences TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы продуктов
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    brand VARCHAR(255) NOT NULL,
    product_title VARCHAR(500) NOT NULL,
    image TEXT,
    ingredients TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(brand, product_title)
);

-- Создание таблицы продуктов пользователей
CREATE TABLE IF NOT EXISTS user_products (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    brand VARCHAR(255) NOT NULL,
    product_title VARCHAR(500) NOT NULL,
    image TEXT,
    ingredients TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, brand, product_title),
    FOREIGN KEY (user_id) REFERENCES user_states(user_id) ON DELETE CASCADE
);

-- Создание индексов для оптимизации
CREATE INDEX IF NOT EXISTS idx_user_states_user_id ON user_states(user_id);
CREATE INDEX IF NOT EXISTS idx_products_brand ON products(brand);
CREATE INDEX IF NOT EXISTS idx_user_products_user_id ON user_products(user_id);
CREATE INDEX IF NOT EXISTS idx_user_products_brand ON user_products(brand);

-- Создание функции для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Создание триггеров для автоматического обновления updated_at
CREATE TRIGGER update_user_states_updated_at 
    BEFORE UPDATE ON user_states 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at 
    BEFORE UPDATE ON products 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_products_updated_at 
    BEFORE UPDATE ON user_products 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
