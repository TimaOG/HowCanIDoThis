CREATE TABLE Users (
	id SERIAL PRIMARY KEY,
	name character varying(255) NOT NULL,
	email character varying(255) NOT NULL,
	password character varying(255) NOT NULL,
	discribtion character varying(1023),
    profileImg character varying(1023),
	rating smallint,
    balance int,
    isPremiumUser BOOLEAN,
	isActiveUser BOOLEAN,
	historyCount INT,
    responsibility smallint, --ответсвенность
    doneOnTime smallint,
    answerSpead smallint,
    registrationDate timestamp without time zone
);


CREATE TABLE TaskTypeFirst (
	id SERIAL PRIMARY KEY,
	name character varying(255) NOT NULL
);

CREATE TABLE taskTypeSecond (
	id SERIAL PRIMARY KEY,
	name character varying(255) NOT NULL,
	fkFirstType INT NOT NULL
);

CREATE TABLE Orders (
	id SERIAL PRIMARY KEY,
	name character varying(255) NOT NULL,
	fkUserOwner INT NOT NULL,
	discribtion TEXT NOT NULL,
	price INT NOT NULL,
	deadline timestamp without time zone NOT NULL,
    urgency smallint, --срочность,
    workType smallint, 
	taskTypeFirst smallint NOT NULL,
	taskTypeSecond smallint NOT NULL,
    tags text,
	isActive BOOLEAN,
	tzPath character varying(1023)
);
CREATE TABLE Offers (
	id SERIAL PRIMARY KEY,
	name character varying(255) NOT NULL,
	fkUserOwner INT NOT NULL,
	discribtion TEXT NOT NULL,
	price INT NOT NULL,
	daysToComplite integer,
    workType smallint, 
	taskTypeFirst smallint NOT NULL,
	taskTypeSecond smallint NOT NULL,
    tags text,
	isActive BOOLEAN,
	rating smallint,
	historyCount SERIAL,
	coverPath character varying(1023)
);

CREATE TABLE Mathes (
    id SERIAL PRIMARY KEY,
    fkUserTaskOwner INT NOT NULL,
	fkUserWhoDo INT NOT NULL,
	fkWhatTaskId INT NOT NULL,
	userOwnerConfim BOOLEAN,
	userWhoDoConfim BOOLEAN,
	isOrder BOOLEAN,
    startTime timestamp without time zone,
    closeTime timestamp without time zone
);

CREATE TABLE Chats (
    id SERIAL PRIMARY KEY,
    fkUserFirst INT NOT NULL,
	fkUserSecond INT NOT NULL,
	fkLastMessage INT
);

CREATE TABLE Messages (
    id SERIAL PRIMARY KEY,
    fkChatId INT,
	messageText TEXT,
	fkUserId INT,
	sendTime timestamp without time zone 
);

