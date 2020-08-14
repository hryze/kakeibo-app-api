-- users table test data
INSERT INTO users
  (user_id, name, email, password)
VALUES
  ('taira', '平侑祐', 'taira@icloud.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK'),
  ('anraku', '安楽りょうすけ', 'anraku@icloud.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK'),
  ('ito', '伊藤りょうご', 'ito@icloud.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK'),
  ('furusawa', '古澤ひろや', 'furusawa@icloud.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK'),
  ('tati1', '館ひろし', 'tati@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('go2', '郷ひろみ', 'go@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('saigo3', '西郷隆盛', 'saigo@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('test4', 'test4', 'test4@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('test5', 'test5', 'test5@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('test6', 'test6', 'test6@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu');

-- group_users table test data
INSERT INTO group_names
  (group_name)
VALUES
  ('伊藤家家計簿'),
  ('平家家計簿1'),
  ('平家家計簿2'),
  ('平家家計簿3'),
  ('源家家計簿');

-- group_users table test data
INSERT INTO group_users
  (group_id, user_id)
VALUES
  (1, 'taira'),
  (1, 'tati1'),
  (1, 'test4'),
  (2, 'taira'),
  (2, 'go2'),
  (2, 'test5'),
  (2, 'saigo3'),
  (3, 'taira'),
  (3, 'test6'),
  (4, 'taira'),
  (4, 'anraku'),
  (4, 'ito'),
  (4, 'furusawa'),
  (5, 'taira'),
  (5, 'anraku'),
  (5, 'ito'),
  (5, 'furusawa'),
  (5, 'test4'),
  (5, 'test5');

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
  (5, 'tati1');
