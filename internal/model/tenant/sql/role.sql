CREATE TABLE `role`(
    id BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '角色名称',
    `role` VARCHAR(64) NOT NULL COMMENT '标识',
    `remark` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '备注',
    `status` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '状态, 0: 禁用, 1: 启用',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `role_key` (`role`)
) COMMENT '系统角色表';

INSERT INTO `role` (`name`, `role`, `remark`, `status`) VALUES ('酶好天使','mei_angel', '用户酶好天使身份', 1);