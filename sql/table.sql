-- create database hinagane_db;

USE hinagane_db;

DROP TABLE member_counter;
DROP TABLE word_counter;
DROP TABLE article_member_link;
DROP TABLE article_word_link;
DROP TABLE member;
DROP TABLE article;
DROP TABLE site;
DROP TABLE word;

CREATE TABLE member (
    id INT AUTO_INCREMENT NOT NULL,
    name VARCHAR(30) NOT NULL,
    PRIMARY KEY(id)
)ENGINE=INNODB DEFAULT CHARSET=utf8;

CREATE TABLE article (
    id INT AUTO_INCREMENT NOT NULL,
    name VARCHAR(1000) NOT NULL,
    url VARCHAR(1000) NOT NULL,
    date_time DATETIME NOT NULL,
    site_id INT NOT NULL,
    pic_url VARCHAR(1000) NOT NULL,
    PRIMARY KEY(id)
)ENGINE=INNODB DEFAULT CHARSET=utf8;

CREATE TABLE site (
    id INT AUTO_INCREMENT NOT NULL,
    name VARCHAR(50) NOT NULL,
    url VARCHAR(300) NOT NULL,
    PRIMARY KEY(id)
)ENGINE=INNODB DEFAULT CHARSET=utf8;

CREATE TABLE word (
    id INT AUTO_INCREMENT NOT NULL,
    name VARCHAR(50) NOT NULL,
    PRIMARY KEY(id)
)ENGINE=INNODB DEFAULT CHARSET=utf8;

CREATE TABLE article_member_link (
    article_id INT NOT NULL,
    member_id INT NOT NULL,
    PRIMARY KEY(article_id,member_id),
    FOREIGN KEY(article_id)
    REFERENCES article(id),
    FOREIGN KEY(member_id)
    REFERENCES member(id)
)ENGINE=INNODB DEFAULT CHARSET=utf8;

CREATE TABLE member_counter (
    member_id INT NOT NULL,
    date_time DATETIME NOT NULL,
    counter INT NULL,
    PRIMARY KEY(member_id,date_time),
    FOREIGN KEY(member_id)
    REFERENCES member(id)
)ENGINE=INNODB DEFAULT CHARSET=utf8;

CREATE TABLE article_word_link (
    article_id INT NOT NULL,
    word_id INT NOT NULL,
    PRIMARY KEY(article_id,word_id),
    FOREIGN KEY(article_id)
    REFERENCES article(id),
    FOREIGN KEY(word_id)
    REFERENCES word(id)
)ENGINE=INNODB DEFAULT CHARSET=utf8;

CREATE TABLE word_counter (
    word_id INT NOT NULL,
    date_time DATETIME NOT NULL,
    counter INT NULL,
    PRIMARY KEY(word_id,date_time),
    FOREIGN KEY(word_id)
    REFERENCES word(id)
)ENGINE=INNODB DEFAULT CHARSET=utf8;

INSERT INTO member(name) VALUES('井口眞緒');
INSERT INTO member(name) VALUES('潮紗理菜');
INSERT INTO member(name) VALUES('影山優佳');
INSERT INTO member(name) VALUES('加藤史帆');
INSERT INTO member(name) VALUES('齊藤京子');
INSERT INTO member(name) VALUES('佐々木久美');
INSERT INTO member(name) VALUES('佐々木美玲');
INSERT INTO member(name) VALUES('高本彩花');
INSERT INTO member(name) VALUES('東村芽依');
INSERT INTO member(name) VALUES('金村美玖');
INSERT INTO member(name) VALUES('河田陽菜');
INSERT INTO member(name) VALUES('小坂菜緒');
INSERT INTO member(name) VALUES('富田鈴花');
INSERT INTO member(name) VALUES('丹生明里');
INSERT INTO member(name) VALUES('濱岸ひより');
INSERT INTO member(name) VALUES('松田好花');
INSERT INTO member(name) VALUES('渡邉美穂');
INSERT INTO member(name) VALUES('上村ひなの');
INSERT INTO member(name) VALUES('高橋未来虹');
INSERT INTO member(name) VALUES('森本茉莉');
INSERT INTO member(name) VALUES('山口陽世');
INSERT INTO member(name) VALUES('柿崎芽実');

INSERT INTO site(name,url) VALUES('日向坂46まとめ速報','http://hiraganakeyaki.blog.jp/');
INSERT INTO site(name,url) VALUES('日向坂46まとめきんぐだむ','http://hiragana46matome.com/');
INSERT INTO site(name,url) VALUES('日向速報','http://hinatasoku.blog.jp/');