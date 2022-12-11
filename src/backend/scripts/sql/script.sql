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
INSERT INTO users(userlogin, userpassword, userrole) VALUES ('artem', '123456', 1);
INSERT INTO users(userlogin, userpassword, userrole) VALUES ('sofa', '123456', 1);

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

INSERT INTO student(studentname, studentsurname, studentgroup, studentnumber, settledate, webaccid)
VALUES ('Александр', 'Прянишников', 'ИУ7-75Б', '19У609', current_date, 3);

INSERT INTO student(studentname, studentsurname, studentgroup, studentnumber, settledate, webaccid)
VALUES ('Артем', 'Богаченко', 'ИУ7-55Б', '18У712', current_date, 4);

INSERT INTO student(studentname, studentsurname, studentgroup, studentnumber, settledate, webaccid)
VALUES ('София', 'Шелия', 'ИУ7-65Б', '19У709', current_date, 5);

CREATE TABLE thing
(
    thingid SERIAL PRIMARY KEY,
    marknumber INT UNIQUE,
    creationdate DATE,
    thingtype TEXT
);

INSERT INTO thing(marknumber, creationdate, thingtype) VALUES (1, current_date, 'Шкаф');
INSERT INTO thing(marknumber, creationdate, thingtype) VALUES (2, current_date, 'Стул');
INSERT INTO thing(marknumber, creationdate, thingtype) VALUES (3, current_date, 'Стол');
INSERT INTO thing(marknumber, creationdate, thingtype) VALUES (4, current_date, 'Мультиварка');
INSERT INTO thing(marknumber, creationdate, thingtype) VALUES (5, current_date, 'Чайник');
INSERT INTO thing(marknumber, creationdate, thingtype) VALUES (6, current_date, 'Стол');

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

INSERT INTO studentroomhistory(studentid, roomid, direction, transferdate) VALUES (1, 2, 1, current_date);
INSERT INTO studentroomhistory(studentid, roomid, direction, transferdate) VALUES (1, 2, 0, current_date);
INSERT INTO studentroomhistory(studentid, roomid, direction, transferdate) VALUES (1, 2, 1, current_date);
INSERT INTO studentroomhistory(studentid, roomid, direction, transferdate) VALUES (2, 2, 1, current_date);
INSERT INTO studentroomhistory(studentid, roomid, direction, transferdate) VALUES (3, 3, 1, current_date);

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

INSERT INTO studentthinghistory(studentid, thingid, direction, transferdate) VALUES (1, 1, 1, current_date);
INSERT INTO studentthinghistory(studentid, thingid, direction, transferdate) VALUES (1, 1, 0, current_date);
INSERT INTO studentthinghistory(studentid, thingid, direction, transferdate) VALUES (1, 1, 1, current_date);
INSERT INTO studentthinghistory(studentid, thingid, direction, transferdate) VALUES (2, 3, 1, current_date);
INSERT INTO studentthinghistory(studentid, thingid, direction, transferdate) VALUES (1, 5, 1, current_date);
INSERT INTO studentthinghistory(studentid, thingid, direction, transferdate) VALUES (3, 6, 1, current_date);

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

INSERT INTO thingroomhistory(srcroomid, dstroomid, thingid, transferdate) VALUES (1, 2, 1, current_date);
INSERT INTO thingroomhistory(srcroomid, dstroomid, thingid, transferdate) VALUES (1, 2, 3, current_date);
INSERT INTO thingroomhistory(srcroomid, dstroomid, thingid, transferdate) VALUES (2, 3, 3, current_date);
INSERT INTO thingroomhistory(srcroomid, dstroomid, thingid, transferdate) VALUES (1, 3, 6, current_date);

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