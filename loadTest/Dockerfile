FROM base/archlinux:2018.05.01

RUN pacman --noconfirm -Sy git go

ENV GOPATH=/root/go
ENV PATH=$PATH:$GOPATH/bin

COPY config.js install.sh run.sh ./
RUN ./install.sh

CMD ./run.sh