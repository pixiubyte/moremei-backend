
CREATE TABLE `user_archive_level`(
    `id` BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `level` VARCHAR(32) NOT NULL COMMENT '等级',
    `title` VARCHAR(64) NOT NULL COMMENT '名称',
    `title_sub` VARCHAR(64) NOT NULL COMMENT '副标题',
    `min_score` DECIMAL(5,2) NOT NULL COMMENT '积分下限',
    `max_score` DECIMAL(5,2) NOT NULL COMMENT '积分下限',
    `icon` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '图标',
    `background` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '背景图',
    `text` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '提醒内容',
    `extra` JSON COMMENT "拓展字段",
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT '用户档案等级';

INSERT INTO `user_archive_level` (`level`, `title`, `title_sub`, `min_score`, `max_score`, `text`, `icon`, `background`) VALUES
('0', '关护', '健康指数', 0, 30, '健康状况需要特别关护', 'https://f.tongbotangkg.cn/static/mp/home/bg_health_badge_one.png', 'https://f.tongbotangkg.cn/static/mp/home/icon_badge_health_one.png'),
('1', '调理', '健康指数', 31, 59, '健康需要进一步调理', 'https://f.tongbotangkg.cn/static/mp/home/bg_health_badge_two.png', 'https://f.tongbotangkg.cn/static/mp/home/icon_badge_health_two.png'),
('2', '稳健', '健康指数', 60, 74, '健康基本稳固', 'https://f.tongbotangkg.cn/static/mp/home/bg_health_badge_three.png', 'https://f.tongbotangkg.cn/static/mp/home/icon_badge_health_three.png'),
('3', '良好', '健康指数', 75, 89, '健康状态良好', 'https://f.tongbotangkg.cn/static/mp/home/bg_health_badge_four.png', 'https://f.tongbotangkg.cn/static/mp/home/icon_badge_health_four.png'),
('4', '优秀', '健康指数', 90, 100, '健康状况非常理想', 'https://f.tongbotangkg.cn/static/mp/home/bg_health_badge_five.png', 'https://f.tongbotangkg.cn/static/mp/home/icon_badge_health_five.png');


INSERT INTO `system_config` (`name`, `key`, `value_type`, `group`, `value`) VALUES
('banner排行卡片背景', 'content_bc_rank_bg','string', 'content_conf', 'https://f.tongbotangkg.cn/static/mp/home/bg_ranking.png'),
('banner排行卡片头图', 'content_bc_rank_cover','string', 'content_conf', 'https://f.tongbotangkg.cn/static/mp/home/icon_ranking.png '),
('banner天使卡片背景', 'content_bc_mei_bg','string', 'content_conf', 'https://f.tongbotangkg.cn/static/mp/home/bg_partner.png'),
('banner天使卡片默认徽章', 'content_bc_mei_st_cover','string', 'content_conf', 'https://f.tongbotangkg.cn/static/mp/home/icon_partner_empty.png'),
('banner天使卡片天使徽章', 'content_bc_mei_angel_st_cover','string', 'content_conf', 'https://f.tongbotangkg.cn/static/mp/home/icon_partner.png'),
('banner每日金句卡片背景', 'content_bc_daily_words_Bg','string', 'content_conf', 'https://f.tongbotangkg.cn/static/mp/home/bg_daily_words.png'),
('banner每日金句卡片头像', 'content_bc_daily_words_cover','string', 'content_conf', 'https://f.tongbotangkg.cn/static/mp/home/icon_zbs.png'),
('aiManager卡片默认背景', 'content_bc_ai_manager_default','string', 'content_conf', 'https://f.tongbotangkg.cn/static/mp/home/add/health_manager_ai_default.png'),
('aiManager卡片默认背景', 'content_bc_ai_manager','string', 'content_conf', 'https://f.tongbotangkg.cn/static/mp/home/add/health_manager_ai.png')
;

INSERT INTO `user_contents_template_modules`(`type`,`class`,`content`,`is_hide`,`is_public`) VALUES
('home_page','home_banner','{"type":"banner","sort":0,"contents":[{"type":"banner_card_tp1","sort":0},{"type":"banner_card_tp2","sort":1},{"type":"banner_card_tp3","sort":2},{"type":"banner_card_tp4","sort":3}]}',0,1),
('home_page','home_to_yihao','{"type": "easy_card","sort":10, "title_hide": true, "contents": [{"type":"image_card","sort": 0, "background": "https://f.tongbotangkg.cn/static/mp/home/bg_yihao_link.png"}], "more_hide": true}',0,1),
('home_page','health_special','{"type":"func_list","sort":20,"title":"健康专题与评测","more":"硬件测一测","more_to":"","contents":[{"type":"health_special_card","sort":0,"title":"nxjk"},{"type":"health_special_card","sort":1,"title":"nxjk"}]}',0,1),
('home_page','health_funcs','{"type":"func_list","sort":40,"title_hide":true,"more_hide":true,"contents":[{"type":"vital_card","sort":0,"title":""},{"type":"health_special_card","sort":1,"title":"nxjk"},{"type":"todo_card","sort":2,"title":""}]}',0,1),
('home_page','product_recommend','{"type":"h_scroll_ai","sort":50,"title":"健康保健产品","more_hide":true,"contents":[]}',0,1),
('home_page','article_meihao','{"type":"h_scroll_content","sort":60,"title":"酶好新闻","more_hide":true,"contents":[]}',0,1),
('home_page','articles','{"type":"articles","sort":70,"title":"健康保健","more_hide":true,"contents":[]}',0,1),
('home_page','ai_manager','{"sort": 30, "type": "ai_manager", "title": "健康管理", "contents": [{"text": "多多锻炼，注意饮食健康，继续加油！", "cover": "", "main_bt": "问问AI", "background": "https://f.tongbotangkg.cn/static/mp/home/add/health_manager_ai.png", "main_bt_bg": "https://f.tongbotangkg.cn/static/mp/home/add/health_manager_ai_btn_bg.png", "main_bt_to": "/pages/chat/chat", "main_bt_icon": "https://f.tongbotangkg.cn/static/mp/common/icon-ai.png"}], "more_hide": true}',0,1)
;

UPDATE `mei_health_special` SET `extra` = '{"tags": ["生理周期调理", "心里健康管理"], "cover": "https://f.tongbotangkg.cn/static/mp/home/ic_hhealth_card_3.png", "theme": "#D33340", "title": "女性健康专题", "main_bt": "进入专题", "title_sub": "关爱女性健康", "background": "https://f.tongbotangkg.cn/static/mp/home/health_card_bg_04.png", "planable_product": "82f70f450c761604e223e639e3c988eb"}' WHERE `type` = 'nxjk';
