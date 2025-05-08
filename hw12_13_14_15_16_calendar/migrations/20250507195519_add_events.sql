-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id          varchar not null
                constraint events_pk primary key,
    title       varchar not null,
    date_start  timestamp not null,
    date_end    timestamp,
    description text not null,
    user_id     varchar not null,
    notify_days integer
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
