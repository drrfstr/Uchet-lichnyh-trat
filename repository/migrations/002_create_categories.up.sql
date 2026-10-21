CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

INSERT OR IGNORE INTO categories (name) VALUES ('Еда');
INSERT OR IGNORE INTO categories (name) VALUES ('Транспорт');
INSERT OR IGNORE INTO categories (name) VALUES ('Жильё');
INSERT OR IGNORE INTO categories (name) VALUES ('Развлечения');
INSERT OR IGNORE INTO categories (name) VALUES ('Здоровье');
INSERT OR IGNORE INTO categories (name) VALUES ('Одежда');
INSERT OR IGNORE INTO categories (name) VALUES ('Другое');