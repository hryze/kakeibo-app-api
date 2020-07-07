-- big_categories table default data
INSERT INTO big_categories
  (id, category_name)
VALUES
  (1, "食費"),
  (2, "日用品"),
  (3, "趣味・娯楽"),
  (4, "交際費"),
  (5, "交通費"),
  (6, "衣服・美容"),
  (7, "健康・医療"),
  (8, "通信費"),
  (9, "教養・教育"),
  (10, "住宅"),
  (11, "水道・光熱費"),
  (12, "自動車"),
  (13, "保険"),
  (14, "税金・社会保険"),
  (15, "現金・カード"),
  (16, "その他");

-- medium_categories table default data
INSERT INTO medium_categories
  (id, category_name, big_category_id)
VALUES
  (1, "食料品", 1),
  (2, "朝食", 1),
  (3, "昼食", 1),
  (4, "夕食", 1),
  (5, "外食", 1),
  (6, "カフェ", 1),
  (7, "その他食費", 1),
  (8, "消耗品", 2),
  (9, "子育て用品", 2),
  (10, "ペット用品", 2),
  (11, "家具", 2),
  (12, "家電", 2),
  (13, "その他日用品", 2),
  (14, "アウトドア", 3),
  (15, "旅行", 3),
  (16, "イベント", 3),
  (17, "スポーツ", 3),
  (18, "映画・動画", 3),
  (19, "音楽", 3),
  (20, "漫画", 3),
  (21, "書籍", 3),
  (22, "ゲーム", 3),
  (23, "その他趣味・娯楽", 3),
  (24, "飲み会", 4),
  (25, "プレゼント", 4),
  (26, "冠婚葬祭", 4),
  (27, "その他交際費", 4),
  (28, "電車", 5),
  (29, "バス", 5),
  (30, "タクシー", 5),
  (31, "新幹線", 5),
  (32, "飛行機", 5),
  (33, "その他交通費", 5),
  (34, "衣服", 6),
  (35, "アクセサリー", 6),
  (36, "クリーニング", 6),
  (37, "美容院・理髪", 6),
  (38, "化粧品", 6),
  (39, "エステ・ネイル", 6),
  (40, "その他衣服・美容", 6),
  (41, "病院", 7),
  (42, "薬", 7),
  (43, "ボディケア", 7),
  (44, "フィットネス", 7),
  (45, "その他健康・医療", 7),
  (46, "携帯電話", 8),
  (47, "固定電話", 8),
  (48, "インターネット", 8),
  (49, "放送サービス", 8),
  (50, "情報サービス", 8),
  (51, "宅配・運送", 8),
  (52, "切手・はがき", 8),
  (53, "その他通信費", 8),
  (54, "新聞", 9),
  (55, "参考書", 9),
  (56, "受験料", 9),
  (57, "学費", 9),
  (58, "習い事", 9),
  (59, "塾", 9),
  (60, "その他教養・教育", 9),
  (61, "家賃", 10),
  (62, "住宅ローン", 10),
  (63, "リフォーム", 10),
  (64, "その他住宅", 10),
  (65, "水道", 11),
  (66, "電気", 11),
  (67, "ガス", 11),
  (68, "その他水道・光熱費", 11),
  (69, "自動車ローン", 12),
  (70, "ガソリン", 12),
  (71, "駐車場", 12),
  (72, "高速料金", 12),
  (73, "車検・整備", 12),
  (74, "その他自動車", 12),
  (75, "生命保険", 13),
  (76, "医療保険", 13),
  (77, "自動車保険", 13),
  (78, "住宅保険", 13),
  (79, "学資保険", 13),
  (80, "その他保険", 13),
  (81, "所得税", 14),
  (82, "住民税", 14),
  (83, "年金保険料", 14),
  (84, "自動車税", 14),
  (85, "その他税金・社会保険", 14),
  (86, "現金引き出し", 15),
  (87, "カード引き落とし", 15),
  (88, "電子マネー", 15),
  (89, "立替金", 15),
  (90, "その他現金・カード", 15),
  (91, "仕送り", 16),
  (92, "お小遣い", 16),
  (93, "使途不明金", 16),
  (94, "雑費", 16);

