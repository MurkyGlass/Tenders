CREATE EXTENSION IF NOT EXISTS pg_cron CASCADE;
SELECT cron.schedule(
'Del_refresh',
'* * * * *',
$$DELETE from refresh_tokens WHERE expires_at < NOW()$$
);
SELECT cron.schedule(
'Update_tender_status',
'* * * * *',
$$
UPDATE Tenders SET id_status = 6 WHERE datetime_end < NOW() AND datetime_end + (5 * INTERVAL '1 day') > NOW() AND id_status = 2$$
);
SELECT cron.schedule(
'Revoke_tender_status',
'* * * * *',
$$UPDATE Tenders SET id_status = 4 WHERE datetime_end + (5 * INTERVAL '1 day') < NOW() AND id_status = 6$$
);
CREATE OR REPLACE FUNCTION update_offers_on_tender_complete()
RETURNS TRIGGER AS $$
BEGIN
    -- Проверяем, что статус тендера изменился на 3 (завершен)
    IF NEW.id_status = 3 AND (OLD.id_status IS DISTINCT FROM NEW.id_status) THEN
        -- Обновляем все остальные офферы на статус 7
        UPDATE Offers 
        SET id_status = 7 
        WHERE id_tender = NEW.id_tender
          AND id_status != 5;  -- все, кто не победитель          
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trigger_complete_tender_offers
    AFTER UPDATE OF id_status ON Tenders
    FOR EACH ROW
    EXECUTE FUNCTION update_offers_on_tender_complete();
CREATE OR REPLACE FUNCTION update_offers_on_tender_del()
RETURNS TRIGGER AS $$
BEGIN
    -- Проверяем, что статус тендера изменился на 3 (завершен)
    IF NEW.id_status = 4 AND (OLD.id_status IS DISTINCT FROM NEW.id_status) THEN
        -- Обновляем все остальные офферы на статус 7
        UPDATE Offers 
        SET id_status = 4 
        WHERE id_tender = NEW.id_tender;      
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trigger_del_tender_offers
    AFTER UPDATE OF id_status ON Tenders
    FOR EACH ROW
    EXECUTE FUNCTION update_offers_on_tender_del();