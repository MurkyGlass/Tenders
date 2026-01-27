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
('Docs')
ON CONFLICT (name) DO NOTHING;
INSERT INTO Type_Operation(name)
VALUES
('Create'),
('Update'),
('Delete')
ON CONFLICT (name) DO NOTHING;
-- TODO Insert Rights,Roles,Districts,Categories,Statuses
-- Insert master role by Role in company