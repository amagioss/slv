---
sidebar_position: 1
---

# Maintaining Central Profiles

## Why Centralized Profile Storage?

As you scale SLV usage across a team or organization, managing environments efficiently becomes essential. A **profile** in SLV is a logical collection of environments. When multiple team members each maintain their own SLV environment, sharing and synchronizing environment data (especially public keys) quickly becomes challenging.

Relying on ad-hoc methods (like manually sharing EDS files) leads to scalability and security issues:

- What happens when a new team member joins?
- How do you revoke access when someone leaves?
- Should every team member track all these changes independently?

To solve this, SLV supports **centralized profile management using Git**—making it easier to share, maintain, and govern profiles collaboratively.

## Git Based Profiles
Use GitHub as the single source-of-truth for profile data, storing it in a dedicated repository (or branch) that everyone can discover but only a few trusted maintainers can write to. 

### Benefits of Using Git for Profiles

Using Git repositories to manage SLV profiles brings several key advantages:

- **Native Access Control** : Use Git’s built-in permission model to allow **read-only access for all** while **restricting write access** to selected maintainers. This enforces consistency and security.
- **Version History** : Git automatically tracks every change. You can **audit, review, or roll back** any update, giving you full control over the profile’s lifecycle.
- **CI/CD Integration** : Integrate profile changes with **CI/CD pipelines** for automated validation, deployment, or notifications—helping teams enforce best practices at scale.
- **Reliability & Backup** : A Git-hosted profile provides a **redundant and reliable backup**, accessible from anywhere. Team members always fetch the latest, verified version.

### Setting up a Git Profile

To create a Git based SLV profile, follow the steps as given below: 

1. Create a **Github Repository** if you dont have one already
2. Ensure that the branch you specify to SLV is **non-empty**, meaning there is atleast an empty README.md file. SLV does not proceed with profile creation when the remote repository is empty. This is to avoid issues with non existent branches, or other discrepancies.
3. Once you have the repository set up, run the following command
    ```bash
    # For SSH based Github Access
    slv profile new git --name test --repo git@github.com:username/reponame.git

    # For HTTP based Github Access (PAT)
    slv profile new git --name test --repo https://github.com/username/reponame.git
    ```
    > **Note:** You will have to use `--token` or `--ssh-key` or `--auth-header` to authenticate with the repository.
4. Now that the new profile is created, you can add environments to the profile using the following command
    ```bash
    slv env add --env-def <ENVIRONMENT_DEFINITION_STRING>
    ```
    > **Note:** This will fail to write to the remote repository if write access is not provided.

### Updating and Maintaining a Git Profile

When a large team needs to update the profile frequently, it’s crucial to maintain tight access controls and preserve the profile’s integrity. Granting everyone direct write access to the profile branch is rarely ideal. Instead, manage changes through pull requests, ticket-based reviews, or GitHub Actions workflows. If you need to restrict GitHub access altogether—including read permissions—consider hosting the profile as an HTTP-based resource instead.

## HTTP(S) Based Profiles

Host the profile files on a secure HTTP(S) endpoint and give SLV the URL. This approach delivers read-only, always-up-to-date profile data without exposing your GitHub repository or requiring any credentials.
### Benefits of Using URL Based Profiles
- **Instant, cache-friendly downloads** : A single HTTP GET grabs the profile file immediately, skipping Git’s object negotiation and auth handshakes for sub-second fetch times.
- **No Git footprint on clients** : Consumers only need curl or a browser; works inside ephemeral CI containers and on devices where installing Git or managing keys is impractical.
- **One URL, infinite reach** : Drop the link in docs, Slack, or Terraform modules and every team, script, or service can read the canonical data without onboarding to your repo.
- **Edits locked to maintainers** : Publish automation keeps write access inside a guarded branch; everyone else sees a read-only snapshot, preserving integrity and audit trails.
- **CDN acceleration & rollback** : Fronting the endpoint with CloudFront / Cloudflare gives global low-latency caching, while versioned object keys let you revert instantly if a bad profile ships.
- **Simplified governance** : No SSH keys, PATs, or GitHub role mapping—standard web ACLs and access logs provide clear permission boundaries and compliance reporting.

### Serving Profile Files

There are two profile files that need to be served
- `environments.y(a)ml` - Contains information about the environments that are present in the profile, including the root environment.
- `settings.y(a)ml` - Contains profile settings (A feature that will be released in the coming months)

> **Hosting tip**: Think of the profile URL as a directory, not a single file.
If you register https://example.com/profile with SLV, the loader will automatically fetch:
>- https://example.com/profile/environments.yaml
>- https://example.com/profile/settings.yaml \
>So be sure those two files (or their .yml equivalents) are served at exactly those paths—the base URL itself doesn’t need to return a document, just act as the parent location for the required files.

### Serving a Git Based Profile via URL

You can also use a Git based profile as your source of truth, and use Github Pages to serve these files. 
>**Warning:** Although the profile hosts only benign environment public keys, it still reveals PII such as user email addresses, so place the directory behind a long, random slug—e.g., https://example.com/4f16e8bd/profile so SLV can still fetch /environments.yaml and /settings.yaml while opportunistic crawlers remain unaware; keep in mind this URL-obfuscation offers lightweight privacy, not a full replacement for stronger access controls if compliance later demands them.
