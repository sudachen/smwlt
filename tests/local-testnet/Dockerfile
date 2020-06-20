FROM ubuntu:18.04
LABEL maintainer="Alexey Sudachen <alexey@sudachen.name>"

USER root
WORKDIR "/"
COPY bin/go-spacemesh /tests/local-testnet/bin/
COPY bin/poet /tests/local-testnet/bin/
COPY bin/local-testnet /tests/local-testnet/bin/
COPY genesis_accounts.json /tests/local-testnet/
CMD ["/tests/local-testnet/bin/local-testnet"]
EXPOSE 19090 19190

