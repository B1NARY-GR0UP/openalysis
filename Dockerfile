FROM ubuntu:latest
LABEL authors="justlorain"

ENTRYPOINT ["top", "-b"]