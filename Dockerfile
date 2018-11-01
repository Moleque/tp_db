FROM ubuntu:18.04

MAINTAINER Dmitriy Salman

# Обновление списка пакетов
RUN apt-get -y update


# ===============================
# Установка postgresql

ENV PGVER 10
RUN apt-get install -y postgresql-$PGVER

# Запуск остальных команд под пользователем `postgres`, созданным пакетом `postgres-$PGVER` при установке
USER postgres

COPY database/create.sql database/create.sql

# Создание роли PostgreSQL с именем `docker` и паролем `docker`, 
# затем создание базы данных `docker`, принадлежащей роли `docker`
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql -a -f database/create.sql &&\
    /etc/init.d/postgresql stop

# Регулировка конфигурации PostgreSQL, чтобы можно было удаленно подключаться к базе данных
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# PostgreSQL будет принимать подключения со всех адресов
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Повесить PostgreSQL на порт
EXPOSE 5432

# Добавление VOLUME, для разрешения резервного копирования конфигураций, журналов и баз данных
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Вернуться к пользователю root
USER root


# ===============================
# Сборка проекта

# Установка golang
ENV GOVER 1.10
RUN apt install -y golang-$GOVER git

# Выставляем переменную окружения для сборки проекта
ENV GOROOT /usr/lib/go-$GOVER
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

# Копируем исходный код в Docker-контейнер
WORKDIR $GOPATH/src/github.com/moleque/tp_db
ADD . $GOPATH/src/github.com/moleque/tp_db
RUN go get github.com/Moleque/tp_db/forum/controllers
RUN go build ./forum/main.go

# Собираем генераторы
# RUN go install ./vendor/github.com/go-swagger/go-swagger/cmd/swagger
# RUN go install ./vendor/github.com/jteeuwen/go-bindata/go-bindata

# Собираем и устанавливаем пакет
# RUN go generate -x ./restapi
# RUN go install ./cmd/forum

# Объявлем порт сервера
EXPOSE 5000


# ===============================
# Запускаем PostgreSQL и сервер

CMD service postgresql start && ./main --scheme=http --port=5000 --host=0.0.0.0 --database=postgres://docker:docker@localhost/docker