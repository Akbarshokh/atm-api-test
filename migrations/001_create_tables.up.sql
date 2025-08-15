CREATE EXTENSION IF NOT EXISTS pgcrypto; -- используем для gen_random_uuid()


-- Providers of card networks(исходя из задание-привожу очень примитивное решение)
CREATE TABLE IF NOT EXISTS card_providers (
                                              id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                              code        TEXT NOT NULL UNIQUE,         -- e.g., "UZCARD", "HUMO", "VISA", "MASTERCARD"
                                              name        TEXT NOT NULL,
                                              created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Сразу создаю базовые карты
INSERT INTO card_providers (code, name) VALUES
                                            ('UZCARD', 'UzCard'),
                                            ('HUMO', 'Humo'),
                                            ('VISA', 'Visa'),
                                            ('MASTERCARD', 'Mastercard')
ON CONFLICT (code) DO NOTHING;


-- Cards
CREATE TABLE IF NOT EXISTS cards (
                                     id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                     provider_id     UUID NOT NULL REFERENCES card_providers(id) ON UPDATE CASCADE,
                                     pan             TEXT NOT NULL UNIQUE,          -- буду сохранять masked_pan
                                     hashed_pin      TEXT NOT NULL,
                                     pin_attempts    SMALLINT NOT NULL DEFAULT 0 CHECK (pin_attempts >= 0 AND pin_attempts <= 10),
                                     is_blocked      BOOLEAN NOT NULL DEFAULT FALSE,
                                     balance         NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
                                     currency_code   CHAR(3) NOT NULL DEFAULT 'UZS', -- e.g. UZS, USD
                                     status          TEXT NOT NULL DEFAULT 'active', -- можно будет внедрить и другие статусы: active, closed, stolen,
                                     created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
                                     updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);


-- Helpful indexes
CREATE INDEX IF NOT EXISTS idx_cards_provider ON cards(provider_id);
CREATE INDEX IF NOT EXISTS idx_cards_blocked ON cards(is_blocked);
CREATE INDEX IF NOT EXISTS idx_cards_status ON cards(status);

-- =========================
-- Transactions (audit log)
-- =========================
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type t WHERE t.typname = 'transaction_type') THEN
            CREATE TYPE transaction_type AS ENUM ('withdrawal', 'deposit', 'balance_check', 'pin_change', 'block', 'unblock');
        END IF;
    END $$;

CREATE TABLE IF NOT EXISTS card_transactions (
                                                 id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                                 card_id         UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
                                                 type            transaction_type NOT NULL,
                                                 amount          NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (amount >= 0),
                                                 fee             NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (fee >= 0),
                                                 success         BOOLEAN NOT NULL DEFAULT TRUE,
                                                 message         TEXT,                  -- optional error/info message
                                                 metadata        JSONB NOT NULL DEFAULT '{}'::jsonb, -- e.g., ATM id, location, receipt flags, etc.
                                                 created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_card_tx_card ON card_transactions(card_id);
CREATE INDEX IF NOT EXISTS idx_card_tx_type ON card_transactions(type);
CREATE INDEX IF NOT EXISTS idx_card_tx_created_at ON card_transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_card_tx_success ON card_transactions(success);

-- =========================
-- Optional daily stats table (use in app logic if needed)
-- =========================
CREATE TABLE IF NOT EXISTS card_daily_counters (
                                                   card_id         UUID PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
                                                   day             DATE NOT NULL DEFAULT CURRENT_DATE,
                                                   withdrawn_sum   NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (withdrawn_sum >= 0),
                                                   deposit_sum     NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (deposit_sum >= 0),
                                                   updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_counters_day ON card_daily_counters(day);
