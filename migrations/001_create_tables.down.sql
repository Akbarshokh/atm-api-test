DO $$
    BEGIN
        IF EXISTS (SELECT 1 FROM pg_type t WHERE t.typname = 'transaction_type') THEN
            DROP TYPE transaction_type;
        END IF;
    END $$;

DROP TABLE IF EXISTS cards;

DROP TABLE IF EXISTS card_providers;

