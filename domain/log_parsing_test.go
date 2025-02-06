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
		expected *cloudtrailLogMapping
	}{
		"basic": {
			input: testBasicPainless,
			expected: &cloudtrailLogMapping{
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
			input: testComplexPainless,
			expected: &cloudtrailLogMapping{
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
			clm := &cloudtrailLogMapping{}
			clm.scan(tc.input)

			got := clm.toString()
			want := tc.expected.toString()

			if got != want {
				t.Errorf("Result not as expected:\n%v", diff.LineDiff(got, want))
			}
		})
	}
}
