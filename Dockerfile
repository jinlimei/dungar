FROM gitlab.int.magneato.site:4567/dungar/test-image:latest
LABEL MAINTAINER="Jinli <jinli@jinli.dev>"

COPY . ./

RUN find . -type d -print0 | xargs -0 chmod 0755 && \
    find . -type f -print0 | xargs -0 chmod 0644 && \
    mv secrets.test.ini secrets.ini && \
    mv settings.test.ini settings.ini && \
    chmod +x ./tools/*.sh

ENV IN_CI_ENV 1

CMD ["./dungar", "botinfo"]
