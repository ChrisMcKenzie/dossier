FROM busybox

ADD dossier /usr/local/bin/dossier

ENTRYPOINT ["/usr/local/bin/dossier"]
CMD ["--base=/var"]
