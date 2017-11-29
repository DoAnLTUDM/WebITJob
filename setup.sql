CREATE TABLE  company (
  id INT NOT NULL PRIMARY KEY ,
  idJob INT[] REFERENCES job(id),
  name TEXT,
  address TEXT,
  country TEXT
);

CREATE TABLE  job (
  id INT NOT NULL PRIMARY KEY,
  title TEXT,
  salary MONEY,
  address TEXT,
  time_posted TIME,
  reason TEXT,
  description TEXT,
  skill TEXT,
  qualification TEXT
);

CREATE TABLE programming_language(
  id INT NOT NULL PRIMARY KEY ,
  idJob INT[] REFERENCES job(id),
  name text
);

CREATE TABLE user(
  id INT NOT NULL PRIMARY KEY ,
  idPLangs INT[] REFERENCES programming_language(id),
  idCompanies INT[] REFERENCES company(id),
  name TEXT,
  email TEXT,
  username TEXT,
  passwd TEXT,
  address TEXT,
  phone TEXT
);
