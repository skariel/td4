http conf:
-----------
`api_basepath` and `static_basepath` in `front/src/routes/utils.js`
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
export TD4_CORS_ORIGIN=
export TD4_KEY_FILE_PATH=
export TD4_CERTIFICATE_FILE_PATH=


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


# TODO: mytests
# TODO: mysolutions
# TODO: content in "about" page
# TODO: deletion of tests
# TODO: deletion of solutions
# TODO: global stats
# TODO: nicer ui
# TODO: better date-time
# TODO: pagination for solutions
# TODO: invalidate cache on new test and make cache TTL longer
# TODO: remove old, stopped containers periodically
# TODO: show the "edit test" button only if test can be really editted...
# TODO: show the "edit solution" button only if solution can be really editted...
# TODO: custom 404 page for github pages: https://docs.github.com/en/github/working-with-github-pages/creating-a-custom-404-page-for-your-github-pages-site
# TODO: front-end cache for fetch (with some TTL, say 5 sec)
# TODO: show "restore solution code" button only if modified solution!



