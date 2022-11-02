BEGIN;

CREATE TABLE IF NOT EXISTS users (
  id SERIAL NOT NULL,
  email VARCHAR(128) UNIQUE NOT NULL,
  password VARCHAR(60),
  fullname VARCHAR(128) NOT NULL,
  is_verified BOOLEAN NOT NULL,
  token_id VARCHAR(36),
  token_verification VARCHAR(128),
  token_expiration VARCHAR(10),

  CONSTRAINT users__pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS category (
  id SERIAL NOT NULL,
  name VARCHAR(50) NOT NULL,

  CONSTRAINT category__pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS books (
  id SERIAL NOT NULL,
  title VARCHAR(60) NOT NULL,
  author VARCHAR(60) NOT NULL,
  synopsis TEXT NOT NULL,
  cover VARCHAR(150) NOT NULL,
  reader INT NOT NULL DEFAULT 0,
  category_id SERIAL NOT NULL,

  CONSTRAINT books__pkey PRIMARY KEY (id),
  CONSTRAINT books__category__fk FOREIGN KEY (category_id) REFERENCES category(id)
);
CREATE INDEX IF NOT EXISTS book_title ON books(title);
CREATE INDEX IF NOT EXISTS books__category__idx ON books(category_id);

CREATE TABLE IF NOT EXISTS book_rating (
  book_id SERIAL NOT NULL,
  user_id SERIAL NOT NULL,
  rating REAL NOT NULL,

  CONSTRAINT book_rating__pkey PRIMARY KEY (book_id, user_id),
  CONSTRAINT book_rating__books__fk FOREIGN KEY (book_id) REFERENCES books(id),
  CONSTRAINT book_rating__users__fk FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX IF NOT EXISTS book_rating__books__idx ON book_rating(book_id);
CREATE INDEX IF NOT EXISTS book_rating__users__idx ON book_rating(user_id);

COMMIT;