-- Preview duplicate local/CAS accounts before running the merge migration.
-- Rule:
--   1. A CAS shadow account is a user currently bound by provider = 'cas'
--      whose username starts with 'cas_'.
--   2. Prefer a non-shadow local account with the same CAS email.
--   3. Fall back to a non-shadow local account whose username equals the CAS
--      provider_subject.
--   4. Only identities with exactly one best target are merged by the apply
--      script. Ambiguous rows must be handled manually.

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
  c.cas_subject,
  c.identity_email,
  c.shadow_user_id,
  c.shadow_username,
  c.shadow_email,
  MIN(c.target_user_id) AS target_user_id,
  MIN(c.target_username) AS target_username,
  MIN(c.target_email) AS target_email,
  MIN(c.matched_by) AS matched_by
FROM cas_user_merge_candidates c
JOIN (
  SELECT identity_id, MIN(match_rank) AS match_rank
  FROM cas_user_merge_candidates
  GROUP BY identity_id
) best ON best.identity_id = c.identity_id
      AND best.match_rank = c.match_rank
GROUP BY
  c.identity_id,
  c.cas_subject,
  c.identity_email,
  c.shadow_user_id,
  c.shadow_username,
  c.shadow_email
HAVING COUNT(DISTINCT c.target_user_id) = 1;

SELECT *
FROM cas_user_merge_pairs
ORDER BY identity_id;

SELECT c.*
FROM cas_user_merge_candidates c
LEFT JOIN cas_user_merge_pairs p ON p.identity_id = c.identity_id
WHERE p.identity_id IS NULL
ORDER BY c.identity_id, c.match_rank, c.target_user_id;
