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