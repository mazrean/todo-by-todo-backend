# syntax = docker/dockerfile:1

FROM bitnami/minideb:bookworm

RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

ENV TIDB_VERSION=v8.1.1

RUN curl --proto '=https' --tlsv1.2 -sSf https://tiup-mirrors.pingcap.com/install.sh | sh && ln -s /root/.tiup/bin/tiup /bin/tiup

RUN tiup install tidb:${TIDB_VERSION} pd:${TIDB_VERSION} tikv:${TIDB_VERSION} playground

ENTRYPOINT [ "tiup" ,"playground", "v8.1.1", "--db.host", "0.0.0.0", "--pd.host", "0.0.0.0" ]
CMD [ "--tiflash", "0" ]

EXPOSE 4000
