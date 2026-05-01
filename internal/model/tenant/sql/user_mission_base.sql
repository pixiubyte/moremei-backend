-- Active: 1717141082295@@192.168.1.27@3306@moremei-ai-saas-tenant-01
CREATE TABLE user_mission_base(
    id BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    title VARCHAR(128) NOT NULL COMMENT '任务名',
    label VARCHAR(64) NOT NULL COMMENT '标签',
    points INT NOT NULL COMMENT '任务积分',
    single TINYINT NOT NULL COMMENT '单次任务',
    daily_limit INT NOT NULL COMMENT '每日积分限制',
    disabled TINYINT NOT NULL DEFAULT false COMMENT '是否禁用',
    description VARCHAR(512) COMMENT '任务描述',
    icon varchar(1024) NOT NULL COMMENT '图标',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT '任务列表';

-- 向表中插入基础数据
INSERT INTO user_mission_base (title, label, points, single, daily_limit, description, icon) VALUES
('关注公众号', 'first_follow', 50, 1, 50, '关注公众号', 'https://f.tongbotangkg.cn/static/mp/mine/icon-profile_task_follow.png'),
('首次健康评估', 'first_assessment', 100, 1, 100, '首次健康问卷评估', 'https://f.tongbotangkg.cn/static/mp/mine/icon-profile_task_first_health_data.png'),
('线下资料绑定', 'update_profile', 100, 1, 100, '填写个人资料', 'https://f.tongbotangkg.cn/static/mp/mine/icon-profile_task_link.png'),
('添加健康管理师', 'add_health_manager', 50, 1, 50, '添加健康管理师', 'https://f.tongbotangkg.cn/static/mp/mine/icon-profile_task_health_teacher.png'),
('每日签到', 'check_in', 10, 0, 10, '每日签到', 'https://f.tongbotangkg.cn/static/mp/mine/icon-profile_task_sign_in.png'),
('点赞文章', 'article_liked', 2, 0, 20, '每次点赞文章', 'https://f.tongbotangkg.cn/static/mp/mine/icon-profile_task_upvote.png'),
('阅读文章', 'article_read', 5, 0, 25, '每次阅读文章', 'https://f.tongbotangkg.cn/static/mp/mine/icon-profile_task_reading.png'),
('每日健康指标记录', 'health_metric_record', 10, 0, 0, '每日健康指标记录', 'https://f.tongbotangkg.cn/static/mp/mine/icon-profile_task_log.png'),
('和健康管理师互动', 'interact_with_manager', 20, 0, 0, '每日健康管理师互动', 'https://f.tongbotangkg.cn/static/mp/mine/icon-profile_task.png'),
('邀请好友加入', 'invite_friend', 20, 0, 0, '每邀请一位好友', 'https://f.tongbotangkg.cn/static/mp/mine/icon-profile_task_invite.png');
