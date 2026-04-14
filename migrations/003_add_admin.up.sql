INSERT INTO users (login, password_hash, role) VALUES
    (
        'admin',
        '$2a$10$N9qo8uLOickgx2ZMRZoMy.MrYV6Jg2o5LHeD7XlJ.T7JQ4QbQ5qKq',
        'admin'
    )
    ON CONFLICT (login) DO NOTHING;