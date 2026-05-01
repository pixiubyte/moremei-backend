CREATE TABLE article_category(
    id BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    name VARCHAR(128) NOT NULL COMMENT '分类名称',
    type VARCHAR(128) NOT NULL COMMENT '分类标签',
    status TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '状态: 0 禁用; 1 启用',
    description TEXT NOT NULL COMMENT '描述',
    extra JSON COMMENT '拓展字段',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY unique_name_type (name, type)
) COMMENT '文章分类类型表';

-- 文章分类
INSERT INTO article_category (name, type, status, description) VALUES
('酶好新闻', 'home-company-news', 1, '公司相关新闻'),
('健康科普', 'home', 1, '健康科普相关文章'),
('中医养生', 'home', 1, '中医养生相关文章'),
('女性健康', 'home', 1, '女性健康相关文章'),
('睡眠健康', 'home', 1, '睡眠健康相关文章'),
('体脂体重', 'home', 1, '体脂体重相关文章'),
('血糖健康', 'home', 1, '血糖健康相关文章'),
('血脂健康', 'home', 1, '血脂健康相关文章'),
('血压健康', 'home', 1, '血压健康相关文章'),
('心脏健康', 'home', 1, '心脏健康相关文章'),
('脾胃健康', 'home', 1, '脾胃健康相关文章'),
('情绪健康', 'home', 1, '情绪健康相关文章');
