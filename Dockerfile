FROM ubuntu:latest

RUN apt-get update && apt-get upgrade -y
RUN apt-get install curl git bzr -y

# install golang
RUN curl -s https://storage.googleapis.com/golang/go1.3.linux-amd64.tar.gz | tar -v -C /usr/local/ -xz

# path config
ENV PATH  $PATH:/usr/local/go/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin:/go/bin
ENV GOPATH  /go
ENV GOROOT  /usr/local/go

RUN cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime 
ADD . /go/src/github.com/ernado/cyvisor
RUN cd /go/src/github.com/ernado/cyvisor && go get .

ENTRYPOINT ["cyvisor"]
