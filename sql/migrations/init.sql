BEGIN;


CREATE TABLE status (
    service TEXT PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    status TEXT NOT NULL
);

CREATE TABLE snapshots (
  service TEXT PRIMARY KEY,
  timestamp TIMESTAMP NOT NULL,
  total_data_points INT NOT NULL,
  down_data_points INT NOT NULL,
  uptime_percentage FLOAT NOT NULL
);

COMMIT;