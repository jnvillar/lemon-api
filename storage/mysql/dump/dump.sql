INSERT INTO user (alias, date_created, email, firstname, id, lastname)
values ('juaninv', Now(), 'juan@gmail.com', 'juan', '2c909e4e-88c5-4cfb-8a4e-b41bdbacd55c', 'noli');

INSERT INTO wallet (id, user_id, currency, balance)
values ('80f9eff4-41e2-408b-84c9-4fac97120bf6', '2c909e4e-88c5-4cfb-8a4e-b41bdbacd55c', 'ARS', 1000),
       ('9edbc638-4c9f-4b3a-9b54-25299206993d', '2c909e4e-88c5-4cfb-8a4e-b41bdbacd55c', 'USDT', 1000),
       ('78a10bc9-970b-48ea-8e77-2fed24e73301', '2c909e4e-88c5-4cfb-8a4e-b41bdbacd55c', 'BTC', 1000);

INSERT INTO user (alias, date_created, email, firstname, id, lastname)
values ('leila', Now(), 'lei@gmail.com', 'lei', '2c909e4e-88c5-4cfb-8a4e-b41bdbacd55d', 'apellido');

INSERT INTO wallet (id, user_id, currency, balance)
values ('80f9eff4-41e2-408b-84c9-4fac97120bf7', '2c909e4e-88c5-4cfb-8a4e-b41bdbacd55d', 'ARS', 1000),
       ('9edbc638-4c9f-4b3a-9b54-25299206993f', '2c909e4e-88c5-4cfb-8a4e-b41bdbacd55d', 'USDT', 1000),
       ('78a10bc9-970b-48ea-8e77-2fed24e73302', '2c909e4e-88c5-4cfb-8a4e-b41bdbacd55d', 'BTC', 1000);