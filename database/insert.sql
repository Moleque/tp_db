INSERT INTO users (id, email, nickname, fullname, about)
VALUES (1, 'moleque@mail.ru', 'Moleque', 'Salman Dima', 'Hi!');

INSERT INTO users (id, email, nickname, fullname, about)
VALUES (2, 'capoeira@mail.ru', 'Capoeira', 'Capoeira de Angola', '');

INSERT INTO users (id, email, nickname, fullname, about)
VALUES (3, 'psurdo@yandex.ru', 'Surdo', 'Surdo Pavel', 'My name is Surdo');

INSERT INTO users (id, email, nickname, fullname, about)
VALUES (4, 'espiao@yandex.ru', 'Espiao', 'Pasha Salman', 'Hello');

INSERT INTO users (id, email, nickname, fullname, about)
VALUES (5, 'cueca@gmail.com', 'Cueca', 'Timor', '');

-- ///////////////////////////////////
INSERT INTO forums (id, slug, title, user_id, threads, posts)
VALUES (2, 'capoeira', 'Capoeirando', 1, 0, 0);

INSERT INTO forums (slug, title, user_id)
VALUES ('sdf', 'Capoeirando', 1);