# AIPIM - AWS Interfaced Painless Identifiers Mapper


## Todo

- [x] From data structure to painless
- [x] From data painless to da structure
- [x] Create API
- [x] Create Front Endx
- [ ] Add Testing capabilities

## Known problems

- Arrays
    - Multiple array levels is complex to parse/convert in painless
    - For performance measures, ideally we group nested fields in a same loop, for example, we might want
    to extract `a[].b` and `a[].c`. Instead of having two loops on `a`, ideally we have one loop on `a` 
    extracting `b` and `c`. This is important from a performance standpoint
    - read Loops with multilevel access after loop (e.g. `json.responseElements.securityGroupRuleSet.items[].referencedGroupInfo.groupId`)
    - read Loops with multilevel access after loop (e.g. `json.responseElements.securityGroupRuleSet.items[].referencedGroupInfo.groupId`) because it needs to have optional markers, it needs to become `referencedGroupInfo?.groupId`
    - arrays without further values, e.g `json.requestParameters.alarmNames[]`

- Special Cases
    - dash containing sources e.g `ec2-instance-connect`
    - lambda that has versioning at the end of the event, it's is not `AddPermission` but `AddPermission20250101`
    - `SendCommand` logic of if pointed to all instances (`*`), add full account as target
    - `AssumeRole` has different target based on type of user type

- Actor per event
    - not yet implemented