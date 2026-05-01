-- Active: 1717141082295@@192.168.1.27@3306@moremei-ai-saas-tenant-01
CREATE TABLE user_points_level(
    id BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    level VARCHAR(32) NOT NULL COMMENT '等级',
    title VARCHAR(64) NOT NULL COMMENT '名称',
    min_points INT NOT NULL COMMENT '积分下限',
    max_points INT NOT NULL COMMENT '积分下限',
    icon VARCHAR(1024) COMMENT '图标',
    remark VARCHAR(1024) COMMENT "备注",
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT '用户积分等级';

INSERT INTO user_points_level (level, title, min_points, max_points, icon, remark) VALUES
('0', '酶动新芽', 0, 199, null, '预计一周'),
('1', '酶力萌芽', 200, 499, null, '预计2-3周'),
('2', '酶动初成', 500, 999, null, '预计1-2个月'),
('3', '酶力达人', 1000, 1999, null, '预计2-3个月'),
('4', '酶动专家', 2000, 3499, null, '预计4-5个月'),
('5', '酶动大师', 3500, 4999, null, '预计6-8个月'),
('6', '酶力宗师', 5000, 5999, null, '预计9-12个月'),
('7', '酶动传奇', 6000, 0, null, '预计1年以上');

