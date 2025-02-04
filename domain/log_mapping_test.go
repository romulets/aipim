package domain

import (
	"os"
	"testing"

	"github.com/andreyvit/diff"
)

func TestToString(t *testing.T) {
	tcs := map[string]struct {
		overwrite bool
		in        cloudtrailLogMapping
		outFile   string
	}{
		"basic": {
			overwrite: false,
			outFile:   "basic.painless",
			in: cloudtrailLogMapping{
				defaultActor: "json.userIdentity.arn",
				defaultRelatedEntities: []string{
					"json.userIdentity.accessKeyId",
					"json.userIdentity.arn",
					"json.userIdentity.userName",
					"json.userIdentity.sessionContext.sessionIssuer.arn",
					"json.userIdentity.sessionContext.sessionIssuer.userName",
				},

				sources: []mappedSource{
					{
						sourceName: "iam",
						relatedEntityFields: []string{
							"json.requestParameters.userName",
							"json.requestParameters.roleName",
						},
						events: []mappedEvent{
							{
								eventName: "CreateUser",
								targetFields: []string{
									"json.requestParameters.userName",
									"json.requestParameters.roleName",
								},
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			if tc.overwrite {
				os.WriteFile("./testdata/tostring_output/"+tc.outFile, []byte(tc.in.toString()), 0o644)
				return
			}

			data, err := os.ReadFile("./testdata/tostring_output/" + tc.outFile)
			if err != nil {
				t.Fatal(err)
			}

			got := tc.in.toString()
			want := string(data)

			if got != want {
				t.Errorf("Result not as expected:\n%v", diff.LineDiff(got, want))
			}
		})
	}
}
