CREATE KEYSPACE horaeapi WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};

USE horaeapi;

CREATE TABLE company (id uuid, email text, name text, password text, city text, state text, PRIMARY KEY(email));
CREATE INDEX ON company (id);

CREATE TABLE contact (id uuid, company_id uuid, name text, email text, phone text, PRIMARY KEY(id));
CREATE INDEX ON contact (email);
CREATE INDEX ON contact (company_id);

CREATE TABLE calendar (id uuid, company_id uuid, contact_id uuid, start_at timestamp, end_at timestamp, description text, value float, PRIMARY KEY(id));
CREATE INDEX ON calendar (company_id);
CREATE INDEX ON calendar (contact_id);

ALTER TABLE contact ADD "token" int;