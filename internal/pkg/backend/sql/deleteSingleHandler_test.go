package sql

import (
	"testing"
)

func TestDeleteSingle(t *testing.T) {
	for _, tc := range TestCasesDeleteSingle {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.Kind == "success" {
				tc.Run = itSucceeds
				tc.Run(t, tc)
			} else if tc.Kind == "failure" {
				tc.Run = itFails
				tc.Run(t, tc)
			} else {
				t.Errorf("testcase should either be success or failure kind")
			}
		})
	}
}
