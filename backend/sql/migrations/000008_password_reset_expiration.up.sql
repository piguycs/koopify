alter table password_resets add column token_expiration timestamptz not null default (now() + interval '30 minutes');
