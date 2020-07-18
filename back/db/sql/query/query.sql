
-- name: UpsertUser :exec
INSERT INTO td4.users(id, display_name, email, avatar)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id)
DO UPDATE SET display_name=$2, email=$3, avatar=$4;

-- name: CleanPendingRunsPerUSer :exec
DELETE FROM td4.pending_runs_per_user WHERE total = 0;

-- name: DeleteTestByID :exec
DELETE FROM td4.test_codes
WHERE id = $1;

-- name: DeleteSolutionByID :exec
DELETE FROM td4.solution_codes
WHERE id = $1;

-- name: GetTestCodeByID :one
SELECT
    t.*,
    u.display_name,
    u.avatar
FROM td4.test_codes t
JOIN td4.users u
ON t.created_by = u.id
WHERE t.id = $1
LIMIT 1;

-- name: GetTestCodes :many
SELECT
    t.*,
    u.display_name,
    u.avatar
FROM td4.test_codes t
JOIN td4.users u
ON t.created_by = u.id
ORDER BY t.ts_updated DESC
LIMIT 10 OFFSET $1;

-- name: GetTestCodesByUser :many
SELECT
    t.*,
    u.display_name,
    u.avatar
FROM td4.test_codes t
JOIN td4.users u
ON t.created_by = u.id
WHERE t.created_by = (SELECT id FROM td4.users u WHERE u.display_name = $1) OR EXISTS (SELECT * FROM td4.solution_codes s WHERE s.created_by = $1)
ORDER BY t.ts_updated DESC
LIMIT 10 OFFSET $2;

-- name: GetSolutionCodesByTest :many
SELECT
    s.*,
    u.display_name,
    u.avatar,
    r.id as run_id,
    r.ts_start,
    r.ts_end,
    r.status
FROM td4.solution_codes s
JOIN td4.users u
ON s.created_by = u.id
JOIN td4.runs r
ON s.id = r.solution_code_id
WHERE s.test_code_id = $1
ORDER BY s.ts_updated DESC
LIMIT 10 OFFSET $2;

-- name: GetSolutionCodeByID :one
SELECT
    s.*,
    u.display_name,
    u.avatar,
    r.id as run_id,
    r.ts_start,
    r.ts_end,
    r.status
FROM td4.solution_codes s
JOIN td4.test_codes t
ON s.test_code_id = t.id
JOIN td4.users u
ON t.created_by = u.id
JOIN td4.runs r
ON s.id = r.solution_code_id
WHERE s.id = $1;

-- name: GetResultsByRun :many
SELECT * FROM td4.run_results WHERE run_id = $1;

-- name: GetConfByDisplayName :one
SELECT * FROM td4.run_configs
WHERE display_name = $1;

-- name: RAWGetSolutionCodeByID :one
SELECT * FROM td4.solution_codes
WHERE id = $1;

-- name: FetchSomeRun :many
WITH pending_run AS (
UPDATE td4.pending_runs_per_user SET total = total - 1
WHERE user_id = (
    SELECT user_id
    FROM td4.pending_runs_per_user
    WHERE total > 0
    ORDER BY random()
    LIMIT 1
    FOR UPDATE
    SKIP LOCKED
)
RETURNING *)
UPDATE td4.runs
SET status = 'wip', ts_start = NOW()
WHERE id = (
    SELECT id
    FROM td4.runs t, pending_run
    WHERE t.created_by = pending_run.user_id AND status = 'pending'
    ORDER BY t.ts_updated ASC
    LIMIT 1
    FOR UPDATE
    SKIP LOCKED
)
RETURNING *;

-- name: InsertTestCode :one
INSERT INTO td4.test_codes(created_by, updated_by, title, descr, code)
VALUES ($1, $1, $2, $3, $4)
RETURNING *;

-- name: InsertSolutionCode :one
INSERT INTO td4.solution_codes(created_by, updated_by, test_code_id, code)
VALUES ($1, $1, $2, $3)
RETURNING *;

-- name: InsertRunResult :one
INSERT INTO td4.run_results(run_id, status, title, output)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: EndRunByID :exec
UPDATE td4.runs
SET
    ts_end=NOW(),
    status=$2
WHERE
    id=$1;


-- name: FailLongRuns :exec
UPDATE td4.runs r
SET status = 'stop', ts_end = NOW()
WHERE r.status = 'wip' AND NOW() - r.ts_start > (
    SELECT max_time_secs
    FROM td4.run_configs rc
    WHERE r.run_config = rc.display_name
) * '1 sec'::interval;



-- name: UpdateSolutionCode :exec
UPDATE td4.solution_codes SET code = $2 WHERE id = $1;

-- name: UpdateTestCode :exec
UPDATE td4.test_codes SET title = $2, descr = $3, code = $4 WHERE id = $1;
