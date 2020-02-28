DROP SCHEMA IF EXISTS td4 CASCADE;
CREATE SCHEMA td4;


CREATE TABLE td4.users (
    ts_created timestamptz DEFAULT now() NOT NULL,
    ts_updated timestamptz DEFAULT now() NOT NULL,
    
    id text PRIMARY KEY,
    display_name text NOT NULL UNIQUE,
    email text NOT NULL UNIQUE,
    avatar text NOT NULL UNIQUE
);
INSERT INTO td4.users(id, display_name, email, avatar)
VALUES ('admin', 'admin', '', '');


CREATE TABLE td4.test_codes (
    id SERIAL PRIMARY KEY,
    ts_created timestamptz DEFAULT now() NOT NULL,
    ts_updated timestamptz DEFAULT now() NOT NULL,
    created_by text NOT NULL REFERENCES td4.users(id) ON DELETE CASCADE,
    updated_by text NOT NULL REFERENCES td4.users(id) ON DELETE CASCADE,

    title text NOT NULL,
    descr text NOT NULL,
    code text NOT NULL,
    is_private bool NOT NULL,
    is_draft bool DEFAULT false NOT NULL
);
CREATE INDEX upserted_by_test_codes_index ON td4.test_codes (created_by, updated_by);


CREATE TABLE td4.solution_codes (
    id SERIAL PRIMARY KEY,
    ts_created timestamptz DEFAULT now() NOT NULL,
    ts_updated timestamptz DEFAULT now() NOT NULL,
    created_by text NOT NULL REFERENCES td4.users(id) ON DELETE CASCADE,
    updated_by text NOT NULL REFERENCES td4.users(id) ON DELETE CASCADE,

    test_code_id integer NOT NULL REFERENCES td4.test_codes(id) ON DELETE CASCADE,
    code text NOT NULL,
    is_private bool NOT NULL,
    is_draft bool DEFAULT false NOT NULL
);
CREATE INDEX upserted_by_solution_codes_index ON td4.test_codes (created_by, updated_by);
CREATE INDEX test_code_id_solution_codes_index ON td4.solution_codes (test_code_id, is_private, is_draft);


CREATE TABLE td4.run_configs (
    display_name text PRIMARY KEY,
    ts_created timestamptz DEFAULT now() NOT NULL,
    ts_updated timestamptz DEFAULT now() NOT NULL,
    created_by text NOT NULL REFERENCES td4.users(id) ON DELETE CASCADE,
    updated_by text NOT NULL REFERENCES td4.users(id) ON DELETE CASCADE,

    cpu_period integer NOT NULL DEFAULT 100000,
    cpu_quota integer NOT NULL DEFAULT   10000,
    memory_usage_mb integer NOT NULL DEFAULT 64,
    disk_usage_mb integer NOT NULL DEFAULT 64,
    max_time_secs integer NOT NULL DEFAULT 10
);
CREATE INDEX upserted_by_run_configs_index ON td4.run_configs (created_by, updated_by);
INSERT INTO td4.run_configs(display_name, created_by, updated_by)
VALUES ('default', 'admin', 'admin');


CREATE TYPE td4.type_run_status
AS ENUM ('pending', 'wip', 'pass', 'fail', 'problem', 'stop');


CREATE TABLE td4.runs (
    id SERIAL PRIMARY KEY,
    ts_created timestamptz DEFAULT now() NOT NULL,
    ts_updated timestamptz DEFAULT now() NOT NULL,
    created_by text NOT NULL REFERENCES td4.users(id) ON DELETE CASCADE,
    updated_by text NOT NULL REFERENCES td4.users(id) ON DELETE CASCADE,

    ts_start timestamptz,
    ts_end timestamptz,
    solution_code_id integer NOT NULL REFERENCES td4.solution_codes(id) ON DELETE CASCADE,
    run_config text NOT NULL REFERENCES td4.run_configs(display_name) ON DELETE CASCADE,
    status td4.type_run_status NOT NULL default 'pending'::td4.type_run_status NOT NULL,
    UNIQUE(solution_code_id)
);
CREATE INDEX status_runs_index ON td4.runs (status);
CREATE INDEX created_by_runs_index ON td4.runs (created_by, status);
CREATE INDEX upserted_by_runs_index ON td4.runs (created_by, updated_by);


-- results can be used to store unit test metadata, and later on also actual results
CREATE TYPE td4.type_run_result_status
AS ENUM ('pass', 'fail', 'skip', 'stop', 'todo');


CREATE TABLE td4.run_results (
    id SERIAL PRIMARY KEY,
    ts_created timestamptz DEFAULT now() NOT NULL,
    ts_updated timestamptz DEFAULT now() NOT NULL,

    run_id integer NOT NULL REFERENCES td4.runs(id) ON DELETE CASCADE,

	status td4.type_run_result_status NOT NULL DEFAULT 'todo',
	title text,
	output text
);
CREATE INDEX run_id_run_results_index ON td4.run_results (run_id);


-- Table to support fair work stealing by workers
-- A user is picked at random and then a test is picked in order of creation
CREATE TABLE td4.pending_runs_per_user (
    user_id text NOT NULL PRIMARY KEY,
    ts_created timestamptz NOT NULL DEFAULT now() NOT NULL,
    ts_updated timestamptz NOT NULL DEFAULT now() NOT NULL,

    total integer NOT NULL default 1
);
CREATE INDEX user_id_pending_runs_per_user_index ON td4.pending_runs_per_user (user_id);


CREATE FUNCTION td4.function_insert_run() RETURNS trigger
LANGUAGE plpgsql
AS $$BEGIN
    INSERT INTO td4.pending_runs_per_user(user_id)
    VALUES (NEW.created_by)
    ON CONFLICT (user_id)
    DO UPDATE
    SET user_id = NEW.created_by, total = EXCLUDED.total + 1;
    RETURN NEW;
END$$;

DO LANGUAGE plpgsql
$$BEGIN
CREATE TRIGGER trigger_insert_run
AFTER INSERT ON td4.runs
FOR EACH ROW EXECUTE FUNCTION td4.function_insert_run();
END$$;


-- automatic updating `ts_updated` columns

CREATE FUNCTION td4.function_set_timestamp() RETURNS trigger
LANGUAGE plpgsql
AS $$BEGIN
  NEW.ts_updated = NOW();
  RETURN NEW;
END;
$$;

DO LANGUAGE plpgsql
$$
DECLARE
    t text;
BEGIN
    FOR t IN
    SELECT table_name
    FROM information_schema.tables
    WHERE table_schema='td4'
    LOOP
        EXECUTE 'CREATE TRIGGER trigger_set_timestamp BEFORE UPDATE ON td4.' || t || ' FOR EACH ROW EXECUTE FUNCTION td4.function_set_timestamp();';
    END LOOP;
END$$;


