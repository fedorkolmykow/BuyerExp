\connect postgres

CREATE DATABASE avitoexp;

\connect avitoexp

CREATE TABLE Users(
	user_id serial NOT NULL,
	mail varchar(255) NOT NULL UNIQUE,
	confirmed BOOLEAN NOT NULL,
	hash  bytea,
	CONSTRAINT User_pk PRIMARY KEY (user_id)
) WITH (
  OIDS=FALSE
);



CREATE TABLE Notices (
	notice_id serial NOT NULL,
	url varchar(255) NOT NULL,
	price FLOAT NOT NULL,
	CONSTRAINT Notice_pk PRIMARY KEY (notice_id)
) WITH (
  OIDS=FALSE
);



CREATE TABLE Subscription (
	notice_id serial NOT NULL,
	user_id serial NOT NULL,
	CONSTRAINT Subscription_pk PRIMARY KEY (notice_id,user_id)
) WITH (
  OIDS=FALSE
);





ALTER TABLE Subscription ADD CONSTRAINT Subscription_fk0 FOREIGN KEY (notice_id) REFERENCES Notices(notice_id);
ALTER TABLE Subscription ADD CONSTRAINT Subscription_fk1 FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE;
