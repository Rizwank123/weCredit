-- +goose Up
-- +goose StatementBegin
DROP TYPE IF EXISTS "public"."login_code_status";

CREATE TYPE "public"."login_code_status" AS ENUM ('SUCCESS', 'FAILED', 'PENDING');

-- Table Definition
CREATE TABLE "public"."login_codes" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "username" varchar NOT NULL,
    "code" varchar NOT NULL,
    "expiry_time" timestamptz NOT NULL,
    "status" "public"."login_code_status" NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "response_meta" text
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."login_codes";

DROP TYPE IF EXISTS "public"."login_code_status";

-- +goose StatementEnd