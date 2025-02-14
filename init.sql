CREATE USER IF NOT EXISTS 'root'@'%' IDENTIFIED BY '27052002';
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;

-- Đảm bảo bảng roles tồn tại trước khi chèn dữ liệu
CREATE TABLE IF NOT EXISTS roles (
                                     id INT PRIMARY KEY,
                                     name VARCHAR(50) NOT NULL
    );

-- Thêm dữ liệu nếu chưa tồn tại
INSERT INTO roles (id, name) VALUES
                                 (1, 'Admin'),
                                 (2, 'Member')
    ON DUPLICATE KEY UPDATE name = VALUES(name);
