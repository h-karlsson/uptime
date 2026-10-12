-- Demo account for local testing: demo@uptime.io / demo1234
-- The hash is bcrypt (cost 10), generated with golang.org/x/crypto/bcrypt.

-- +goose Up
INSERT INTO users (email, password_hash)
VALUES ('demo@uptime.io', '$2a$10$Nn6lG1Rsyu75/wIvNIGdnOW/JfkiRmPjyRDG5ar6/xURxG1wLjIPe')
ON CONFLICT (email) DO NOTHING;

-- +goose Down
DELETE FROM users WHERE email = 'demo@uptime.io';
