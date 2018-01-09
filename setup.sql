CREATE TABLE skill(
  id SERIAL UNIQUE NOT NULL PRIMARY KEY ,
  nameSkill TEXT UNIQUE
);

CREATE TABLE  company (
  id SERIAL UNIQUE NOT NULL PRIMARY KEY,
  idJobs INT ARRAY,
  idSkill TEXT ARRAY,
  nameComp TEXT UNIQUE ,
  address TEXT,
  country TEXT,
  logo TEXT,
  banner TEXT,
  intro TEXT ARRAY
);

CREATE TABLE  job (
  id SERIAL UNIQUE NOT NULL  PRIMARY KEY,
  idSkill TEXT ARRAY,
  idComp INT REFERENCES company(id),
  title TEXT UNIQUE,
  salary TEXT,
  address TEXT,
  time_posted DATE,
  reason TEXT ARRAY,
  description TEXT ARRAY,
  skill TEXT ARRAY
);

CREATE TABLE user_account(
  id SERIAL UNIQUE NOT NULL PRIMARY KEY ,
  idSkill INT ARRAY,
  idComps INT ARRAY,
  fullname TEXT,
  email TEXT,
  username TEXT UNIQUE,
  passwd TEXT,
  address TEXT,
  phone TEXT
);


