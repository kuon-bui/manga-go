BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- roles
INSERT INTO roles (id, name, deleted_at)
VALUES
  ('00000000-0000-0000-0000-000000000101', 'admin', NULL),
  ('00000000-0000-0000-0000-000000000102', 'editor', NULL),
  ('00000000-0000-0000-0000-000000000103', 'translator', NULL),
  ('00000000-0000-0000-0000-000000000104', 'user', NULL)
ON CONFLICT (id) DO UPDATE
SET
  name = EXCLUDED.name,
  deleted_at = NULL,
  updated_at = NOW();

-- permissions
INSERT INTO permissions (id, name, deleted_at)
VALUES
  ('00000000-0000-0000-0000-000000000201', 'users.manage', NULL),
  ('00000000-0000-0000-0000-000000000202', 'roles.manage', NULL),
  ('00000000-0000-0000-0000-000000000203', 'comics.read', NULL),
  ('00000000-0000-0000-0000-000000000204', 'comics.write', NULL),
  ('00000000-0000-0000-0000-000000000205', 'chapters.write', NULL),
  ('00000000-0000-0000-0000-000000000206', 'comments.moderate', NULL),
  ('00000000-0000-0000-0000-000000000207', 'translation_groups.manage', NULL)
ON CONFLICT (id) DO UPDATE
SET
  name = EXCLUDED.name,
  deleted_at = NULL,
  updated_at = NOW();

-- role-permission mapping
INSERT INTO roles_permissions (role_id, permission_id)
VALUES
  ('00000000-0000-0000-0000-000000000101', '00000000-0000-0000-0000-000000000201'),
  ('00000000-0000-0000-0000-000000000101', '00000000-0000-0000-0000-000000000202'),
  ('00000000-0000-0000-0000-000000000101', '00000000-0000-0000-0000-000000000203'),
  ('00000000-0000-0000-0000-000000000101', '00000000-0000-0000-0000-000000000204'),
  ('00000000-0000-0000-0000-000000000101', '00000000-0000-0000-0000-000000000205'),
  ('00000000-0000-0000-0000-000000000101', '00000000-0000-0000-0000-000000000206'),
  ('00000000-0000-0000-0000-000000000101', '00000000-0000-0000-0000-000000000207'),

  ('00000000-0000-0000-0000-000000000102', '00000000-0000-0000-0000-000000000203'),
  ('00000000-0000-0000-0000-000000000102', '00000000-0000-0000-0000-000000000204'),
  ('00000000-0000-0000-0000-000000000102', '00000000-0000-0000-0000-000000000205'),
  ('00000000-0000-0000-0000-000000000102', '00000000-0000-0000-0000-000000000206'),

  ('00000000-0000-0000-0000-000000000103', '00000000-0000-0000-0000-000000000203'),
  ('00000000-0000-0000-0000-000000000103', '00000000-0000-0000-0000-000000000205'),
  ('00000000-0000-0000-0000-000000000103', '00000000-0000-0000-0000-000000000207'),

  ('00000000-0000-0000-0000-000000000104', '00000000-0000-0000-0000-000000000203')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- users (password plaintext for test note: Admin@123456 / Editor@123456 / User@123456)
INSERT INTO users (id, name, email, password, translation_group_id, reset_password_token, reset_password_expiry_at, deleted_at)
VALUES
  ('00000000-0000-0000-0000-000000000301', 'Seed Admin User', 'seed-admin@example.com', crypt('Admin@123456', gen_salt('bf', 12)), NULL, NULL, NULL, NULL),
  ('00000000-0000-0000-0000-000000000302', 'Seed Editor User', 'seed-editor@example.com', crypt('Editor@123456', gen_salt('bf', 12)), NULL, NULL, NULL, NULL),
  ('00000000-0000-0000-0000-000000000303', 'Seed Translator User', 'seed-translator@example.com', crypt('Translator@123456', gen_salt('bf', 12)), NULL, NULL, NULL, NULL),
  ('00000000-0000-0000-0000-000000000304', 'Seed Reader User', 'seed-reader@example.com', crypt('User@123456', gen_salt('bf', 12)), NULL, NULL, NULL, NULL)
ON CONFLICT (id) DO UPDATE
SET
  name = EXCLUDED.name,
  email = EXCLUDED.email,
  password = EXCLUDED.password,
  deleted_at = NULL,
  updated_at = NOW();

-- user-role mapping
INSERT INTO users_roles (user_id, role_id)
VALUES
  ('00000000-0000-0000-0000-000000000301', '00000000-0000-0000-0000-000000000101'),
  ('00000000-0000-0000-0000-000000000302', '00000000-0000-0000-0000-000000000102'),
  ('00000000-0000-0000-0000-000000000303', '00000000-0000-0000-0000-000000000103'),
  ('00000000-0000-0000-0000-000000000304', '00000000-0000-0000-0000-000000000104')
ON CONFLICT (user_id, role_id) DO NOTHING;

-- authors
INSERT INTO authors (id, name, deleted_at)
VALUES
  ('00000000-0000-0000-0000-000000000401', 'Eiichiro Oda', NULL),
  ('00000000-0000-0000-0000-000000000402', 'Yusuke Murata', NULL),
  ('00000000-0000-0000-0000-000000000403', 'Aka Akasaka', NULL)
ON CONFLICT (id) DO UPDATE
SET
  name = EXCLUDED.name,
  deleted_at = NULL,
  updated_at = NOW();

-- genres
INSERT INTO genres (id, name, slug, description, thumbnail, deleted_at)
VALUES
  ('00000000-0000-0000-0000-000000000501', 'Action', 'action', 'Fast-paced action stories', 'https://picsum.photos/seed/action/400/240', NULL),
  ('00000000-0000-0000-0000-000000000502', 'Romance', 'romance', 'Romance and relationship stories', 'https://picsum.photos/seed/romance/400/240', NULL),
  ('00000000-0000-0000-0000-000000000503', 'Fantasy', 'fantasy', 'Magic and fantasy worlds', 'https://picsum.photos/seed/fantasy/400/240', NULL)
ON CONFLICT (id) DO UPDATE
SET
  name = EXCLUDED.name,
  slug = EXCLUDED.slug,
  description = EXCLUDED.description,
  thumbnail = EXCLUDED.thumbnail,
  deleted_at = NULL,
  updated_at = NOW();

-- tags
INSERT INTO tags (id, name, slug, deleted_at)
VALUES
  ('00000000-0000-0000-0000-000000000601', 'adventure', 'adventure', NULL),
  ('00000000-0000-0000-0000-000000000602', 'school-life', 'school-life', NULL),
  ('00000000-0000-0000-0000-000000000603', 'comedy', 'comedy', NULL),
  ('00000000-0000-0000-0000-000000000604', 'drama', 'drama', NULL)
ON CONFLICT (id) DO UPDATE
SET
  name = EXCLUDED.name,
  slug = EXCLUDED.slug,
  deleted_at = NULL,
  updated_at = NOW();

-- comics
INSERT INTO comics (
  id, title, slug, alternative_titles, description, thumbnail, banner,
  type, status, is_published, is_hot, is_featured, artist_id, published_year,
  age_rating, last_chapter_at, deleted_at
)
VALUES
  (
    '00000000-0000-0000-0000-000000000701', 'Pirate Legend', 'pirate-legend',
    '["The Sea King"]'::jsonb,
    'A long-running pirate adventure.',
    'https://picsum.photos/seed/pirate-thumb/600/800',
    'https://picsum.photos/seed/pirate-banner/1200/480',
    'manga', 'ongoing', true, true, true,
    '00000000-0000-0000-0000-000000000401', 2010,
    '13+', NOW(), NULL
  ),
  (
    '00000000-0000-0000-0000-000000000702', 'City Hunter X', 'city-hunter-x',
    '["Urban Hunter"]'::jsonb,
    'A modern action story in a mega city.',
    'https://picsum.photos/seed/city-thumb/600/800',
    'https://picsum.photos/seed/city-banner/1200/480',
    'manhwa', 'ongoing', true, false, true,
    '00000000-0000-0000-0000-000000000402', 2020,
    '16+', NOW(), NULL
  ),
  (
    '00000000-0000-0000-0000-000000000703', 'Love in Orbit', 'love-in-orbit',
    '["Orbit Love Story"]'::jsonb,
    'Romance drama set on a space station.',
    'https://picsum.photos/seed/love-thumb/600/800',
    'https://picsum.photos/seed/love-banner/1200/480',
    'comic', 'completed', false, false, false,
    '00000000-0000-0000-0000-000000000403', 2018,
    'all', NOW(), NULL
  )
ON CONFLICT (id) DO UPDATE
SET
  title = EXCLUDED.title,
  slug = EXCLUDED.slug,
  alternative_titles = EXCLUDED.alternative_titles,
  description = EXCLUDED.description,
  thumbnail = EXCLUDED.thumbnail,
  banner = EXCLUDED.banner,
  type = EXCLUDED.type,
  status = EXCLUDED.status,
  is_published = EXCLUDED.is_published,
  is_hot = EXCLUDED.is_hot,
  is_featured = EXCLUDED.is_featured,
  artist_id = EXCLUDED.artist_id,
  published_year = EXCLUDED.published_year,
  age_rating = EXCLUDED.age_rating,
  last_chapter_at = EXCLUDED.last_chapter_at,
  deleted_at = NULL,
  updated_at = NOW();

-- comic-authors mapping
INSERT INTO comic_authors (comic_id, author_id)
VALUES
  ('00000000-0000-0000-0000-000000000701', '00000000-0000-0000-0000-000000000401'),
  ('00000000-0000-0000-0000-000000000702', '00000000-0000-0000-0000-000000000402'),
  ('00000000-0000-0000-0000-000000000703', '00000000-0000-0000-0000-000000000403')
ON CONFLICT (comic_id, author_id) DO NOTHING;

-- comic-genres mapping
INSERT INTO comic_genres (comic_id, genre_id)
VALUES
  ('00000000-0000-0000-0000-000000000701', '00000000-0000-0000-0000-000000000501'),
  ('00000000-0000-0000-0000-000000000701', '00000000-0000-0000-0000-000000000503'),
  ('00000000-0000-0000-0000-000000000702', '00000000-0000-0000-0000-000000000501'),
  ('00000000-0000-0000-0000-000000000703', '00000000-0000-0000-0000-000000000502')
ON CONFLICT (comic_id, genre_id) DO NOTHING;

-- comic-tags mapping
INSERT INTO comic_tags (comic_id, tag_id)
VALUES
  ('00000000-0000-0000-0000-000000000701', '00000000-0000-0000-0000-000000000601'),
  ('00000000-0000-0000-0000-000000000701', '00000000-0000-0000-0000-000000000603'),
  ('00000000-0000-0000-0000-000000000702', '00000000-0000-0000-0000-000000000603'),
  ('00000000-0000-0000-0000-000000000703', '00000000-0000-0000-0000-000000000602'),
  ('00000000-0000-0000-0000-000000000703', '00000000-0000-0000-0000-000000000604')
ON CONFLICT (comic_id, tag_id) DO NOTHING;

-- chapters
INSERT INTO chapters (id, comic_id, number, title, slug, is_published, chapter_idx, deleted_at)
VALUES
  ('00000000-0000-0000-0000-000000000801', '00000000-0000-0000-0000-000000000701', '1', 'Romance Dawn', 'chapter-1', true, 1, NULL),
  ('00000000-0000-0000-0000-000000000802', '00000000-0000-0000-0000-000000000701', '2', 'Into The Sea', 'chapter-2', true, 2, NULL),
  ('00000000-0000-0000-0000-000000000803', '00000000-0000-0000-0000-000000000702', '1', 'Neon Alley', 'chapter-1', true, 1, NULL),
  ('00000000-0000-0000-0000-000000000804', '00000000-0000-0000-0000-000000000703', '1', 'Docking Hearts', 'chapter-1', false, 1, NULL)
ON CONFLICT (id) DO UPDATE
SET
  comic_id = EXCLUDED.comic_id,
  number = EXCLUDED.number,
  title = EXCLUDED.title,
  slug = EXCLUDED.slug,
  is_published = EXCLUDED.is_published,
  chapter_idx = EXCLUDED.chapter_idx,
  deleted_at = NULL,
  updated_at = NOW();

-- pages
INSERT INTO pages (id, chapter_id, page_number, image_url, deleted_at)
VALUES
  ('00000000-0000-0000-0000-000000000901', '00000000-0000-0000-0000-000000000801', 1, 'https://picsum.photos/seed/p801-1/1080/1920', NULL),
  ('00000000-0000-0000-0000-000000000902', '00000000-0000-0000-0000-000000000801', 2, 'https://picsum.photos/seed/p801-2/1080/1920', NULL),
  ('00000000-0000-0000-0000-000000000903', '00000000-0000-0000-0000-000000000802', 1, 'https://picsum.photos/seed/p802-1/1080/1920', NULL),
  ('00000000-0000-0000-0000-000000000904', '00000000-0000-0000-0000-000000000803', 1, 'https://picsum.photos/seed/p803-1/1080/1920', NULL),
  ('00000000-0000-0000-0000-000000000905', '00000000-0000-0000-0000-000000000804', 1, 'https://picsum.photos/seed/p804-1/1080/1920', NULL)
ON CONFLICT (id) DO UPDATE
SET
  chapter_id = EXCLUDED.chapter_id,
  page_number = EXCLUDED.page_number,
  image_url = EXCLUDED.image_url,
  deleted_at = NULL,
  updated_at = NOW();

-- translation groups (requires users existing)
INSERT INTO translation_groups (id, name, owner_id, slug, deleted_at)
VALUES
  ('00000000-0000-0000-0000-000000001001', 'Luminous Scans', '00000000-0000-0000-0000-000000000303', 'luminous-scans', NULL),
  ('00000000-0000-0000-0000-000000001002', 'Night Owl Team', '00000000-0000-0000-0000-000000000302', 'night-owl-team', NULL)
ON CONFLICT (id) DO UPDATE
SET
  name = EXCLUDED.name,
  owner_id = EXCLUDED.owner_id,
  slug = EXCLUDED.slug,
  deleted_at = NULL,
  updated_at = NOW();

UPDATE users
SET translation_group_id = '00000000-0000-0000-0000-000000001001'
WHERE id = '00000000-0000-0000-0000-000000000303';

UPDATE users
SET translation_group_id = '00000000-0000-0000-0000-000000001002'
WHERE id = '00000000-0000-0000-0000-000000000302';

-- reading histories (unique user_id + chapter_id)
INSERT INTO reading_histories (id, user_id, chapter_id, comic_id, last_read_at, deleted_at)
VALUES
  ('00000000-0000-0000-0000-000000001101', '00000000-0000-0000-0000-000000000304', '00000000-0000-0000-0000-000000000801', '00000000-0000-0000-0000-000000000701', NOW() - INTERVAL '1 day', NULL),
  ('00000000-0000-0000-0000-000000001102', '00000000-0000-0000-0000-000000000304', '00000000-0000-0000-0000-000000000803', '00000000-0000-0000-0000-000000000702', NOW() - INTERVAL '2 hours', NULL),
  ('00000000-0000-0000-0000-000000001103', '00000000-0000-0000-0000-000000000302', '00000000-0000-0000-0000-000000000802', '00000000-0000-0000-0000-000000000701', NOW() - INTERVAL '30 minutes', NULL)
ON CONFLICT (user_id, chapter_id) DO UPDATE
SET
  comic_id = EXCLUDED.comic_id,
  last_read_at = EXCLUDED.last_read_at,
  deleted_at = NULL,
  updated_at = NOW();

-- reading progresses
INSERT INTO reading_progresses (id, user_id, comic_id, chapter_id, scroll_percent, deleted_at)
VALUES
  ('00000000-0000-0000-0000-000000001201', '00000000-0000-0000-0000-000000000304', '00000000-0000-0000-0000-000000000701', '00000000-0000-0000-0000-000000000801', 45, NULL),
  ('00000000-0000-0000-0000-000000001202', '00000000-0000-0000-0000-000000000304', '00000000-0000-0000-0000-000000000702', '00000000-0000-0000-0000-000000000803', 88, NULL),
  ('00000000-0000-0000-0000-000000001203', '00000000-0000-0000-0000-000000000302', '00000000-0000-0000-0000-000000000701', '00000000-0000-0000-0000-000000000802', 12, NULL)
ON CONFLICT (id) DO UPDATE
SET
  user_id = EXCLUDED.user_id,
  comic_id = EXCLUDED.comic_id,
  chapter_id = EXCLUDED.chapter_id,
  scroll_percent = EXCLUDED.scroll_percent,
  deleted_at = NULL,
  updated_at = NOW();

-- comments (unique user_id + chapter_id)
INSERT INTO comments (id, user_id, chapter_id, comic_id, page_index, content, last_read_at, deleted_at)
VALUES
  ('00000000-0000-0000-0000-000000001301', '00000000-0000-0000-0000-000000000304', '00000000-0000-0000-0000-000000000801', '00000000-0000-0000-0000-000000000701', 1, 'This opening chapter is amazing.', NOW(), NULL),
  ('00000000-0000-0000-0000-000000001302', '00000000-0000-0000-0000-000000000302', '00000000-0000-0000-0000-000000000802', '00000000-0000-0000-0000-000000000701', 1, 'Needs better pacing, but good art.', NOW(), NULL),
  ('00000000-0000-0000-0000-000000001303', '00000000-0000-0000-0000-000000000303', '00000000-0000-0000-0000-000000000803', '00000000-0000-0000-0000-000000000702', 1, 'Translation notes are ready.', NOW(), NULL)
ON CONFLICT (user_id, chapter_id) DO UPDATE
SET
  comic_id = EXCLUDED.comic_id,
  page_index = EXCLUDED.page_index,
  content = EXCLUDED.content,
  last_read_at = EXCLUDED.last_read_at,
  deleted_at = NULL,
  updated_at = NOW();

-- reactions
INSERT INTO reactions (id, user_id, comment_id, type, last_read_at, deleted_at)
VALUES
  ('00000000-0000-0000-0000-000000001401', '00000000-0000-0000-0000-000000000301', '00000000-0000-0000-0000-000000001301', 'like', NOW(), NULL),
  ('00000000-0000-0000-0000-000000001402', '00000000-0000-0000-0000-000000000304', '00000000-0000-0000-0000-000000001302', 'like', NOW(), NULL),
  ('00000000-0000-0000-0000-000000001403', '00000000-0000-0000-0000-000000000302', '00000000-0000-0000-0000-000000001303', 'love', NOW(), NULL)
ON CONFLICT (id) DO UPDATE
SET
  user_id = EXCLUDED.user_id,
  comment_id = EXCLUDED.comment_id,
  type = EXCLUDED.type,
  last_read_at = EXCLUDED.last_read_at,
  deleted_at = NULL,
  updated_at = NOW();

-- user comic reads (read_data as bytea)
INSERT INTO user_comic_reads (id, user_id, comic_id, read_data, deleted_at)
VALUES
  ('00000000-0000-0000-0000-000000001501', '00000000-0000-0000-0000-000000000304', '00000000-0000-0000-0000-000000000701', decode('0f', 'hex'), NULL),
  ('00000000-0000-0000-0000-000000001502', '00000000-0000-0000-0000-000000000304', '00000000-0000-0000-0000-000000000702', decode('03', 'hex'), NULL),
  ('00000000-0000-0000-0000-000000001503', '00000000-0000-0000-0000-000000000302', '00000000-0000-0000-0000-000000000701', decode('01', 'hex'), NULL)
ON CONFLICT (id) DO UPDATE
SET
  user_id = EXCLUDED.user_id,
  comic_id = EXCLUDED.comic_id,
  read_data = EXCLUDED.read_data,
  deleted_at = NULL,
  updated_at = NOW();

COMMIT;
