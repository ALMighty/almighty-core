insert into spaces (id, name) values ('00000000-1111-2222-3333-000000000000', 'test');
insert into areas (id, name, space_id) values ('00000000-1111-2222-3333-444444444444', 'test area', '00000000-1111-2222-3333-000000000000');
insert into work_items (id, space_id, type, fields) values (12345, '00000000-1111-2222-3333-000000000000', '26787039-b68f-4e28-8814-c2f93be1ef4e', '{"system.title":"Title"}'::json);