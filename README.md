# AIPIM - AWS Interfaced Painless Identifiers Mapper


## Todo

- [x] From data structure to painless
- [ ] From data painless to da structure - in dev @kubasobon
- [ ] Create API
- [ ] Create Front End - in dev @romulets
- [ ] Add Testing capabilities

## Known problems

- Arrays
    - Multiple array levels is complex to parse/convert in painless
    - For performance measures, ideally we group nested fields in a same loop, for example, we might want
    to extract `a[].b` and `a[].c`. Instead of having two loops on `a`, ideally we have one loop on `a` 
    extracting `b` and `c`. This is important from a performance standpoint