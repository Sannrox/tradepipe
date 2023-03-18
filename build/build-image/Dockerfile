ARG BASE_IMAGE

FROM $BASE_IMAGE 

RUN touch /tradepipe-builder 

RUN chmod -R a+rwx /usr/local/go/pkg


ENV HOME=/go/src/github.com/Sannrox/tradepipe
WORKDIR ${HOME}


ENV OUPUT_SUBPATH=_output/dockerized 

ENV GIT_VERSION_FIEL=${HOME}/.dockerized-versions 

RUN git config --system user.name "tradepipe-builder" && \
    git config --system user.email "tradepipe-builder@tradepipe-builder"

RUN chmod -R a+rwx $GOPATH

# ADD localtime /etc/localtime
# RUN chmod 644 /etc/localtime
