FROM fedora:29

ENV PACKAGES="iproute can-utils"

RUN dnf install -y ${PACKAGES}

USER root

ENTRYPOINT exec /bin/bash -c "ip link show"
#CMD exec /bin/bash -c "trap : TERM INT; sleep infinity & wait"