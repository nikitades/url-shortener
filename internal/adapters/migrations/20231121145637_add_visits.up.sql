CREATE TABLE public.visits (
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    user_agent text NOT NULL,
    url_id bigint NOT NULL,
    url_source text NOT NULL,
    url_code text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    CONSTRAINT fk_url_id FOREIGN KEY(url_id) REFERENCES urls(id) ON DELETE RESTRICT ON UPDATE RESTRICT
);