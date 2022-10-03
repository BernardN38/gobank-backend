CREATE TABLE transactions (
                        id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
                       amount BIGINT NOT NULL,
                       to_account uuid      NOT NULL ,
                       from_account uuid NOT NULL,
                       created_at timestamp NOT NULL
);