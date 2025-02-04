# AIPIM - AWS Interfaced Painless Identifiers Mapper


## Todo

- [ ] From data structure to painless
    - Handle multiple array levels
- [ ] From data painless to da structure
- [ ] Create API
- [ ] Create Front End
- [ ] Add Testing capabilities

## Know problems

- Arrays
    - Multiple array levels is complex to parse/convert in painless
    - For performance measures, ideally we group nested fields in a same loop, for example, we might want
    to extract `a[].b` and `a[].c`. Instead of having two loops on `a`, ideally we have one loop on `a` 
    extracting `b` and `c`. This is important from a performance standpoint