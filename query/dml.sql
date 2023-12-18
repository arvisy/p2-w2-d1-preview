INSERT INTO games (title, genre, price, stock)
VALUES 
    ('Final Fantasy', 'RPG', 59.99, 100),
    ('FIFA 2023', 'Sports', 49.99, 120),
    ('Doom Eternal', 'FPS', 49.99, 80);

INSERT INTO branches (name, location)
VALUES 
    ('Downtown Branch', '123 Downtown St'),
    ('Uptown Branch', '456 Uptown Ave');

INSERT INTO sales (game_id, branch_id, sale_date, quantity)
VALUES
    (1, 1, '2023-08-26', 2),
    (2, 2, '2023-08-25', 3);