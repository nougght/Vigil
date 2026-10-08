CREATE TABLE agent_groups (
        id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
        name TEXT NOT NULL,
        description TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ 
    );

-- unique name if not deleted
CREATE UNIQUE INDEX agent_groups_name_lower_uniq
    ON agent_groups (lower(name))
    WHERE deleted_at IS NULL;
    
ALTER TABLE agents ADD COLUMN group_id UUID REFERENCES groups (id) ON DELETE SET NULL;