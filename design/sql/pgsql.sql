CREATE DATABASE go_project;
SET search_path TO go_project;

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

CREATE TABLE system_action (
    id serial PRIMARY KEY,
    key_name character varying(64) NOT NULL UNIQUE,
    title character varying(16) NOT NULL DEFAULT '',
    sort integer NOT NULL DEFAULT 0,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE admin_role (
    id serial PRIMARY KEY,
    name character varying(32) NOT NULL DEFAULT '',
    actions jsonb,
    menus jsonb,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

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

INSERT INTO admin_user (username,password,role_id) VALUES
('super','84BHqjrPHK-07b8xirALIuU7iwjQgzdov-dGqlfWegFueqaSLarejC1bN9Fw',-1),
('admin','84BHqjrPHK-07b8xirALIuU7iwjQgzdov-dGqlfWegFueqaSLarejC1bN9Fw',-1);
-- 受保护的超管账号，初始密码: 123456
