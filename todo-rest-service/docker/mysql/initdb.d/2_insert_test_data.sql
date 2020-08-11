-- todo_list table test data
INSERT INTO todo_list
  (id, implementation_date, due_date, todo_content, complete_flag, user_id)
VALUES
  (1, "2020-07-05", "2020-07-05", "今月の予算を立てる", true, "taira7"),
  (2, "2020-07-09", "2020-07-10", "コストコ鶏肉セール 5パック購入", true, "taira7"),
  (3, "2020-07-10", "2020-07-10", "電車定期券更新", true, "taira7"),
  (4, "2020-07-10", "2020-07-12", "醤油購入", false, "taira7"),
  (5, "2020-08-01", "2020-08-10", "水道代支払い", false, "taira7"),
  (6, "2020-08-01", "2020-08-15", "国保支払い", false, "taira7"),
  (7, "2020-07-20", "2020-07-20", "給料日 飲みに行く", true , "anraku8"),
  (8, "2020-07-21", "2020-07-22", "牛乳購入", false , "anraku8"),
  (9, "2020-07-25", "2020-07-25", "自分用におむつ購入", false , "anraku8"),
  (10, "2020-07-27", "2020-07-30", "牛肉2パック購入", false , "anraku8");

-- group_todo_list table test data
INSERT INTO group_todo_list
  (id, implementation_date, due_date, todo_content, complete_flag, user_id, group_id)
VALUES
  (1, "2020-07-05", "2020-07-05", "今月の予算を立てる", true, "taira7", 4),
  (2, "2020-07-09", "2020-07-10", "コストコ鶏肉セール 5パック購入", true, "taira7", 4),
  (3, "2020-07-10", "2020-07-12", "醤油購入", false, "taira7", 4),
  (4, "2020-07-20", "2020-07-20", "給料日 みんなで飲みに行く", true , "anraku8", 4),
  (5, "2020-07-21", "2020-07-22", "牛乳購入", false , "anraku8", 4),
  (6, "2020-07-27", "2020-07-30", "牛肉2パック購入", false , "anraku8", 4);
