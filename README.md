http conf:
-----------
`basepath` in `front/src/utils.js`
port and CORS in `back/serve.go`
redirect and homepage in github oauth app
POSTGRES conf in `sql/db/dbconn.go`

Frontend
-----------
npm run dev
npm build
npx sapper export
npx serve __sapper__/export

SQLC
------------
sqlc generate

server
----------
./dev_run.py

worker
-----------
currently run manually: `go run .` or `go run worker` etc.

