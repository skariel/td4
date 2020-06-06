module test_worker

go 1.14

replace sql => /home/skariel/prog/td4/sql

require (
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/joshdk/go-junit v0.0.0-20200312181801-e5d93c0f31a8
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	sql v0.0.0-00010101000000-000000000000
)
