CREATE TABLE clients (
  walletId UUID PRIMARY KEY,
  amount INTEGER NOT NULL check (amount >= 0)
);

insert into clients values ('a0eebc999c0b4ef8bb6d6bb9bd38-0a11', 1000);
insert into clients values ('29f99fe4-f355-414c-89e8-824b04039c14', 2000);
insert into clients values ('e0eebc999c0b4ef8bb6d6bb9bd38-0a11', 0);
insert into clients values ('99f99fe4-f355-414c-89e8-824b04039c14', 1200);
