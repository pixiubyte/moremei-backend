
INSERT INTO `physical_exam_item` (`types`, `name`, `field`, `unit`, `component`, `value_type`, `required`, value_options) VALUES
('blood-sugar', '血糖校准', 'calibration_value', 'mmol/L', '', 'float64', 0, NULL),
('blood-sugar', '血糖报警','alarm', '', '','string', 0, '[{"label":"降糖过快","value":"overly_reduced"},{"label":"升糖过快","value":"overly_rise"},{"label":"血糖过低","value":"low_glucose"}],{"label":"血糖过高","value":"high_glucose"}'),
('blood-sugar', '报警取消','cancel_alarm', '', '','string', 0, '[{"label":"降糖过快","value":"overly_reduced"},{"label":"升糖过快","value":"overly_rise"},{"label":"血糖过低","value":"low_glucose"}],{"label":"血糖过高","value":"high_glucose"}'),
('blood-sugar', '设备SN', 'device_sn', '', '', 'string', 0, NULL),
('blood-sugar', '是否生成体检报告', 'if_exam_advice', '', '', 'bool', 0, NULL);

-- system config, 导入酶好生活子公司列表
INSERT INTO `system_config` (`name`, `key`, `value_type`, `group`, `value`) VALUES
('酶好生活子公司列表', 'mei_subcompany_list','list', 'app', '["华北-酶好幸福", "华北-酶好道为", "华北-酶好大诚", "华东-酶好洋阳", "华东-酶好意祥", "华东-酶好凤呈", "华东-酶好印象", "华东-酶好时光", "华东-酶好吴尚", "华东-酶好安悦", "华中-酶好芳华", "华中-酶好蓉悦", "华中-酶好和礼", "华中-酶好玖悦", "东北-酶好盛世", "华南-酶好领悦", "西南-酶好天悦", "西南-酶好成悦", "西南-酶好原悦", "西北-酶好太享", "西北-酶好嘉悦", "西北-酶好宁悦", "其他-1部", "其他-7部", "其他-8部", "其他-9部", "其他-15部", "其他-19部", "其他-21部", "其他-26部"]');

-- questionnaire_mei_product, 新增方案字段
ALTER TABLE `questionnaire_mei_product` ADD COLUMN `scheme` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '方案';

-- questionnaire_mei_product, 新增方案数据, 
-- A: 清徽悦、木瓜咀嚼、沙棘VC片; B: 清徽悦、山柚虾青素、木瓜咀嚼、香橙味益生元; C: 清徽悦、山柚虾青素、木瓜咀嚼、清徽悦每美四臻丸;
INSERT INTO `questionnaire_mei_product` (`scheme`, `mei_product_id`, `questionnaire_id`) VALUES
('nxjk_ps_a', '16', '6'),
('nxjk_ps_a', '7', '6'),
('nxjk_ps_a', '10', '6'),
('nxjk_ps_b', '16', '6'),
('nxjk_ps_b', '13', '6'),
('nxjk_ps_b', '7', '6'),
('nxjk_ps_b', '14', '6'),
('nxjk_ps_c', '16', '6'),
('nxjk_ps_c', '13', '6'),
('nxjk_ps_c', '7', '6');
