create table if not exists entries (
  id integer primary key,
  word text not null,
  pos text not null,
  defs text not null
);
create unique index if not exists ux_entries_word_pos on entries(word, pos);
create index if not exists ix_entries_word on entries(word);