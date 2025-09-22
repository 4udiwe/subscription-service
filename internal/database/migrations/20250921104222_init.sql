-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS offer (
    id UUID DEFAULT gen_random_uuid() NOT NULL,
    name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    duration_months INTEGER NOT NULL, 
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()

    PRIMARY(id)
    UNIQUE(name, price)
);

CREATE TABLE IF NOT EXISTS subscription (
    id UUID DEFAULT gen_random_uuid() NOT NULL,
    user_id UUID NOT NULL,
    offer_id UUID NOT NULL REFERENCES offer(id) ON DELETE RESTRICT,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    CHECK (end_date IS NULL OR end_date >= start_date)
);

CREATE INDEX IF NOT EXISTS idx_subscription_user_id ON subscription(user_id);
CREATE INDEX IF NOT EXISTS idx_subscription_start_date ON subscription(start_date);
CREATE INDEX IF NOT EXISTS idx_subscription_end_date ON subscription(end_date);
CREATE INDEX IF NOT EXISTS idx_subscription_offer_id ON subscription(offer_id);
CREATE INDEX IF NOT EXISTS idx_offer_name_price ON offer(name, price);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_subscription_user_id;
DROP INDEX IF EXISTS idx_subscription_start_date;
DROP INDEX IF EXISTS idx_subscription_end_date;
DROP INDEX IF EXISTS idx_offer_name_price;

DROP TABLE IF EXISTS subscription;
DROP TABLE IF EXISTS offer;
-- +goose StatementEnd
