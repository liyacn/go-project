CREATE DATABASE go_project;
SET search_path TO go_project;

CREATE OR REPLACE FUNCTION on_update()
RETURNS TRIGGER AS $$
BEGIN
    IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
        NEW.updated_at = CURRENT_TIMESTAMP;
        RETURN NEW;
    ELSE
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE "user" (
    id serial PRIMARY KEY,
    openid character varying(50) NOT NULL UNIQUE,
    unionid character varying(50) NOT NULL DEFAULT '',
    phone_number character varying(20) NOT NULL DEFAULT '',
    nickname character varying(10) NOT NULL DEFAULT '',
    avatar_url character varying(150) NOT NULL DEFAULT '',
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON "user" (phone_number);
CREATE TRIGGER user_update BEFORE UPDATE ON "user" FOR EACH ROW EXECUTE FUNCTION on_update();

CREATE TABLE system_action (
    id serial PRIMARY KEY,
    key_name character varying(64) NOT NULL UNIQUE,
    title character varying(16) NOT NULL DEFAULT '',
    sort integer NOT NULL DEFAULT 0,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER system_action_update BEFORE UPDATE ON system_action FOR EACH ROW EXECUTE FUNCTION on_update();

CREATE TABLE system_config (
    id serial PRIMARY KEY,
    key_name character varying(64) NOT NULL UNIQUE,
    content jsonb,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER system_config_update BEFORE UPDATE ON system_config FOR EACH ROW EXECUTE FUNCTION on_update();
INSERT INTO system_config (key_name,content) VALUES
('sys_menu_keys', '[]'),
('sys_menu_trees','[]');

CREATE TABLE admin_role (
    id serial PRIMARY KEY,
    name character varying(32) NOT NULL DEFAULT '',
    actions jsonb,
    menus jsonb,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER admin_role_update BEFORE UPDATE ON admin_role FOR EACH ROW EXECUTE FUNCTION on_update();

CREATE TABLE admin_user (
    id serial PRIMARY KEY,
    username character varying(32) NOT NULL UNIQUE,
    password character varying(64) NOT NULL DEFAULT '',
    pwd_exp integer NOT NULL DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::integer + 7862400),
    role_id integer NOT NULL DEFAULT 0,
    status smallint NOT NULL DEFAULT 1,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON admin_user (role_id);
CREATE TRIGGER admin_user_update BEFORE UPDATE ON admin_user FOR EACH ROW EXECUTE FUNCTION on_update();

INSERT INTO admin_user (username,password,role_id) VALUES
('super','84BHqjrPHK-07b8xirALIuU7iwjQgzdov-dGqlfWegFueqaSLarejC1bN9Fw',-1),
('admin','84BHqjrPHK-07b8xirALIuU7iwjQgzdov-dGqlfWegFueqaSLarejC1bN9Fw',-1);
-- 受保护的超管账号，初始密码: 123456
