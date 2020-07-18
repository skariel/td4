http conf:
-----------
`api_basepath`, `static_basepath`, and `invalidate_cache_ttl_ms` in `front/src/routes/utils.js`
consts in `back/server_api/serve.go`
redirect and homepage in github oauth app
POSTGRES conf in `back/db/dbconn.go`
worker conf consts in `back/worker_test/worker.go`
env variables: `TD4_*`:
export TD4_github_client_id=
export TD4_github_client_secret=
export TD4_JWT_SECRET=
export TD4_API_PORT=
export TD4_SOCIAL_AUTH_REDIRECT="https://api.solvemytest.dev/auth/github/callback"
export TD4_SOCIAL_AUTH_FINAL_DEST=
export TD4_CORS_ORIGIN=https://solvemytest.dev
export TD4_KEY_FILE_PATH=/etc/letsencrypt/live/solvemytest.dev/privkey.pem
export TD4_CERTIFICATE_FILE_PATH=/etc/letsencrypt/live/solvemytest.dev/fullchain.pem
export TD4_CACHE_TTL_SECONDS=7

(key and cert file paths above are the default for lets encrypt cert-bot)

to run local frontend:
-------------------------
in the cloud:
export TD4_CORS_ORIGIN=http://localhost:3000 

Frontend
-----------
npm run dev
npm build
npx sapper export
npx serve __sapper__/export

SQLC
------------
sqlc generate

useful stuff:
----------------
check which process is using a port: `lsof -i :8081`
go linter being used: golangci-lint: https://github.com/golangci/golangci-lint#binary
screen: `gtk-redshift -l 32.08:34.78`


tech used:
--------------
pytest
postgresql
golang
sapper (...svelte)
docker
sqlc
api + static rendered site
github pages
vultr


update packages
---------------------
go:

in the package root folder:
go get -u

upgrade sqlc:
go get -u github.com/kyleconroy/sqlc/cmd/sqlc

npm:

edit package.json for new versions
npm update --save
npm update --save-dev


# TODO: content in "about" page
# TODO: nicer ui
# TODO: better date-time
# TODO: make sure it looks good on mobile
# TODO: test on Chrome

For later:

# TODO: feedback on actions status ("test was added, press x to close" etc.)
# TODO: syntax highlighting
# TODO: add different languages (go, node (js, ts), etc.)
# TODO: global stats (total # tests pending, running, passing, failing, ETA for running a test etc.)
# TODO: custom 404 page for github pages: https://docs.github.com/en/github/working-with-github-pages/creating-a-custom-404-page-for-your-github-pages-site

Bugs:

# BUG: when updating test code: not all new runs get added to td4.pending_runs_per_user_table
# BUG: rate limiting for new solutions used to work, but apparently not anymore...



