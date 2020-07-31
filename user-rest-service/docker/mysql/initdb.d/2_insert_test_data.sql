-- users table test data
INSERT INTO users
  (user_id, name, email, password)
VALUES
  ('tati1', '館ひろし', 'tati@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('go2', '郷ひろみ', 'go@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('saigo3', '西郷隆盛', 'saigo@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('test4', 'test4', 'test4@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('test5', 'test5', 'test5@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu'),
  ('test6', 'test6', 'test6@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu');

-- group_names table test data
INSERT INTO group_names
  (id, group_name)
VALUES
  (1, '館家'),
  (2, '郷家'),
  (3, '西郷家');

-- group_users table test data
INSERT INTO group_users
  (group_id, user_id)
VALUES
  (1, 'tati1'),
  (2, 'go2'),
  (3, 'saigo3');

-- group_unapproved_users table test data
INSERT INTO group_unapproved_users
  (group_id, user_id)
VALUES
  (1, 'test4'),
  (2, 'test5'),
  (3, 'test6');
