-- +goose Up
-- +goose StatementBegin
DROP TYPE IF EXISTS "public"."user_role";

CREATE TYPE "public"."user_role" AS ENUM ('USER');

-- Table Definition
CREATE TABLE "public"."users" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "full_name" varchar NOT NULL,
    "user_name" varchar NOT NULL,
    "role" "public"."user_role" NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."users";

DROP TYPE IF EXISTS "public"."user_role";

-- +goose StatementEnd