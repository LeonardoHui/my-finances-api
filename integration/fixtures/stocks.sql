INSERT INTO stocks (id, name) VALUES (01, 'AAAA3');
INSERT INTO stocks (id, name) VALUES (02, 'BBBB4');
INSERT INTO stocks (id, name) VALUES (03, 'CCCC11');


INSERT INTO stock_events (id, user_id, broker_id, created_at, event, stock_id, quantity, price, total_price) VALUES (01, 1, 1, '2020-10-07T00:00:00Z', 'BUY',   	1	,5	 , 13403 ,	67014	);
INSERT INTO stock_events (id, user_id, broker_id, created_at, event, stock_id, quantity, price, total_price) VALUES (02, 1, 1, '2020-10-07T00:00:00Z', 'BUY',   	2	,5	 , 9049  ,	45247	);
INSERT INTO stock_events (id, user_id, broker_id, created_at, event, stock_id, quantity, price, total_price) VALUES (03, 1, 1, '2020-10-07T00:00:00Z', 'BUY',	    3	,5	 , 14236 ,	71179	);
INSERT INTO stock_events (id, user_id, broker_id, created_at, event, stock_id, quantity, price, total_price) VALUES (04, 1, 1, '2020-10-07T00:00:00Z', 'DIVIDEND',	1	,5	 , 1 ,  	5	);
INSERT INTO stock_events (id, user_id, broker_id, created_at, event, stock_id, quantity, price, total_price) VALUES (05, 1, 1, '2020-10-07T00:00:00Z', 'DIVIDEND',	2	,5	 , 2 ,  	10	);
INSERT INTO stock_events (id, user_id, broker_id, created_at, event, stock_id, quantity, price, total_price) VALUES (06, 1, 1, '2020-10-07T00:00:00Z', 'DIVIDEND',	3	,5	 , 3 ,  	15	);
INSERT INTO stock_events (id, user_id, broker_id, created_at, event, stock_id, quantity, price, total_price) VALUES (07, 1, 1, '2020-11-07T00:00:00Z', 'BUY',   	1	,5	 , 13403 ,	67014	);
INSERT INTO stock_events (id, user_id, broker_id, created_at, event, stock_id, quantity, price, total_price) VALUES (08, 1, 1, '2020-12-07T00:00:00Z', 'BUY',   	2	,5	 , 9049  ,	45247	);
INSERT INTO stock_events (id, user_id, broker_id, created_at, event, stock_id, quantity, price, total_price) VALUES (09, 1, 1, '2020-12-07T00:00:00Z', 'BUY',	    3	,5	 , 14236 ,	71179	);
