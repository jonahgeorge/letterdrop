FROM golang:onbuild

RUN apt-get update
RUN apt-get install -y curl apt-transport-https

RUN curl -L https://packagecloud.io/mattes/migrate/gpgkey | apt-key add -
RUN echo "deb https://packagecloud.io/mattes/migrate/ubuntu/ xenial main" > /etc/apt/sources.list.d/migrate.list
RUN apt-get update
RUN apt-get install -y migrate
