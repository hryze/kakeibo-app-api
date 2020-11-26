-- users table test data
INSERT INTO users
  (user_id, name, email, password)
VALUES
  ('taira', '平侑祐', 'taira@icloud.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK'),
  ('anraku', '安楽りょうすけ', 'anraku@icloud.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK'),
  ('ito', '伊藤りょうご', 'ito@icloud.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK'),
  ('furusawa', '古澤ひろや', 'furusawa@icloud.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK'),
  ('tati1', '館ひろし', 'tati@developer.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK'),
  ('go2', '郷ひろみ', 'go@developer.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK'),
  ('saigo3', '西郷隆盛', 'saigo@developer.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK'),
  ('test4', 'test4', 'test4@developer.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK'),
  ('test5', 'test5', 'test5@developer.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK'),
  ('test6', 'test6', 'test6@developer.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK');

-- group_users table test data
INSERT INTO group_names
  (group_name)
VALUES
  ('伊藤家家計簿'),
  ('平家家計簿1'),
  ('平家家計簿2'),
  ('平家家計簿3'),
  ('源家家計簿'),
  ('安楽家家計簿');

-- group_users table test data
INSERT INTO group_users
  (group_id, user_id, color_code)
VALUES
  (1, 'taira', '#FF0000'),
  (1, 'tati1', '#00FFFF'),
  (1, 'test4', '#80FF00'),
  (2, 'taira', '#FF0000'),
  (2, 'go2', '#00FFFF'),
  (2, 'test5', '#80FF00'),
  (2, 'saigo3', '#8000FF'),
  (3, 'taira', '#FF0000'),
  (3, 'test6', '#00FFFF'),
  (4, 'taira', '#FF0000'),
  (4, 'anraku', '#00FFFF'),
  (4, 'ito', '#80FF00'),
  (4, 'furusawa', '#8000FF'),
  (4, 'test4', '#FF8000'),
  (4, 'test5', '#0080FF'),
  (4, 'test6', '#00FF80'),
  (5, 'taira', '#FF0000'),
  (5, 'anraku', '#00FFFF'),
  (5, 'ito', '#80FF00'),
  (5, 'furusawa', '#8000FF'),
  (5, 'test4', '#FF8000'),
  (5, 'test5', '#0080FF'),
  (6, 'anraku', '#FF0000');

-- group_unapproved_users table test data
INSERT INTO group_unapproved_users
  (group_id, user_id)
VALUES
  (1, 'go2'),
  (1, 'saigo3'),
  (1, 'test5'),
  (1, 'test6'),
  (2, 'test4'),
  (2, 'tati1'),
  (2, 'anraku'),
  (3, 'anraku'),
  (4, 'tati1'),
  (5, 'tati1'),
  (6, 'taira');
