-- users table test data
INSERT INTO users
  (user_id, name, email, password)
VALUES
  ('taira7', '平侑祐', 'taira@icloud.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK'),
  ('tati1', '館ひろし', 'tati@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('go2', '郷ひろみ', 'go@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('saigo3', '西郷隆盛', 'saigo@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('test4', 'test4', 'test4@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('test5', 'test5', 'test5@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('test6', 'test6', 'test6@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('anraku8', '安楽亮佑', 'anraku@icloud.com', '$2a$10$6KihsMbh9tfCJEJcgo4fqujK.gbbGaq/v7YY9GZF7snhcMsGRQoVK');

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
  (1, 'taira7'),
  (1, 'tati1'),
  (1, 'test4'),
  (2, 'taira7'),
  (2, 'go2'),
  (2, 'test5'),
  (2, 'saigo3'),
  (3, 'taira7'),
  (3, 'test6'),
  (4, 'taira7'),
  (4, 'anraku8'),
  (5, 'tati1');

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
  (2, 'anraku8'),
  (3, 'anraku8'),
  (4, 'tati1'),
  (5, 'tati1');
