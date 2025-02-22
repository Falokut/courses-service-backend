INSERT INTO users (username, fio, password, role_id) VALUES
    ('user1', 'Иван Иванов', '$2b$12$abc123hashedpassword1', 1),
    ('user2', 'Мария Смирнова', '$2b$12$abc123hashedpassword2', 2),
    ('user3', 'Алексей Кузнецов', '$2b$12$abc123hashedpassword3', 3),
    ('user4', 'Наталья Попова', '$2b$12$abc123hashedpassword4', 2),
    ('user5', 'Дмитрий Соколов', '$2b$12$abc123hashedpassword5', 3),
    ('user6', 'Светлана Морозова', '$2b$12$abc123hashedpassword6', 2),
    ('user7', 'Егор Васильев', '$2b$12$abc123hashedpassword7', 3),
    ('user8', 'Анна Новикова', '$2b$12$abc123hashedpassword8', 2),
    ('user9', 'Сергей Михайлов', '$2b$12$abc123hashedpassword9', 3),
    ('user10', 'Елена Федорова', '$2b$12$abc123hashedpassword10', 2);

INSERT INTO courses (name) VALUES
    ('Основы программирования'),
    ('Математический анализ'),
    ('История России'),
    ('Физика для начинающих'),
    ('Английский язык'),
    ('Экономика предприятия'),
    ('Основы маркетинга'),
    ('Биология человека'),
    ('Философия науки'),
    ('Основы дизайна');

INSERT INTO elective_lessons (elective_id, title, created_at, description, video_url) VALUES
    (1, 'Введение в Python', NOW(), 'Основы синтаксиса и структуры кода', 'https://example.com/video1'),
    (2, 'Пределы и производные', NOW(), 'Основные понятия анализа', 'https://example.com/video2'),
    (3, 'Древняя Русь', NOW(), 'История от Киевской Руси до Московского княжества', 'https://example.com/video3'),
    (4, 'Основы механики', NOW(), 'Законы Ньютона и их применение', 'https://example.com/video4'),
    (5, 'Грамматика английского', NOW(), 'Разбор базовых конструкций', 'https://example.com/video5'),
    (6, 'Маркетинговые исследования', NOW(), 'Методы анализа рынка', 'https://example.com/video6'),
    (7, 'Физиология человека', NOW(), 'Структура и функции организма', 'https://example.com/video7'),
    (8, 'Современная философия', NOW(), 'Ключевые идеи и течения', 'https://example.com/video8'),
    (9, 'Визуальный дизайн', NOW(), 'Основные принципы композиции', 'https://example.com/video9'),
    (10, 'Управление проектами', NOW(), 'Методы и практики управления', 'https://example.com/video10');

INSERT INTO assignments (elective_id, teacher_id, assignment_text, assignment_at) VALUES
    (1, 3, 'Сделать калькулятор на Python', NOW()),
    (2, 5, 'Решить задачи по производным', NOW()),
    (3, 9, 'Подготовить доклад об истории Руси', NOW()),
    (4, 7, 'Смоделировать движение тела под углом', NOW()),
    (5, 9, 'Составить эссе на английском', NOW()),
    (6, 5, 'Проанализировать рынок электроники', NOW()),
    (7, 7, 'Подготовить презентацию о дыхательной системе', NOW()),
    (8, 3, 'Написать эссе о философии Платона', NOW()),
    (9, 9, 'Создать макет веб-страницы', NOW()),
    (10, 5, 'Разработать план проекта', NOW());

INSERT INTO submitted_assignments (student_id, assignment_id, submitted_at, grade) VALUES
    (2, 1, NOW(), 9),
    (4, 2, NOW(), 8),
    (6, 3, NOW(), 10),
    (8, 4, NOW(), 7),
    (10, 5, NOW(), 8),
    (2, 6, NOW(), 9),
    (4, 7, NOW(), 10),
    (6, 8, NOW(), 7),
    (8, 9, NOW(), 9),
    (10, 10, NOW(), 10);
