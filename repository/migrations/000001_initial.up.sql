create table spell (
  id uuid primary key,
  name varchar not null,
  mana int not null,
  damage int not null,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  CHECK (name <> '')
);

create unique index on spell(name);
