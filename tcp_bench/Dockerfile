
# Pull base image.
FROM ubuntu:14.04

RUN \
  sed -i 's/# \(.*multiverse$\)/\1/g' /etc/apt/sources.list && \
  apt-get update && \
  apt-get install -y curl apache2 build-essential software-properties-common php5 libapache2-mod-php5 mercurial

# Install docker
RUN \
  mkdir -p /goroot && \
  curl https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1

# Set go environment variables
ENV GOROOT /goroot
ENV GOPATH /data
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH

# Define working directory.
WORKDIR /data

ADD . /data

# run apache
RUN \
  rm -rf /var/www/html && \
  ln -s /data/web /var/www/html

RUN \
    go get code.google.com/p/go-uuid/uuid && \
    go build tcp_bench && \
    go install tcp_bench

# Define default command.
CMD ["./run.sh"]
