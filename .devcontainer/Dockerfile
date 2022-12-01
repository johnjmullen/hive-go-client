FROM ubuntu:18.04

ENV DEBIAN_FRONTEND=noninteractive
LABEL Name=hive-go-client Version=0.0.1
RUN apt-get update \
    && apt -y upgrade \
    && apt -y install \
              locales \
              rsync \
              openssh-server \
              sudo \
              procps \
              wget \
              unzip \
              mc \
              ca-certificates \
              curl \
              software-properties-common \
              bash-completion \
              shellcheck \
              whois \
              libudev-dev \
              libpam0g-dev \
              golang-any \
              python-dev \
    && add-apt-repository ppa:hio/appliance \
    && add-apt-repository ppa:longsleep/golang-backports \
    && add-apt-repository ppa:git-core/ppa \                                                  
    && apt-get update \
    && apt-get -y install \
                  git \
                  subversion \
                  libvirt-daemon-system \
                  libvirt-dev libsystemd-dev \
                  libvirt-dbus \
                  debhelper \
                  devscripts \
                  openssl \
                  golang-go \
    && rm -rf /var/lib/apt/lists/* \
    && localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8 \
    && chsh -s /bin/bash

SHELL ["/bin/bash", "-c"]
RUN echo "source /usr/share/bash-completion/completions/git" >> ~/.bashrc
RUN source ~/.bashrc


  # hivesetup
RUN mkdir -p /opt/hive/conf /opt/hive/auth \
    && ssh-keygen -A  \
    && mkdir -p /run/sshd  \
    && echo '{"broker": true, "clusterId": "TestClusterId", "cma": "localhost", "cpuGovernor": "performance", "crs": true, "dbName": "TestDb", "firmware": { "active": "hiveio-fabric-v7.4.1-1342", "pendingSwap": false, "software": "hiveio-fabric-v8.3.0-1118" }, "hostname": "TestHostname", "loglevel": "debug", "maxCloneDensity": 119, "ntp": "0.us.pool.ntp.org", "role": "cma", "timezone": "UTC"}' > /opt/hive/conf/appliance.json \
    && echo -n "960d739efb4e" > /opt/hive/conf/hostid

ENV BASE=/usr/local
ENV GO111MODULE=on
ENV GOPRIVATE=github.com/hive-io
ENV GOPATH=$HOME/go
ENV PATH=$GOPATH/bin:$PATH

ENV LC_ALL=en_US.UTF-8 \
    LANG=en_US.UTF-8 \
    LANGUAGE=en_US.UTF-8 \
    SHELL=/bin/bash 

USER root
RUN /usr/sbin/useradd admin1 -s /bin/bash -m \
    && echo admin1:admin | chpasswd

COPY [".devcontainer/testenv.sh","/usr/bin/testenv.sh"]
COPY [".devcontainer/publish-aptly.sh","/usr/bin/publish-aptly.sh"]
WORKDIR /workspaces/hive-go-client
COPY . .

