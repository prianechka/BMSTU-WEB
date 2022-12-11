DROP SCHEMA if exists ppo cascade;
CREATE SCHEMA ppo;

CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    userlogin TEXT unique,
    userpassword TEXT,
    userrole int
);

INSERT INTO users(userlogin, userpassword, userrole) VALUES ('commend', '1234', 3);
INSERT INTO users(userlogin, userpassword, userrole) VALUES ('supp', 'supp', 2);
INSERT INTO users(userlogin, userpassword, userrole) VALUES ('priany', '123456', 1);

CREATE TABLE rooms
(
    roomid SERIAL PRIMARY KEY,
    roomtype TEXT,
    roomnumber INT
);

INSERT INTO rooms(roomtype, roomnumber) VALUES ('Склад', 0);
INSERT INTO rooms(roomtype, roomnumber) VALUES ('Комната', 628);
INSERT INTO rooms(roomtype, roomnumber) VALUES ('Комната', 629);

CREATE TABLE student
(
    studentid SERIAL PRIMARY KEY,
    studentname TEXT,
    studentsurname TEXT,
    studentgroup TEXT,
    studentnumber TEXT UNIQUE,
    settledate DATE,
    webaccid int,
    FOREIGN KEY (webaccid) references users(id)
);

CREATE TABLE thing
(
    thingid SERIAL PRIMARY KEY,
    marknumber INT UNIQUE,
    creationdate DATE,
    thingtype TEXT
);

CREATE TABLE studentroomhistory
(
    id SERIAL PRIMARY KEY,
    studentid int,
    roomid int,
    direction int,
    transferdate date,
    FOREIGN KEY (studentid) references student(studentid),
    FOREIGN KEY (roomid) references rooms(roomid)
);

CREATE TABLE studentthinghistory
(
    id SERIAL PRIMARY KEY,
    studentid int,
    thingid int,
    direction int,
    transferdate date,
    FOREIGN KEY (studentid) references student(studentid),
    FOREIGN KEY (thingid) references thing(thingid)
);

CREATE TABLE thingroomhistory
(
    id SERIAL PRIMARY KEY,
    srcroomid int,
    dstroomid int,
    thingid int,
    transferdate date,
    FOREIGN KEY (thingid) references thing(thingid),
    FOREIGN KEY (srcroomid) references rooms(roomid),
    FOREIGN KEY (dstroomid) references rooms(roomid)
);

create function findroom(idthing integer) returns integer
    language plpgsql
as
$$
    DECLARE tmp int = 1;
    BEGIN
        SELECT dstroomid into tmp
        FROM    (SELECT STH.ThingID, max(STH.ID) as "last"
                FROM ThingRoomHistory as STH
                GROUP BY STH.ThingID) as LST
            LEFT JOIN ThingRoomHistory as TRH on (LST."last" = TRH.ID)
        WHERE LST.thingid = idThing;
    if tmp is null
        then tmp = 1;
    end if;
    return tmp;
    END
$$;

create function findstudent(idthing integer) returns integer
    language plpgsql
as
$$
    DECLARE tmp int = -1;
    BEGIN
        SELECT studentid into tmp
            FROM (SELECT STH.ThingID, max(STH.ID) as "last"
                  FROM StudentThingHistory as STH
                   GROUP BY STH.ThingID) as LH
            LEFT JOIN StudentThingHistory as STH on (LH."last" = STH.ID)
        WHERE LH.thingid = idThing and STH.direction = 1;
    if tmp is null
        then tmp = -1;
    end if;
    return tmp;
    END
$$;

create function findstudentroom(idtstudent integer) returns integer
    language plpgsql
as
$$
    DECLARE tmp int = 0;
    BEGIN
        SELECT SRH.roomid into tmp
            FROM (SELECT SRH1.studentid, max(SRH1.ID) as "last"
                  FROM StudentRoomHistory as SRH1
                   GROUP BY SRH1.StudentID) as LH
            LEFT JOIN StudentRoomHistory as SRH on (LH."last" = SRH.ID)
        WHERE LH.studentid = idTStudent and SRH.direction = 1;
    if tmp is null
        then tmp = 0;
    end if;
    return tmp;
    END
$$;