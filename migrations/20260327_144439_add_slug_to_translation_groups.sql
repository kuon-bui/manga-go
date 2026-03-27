-- +migrate Up
ALTER TABLE translation_groups
    ADD COLUMN slug VARCHAR(255);

WITH normalized AS (
    SELECT
        id,
        TRIM(BOTH '-' FROM REGEXP_REPLACE(LOWER(COALESCE(name, '')), '[^a-z0-9]+', '-', 'g')) AS base_slug
    FROM translation_groups
),
ranked AS (
    SELECT
        id,
        CASE
            WHEN base_slug = '' THEN NULL
            ELSE base_slug
        END AS base_slug,
        ROW_NUMBER() OVER (
            PARTITION BY CASE
                WHEN base_slug = '' THEN NULL
                ELSE base_slug
            END
            ORDER BY id
        ) AS rn
    FROM normalized
)
UPDATE translation_groups tg
SET slug = CASE
    WHEN ranked.base_slug IS NULL THEN tg.id::text
    WHEN ranked.rn = 1 THEN ranked.base_slug
    ELSE ranked.base_slug || '-' || ranked.rn::text
END
FROM ranked
WHERE tg.id = ranked.id;

ALTER TABLE translation_groups
    ALTER COLUMN slug SET NOT NULL;

CREATE UNIQUE INDEX idx_translation_groups_slug_unique
    ON translation_groups (slug)
    WHERE deleted_at IS NULL;

-- +migrate Down
DROP INDEX IF EXISTS idx_translation_groups_slug_unique;

ALTER TABLE translation_groups
    DROP COLUMN IF EXISTS slug;
