CREATE EXTENSION IF NOT EXISTS pg_cron CASCADE;
SELECT cron.schedule(
'Del_refresh',
'* * * * *',
$$DELETE from refresh_tokens WHERE expires_at < NOW()$$
);
INSERT INTO Entities(name)
VALUES
('Companies'),
('Tenders'),
('Offers'),
('Docs'),
('Users')
ON CONFLICT (name) DO NOTHING;
INSERT INTO Type_Operation(name)
VALUES
('Create'),
('Update'),
('Delete'),
('Login'),
('Revoke'),
('Refresh'),
('Authenfication')
ON CONFLICT (name) DO NOTHING;
-- TODO Insert Rights,Districts,Categories,Statuses
INSERT INTO Role_in_Company(id_role,name, is_creater)
VALUES
(1,'Директор',TRUE)
ON CONFLICT (id_role) DO NOTHING;
INSERT INTO Role(name)
VALUES
('User'),
('Admin'),
('Moderator')
ON CONFLICT (name) DO NOTHING;
INSERT INTO Districts (name) VALUES
('Центральный федеральный округ'),
('Северо-Западный федеральный округ'),
('Южный федеральный округ'),
('Северо-Кавказский федеральный округ'),
('Приволжский федеральный округ'),
('Уральский федеральный округ'),
('Сибирский федеральный округ'),
('Дальневосточный федеральный округ')
ON CONFLICT (name) DO NOTHING;
INSERT INTO Statuses (name) VALUES
('Черновик'),
('Активный'),
('Завершен'),
('Отменен'),
('Приостановлен'),
('На рассмотрении'),
('Отказ')
ON CONFLICT (name) DO NOTHING;
--Testing data
INSERT INTO Companies (
    name, 
    email, 
    address, 
    inn, 
    egrul, 
    description
) VALUES (
    'ООО "Медицинские технологии"',
    'info@medtech.ru',
    'г. Москва, ул. Ленина, д. 10, офис 305',
    '7701234567',
    '1234567890123',
    'Компания специализируется на поставках медицинского оборудования и расходных материалов для государственных и частных медицинских учреждений. На рынке более 10 лет.'
) ON CONFLICT (name) DO NOTHING;
INSERT INTO Tenders (name, description, datetime_start, datetime_end, id_company, id_status, id_district)
VALUES 
-- Тендер 1: Поставка медицинского оборудования
(
    'Поставка медицинского оборудования',
    'Поставка медицинского оборудования для городской больницы №1 в соответствии с техническим заданием. В комплект входит: аппараты УЗИ (2 шт.), рентген-аппараты (1 шт.), оборудование для реанимации.',
    '2026-03-25 10:00:00',
    '2026-04-25 18:00:00',
    1, -- id_company (предполагаем, что компания с id=1 существует)
    2, -- id_status (1 = Активный)
    1  -- id_district (1 = Москва)
),

-- Тендер 2: Строительство спортивного комплекса
(
    'Строительство спортивного комплекса',
    'Строительство спортивного комплекса с бассейном и игровыми залами в районе Северное Бутово. Общая площадь 5000 кв.м. Включает: бассейн 25м, спортивные залы, раздевалки, административные помещения.',
    '2026-04-01 09:00:00',
    '2026-06-01 18:00:00',
    1, -- id_company (Администрация г. Москвы)
    1, -- id_status (1 = Активный)
    1  -- id_district (1 = Москва)
),
-- Тендер 3: Строительство спортивного комплекса
(
    'Закупка кирпичей',
    'закупка качественных кирпичей 1000шт одной поставкой',
    '2026-05-01 09:00:00',
    '2026-08-01 18:00:00',
    1, -- id_company (Администрация г. Москвы)
    2, -- id_status (1 = Активный)
    1  -- id_district (1 = Москва)
);