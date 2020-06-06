# rootless mode docker:

systemctl --user start docker
https://docs.docker.com/engine/security/rootless/

# useful docker:

docker run --rm -it ipython
docker container prune
docker ps -a

## build td4 image:

in /dock:
docker build -t td4:v1 .

# TODO: limit docker size
