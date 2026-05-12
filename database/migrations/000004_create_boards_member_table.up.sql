CREATE TABLE boards_members (
    board_internal_id BIGINT NOT NULL REFERENCES boards(internal_id) ON DELETE CASCADE,
    user_internal_id BIGINT NOT NULL REFERENCES users(internal_id) ON DELETE CASCADE,
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (board_internal_id, user_internal_id)
);

ALTER TABLE boards_members RENAME TO board_members;