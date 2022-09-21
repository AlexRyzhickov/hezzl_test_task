CREATE TABLE campaigns
(
    id   int PRIMARY KEY CHECK (id > 0),
    name varchar
);

INSERT INTO campaigns (id, name) VALUES (1, 'admin');

CREATE TABLE items
(
    id          int CHECK (id > 0),
    campaign_id int CHECK (campaign_id > 0),
    name        varchar,
    description varchar,
    priority    int,
    removed     bool,
    created_at  timestamp,
    PRIMARY KEY (id, campaign_id)
);