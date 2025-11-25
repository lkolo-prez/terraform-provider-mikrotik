---
name: Feature Request
about: Suggest a RouterOS v7 feature for the provider
title: '[FEATURE] '
labels: enhancement, routeros-v7
assignees: ''
---

## Feature Description

**RouterOS Path:** `/path/to/feature`  
**RouterOS Version:** 7.x+  
**Priority:** P0 / P1 / P2 / P3

## Use Case

Describe why this feature is needed and what problem it solves.

## Proposed Resources

List the Terraform resources that should be created:

- [ ] `mikrotik_resource_name` - Description
- [ ] `mikrotik_other_resource` - Description

## Example Configuration

```hcl
resource "mikrotik_example" "test" {
  name = "example"
  # Add attributes
}
```

## RouterOS CLI Example

```
/path/to/feature
add name=example attribute=value
```

## Implementation Complexity

- **Estimated Effort:** X weeks
- **Number of Attributes:** ~X
- **Dependencies:** List any dependent features
- **Testing Requirements:** Low / Medium / High

## Additional Context

Add any other context, screenshots, or RouterOS documentation links.
