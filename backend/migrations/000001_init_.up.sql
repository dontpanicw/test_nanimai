-- Сервисы (внешние системы, которые могут работать с балансом)
CREATE TABLE IF NOT EXISTS services(
id         BIGSERIAL PRIMARY KEY,
name       TEXT UNIQUE NOT NULL,
api_key    TEXT NOT NULL UNIQUE, -- UUID
);

-- Аккаунты пользователей
CREATE TABLE IF NOT EXISTS users (
id              BIGSERIAL PRIMARY KEY,
current_amount  NUMERIC(20,2) NOT NULL DEFAULT 0,
reserved_amount NUMERIC(20,2) NOT NULL DEFAULT 0, --сумма, которая уже зарезервирована в активных транзакциях, но ещё не списана
max_amount      NUMERIC(20,2) NOT NULL DEFAULT 0,
CHECK (current_amount >= 0),
CHECK (reserved_amount >= 0),
CHECK (max_amount >= 0),
CHECK (reserved_amount <= current_amount),
CHECK (current_amount <= max_amount)
);

-- Транзакции (резервы средств)
CREATE TYPE IF NOT EXISTS reservation_status AS ENUM ('ACTIVE', 'CONFIRMED', 'CANCELLED', 'EXPIRED');

CREATE TABLE IF NOT EXISTS reservations (
id                BIGSERIAL PRIMARY KEY,
account_id        BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
owner_service_id  BIGINT NOT NULL REFERENCES services(id) ON DELETE CASCADE,
amount            NUMERIC(20,2) NOT NULL CHECK (amount > 0),
status            reservation_status NOT NULL DEFAULT 'ACTIVE',
idempotency_key   TEXT NOT NULL,
expires_at        TIMESTAMPTZ NOT NULL,
created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
confirmed_at      TIMESTAMPTZ,
cancelled_at      TIMESTAMPTZ,
UNIQUE (owner_service_id, idempotency_key) -- идемпотентность по сервису
);

-- Журнал операций (для аудита)
CREATE TYPE IF NOT EXISTS ledger_op AS ENUM (
    'LIMIT_INCREASE', 'LIMIT_DECREASE',
    'BALANCE_INCREASE', 'BALANCE_DECREASE',
    'RESERVE_OPEN', 'RESERVE_CONFIRM', 'RESERVE_CANCEL', 'RESERVE_EXPIRE'
);

CREATE TABLE IF NOT EXISTS ledger (
id              BIGSERIAL PRIMARY KEY,
account_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
reservation_id  BIGINT REFERENCES reservations(id) ON DELETE SET NULL,
actor_service_id BIGINT REFERENCES services(id), -- кто сделал операцию
operation       ledger_op NOT NULL,
delta_current   NUMERIC(20,2) NOT NULL DEFAULT 0,
delta_reserved  NUMERIC(20,2) NOT NULL DEFAULT 0,
delta_max       NUMERIC(20,2) NOT NULL DEFAULT 0,
created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO services (name, api_key) VALUES ('payments', '2d9a5f20-16ac-4b47-85f4-1b62b2675c8f'),
                                            ('shop',     'cd0fbe13-7541-4fa7-94c8-774a9f9a0e01');
