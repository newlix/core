-- Item is a to-do item.
CREATE TABLE IF NOT EXISTS "item" (
  id   text PRIMARY KEY
);
ALTER TABLE "item" ADD COLUMN IF NOT EXISTS "text" TEXT;
ALTER TABLE "item" ALTER COLUMN "text" SET NOT NULL;
ALTER TABLE "item" ADD COLUMN IF NOT EXISTS "created_at" Timestamp;
