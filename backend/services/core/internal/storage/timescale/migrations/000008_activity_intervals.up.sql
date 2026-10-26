
-- интервалы активности
CREATE TABLE activity_intervals (
  agent_id   UUID NOT NULL,
  started_at TIMESTAMPTZ NOT NULL,
  ended_at   TIMESTAMPTZ,
  kind       SMALLINT NOT NULL,
  app_id     INT,             -- for kind == "focus" 
  category   SMALLINT,        -- for kind == "focus" 
  title      TEXT,            -- for kind == "focus"  
  meta       JSONB
);

SELECT create_hypertable('activity_intervals', 'started_at',
                         chunk_time_interval => INTERVAL '1 day');