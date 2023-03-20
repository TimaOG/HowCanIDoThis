CREATE TABLE Users (
	id SERIAL PRIMARY KEY,
	userName varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	userPassword varchar(255) NOT NULL,
	discribtion varchar(1023),
    profileImg varchar,
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
	taskTypeName character varying(255) NOT NULL
);

CREATE TABLE taskTypeSecond (
	id SERIAL PRIMARY KEY,
	taskTypeName character varying(255) NOT NULL,
	fkFirstType INT NOT NULL,
	FOREIGN KEY (fkFirstType) REFERENCES TaskTypeFirst (Id)
);

CREATE TABLE Orders (
	id SERIAL PRIMARY KEY,
	orderName varchar(255) NOT NULL,
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
	tzPath varchar,
	FOREIGN KEY (fkUserOwner) REFERENCES Users (Id),
	FOREIGN KEY (taskTypeFirst) REFERENCES TaskTypeFirst (Id),
	FOREIGN KEY (taskTypeSecond) REFERENCES taskTypeSecond (Id)
);
CREATE TABLE Offers (
	id SERIAL PRIMARY KEY,
	offerName varchar(255) NOT NULL,
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
	coverPath varchar,
	FOREIGN KEY (fkUserOwner) REFERENCES Users (Id),
	FOREIGN KEY (taskTypeFirst) REFERENCES TaskTypeFirst (Id),
	FOREIGN KEY (taskTypeSecond) REFERENCES taskTypeSecond (Id)
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
    closeTime timestamp without time zone,
	FOREIGN KEY (fkUserTaskOwner) REFERENCES Users (Id),
	FOREIGN KEY (fkUserWhoDo) REFERENCES Users (Id)
);

CREATE TABLE Chats (
    id SERIAL PRIMARY KEY,
    fkUserFirst INT NOT NULL,
	fkUserSecond INT NOT NULL,
	lastMessage varchar(20),
	lastTime timestamp without time zone,
	FOREIGN KEY (fkUserFirst) REFERENCES Users (Id),
	FOREIGN KEY (fkUserSecond) REFERENCES Users (Id)
);


CREATE TABLE Messages (
    id SERIAL PRIMARY KEY,
    fkChatId INT,
	messageText TEXT,
	fkUserId INT,
	messageType varchar(10),
	sendTime timestamp without time zone,
	FOREIGN KEY (fkChatId) REFERENCES Chats (Id),
	FOREIGN KEY (fkUserId) REFERENCES Users (Id) 
);

