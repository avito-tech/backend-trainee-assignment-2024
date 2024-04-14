create table banner (
	id			bigserial	primary key,
	content		jsonb		not null,
	is_active	boolean		not null,
	created_at	timestamptz not null	default now(),
	updated_at	timestamptz not null	default now()
);

create table feature_tag_banner (
	feature_id	bigint		not null,
	tag_id		bigint		not null,
	banner_id	bigint		not null	references banner (id),
	created_at	timestamptz	not null	default now(),
	primary key (feature_id, tag_id)
);
