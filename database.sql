CREATE KEYSPACE horaeapi WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};

CREATE TABLE company (id uuid, email text, name text, password text, city text, state text, PRIMARY KEY(email));
CREATE INDEX ON company (id);

CREATE TABLE contact (id uuid, company_id uuid, name text, email text, phone text, PRIMARY KEY(id));
CREATE INDEX ON contact (email);

CREATE TABLE calendar (id uuid, company_id uuid, contact_id uuid, date timestamp, duration timestamp, description text, value decimal, PRIMARY KEY(id));
