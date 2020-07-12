-- custom_categories table test data
INSERT INTO custom_categories
  (id, category_name, big_category_id, user_id)
VALUES
  (1, "調味料", 1, "taira7"),
  (2, "パン", 1, "taira7"),
  (3, "外食",1 , "taira7"),
  (4, "洗剤", 2, "taira7"),
  (5, "トイレットペーパー", 2, "taira7"),
  (6, "歯磨き粉", 2, "taira7"),
  (7, "調味料", 1, "anraku8"),
  (8, "パン", 1, "anraku8"),
  (9, "外食",1 , "anraku8"),
  (10, "洗剤", 2, "anraku8"),
  (11, "トイレットペーパー", 2, "anraku8"),
  (12, "歯磨き粉", 2, "anraku8"),
  (13, "おむつ", 2, "anraku8");


-- expenses table test data
INSERT INTO expenses
  (payment_date, payment_shop, memo, amount, user_id, big_category_id, medium_category_id, custom_category_id)
VALUES
  ("2020-07-01" ,"コストコ", "セールで牛肉購入", 4500, "taira7", 1, 1, NULL),
  ("2020-07-02" ,"ニトリ", "ベッド購入", 15000, "taira7", 2, 11, NULL),
  ("2020-07-02" , NULL, NULL, 1300, "taira7", 1, NULL, 3),
  ("2020-07-01" , NULL, "電車定期代", 12000, "taira7", 5, 28, NULL),
  ("2020-07-03" , NULL, NULL, 65000, "taira7", 10, 61, NULL),
  ("2020-07-04" , NULL, NULL, 500, "taira7", 1, 6, NULL),
  ("2020-07-05" , NULL, NULL, 4800, "taira7", 8, 46, NULL),
  ("2020-07-05" , NULL, "みんなのGo言語", 2500, "taira7", 9, 55, NULL),
  ("2020-07-06" ,"コンビニ" , NULL, 120, "taira7", 1, NULL, 2),
  ("2020-07-07" , NULL, "歯磨き粉3つ購入", 300, "taira7", 2, NULL, 6);
