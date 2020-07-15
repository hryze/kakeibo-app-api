-- custom_categories table test data
INSERT INTO custom_categories
  (id, category_name, big_category_id, user_id)
VALUES
  (1, "調味料", 2, "taira7"),
  (2, "パン", 2, "taira7"),
  (3, "外食", 2, "taira7"),
  (4, "洗剤", 3, "taira7"),
  (5, "トイレットペーパー", 3, "taira7"),
  (6, "歯磨き粉", 3, "taira7"),
  (7, "調味料", 2, "anraku8"),
  (8, "パン", 2, "anraku8"),
  (9, "外食", 2, "anraku8"),
  (10, "洗剤", 3, "anraku8"),
  (11, "トイレットペーパー", 3, "anraku8"),
  (12, "歯磨き粉", 3, "anraku8"),
  (13, "おむつ", 3, "anraku8"),
  (14 , "株配当金", 1, "taira7"),
  (15 , "株配当金", 1, "anraku8");


-- expenses table test data
INSERT INTO transactions
  (transaction_type, transaction_date, shop, memo, amount, user_id, big_category_id, medium_category_id, custom_category_id)
VALUES
  ("expense", "2020-07-01", "コストコ", "セールで牛肉購入", 4500, "taira7", 2, 6, NULL),
  ("expense", "2020-07-02", "ニトリ", "ベッド購入", 15000, "taira7", 3, 16, NULL),
  ("expense", "2020-07-02", NULL, NULL, 1300, "taira7", 2, NULL, 3),
  ("expense", "2020-07-01", NULL, "電車定期代", 12000, "taira7", 6, 33, NULL),
  ("expense", "2020-07-03", NULL, NULL, 65000, "taira7", 11, 66, NULL),
  ("expense", "2020-07-04", NULL, NULL, 500, "taira7", 2, 11, NULL),
  ("expense", "2020-07-05", NULL, NULL, 4800, "taira7", 8, 49, NULL),
  ("expense", "2020-07-05", NULL, "みんなのGo言語", 2500, "taira7", 10, 60, NULL),
  ("expense", "2020-07-06", "コンビニ", NULL, 120, "taira7", 2, NULL, 2),
  ("expense", "2020-07-07", NULL, "歯磨き粉3つ購入", 300, "taira7", 3, NULL, 6),
  ("income", "2020-07-10", NULL, "給料日", 450000, "taira7", 1, 1, NULL),
  ("income", "2020-07-20", NULL, "賞与", 1000000, "taira7", 1, 2, NULL),
  ("income", "2020-07-20", NULL, "株配当金", 200000, "taira7", 1, NULL, 14),
  ("expense", "2020-07-01", "コストコ", "セールで牛肉購入", 4500, "anraku8", 2, 6, NULL),
  ("expense", "2020-07-02", "ニトリ", "ベッド購入", 15000, "anraku8", 3, 16, NULL),
  ("expense", "2020-07-02", NULL, "醤油", 1300, "anraku8", 2, NULL, 7),
  ("expense", "2020-07-01", NULL, "電車定期代", 12000, "anraku8", 6, 33, NULL),
  ("expense", "2020-07-03", NULL, NULL, 65000, "anraku8", 11, 66, NULL),
  ("expense", "2020-07-04", NULL, NULL, 500, "anraku8", 2, 11, NULL),
  ("expense", "2020-07-05", NULL, "携帯", 4800, "anraku8", 9, 51, NULL),
  ("expense", "2020-07-05", NULL, "React参考書", 2500, "anraku8", 10, 60, NULL),
  ("expense", "2020-07-06", "クリエイト" , NULL, 340, "anraku8", 3, NULL, 11),
  ("expense", "2020-07-07", NULL, "自分用におむつ3つ購入", 1200, "anraku8", 3, NULL, 13),
  ("income", "2020-07-10", NULL, "給料日", 140000, "anraku8", 1, 1, NULL),
  ("income", "2020-07-20", NULL, "賞与", 30000, "anraku8", 1, 2, NULL);
