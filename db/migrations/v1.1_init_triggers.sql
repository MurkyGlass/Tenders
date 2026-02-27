CREATE EXTENSION IF NOT EXISTS pg_cron CASCADE;
SELECT cron.schedule(
'Del_refresh',
'* * * * *',
$$DELETE from refresh_tokens WHERE expires_at < NOW()$$
);
