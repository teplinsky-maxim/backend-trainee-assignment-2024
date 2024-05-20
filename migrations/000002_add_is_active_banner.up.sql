BEGIN;

ALTER TABLE "banners"
    ADD COLUMN "is_active" boolean NOT NULL DEFAULT true;

COMMIT;