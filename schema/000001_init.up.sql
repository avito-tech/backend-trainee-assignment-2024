CREATE TABLE tag (
	id bigserial PRIMARY KEY,
	created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE feature (
	id bigserial PRIMARY KEY,
	created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE banner (
	id bigserial PRIMARY KEY,
	feature_id bigint NOT NULL,
	content jsonb NOT NULL,
	is_active boolean NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (feature_id) REFERENCES feature (id)
);

CREATE TABLE banner_tag (
	tag_id bigint NOT NULL,
	banner_id bigint NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (tag_id, banner_id),
    FOREIGN KEY (tag_id) REFERENCES tag (id),
    FOREIGN KEY (banner_id) REFERENCES banner (id)
);
