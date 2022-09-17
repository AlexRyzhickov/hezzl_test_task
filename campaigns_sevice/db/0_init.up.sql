CREATE TABLE campaigns (
    id   int PRIMARY KEY,
    name varchar
);

INSERT INTO campaigns (id, name) VALUES (0, 'admin');

CREATE TABLE items (
    id          int,
    campaign_id int,
    name        varchar,
    description varchar,
    priority    int,
    removed     bool,
    created_at  timestamp,
    PRIMARY KEY (id, campaign_id)
);