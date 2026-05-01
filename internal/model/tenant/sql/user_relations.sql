CREATE TABLE IF NOT EXISTS user_relations (
    id BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    user_id  BIGINT(20) UNSIGNED NOT NULL COMMENT '用户ID',
    follower_id  BIGINT(20) UNSIGNED NOT NULL COMMENT '被关联者ID',
    relation_type VARCHAR(50) NOT NULL COMMENT '关系类型',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `user_relations_key` (`user_id`, `follower_id`),
    UNIQUE KEY `follower_id_key` (`follower_id`)
) COMMENT '用户关系表';
