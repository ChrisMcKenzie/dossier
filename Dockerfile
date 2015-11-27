FROM scratch

RUN mkdir -p /usr/local/bin /var/dossier

ADD dossier /usr/local/bin

ENTRYPOINT dossier
