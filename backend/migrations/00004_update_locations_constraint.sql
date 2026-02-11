-- +goose Up
-- +goose StatementBegin
ALTER TABLE Locations DROP CONSTRAINT locations_name_address_key;
CREATE UNIQUE INDEX locations_name_address_unique
  ON Locations (name, address, CAST(point AS text)) NULLS NOT DISTINCT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX locations_name_address_unique;
ALTER TABLE Locations ADD CONSTRAINT locations_name_address_key UNIQUE(name, address);
-- +goose StatementEnd
