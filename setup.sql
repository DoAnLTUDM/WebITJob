CREATE TABLE  company (
  id SERIAL UNIQUE NOT NULL PRIMARY KEY,
  idJobs INT ARRAY,
  company_name TEXT UNIQUE ,
  address TEXT,
  country TEXT,
  logo TEXT,
  urlC TEXT
);

CREATE TABLE  job (
  id SERIAL UNIQUE NOT NULL  PRIMARY KEY,
  idPLang INT,
  title TEXT,
  salary MONEY,
  address TEXT,
  time_posted TIME,
  reason TEXT,
  description TEXT,
  skill TEXT,
  qualification TEXT,
  company_name TEXT
);

CREATE TABLE programming_language(
  id SERIAL UNIQUE NOT NULL PRIMARY KEY ,
  pl_name TEXT
);

CREATE TABLE user_account(
  id SERIAL UNIQUE NOT NULL PRIMARY KEY ,
  idPLangs INT ARRAY,
  idCompanies INT ARRAY,
  usr_name TEXT,
  email TEXT,
  username TEXT,
  passwd TEXT,
  address TEXT,
  phone TEXT
);
