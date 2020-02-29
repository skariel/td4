-- name: UpsertUser :exec
INSERT INTO td4.users(id, display_name, email, avatar)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id)
DO UPDATE SET display_name=$2, email=$3, avatar=$4;

-- name: CleanPendingRunsPerUSer :exec
DELETE FROM td4.pending_runs_per_user WHERE total == 0;

-- name: GetTestCodeByID :one
SELECT *
FROM td4.test_codes
WHERE id = $1;

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
SELECT *
FROM td4.test_codes t
WHERE created_by = $1
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
ORDER BY s.ts_updated DESC
LIMIT 10 OFFSET $1;

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
ON s.id = r.solution_code_id;

-- name: GetConfByDisplayName :one
SELECT * FROM td4.run_configs
WHERE display_name=$1;

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
INSERT INTO td4.test_codes(created_by, updated_by, title, descr, code, is_private)
VALUES ($1, $1, $2, $3, $4, false)
RETURNING *;

-- name: InsertSolutionCode :one
INSERT INTO td4.solution_codes(created_by, updated_by, test_code_id, code, is_private)
VALUES ($1, $1, $2, $3, false)
RETURNING *;

-- name: InsertRunResult :one
INSERT INTO td4.run_results(run_id, status, title, output)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: UpdateRunStatusByID :exec
UPDATE td4.runs
SET
    ts_end=NOW(),
    status=$2
WHERE
    id=$1;



