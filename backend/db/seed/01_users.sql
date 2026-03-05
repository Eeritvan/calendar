-- +goose Up
-- +goose StatementBegin
INSERT INTO Users (name, password_hash)
VALUES ('user', '$2a$12$nIlSkV//4Kukvjt1Mt0f2OU4/KC4dAAl.6OXd1TQPrk2dih57C60q'); -- pw: Testtest

INSERT INTO Users (name, password_hash)
VALUES ('tester', '$2a$12$ZnVcm.m5IfejG7GEjAw0x.205f4ZSoQfd4iRwPA8rQs0TmcWCD34y'); -- pw: Testtest
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE Users CASCADE;
-- +goose StatementEnd
