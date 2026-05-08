-- Заполняем базу 

INSERT INTO users (email)
SELECT 'user_' || i || '@example.com'
FROM generate_series(1, 100) AS i;

INSERT INTO tasks (title, user_id)
SELECT
    'Task ' || t || ' for user ' || u.id,
    u.id
FROM users u
CROSS JOIN generate_series(1, 100) AS t;