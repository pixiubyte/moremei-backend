CREATE TABLE `mei_health_expert`(
    id BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `avatar` VARCHAR(1024) NOT NULL COMMENT '头像',
    `name` VARCHAR(64) NOT NULL COMMENT '姓名',
    `profile` TEXT NOT NULL COMMENT '人物简介',
    `cover` VARCHAR(1024) NOT NULL COMMENT '封面图',
    `tags` JSON COMMENT '标签',
    `short_titles` JSON COMMENT '简称（博士、专家）',
    `titles` JSON COMMENT '头衔/职务（企业领域）',
    `research_roles` JSON COMMENT '学术界角色、头衔',
    `extra` JSON COMMENT '扩展字段',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT '健康专家';

INSERT INTO `mei_health_expert` (`avatar`, `name`, `profile`, `cover`, `tags`, `short_titles`, `titles`, `research_roles`) VALUES
('', '徐明', '酶好生活集团/禾呈美科学家\\n清徽悦抑菌片研发博士', '', null, '["医学博士", "药学专家"]', '[]', '["中南大学分子药物与治疗研究所创始所长","美国癌症研究所MSKCC 研究员","北京巴普国际医学研究院院长","润界智库管理委员会主席"]')
