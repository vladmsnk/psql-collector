package collector

const (
	SelectTablesInfo = `
SELECT
    relid,
    relname AS tablename,
    n_live_tup,--the number of live tuples
    n_dead_tup,--the number of dead tuples
    seq_scan,
    idx_scan,
    n_tup_ins AS inserts,-- the number of inserts
    n_tup_upd AS updates,-- the number of updates
    n_tup_del AS deletes,-- the number of deletes
    last_vacuum,--the last vacuum times 
    last_autovacuum,
    last_analyze,--the last analyze times 
    last_autoanalyze
FROM
    pg_stat_user_tables
WHERE
    n_live_tup + seq_scan + idx_scan > 0
ORDER BY
    n_live_tup DESC,
    seq_scan + idx_scan DESC;
`

	SelectQueryTypesDistribution = `
SELECT
  query_type,
  count(*) AS total
FROM (
  SELECT
    CASE
      WHEN query LIKE 'INSERT%' THEN 'INSERT'
      WHEN query LIKE 'UPDATE%' THEN 'UPDATE'
      WHEN query LIKE 'DELETE%' THEN 'DELETE'
      WHEN query LIKE 'SELECT%' THEN 'SELECT'
      ELSE 'OTHER'
    END AS query_type
  FROM pg_stat_statements
) AS categorized_queries
GROUP BY query_type
ORDER BY total DESC;
`
)
