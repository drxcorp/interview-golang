-- +goose Up
-- +goose StatementBegin
INSERT INTO wallets (user_id, balance, currency) VALUES
    (1, 1000.00, 'USD');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM wallets WHERE user_id = 1;
-- +goose StatementEnd
