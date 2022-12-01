#!/bin/bash

echo "127.0.0.1 cma.hiveio.internal" >> /etc/hosts

cp -r /home/admin1/.ssh /root/
chown -R root:root /root/.ssh

git config --add --global url."git@github.com:".insteadOf https://github.com

go install -v github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest
go install -v github.com/ramya-rao-a/go-outline@latest
go install -v github.com/cweill/gotests/gotests@latest
go install -v github.com/fatih/gomodifytags@latest
go install -v github.com/josharian/impl@latest
go install -v github.com/haya14busa/goplay/cmd/goplay@latest
go install -v github.com/go-delve/delve/cmd/dlv@latest
go install -v honnef.co/go/tools/cmd/staticcheck@latest
go install -v golang.org/x/tools/gopls@latest
go install -v golang.org/x/tools/cmd/goimports@latest
go install -v github.com/rogpeppe/godef@latest
go install -v github.com/stamblerre/gocode@latest

go mod download
go install -v github.com/spf13/cobra/cobra@latest
go mod tidy -compat=1.17