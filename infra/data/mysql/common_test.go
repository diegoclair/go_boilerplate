package mysql

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithCount(t *testing.T) {
	tests := []struct {
		name          string
		inputQuery    string
		expectedCount string
		shouldContain bool
	}{
		{
			name: "Should add COUNT to query with tabs formatting",
			inputQuery: `
				SELECT
					id,
					name
				FROM tab_test
			`,
			expectedCount: "COUNT(*) OVER() as total_count",
			shouldContain: true,
		},
		{
			name:          "Should add COUNT to inline query",
			inputQuery:    `SELECT id, name FROM tab_test`,
			expectedCount: "COUNT(*) OVER() as total_count",
			shouldContain: true,
		},
		{
			name: "Should add COUNT to query with spaces instead of tabs",
			inputQuery: `
  SELECT
    id,
    name
  FROM tab_test
  `,
			expectedCount: "COUNT(*) OVER() as total_count",
			shouldContain: true,
		},
		{
			name: "Should work with lowercase FROM",
			inputQuery: `
	SELECT id, name
	from tab_test
	`,
			expectedCount: "COUNT(*) OVER() as total_count",
			shouldContain: true,
		},
		{
			name: "Should add COUNT only to first FROM (with JOINs)",
			inputQuery: `
	SELECT
		a.id,
		b.name
	FROM tab_a a
	INNER JOIN tab_b b ON a.id = b.a_id
	`,
			expectedCount: "COUNT(*) OVER() as total_count",
			shouldContain: true,
		},
		{
			name: "Should work with mixed case FROM",
			inputQuery: `
	SELECT id, name
	From tab_test
	`,
			expectedCount: "COUNT(*) OVER() as total_count",
			shouldContain: true,
		},
		{
			name: "Should work with inline from",
			inputQuery: `
				SELECT id, name from tab_test
			`,
			expectedCount: "COUNT(*) OVER() as total_count",
			shouldContain: true,
		},
		{
			name:          "Should return original query if no FROM found",
			inputQuery:    `SELECT 1 + 1`,
			expectedCount: "",
			shouldContain: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := withCount(tt.inputQuery)

			if tt.shouldContain {
				require.Contains(t, result, tt.expectedCount)

				// Verify COUNT appears before FROM
				countIndex := strings.Index(strings.ToLower(result), "count(*) over()")
				fromIndex := strings.Index(strings.ToLower(result), "from")
				require.NotEqual(t, -1, countIndex, "COUNT(*) OVER() should be present")
				require.NotEqual(t, -1, fromIndex, "FROM should be present")
				require.Less(t, countIndex, fromIndex, "COUNT should appear before FROM")

				// Verify only one COUNT was added
				countOccurrences := strings.Count(strings.ToLower(result), "count(*) over()")
				require.Equal(t, 1, countOccurrences, "Should add exactly one COUNT")
			} else {
				require.Equal(t, tt.inputQuery, result, "Should return original query unchanged")
				require.NotContains(t, result, "COUNT(*) OVER()")
			}
		})
	}
}

func TestBuildInPlaceholders(t *testing.T) {
	t.Run("String slices", func(t *testing.T) {
		t.Run("Multiple strings", func(t *testing.T) {
			placeholders, newArgs := buildInPlaceholders([]any{1, "test"}, []string{"a", "b", "c"})

			require.Equal(t, "?, ?, ?", placeholders)
			require.Len(t, newArgs, 5)
			require.Equal(t, []any{1, "test", "a", "b", "c"}, newArgs)
		})

		t.Run("Single string", func(t *testing.T) {
			placeholders, newArgs := buildInPlaceholders([]any{42}, []string{"single"})

			require.Equal(t, "?", placeholders)
			require.Equal(t, []any{42, "single"}, newArgs)
		})

		t.Run("Empty slice", func(t *testing.T) {
			placeholders, newArgs := buildInPlaceholders([]any{1, 2}, []string{})

			require.Empty(t, placeholders)
			require.Equal(t, []any{1, 2}, newArgs)
		})
	})

	t.Run("Int slices", func(t *testing.T) {
		placeholders, newArgs := buildInPlaceholders([]any{"test"}, []int{10, 20, 30})

		require.Equal(t, "?, ?, ?", placeholders)
		require.Equal(t, []any{"test", 10, 20, 30}, newArgs)
	})

	t.Run("Int64 slices", func(t *testing.T) {
		placeholders, newArgs := buildInPlaceholders([]any{}, []int64{100, 200})

		require.Equal(t, "?, ?", placeholders)
		require.Equal(t, []any{int64(100), int64(200)}, newArgs)
	})

	t.Run("Preserves initial args", func(t *testing.T) {
		initial := []any{"preserved", 123}
		_, newArgs := buildInPlaceholders(initial, []string{"new"})

		require.Equal(t, []any{"preserved", 123, "new"}, newArgs)
	})
}
