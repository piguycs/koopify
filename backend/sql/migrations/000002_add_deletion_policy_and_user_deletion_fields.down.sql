alter table users drop column deletion_scheduled_at;
alter table users drop column requested_deletion_at;

drop table deletion_policy;
