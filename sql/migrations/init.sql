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
CREATE TABLE accounts (
  id UUID PRIMARY KEY,
  external_id TEXT
);

CREATE TABLE users (
    id UUID PRIMARY KEY, -- internal anayltics id
    external_id TEXT, -- id from user application
    created_at TIMESTAMP,
    email TEXT,
    name TEXT,
    properties JSON  -- any json
);
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    device TEXT,
    os TEXT,
    browser TEXT,
    location TEXT,
    user_agent TEXT,
    properties JSON
);
CREATE TABLE projects (
  id UUID PRIMARY KEY,
  url TEXT,
  data_tags JSON
);

CREATE TABLE events (
  id UUID PRIMARY KEY,
  project_id UUID REFERENCES projects(id),
  -- user_id UUID REFERENCES users(id),
  -- session_id UUID REFERENCES sessions(id),
  name TEXT NOT NULL, -- click, submit etc..
  timestamp TIMESTAMP NOT NULL,
  metadata JSON -- store json like info about btn curr page
);

CREATE TABLE page_views (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  session_id UUID REFERENCES sessions(id),
  url TEXT,
  title TEXT,
  timestamp TIMESTAMP,
  duration_seconds INT,
  referrer TEXT, -- where did the user come from
  metadata JSON -- store any extra info heere
);

COMMIT;