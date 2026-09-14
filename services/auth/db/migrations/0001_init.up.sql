CREATE TABLE users (
    userid UUID NOT NULL DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT users_pkey PRIMARY KEY (userid),
    CONSTRAINT users_email_key UNIQUE (email)
);

CREATE TABLE refresh_tokens (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    userid UUID NOT NULL,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    revoked_at TIMESTAMPTZ,

    CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id),
    CONSTRAINT refresh_tokens_token_hash_key UNIQUE (token_hash),

    CONSTRAINT refresh_tokens_userid_fkey
        FOREIGN KEY (userid)
        REFERENCES users(userid)
        ON DELETE CASCADE
);
