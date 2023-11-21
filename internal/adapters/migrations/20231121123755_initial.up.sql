CREATE TABLE public.urls (
	id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
	source_url text NOT NULL,
	short_url text NOT NULL,
	created_at timestamp with time zone NOT NULL,
	CONSTRAINT urls_pk PRIMARY KEY (id),
	CONSTRAINT urls_source_url_un UNIQUE (source_url),
	CONSTRAINT urls_short_url_un UNIQUE (short_url)
);