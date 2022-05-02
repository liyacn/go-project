CREATE DATABASE go_project DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE go_project;

CREATE TABLE `user` (
    id bigint AUTO_INCREMENT PRIMARY KEY,
    openid varchar(50) NOT NULL UNIQUE,
    unionid varchar(50) NOT NULL DEFAULT '',
    phone_number varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
    nickname varchar(10) NOT NULL DEFAULT '' COMMENT '昵称',
    avatar_url varchar(150) NOT NULL DEFAULT '' COMMENT '头像链接',
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY (phone_number)
) COMMENT='用户信息';

CREATE TABLE system_action (
    id bigint AUTO_INCREMENT PRIMARY KEY,
    key_name varchar(64) NOT NULL UNIQUE,
    title varchar(16) NOT NULL DEFAULT '' COMMENT '备注标题',
    sort int NOT NULL DEFAULT 0 COMMENT '排序值0~9999',
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) COMMENT='系统接口';

CREATE TABLE admin_role (
    id bigint AUTO_INCREMENT PRIMARY KEY,
    name varchar(32) NOT NULL DEFAULT '',
    actions json NOT NULL DEFAULT ('[]') COMMENT '数组保存system_action.key_name',
    menus json NOT NULL DEFAULT ('[]') COMMENT '数组保存前端菜单的唯一标识',
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) COMMENT='管理员角色';

CREATE TABLE admin_user (
    id bigint AUTO_INCREMENT PRIMARY KEY,
    username varchar(32) NOT NULL UNIQUE,
    password varchar(64) NOT NULL DEFAULT '',
    pwd_exp int NOT NULL DEFAULT (UNIX_TIMESTAMP()+7862400) COMMENT '密码过期时间戳',
    role_id int NOT NULL DEFAULT 0 COMMENT 'admin_role.id;-1:super',
    status tinyint NOT NULL DEFAULT 0 COMMENT '1:enabled,2:disabled',
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY(role_id)
) COMMENT='管理员账号';

INSERT INTO admin_user (username,password,role_id,status) VALUES
('super','84BHqjrPHK-07b8xirALIuU7iwjQgzdov-dGqlfWegFueqaSLarejC1bN9Fw',-1,1),
('admin','84BHqjrPHK-07b8xirALIuU7iwjQgzdov-dGqlfWegFueqaSLarejC1bN9Fw',-1,1);
-- 超管账号，初始密码: 123456
