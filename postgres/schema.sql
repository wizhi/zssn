CREATE TABLE survivor (
    id text,
    name text,
    gender text,
    location point,
    flags integer DEFAULT 0,

    PRIMARY KEY (id)
);

CREATE TABLE resource (
    survivor_id text,
    kind text,
    quantity int,

    CONSTRAINT positive_quantity CHECK (quantity > 0),
    CONSTRAINT singular UNIQUE (survivor_id, kind),
    CONSTRAINT ownership FOREIGN KEY (survivor_id) REFERENCES survivor (id)
);
