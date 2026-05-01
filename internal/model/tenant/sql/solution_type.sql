-- 方案（报告）分类信息
CREATE TABLE `health_solution_type` (
    id BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '分类名称',
    `alias` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '分类别名',
    `key` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '分类键值',
    `type` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '分类类型',
    `icon` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '分类图标',
    `description` TEXT NOT NULL COMMENT '分类描述',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态, 0: 禁用, 1: 启用',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `key_type_key` (`key`,`type`)
) COMMENT '报告分类信息表';

INSERT INTO `health_solution_type` (`name`, `alias`, `key`, `type`, `icon`, `description`, `status`) VALUES
('问卷血压', '', 'blood-pressure', 'questionnaire', 'https://f.tongbotangkg.cn/static/mp/home/todolist/blood_pressure_36.png.png', '', 1),
('问卷血糖', '', 'blood-sugar', 'questionnaire', 'https://f.tongbotangkg.cn/static/mp/home/todolist/blood_glucose_36.png.png', '', 1),
('问卷尿酸', '', 'acid', 'questionnaire', 'https://f.tongbotangkg.cn/static/mp/home/todolist/UA_36.png.png', '', 1),
('问卷体脂', '', 'bmi', 'questionnaire', 'https://f.tongbotangkg.cn/static/mp/home/todolist/blood_pressure_36-1.png.png', '', 1),
('问卷血脂', '', 'blood-fat', 'questionnaire', 'https://f.tongbotangkg.cn/static/mp/home/todolist/blood_fat_36.png.png', '', 1),
('女性健康', '', 'female-reproduction', 'questionnaire', 'https://f.tongbotangkg.cn/static/mp/questionnaire/ic_quiz_entry_health_women_60.png', '', 1),
('糖尿病测评', '', 'diabetes-special', 'questionnaire', 'https://f.tongbotangkg.cn/static/mp/home/todolist/blood_glucose_36.png.png', '', 1),
('登记血压', '', 'blood-pressure', 'physical', 'https://f.tongbotangkg.cn/static/mp/home/todolist/blood_pressure_36.png.png', '', 1),
('登记血糖', '', 'blood-sugar', 'physical', 'https://f.tongbotangkg.cn/static/mp/home/todolist/blood_glucose_36.png.png', '', 1),
('登记尿酸', '', 'acid', 'physical', 'https://f.tongbotangkg.cn/static/mp/home/todolist/UA_36.png.png', '', 1),
('登记体脂', '', 'bmi', 'physical', 'https://f.tongbotangkg.cn/static/mp/home/todolist/blood_pressure_36-1.png.png', '', 1),
('登记血脂', '', 'blood-fat', 'physical', 'https://f.tongbotangkg.cn/static/mp/home/todolist/blood_fat_36.png.png', '', 1),
('膳食分析', '', 'diet', 'vision', 'https://f.tongbotangkg.cn/static/mp/vision/icon_rp_diet.png', '', 1),
('体检报告分析', '', 'exam', 'vision', 'https://f.tongbotangkg.cn/static/mp/vision/icon_rp_exam.png', '', 1);
