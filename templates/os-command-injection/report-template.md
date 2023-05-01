# SSRF Vulnerability Report

## Score (CWE + CVSS)

CWE-918: Server-Side Request Forgery (SSRF)
CVSSv3 Score: X.X

## Summary

During our testing, we discovered a Server-Side Request Forgery (SSRF) vulnerability. The vulnerability allows an attacker to send arbitrary HTTP requests, bypassing any security measures in place and potentially accessing sensitive information.

## Proof of Concept

To reproduce the vulnerability, reproduce the following request:

```http
{{.request}}
```

## Impact

If exploited, this vulnerability could allow an attacker to perform actions on behalf of the application or access sensitive information that should not be accessible to them. Additionally, this vulnerability can be used to bypass security controls like firewalls and access internal resources that should not be publicly accessible.

## Remediation

To remediate this vulnerability, we recommend the following:

- Implement input validation on all user-supplied data, including URLs and IP addresses.
- Ensure that any external requests are made using allowlisted URLs or IP addresses only.
