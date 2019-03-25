# Release History

### 0.0.15 / 2019-03-25
Fixing path templates to make sure curly braces match
Use Go modules

### 0.0.14 / 2019-03-25
Serve Testing service CLI service
Ensure all path templates start with `/`

### 0.0.13 / 2019-03-01
Fix issue which tombstones users.

### 0.0.12 / 2019-02-20
Remove google.api.client_package proto annotations.

### 0.0.11 / 2019-02-19
Update GAPIC config proto annotations.

### 0.0.10 / 2019-01-29
- Expose messaging and identity services when running `gapic-showcase run`.
- Refactor `Echo.WaitRequest` to follow API style for denoting time to live.
- Use GCLI Generated Code for the CLI cmd.
