CREATE EXTENSION moddatetime;

CREATE TABLE admin(
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  name varchar(30) NOT NULL,
  password varchar(128) NOT NULL,
  UNIQUE(name)
);

CREATE TRIGGER mdt_admin
  BEFORE UPDATE ON admin
  FOR EACH ROW
  EXECUTE PROCEDURE moddatetime (updated_at);

CREATE TABLE api_user(
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(), 
  name varchar(30) NOT NULL,
  secret varchar(128) NOT NULL,
  uses integer NOT NULL DEFAULT 0,
  active boolean DEFAULT TRUE,
  UNIQUE(name)
);

CREATE TRIGGER mdt_api_user
  BEFORE UPDATE ON api_user
  FOR EACH ROW
  EXECUTE PROCEDURE moddatetime (updated_at);

