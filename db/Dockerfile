FROM postgres:14

RUN apt-get update && \
    apt-get install -y postgresql-server-dev-14 make gcc

COPY ./extensions/ /tmp/extensions/

WORKDIR /tmp/extensions/pg_raise_error


RUN make && make install

COPY ./config/postgresql.conf /etc/postgresql/

RUN rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
