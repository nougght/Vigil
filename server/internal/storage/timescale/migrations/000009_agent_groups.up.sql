CREATE TABLE agent_groups (
        id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
        name TEXT NOT NULL UNIQUE,
        description TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

ALTER TABLE agents ADD COLUMN group_id UUID REFERENCES groups (id) ON DELETE SET NULL;