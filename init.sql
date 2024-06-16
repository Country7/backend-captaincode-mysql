CREATE USER 'mysqluser'@'%' IDENTIFIED BY 'secretpass';
GRANT ALL PRIVILEGES ON main_db.* TO 'mysqluser'@'%';
-- GRANT ALL PRIVILEGES ON *.* TO 'mysqluser'@'%';
FLUSH PRIVILEGES;
-- ALTER USER 'mysqluser'@'%' IDENTIFIED WITH caching_sha2_password BY 'secretpass';