CREATE TABLE card_position (
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid(),
    list_internal_id BIGINT NOT NULL REFERENCES lists(internal_id) ON DELETE CASCADE,
    card_order UUID[] NOT NULL DEFAULT '{}',
    CONSTRAINT card_position_public_id_unique UNIQUE (public_id),
    CONSTRAINT card_position_list_unique UNIQUE (list_internal_id)
);