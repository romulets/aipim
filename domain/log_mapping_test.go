package domain

import (
	"os"
	"testing"

	"github.com/andreyvit/diff"
)

func TestToString(t *testing.T) {
	tcs := map[string]struct {
		writeGenerated bool
		in             CloudtrailLogMapping
		outFile        string
	}{
		"basic": {
			writeGenerated: true,
			outFile:        "basic.painless",
			in: CloudtrailLogMapping{
				DefaultActor: "json.userIdentity.arn",
				DefaultRelatedEntities: []string{
					"json.userIdentity.accessKeyId",
					"json.userIdentity.arn",
					"json.userIdentity.userName",
					"json.userIdentity.sessionContext.sessionIssuer.arn",
					"json.userIdentity.sessionContext.sessionIssuer.userName",
					"json.resources[].ARN",
				},

				Sources: []MappedSource{
					{
						SourceName: "iam",
						RelatedEntityFields: []string{
							"json.requestParameters.userName",
							"json.requestParameters.roleName",
						},
						Events: []MappedEvent{
							{
								EventName: "CreateUser",
								TargetFields: []string{
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
			in: CloudtrailLogMapping{
				DefaultActor: "json.userIdentity.arn",
				DefaultRelatedEntities: []string{
					"json.userIdentity.accessKeyId",
					"json.userIdentity.arn",
					"json.userIdentity.userName",
					"json.userIdentity.sessionContext.sessionIssuer.arn",
					"json.userIdentity.sessionContext.sessionIssuer.userName",
					"json.resources[].ARN",
				},

				Sources: []MappedSource{
					{
						SourceName: "iam",
						RelatedEntityFields: []string{
							"json.requestParameters.userName",
						},
						Events: []MappedEvent{
							{
								EventName: "CreateUser",
								TargetFields: []string{
									"json.requestParameters.userName",
								},
							},
							{
								EventName: "DeleteUser",
								TargetFields: []string{
									"json.requestParameters.userName",
								},
							},
							{
								EventName: "CreateRole",
								TargetFields: []string{
									"json.requestParameters.roleName",
								},
							},
						},
					},

					{
						SourceName: "ec2",
						Events: []MappedEvent{
							{
								EventName: "StartInstances",
								TargetFields: []string{
									"json.requestParameters.roleName",
									"json.responseElements.instancesSet.items[].instanceId",
								},
							},

							{
								EventName: "StopInstances",
								TargetFields: []string{
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
				os.WriteFile("./testdata/tostring_output/generated_"+tc.outFile, []byte(tc.in.ToString()), 0o644)
			}

			data, err := os.ReadFile("./testdata/tostring_output/" + tc.outFile)
			if err != nil {
				t.Fatal(err)
			}

			got := tc.in.ToString()
			want := string(data)

			if got != want {
				t.Errorf("Result not as expected:\n%v", diff.LineDiff(got, want))
			}
		})
	}
}
