# Contributing Guidelines

The Helm mapkubeapis plugin project accepts contributions via GitHub pull requests. This document outlines the process
to help get your contribution accepted.

## Reporting a Security Issue

Most of the time, when you find a bug in the Helm mapkubeapis plugin, it should be reported using [GitHub
issues](https://github.com/helm/helm-mapkubeapis/issues). However, if you are reporting a _security
vulnerability_, please email a report to
[cncf-helm-security@lists.cncf.io](mailto:cncf-helm-security@lists.cncf.io). This will give us a
chance to try to fix the issue before it is exploited in the wild.

## Signing Your Work

This project requires for you to sign your work. To do so this project will require you to have a valid GPG [key](https://docs.github.com/en/authentication/managing-commit-signature-verification/generating-a-new-gpg-key) and have that key [bound](https://docs.github.com/en/authentication/managing-commit-signature-verification/adding-a-gpg-key-to-your-github-account) to your GitHub account. All checkins will need to have this signing to be approved and merged.

## Sign-off On Your Work

The sign-off is a simple line at the end of the explanation for a commit. All commits need to be
signed. Your signature certifies that you wrote the patch or otherwise have the right to contribute
the material. The rules are pretty simple, if you can certify the below (from
[developercertificate.org](https://developercertificate.org/)):

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
1 Letterman Drive
Suite D4700
San Francisco, CA, 94129

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.

Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```

Then you just add a line to every git commit message:

    Signed-off-by: Joe Smith <joe.smith@example.com>

Use your real name (sorry, no pseudonyms or anonymous contributions.)

If you set your `user.name` and `user.email` git configs, you can sign your commit automatically
with `git commit -s`.

Note: If your git config information is set properly then viewing the `git log` information for your
 commit will look something like this:

```
Author: Joe Smith <joe.smith@example.com>
Date:   Thu Feb 2 11:41:15 2018 -0800

    Update README

    Signed-off-by: Joe Smith <joe.smith@example.com>
```

Notice the `Author` and `Signed-off-by` lines match. If they don't your PR will be rejected by the
automated DCO check.

## Support Channels

Whether you are a user or contributor, official support channels include:

- [Issues](https://github.com/helm/helm-mapkubeapis/issues)
- Slack:
  - User: [#helm-users](https://kubernetes.slack.com/messages/C0NH30761/details/)
  - Contributor: [#helm-dev](https://kubernetes.slack.com/messages/C51E88VDG/)

Before opening a new issue or submitting a new pull request, it's helpful to search the project -
it's likely that another user has already reported the issue you're facing, or it's a known issue
that we're already aware of. It is also worth asking on the Slack channels.

## Semantic Versioning

Helm maintains a strong commitment to backward compatibility. All of our changes to protocols and
formats are backward compatible from one major release to the next. No features, flags, or commands
are removed or substantially modified (unless we need to fix a security issue).

We also remain committed to not changing publicly accessible Go library definitions inside of the `pkg/` directory of our source code in a non-backwards-compatible way.

## Issues

Issues are used as the primary method for tracking anything to do with the Helm project.

### Issue Types

There are 5 types of issues (each with their own corresponding [label](#labels)):

- `question`: These are support or functionality inquiries that we want to have a record of
  for future reference. Generally these are questions that are too complex or large to store in the
  Slack channel or have particular interest to the community as a whole. Depending on the
  discussion, these can turn into `feature` or `bug` issues.
- `enhancement`: These track specific feature requests and ideas until they are complete.
- `bug`: These track bugs with the code
- `docs`: These track problems with the documentation (i.e. missing or incomplete)

### Issue Lifecycle

The issue lifecycle is mainly driven by the core maintainers, but is good information for those
contributing to Helm. All issue types follow the same general lifecycle. Differences are noted
below.

1. Issue creation
2. Triage
    - The maintainer in charge of triaging will apply the proper labels for the issue. This includes
      labels for priority, type, and metadata (such as `good first issue`). The only issue priority
      we will be tracking is whether or not the issue is "critical." If additional levels are needed
      in the future, we will add them.
    - (If needed) Clean up the title to succinctly and clearly state the issue. Also ensure that
      proposals are prefaced with "Proposal: [the rest of the title]".
    - Add the issue to the correct milestone. If any questions come up, don't worry about adding the
      issue to a milestone until the questions are answered.
    - We attempt to do this process at least once per work day.
3. Discussion
    - Issues that are labeled `feature` or `proposal` must write a Helm Improvement Proposal (HIP).
    - Issues that are labeled as `feature` or `bug` should be connected to the PR that resolves it.
    - Whoever is working on a `feature` or `bug` issue (whether a maintainer or someone from the
      community), should either assign the issue to themselves or make a comment in the issue saying
      that they are taking it.
    - `proposal` and `support/question` issues should stay open until resolved or if they have not
      been active for more than 30 days. This will help keep the issue queue to a manageable size
      and reduce noise. Should the issue need to stay open, the `keep open` label can be added.
4. Issue closure

## How to Contribute a Patch

1. Fork the desired repo; develop and test your code changes.
2. Submit a pull request, making sure to sign your work and link the related issue.

Coding conventions and standards are explained in the [official developer
docs](https://helm.sh/docs/developers/).

## Pull Requests

Like any good open source project, we use Pull Requests (PRs) to track code changes.

### Documentation PRs

Documentation PRs will follow the same lifecycle as other PRs. They will also be labeled with the
`docs` label. For documentation, special attention will be paid to spelling, grammar, and clarity
(whereas those things don't matter *as* much for comments in code).

## Labels

The following tables define all label types used for Helm-MapKubeAPIs.

| Label | Description |
| ----- | ----------- |
| `bug` | Marks an issue as a bug or a PR as a bugfix |
| `dependencies` | Pull requests that update a dependency file |
| `documentation` | Improvements or additions to documentation |
| `duplicate` | This issue or pull request already exists |
| `enhancement` | New feature or request |
| `good first issue` | Marks an issue as a good starter issue for someone new to the project |
| `hactoberfest-accepted` | Accept for hactoberfest |
| `help wanted` | Marks an issue needs help from the community to solve |
| `invalid` | This doesn't seem right |
| `question` | Marks an issue as a support request or question |
| `wont fix` | Marks an issue as discussed and will not be implemented (or accepted in the case of a proposal) |

### Size labels

Size labels are used to indicate how "dangerous" a PR is. The guidelines below are used to assign
the labels, but ultimately this can be changed by the maintainers. For example, even if a PR only
makes 30 lines of changes in 1 file, but it changes key functionality, it will likely be labeled as
`size/L` because it requires sign off from multiple people. Conversely, a PR that adds a small
feature, but requires another 150 lines of tests to cover all cases, could be labeled as `size/S`
even though the number of lines is greater than defined below.

| Label | Description |
| ----- | ----------- |
| `size/XS` | Denotes a PR that changes 0-9 lines, ignoring generated files. Very little testing may be required depending on the change. |
| `size/S` | Denotes a PR that changes 10-29 lines, ignoring generated files. Only small amounts of manual testing may be required. |
| `size/M` | Denotes a PR that changes 30-99 lines, ignoring generated files. Manual validation should be required. |
| `size/L` | Denotes a PR that changes 100-499 lines, ignoring generated files. |
| `size/XL` | Denotes a PR that changes 500-999 lines, ignoring generated files. |
| `size/XXL` | Denotes a PR that changes 1000+ lines, ignoring generated files. |
