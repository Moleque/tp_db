FROM ubuntu:18.04

MAINTAINER Dmitriy Salman

# Обновление списка пакетов
RUN apt-get -y update


# ===============================
# Установка postgresql

ENV PGVER 10
RUN apt-get install -y postgresql-$PGVER

# Регулировка конфигурации PostgreSQL, чтобы можно было удаленно подключаться к базе данных
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# PostgreSQL будет принимать подключения со всех адресов
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf


# ===============================
# Установка golang

ENV GOVER 1.10
RUN apt install -y golang-$GOVER git

# Выставляем переменную окружения для сборки проекта
ENV GOROOT /usr/lib/go-$GOVER
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH


# ===============================
# Сборка проекта

# Копируем исходный код в Docker-контейнер
WORKDIR $GOPATH/src/tp_db/
ADD . $GOPATH/src/tp_db/

# Подтягиваем зависимости
RUN go get github.com/julienschmidt/httprouter
RUN go get github.com/lib/pq
RUN go install ./forum

# Объявлем порт сервера
EXPOSE 5000


# ===============================
# Создание базы данных

# Запуск остальных команд под пользователем `postgres`
USER postgres

# Создание роли PostgreSQL с именем `docker` и паролем `docker`, 
# затем создание базы данных `docker`, принадлежащей роли `docker`
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql -d docker -a -f ./database/create.sql &&\
    /etc/init.d/postgresql stop

# Повесить PostgreSQL на порт
EXPOSE 5432

# Добавление VOLUME, для разрешения резервного копирования конфигураций, журналов и баз данных
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]


# ===============================
# Запускаем PostgreSQL и сервер

CMD service postgresql start && forum