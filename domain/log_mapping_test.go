package domain

import (
	"os"
	"testing"

	"github.com/andreyvit/diff"
)

func TestToString(t *testing.T) {
	tcs := map[string]struct {
		writeGenerated bool
		in             cloudtrailLogMapping
		outFile        string
	}{
		"basic": {
			writeGenerated: true,
			outFile:        "basic.painless",
			in: cloudtrailLogMapping{
				defaultActor: "json.userIdentity.arn",
				defaultRelatedEntities: []string{
					"json.userIdentity.accessKeyId",
					"json.userIdentity.arn",
					"json.userIdentity.userName",
					"json.userIdentity.sessionContext.sessionIssuer.arn",
					"json.userIdentity.sessionContext.sessionIssuer.userName",
					"json.resources[].ARN",
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

		"complex": {
			writeGenerated: true,
			outFile:        "complex.painless",
			in: cloudtrailLogMapping{
				defaultActor: "json.userIdentity.arn",
				defaultRelatedEntities: []string{
					"json.userIdentity.accessKeyId",
					"json.userIdentity.arn",
					"json.userIdentity.userName",
					"json.userIdentity.sessionContext.sessionIssuer.arn",
					"json.userIdentity.sessionContext.sessionIssuer.userName",
					"json.resources[].ARN",
				},

				sources: []mappedSource{
					{
						sourceName: "iam",
						relatedEntityFields: []string{
							"json.requestParameters.userName",
						},
						events: []mappedEvent{
							{
								eventName: "CreateUser",
								targetFields: []string{
									"json.requestParameters.userName",
								},
							},
							{
								eventName: "DeleteUser",
								targetFields: []string{
									"json.requestParameters.userName",
								},
							},
							{
								eventName: "CreateRole",
								targetFields: []string{
									"json.requestParameters.roleName",
								},
							},
						},
					},

					{
						sourceName: "ec2",
						events: []mappedEvent{
							{
								eventName: "StartInstances",
								targetFields: []string{
									"json.requestParameters.roleName",
									"json.responseElements.instancesSet.items[].instanceId",
								},
							},

							{
								eventName: "StopInstances",
								targetFields: []string{
									"json.responseElements.instancesSet.items[].instanceId",
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
			if tc.writeGenerated {
				os.WriteFile("./testdata/tostring_output/generated_"+tc.outFile, []byte(tc.in.toString()), 0o644)
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
