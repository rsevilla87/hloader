FROM registry.fedoraproject.org/fedora-minimal

COPY bin/hloader /usr/bin/hloader

ENTRYPOINT ["/usr/bin/hloader"]
