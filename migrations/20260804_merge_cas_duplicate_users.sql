-- One-off migration for duplicate local/CAS accounts.
--
-- Run migrations/20260804_preview_cas_duplicate_users.sql first and take a
-- database backup before applying this file.
--
-- This migration only merges unambiguous CAS shadow accounts into existing
-- local accounts:
--   1. A CAS shadow account is a user currently bound by provider = 'cas'
--      whose username starts with 'cas_'.
--   2. Prefer a non-shadow local account with the same CAS email.
--   3. Fall back to a non-shadow local account whose username equals the CAS
--      provider_subject.
--   4. Ambiguous matches are skipped and must be handled manually.

START TRANSACTION;

DROP TEMPORARY TABLE IF EXISTS cas_user_merge_candidates;

CREATE TEMPORARY TABLE cas_user_merge_candidates AS
SELECT
  ui.id AS identity_id,
  ui.provider_subject AS cas_subject,
  ui.email AS identity_email,
  shadow.id AS shadow_user_id,
  shadow.username AS shadow_username,
  shadow.email AS shadow_email,
  target.id AS target_user_id,
  target.username AS target_username,
  target.email AS target_email,
  'email' AS matched_by,
  1 AS match_rank
FROM user_identities ui
JOIN users shadow ON shadow.id = ui.user_id
JOIN users target
  ON target.id <> shadow.id
 AND target.email = COALESCE(NULLIF(ui.email, ''), NULLIF(shadow.email, ''))
 AND LEFT(COALESCE(target.username, ''), 4) <> 'cas_'
WHERE ui.provider = 'cas'
  AND LEFT(COALESCE(shadow.username, ''), 4) = 'cas_'
  AND COALESCE(NULLIF(ui.email, ''), NULLIF(shadow.email, '')) IS NOT NULL

UNION ALL

SELECT
  ui.id AS identity_id,
  ui.provider_subject AS cas_subject,
  ui.email AS identity_email,
  shadow.id AS shadow_user_id,
  shadow.username AS shadow_username,
  shadow.email AS shadow_email,
  target.id AS target_user_id,
  target.username AS target_username,
  target.email AS target_email,
  'username' AS matched_by,
  2 AS match_rank
FROM user_identities ui
JOIN users shadow ON shadow.id = ui.user_id
JOIN users target
  ON target.id <> shadow.id
 AND target.username = ui.provider_subject
 AND LEFT(COALESCE(target.username, ''), 4) <> 'cas_'
WHERE ui.provider = 'cas'
  AND LEFT(COALESCE(shadow.username, ''), 4) = 'cas_'
  AND ui.provider_subject <> '';

DROP TEMPORARY TABLE IF EXISTS cas_user_merge_pairs;

CREATE TEMPORARY TABLE cas_user_merge_pairs AS
SELECT
  c.identity_id,
  c.shadow_user_id,
  MIN(c.target_user_id) AS target_user_id,
  MIN(c.matched_by) AS matched_by
FROM cas_user_merge_candidates c
JOIN (
  SELECT identity_id, MIN(match_rank) AS match_rank
  FROM cas_user_merge_candidates
  GROUP BY identity_id
) best ON best.identity_id = c.identity_id
      AND best.match_rank = c.match_rank
GROUP BY c.identity_id, c.shadow_user_id
HAVING COUNT(DISTINCT c.target_user_id) = 1;

UPDATE users target
JOIN cas_user_merge_pairs p ON p.target_user_id = target.id
JOIN users shadow ON shadow.id = p.shadow_user_id
SET
  target.email = CASE
    WHEN COALESCE(target.email, '') = '' THEN shadow.email
    ELSE target.email
  END,
  target.info = CASE
    WHEN COALESCE(target.info, '') = ''
     AND COALESCE(shadow.info, '') <> 'cas authenticated user'
    THEN shadow.info
    ELSE target.info
  END,
  target.avatar_url = CASE
    WHEN COALESCE(target.avatar_url, '') = '' THEN shadow.avatar_url
    ELSE target.avatar_url
  END,
  target.personal_blog = CASE
    WHEN COALESCE(target.personal_blog, '') = '' THEN shadow.personal_blog
    ELSE target.personal_blog
  END,
  target.github = CASE
    WHEN COALESCE(target.github, '') = '' THEN shadow.github
    ELSE target.github
  END,
  target.flickr = CASE
    WHEN COALESCE(target.flickr, '') = '' THEN shadow.flickr
    ELSE target.flickr
  END,
  target.weibo = CASE
    WHEN COALESCE(target.weibo, '') = '' THEN shadow.weibo
    ELSE target.weibo
  END,
  target.zhihu = CASE
    WHEN COALESCE(target.zhihu, '') = '' THEN shadow.zhihu
    ELSE target.zhihu
  END,
  target.birthday = CASE
    WHEN COALESCE(target.birthday, '') = '' THEN shadow.birthday
    ELSE target.birthday
  END,
  target.`group` = CASE
    WHEN COALESCE(target.`group`, '') = '' THEN shadow.`group`
    ELSE target.`group`
  END,
  target.hometown = CASE
    WHEN COALESCE(target.hometown, '') = '' THEN shadow.hometown
    ELSE target.hometown
  END,
  target.timejoin = CASE
    WHEN COALESCE(target.timejoin, '') = '' THEN shadow.timejoin
    ELSE target.timejoin
  END;

UPDATE user_identities ui
JOIN cas_user_merge_pairs p ON p.identity_id = ui.id
SET ui.user_id = p.target_user_id;

SET @oauth2_token_exists := (
  SELECT COUNT(*)
  FROM information_schema.tables
  WHERE table_schema = DATABASE()
    AND table_name = 'oauth2_token'
);

SET @rewrite_oauth2_token_sql := IF(
  @oauth2_token_exists > 0,
  'UPDATE oauth2_token t
   JOIN cas_user_merge_pairs p
     ON JSON_VALID(t.data) = 1
    AND JSON_UNQUOTE(JSON_EXTRACT(t.data, ''$.UserID'')) = CAST(p.shadow_user_id AS CHAR)
   SET t.data = JSON_SET(t.data, ''$.UserID'', CAST(p.target_user_id AS CHAR))',
  'SELECT ''oauth2_token table not found; skipped token rewrite'' AS note'
);

PREPARE rewrite_oauth2_token FROM @rewrite_oauth2_token_sql;
EXECUTE rewrite_oauth2_token;
DEALLOCATE PREPARE rewrite_oauth2_token;

DELETE shadow
FROM users shadow
JOIN (
  SELECT DISTINCT shadow_user_id
  FROM cas_user_merge_pairs
) merged ON merged.shadow_user_id = shadow.id
LEFT JOIN user_identities remaining ON remaining.user_id = shadow.id
WHERE remaining.id IS NULL;

SELECT COUNT(*) AS merged_identity_count
FROM cas_user_merge_pairs;

SELECT c.*
FROM cas_user_merge_candidates c
LEFT JOIN cas_user_merge_pairs p ON p.identity_id = c.identity_id
WHERE p.identity_id IS NULL
ORDER BY c.identity_id, c.match_rank, c.target_user_id;

COMMIT;
