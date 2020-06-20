FROM  sudachen/linux:latest
LABEL maintainer="Alexey Sudachen <alexey@sudachen.name>"

USER root
RUN set -ex \
 && apt-get update --fix-missing \
 && apt-get install -qy --no-install-recommends \
     software-properties-common python3 python3-pip\
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* \
 && curl -L --retry 3 --output /usr/bin/jq https://circle-downloads.s3.amazonaws.com/circleci-images/cache/linux-amd64/jq-latest \
 && chmod +x /usr/bin/jq \
 && jq --version \
 && export DOCKER_VERSION=$(curl --silent --fail --retry 3 https://download.docker.com/linux/static/stable/x86_64/ | grep -o -e 'docker-[.0-9]*\.tgz' | sort -r | head -n 1) \
 && DOCKER_URL="https://download.docker.com/linux/static/stable/x86_64/${DOCKER_VERSION}" \
 && echo Docker URL: $DOCKER_URL \
 && curl --silent --show-error --location --fail --retry 3 "${DOCKER_URL}" | tar -xz -C /tmp \
 && mv /tmp/docker/* /usr/bin \
 && rm -rf /tmp/docker  \
 && which docker \
 && (docker version || true)

 #&& COMPOSE_URL="https://circle-downloads.s3.amazonaws.com/circleci-images/cache/linux-amd64/docker-compose-latest" \
 #&& curl --silent --show-error --location --fail --retry 3 --output /usr/bin/docker-compose $COMPOSE_URL \
 #&& chmod +x /usr/bin/docker-compose \
 #&& docker-compose version \
 #&& DOCKERIZE_URL="https://circle-downloads.s3.amazonaws.com/circleci-images/cache/linux-amd64/dockerize-latest.tar.gz" \
 #&& curl --silent --show-error --location --fail --retry 3 $DOCKERIZE_URL | tar -C /usr/local/bin -xzv  \
 #&& dockerize --version

COPY ./app /app
RUN pip3 install -r /app/requirements.txt

CMD ["/usr/bin/python","/app/testnet.py"]

