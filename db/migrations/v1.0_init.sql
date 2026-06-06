CREATE EXTENSION IF NOT EXISTS pg_cron CASCADE;
--Main entities 
CREATE TABLE IF NOT EXISTS Companies(
    id_company SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(50) NOT NULL UNIQUE,
    address VARCHAR(100) NOT NULL,
    inn VARCHAR(12) NOT NULL UNIQUE,
    egrul VARCHAR(13) NOT NULL UNIQUE,
    description VARCHAR(500) NULL
);
CREATE TABLE IF NOT EXISTS Role_in_Company(
    id_role SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    id_company INT NULL REFERENCES Companies(id_company) ON DELETE CASCADE,
    is_creater BOOLEAN DEFAULT false
);
CREATE TABLE IF NOT EXISTS Rights(
    id_right SERIAL PRIMARY KEY,
    rus_name VARCHAR(50) UNIQUE not null,
    name VARCHAR(50) NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS Right_RoleInCompany(
    id_right INT REFERENCES Rights(id_right) ON DELETE CASCADE,
    id_role INT REFERENCES Role_in_Company(id_role) ON DELETE CASCADE,
    PRIMARY KEY(id_right,id_role)
);
CREATE TABLE IF NOT EXISTS Role(
    id_role SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS Users (
    id_user SERIAL PRIMARY KEY,
    login VARCHAR(30) UNIQUE NOT NULL,
	name VARCHAR(30) NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    id_company INT NOT NULL REFERENCES Companies(id_company) ON DELETE CASCADE,
    id_role_in_company INT NOT NULL REFERENCES Role_in_Company(id_role) ON DELETE CASCADE,
    id_role INT NOT NULL REFERENCES Role(id_role) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS Districts(
    id_district SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS Statuses(
    id_status SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS Categories(
    id_category SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS Category_Links(
    id_parent INT REFERENCES Categories(id_category) ON DELETE CASCADE,
    id_children INT REFERENCES Categories(id_category),
    PRIMARY KEY(id_parent,id_children)
);
CREATE TABLE IF NOT EXISTS Tenders(
    id_tender SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(500) NULL,
    datetime_start TIMESTAMP NOT NULL,
    datetime_end TIMESTAMP NOT NULL,
    id_status INT NOT NULL REFERENCES Statuses(id_status) ON DELETE CASCADE,
    id_company INT NOT NULL REFERENCES Companies(id_company) ON DELETE CASCADE,
    id_district INT NOT NULL REFERENCES Districts(id_district) ON DELETE CASCADE   
);
CREATE TABLE IF NOT EXISTS Tender_Category(
    id_tender INT REFERENCES Tenders(id_tender) ON DELETE CASCADE,
    id_category INT REFERENCES Categories(id_category) ON DELETE CASCADE,
    PRIMARY KEY(id_tender,id_category)
);
CREATE TABLE IF NOT EXISTS Offers(
    id_offer SERIAL PRIMARY KEY,
    price FLOAT NOT NULL,
    description VARCHAR(500) NULL,
    datetime_create TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    id_status INT NOT NULL REFERENCES Statuses(id_status) ON DELETE CASCADE,
    id_company INT NOT NULL REFERENCES Companies(id_company) ON DELETE CASCADE,
    id_tender INT NOT NULL REFERENCES Tenders(id_tender) ON DELETE CASCADE  
);
CREATE TABLE IF NOT EXISTS Docs(
    id_doc SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    filename VARCHAR(250) NOT NULL UNIQUE,
    description VARCHAR(500) NULL
);
CREATE TABLE IF NOT EXISTS Doc_Tender(
    id_doc INT REFERENCES Docs(id_doc) ON DELETE CASCADE,
    id_tender INT REFERENCES Tenders(id_tender) ON DELETE CASCADE,
    PRIMARY KEY(id_doc,id_tender)
);
CREATE TABLE IF NOT EXISTS Doc_Company(
    id_doc INT REFERENCES Docs(id_doc) ON DELETE CASCADE,
    id_company INT REFERENCES Companies(id_company) ON DELETE CASCADE,
    PRIMARY KEY(id_doc,id_company)
);
CREATE TABLE IF NOT EXISTS Doc_Offer(
    id_doc INT REFERENCES Docs(id_doc) ON DELETE CASCADE,
    id_offer INT REFERENCES Offers(id_offer) ON DELETE CASCADE,
    PRIMARY KEY(id_doc,id_offer)
);
--Tecnical entities
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL UNIQUE,
    id_user INT NOT NULL REFERENCES Users(id_user) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS reset_tokens (
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL UNIQUE,
    id_user INT NOT NULL REFERENCES Users(id_user) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS Type_Operation(
    id_type SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS Entities(
    id_entity SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS Logs(
    id_log SERIAL PRIMARY KEY,
    id_user INT NOT NULL REFERENCES Users(id_user),
    id_type INT NOT NULL REFERENCES Type_Operation(id_type) ON DELETE CASCADE,
    id_entity INT NOT NULL REFERENCES Entities(id_entity) ON DELETE CASCADE,
    datetime_create TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS Log_Tender(
    id_log INT REFERENCES Logs(id_log) ON DELETE CASCADE,
    id_tender INT REFERENCES Tenders(id_tender),
    PRIMARY KEY(id_log,id_tender)
);
CREATE TABLE IF NOT EXISTS Log_Company(
    id_log INT REFERENCES Logs(id_log) ON DELETE CASCADE,
    id_company INT REFERENCES Companies(id_company),
    PRIMARY KEY(id_log,id_company)
);
CREATE TABLE IF NOT EXISTS Log_Offer(
    id_log INT REFERENCES Logs(id_log) ON DELETE CASCADE,
    id_offer INT REFERENCES Offers(id_offer),
    PRIMARY KEY(id_log,id_offer)
);
CREATE TABLE IF NOT EXISTS Log_Doc(
    id_log INT REFERENCES Logs(id_log) ON DELETE CASCADE,
    id_doc INT REFERENCES Docs(id_doc),
    PRIMARY KEY(id_log,id_doc)
);