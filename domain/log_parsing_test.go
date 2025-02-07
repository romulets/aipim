package domain

import (
	_ "embed" // embed templates
	"testing"

	"github.com/andreyvit/diff"
)

//go:embed testdata/tostring_output/generated_basic.painless
var testBasicPainless string

//go:embed testdata/tostring_output/generated_complex.painless
var testComplexPainless string

func TestSimpleParser(t *testing.T) {
	tcs := map[string]struct {
		input    string
		expected *CloudtrailLogMapping
	}{
		"basic": {
			input: testBasicPainless,
			expected: &CloudtrailLogMapping{
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
			input: testComplexPainless,
			expected: &CloudtrailLogMapping{
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
			clm := &CloudtrailLogMapping{}
			err := clm.Scan(tc.input)
			if err != nil {
				t.Errorf("Expected no errors, got: %s\n", err)
			}

			got := clm.ToString()
			want := tc.expected.ToString()

			if got != want {
				t.Errorf("Result not as expected:\n%v", diff.LineDiff(got, want))
			}
		})
	}
}
