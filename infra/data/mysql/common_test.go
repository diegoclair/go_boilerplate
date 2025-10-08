package mysql

import (
	"strings"
	"testing"
)

func TestWithCount(t *testing.T) {
	tests := []struct {
		name          string
		inputQuery    string
		expectedCount string // What we expect to find in the output
		shouldContain bool   // Should contain COUNT(*) OVER()
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

			// Check if COUNT(*) OVER() was added
			if tt.shouldContain && !strings.Contains(result, tt.expectedCount) {
				t.Errorf("withCount() should contain %q, but got:\n%s", tt.expectedCount, result)
			}

			// For queries without FROM, ensure it returns the original query unchanged
			if !tt.shouldContain {
				if result != tt.inputQuery {
					t.Errorf("withCount() should return original query unchanged when no FROM found, but got:\n%s", result)
				}
				return // Skip remaining checks for this test case
			}

			// Verify COUNT appears before FROM
			countIndex := strings.Index(strings.ToLower(result), "count(*) over()")
			fromIndex := strings.Index(strings.ToLower(result), "from")

			if countIndex == -1 {
				t.Error("withCount() did not add COUNT(*) OVER()")
			}

			if fromIndex == -1 {
				t.Error("withCount() result does not contain FROM keyword")
			}

			if countIndex >= fromIndex {
				t.Errorf("COUNT(*) OVER() should appear before FROM, but got:\n%s", result)
			}

			// Verify only one COUNT was added (should not add multiple)
			countOccurrences := strings.Count(strings.ToLower(result), "count(*) over()")
			if countOccurrences != 1 {
				t.Errorf("Expected exactly 1 COUNT(*) OVER(), but found %d in:\n%s", countOccurrences, result)
			}
		})
	}
}
