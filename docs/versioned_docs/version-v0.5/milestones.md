---
id: milestones
title: Milestones and Roadmap
---

## [v0.7.0-alpha.1](https://github.com/ory/kratos/milestone/9)

_This milestone does not have a description._

### [Bug](https://github.com/ory/kratos/labels/bug)

Something is not working.

#### Issues

- [ ] Do not create system errors on duplicate credentials when linking oidc
      providers ([kratos#694](https://github.com/ory/kratos/issues/694))
- [ ] Typescript ErrorContainer type is incorrect
      ([kratos#782](https://github.com/ory/kratos/issues/782))
- [ ] Refresh Sessions Without Having to Log In Again
      ([kratos#615](https://github.com/ory/kratos/issues/615)) -
      [@hackerman](https://github.com/aeneasr)

### [Feat](https://github.com/ory/kratos/labels/feat)

New feature or request.

#### Issues

- [ ] Selfservice account deletion
      ([kratos#596](https://github.com/ory/kratos/issues/596))
- [ ] Implement Hydra integration
      ([kratos#273](https://github.com/ory/kratos/issues/273))
- [ ] Self-service GDPR identity export
      ([kratos#658](https://github.com/ory/kratos/issues/658))
- [ ] Admin/Selfservice session management
      ([kratos#655](https://github.com/ory/kratos/issues/655))
- [ ] Webhook notification based system
      ([kratos#776](https://github.com/ory/kratos/issues/776))
- [ ] improve multi schema handling in different auth flows
      ([kratos#765](https://github.com/ory/kratos/issues/765))
- [ ] More meta information about the managed identity
      ([kratos#820](https://github.com/ory/kratos/issues/820))
- [ ] Add i18n support to mail templates
      ([kratos#834](https://github.com/ory/kratos/issues/834))
- [ ] Add option for disabling registration
      ([kratos#882](https://github.com/ory/kratos/issues/882))
- [ ] Implement React SPA sample app
      ([kratos#668](https://github.com/ory/kratos/issues/668)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] Rename strategy to method in internal APIs and Documentation
      ([kratos#683](https://github.com/ory/kratos/issues/683)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Double slash in URLs causes CSRF issues
      ([kratos#779](https://github.com/ory/kratos/issues/779))

### [Rfc](https://github.com/ory/kratos/labels/rfc)

A request for comments to discuss and share ideas.

#### Issues

- [ ] improve multi schema handling in different auth flows
      ([kratos#765](https://github.com/ory/kratos/issues/765))

### [Blocking](https://github.com/ory/kratos/labels/blocking)

Blocks milestones or other issues or pulls.

#### Issues

- [ ] Implement Hydra integration
      ([kratos#273](https://github.com/ory/kratos/issues/273))

## [v0.6.0-alpha.1](https://github.com/ory/kratos/milestone/8)

_This milestone does not have a description._

### [Bug](https://github.com/ory/kratos/labels/bug)

Something is not working.

#### Issues

- [ ] Unmable to use Auth0 as a generic OIDC provider
      ([kratos#609](https://github.com/ory/kratos/issues/609))
- [ ] Password reset emails sent twice by each of the two kratos pods in my
      cluster ([kratos#652](https://github.com/ory/kratos/issues/652))
- [ ] Investigate why smtps fails but smtp does not
      ([kratos#781](https://github.com/ory/kratos/issues/781)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] Add randomized constant time to every login request
      ([kratos#832](https://github.com/ory/kratos/issues/832))
- [ ] Fetching a settings request after error is missing identity data
      ([kratos#689](https://github.com/ory/kratos/issues/689)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] Reloading config values does not work
      ([kratos#804](https://github.com/ory/kratos/issues/804)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] Quickstart is failing to mount volume kratos.yml when SELinux is enabled
      using Podman ([kratos#831](https://github.com/ory/kratos/issues/831)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Sending JSON to complete oidc/password strategy flows causes CSRF issues
      ([kratos#378](https://github.com/ory/kratos/issues/378))
- [x] Building From Source fails
      ([kratos#711](https://github.com/ory/kratos/issues/711))

### [Feat](https://github.com/ory/kratos/labels/feat)

New feature or request.

#### Issues

- [ ] Implement Security Questions MFA
      ([kratos#469](https://github.com/ory/kratos/issues/469))
- [ ] Do not send credentials to hooks
      ([kratos#77](https://github.com/ory/kratos/issues/77)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] Implement immutable keyword in JSON Schema for Identity Traits
      ([kratos#117](https://github.com/ory/kratos/issues/117))
- [ ] Add filters to admin api
      ([kratos#249](https://github.com/ory/kratos/issues/249))
- [ ] Feature Request: Webhooks
      ([kratos#271](https://github.com/ory/kratos/issues/271))
- [ ] Support email verification paswordless login
      ([kratos#286](https://github.com/ory/kratos/issues/286))
- [ ] Prevent account enumeration for profile updates
      ([kratos#292](https://github.com/ory/kratos/issues/292)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] Implement identity state and administrative deactivation, deletion of
      identities ([kratos#598](https://github.com/ory/kratos/issues/598)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] SMTP Error spams the server logs
      ([kratos#402](https://github.com/ory/kratos/issues/402))
- [ ] How to sign in with Twitter
      ([kratos#517](https://github.com/ory/kratos/issues/517))
- [ ] Add ability to import user credentials
      ([kratos#605](https://github.com/ory/kratos/issues/605)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] Throttling repeated login requests
      ([kratos#654](https://github.com/ory/kratos/issues/654))
- [ ] Require identity deactivation before administrative deletion
      ([kratos#657](https://github.com/ory/kratos/issues/657))
- [ ] CSRF failure should start a new login/registration flow
      ([kratos#821](https://github.com/ory/kratos/issues/821)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Gracefully handle CSRF errors
      ([kratos#91](https://github.com/ory/kratos/issues/91)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Support remote argon2 execution
      ([kratos#357](https://github.com/ory/kratos/issues/357)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Feature request: adjustable thresholds on how many times a password has
      been in a breach according to haveibeenpwned
      ([kratos#450](https://github.com/ory/kratos/issues/450))
- [x] Add return_to after logout
      ([kratos#702](https://github.com/ory/kratos/issues/702)) -
      [@Patrik](https://github.com/zepatrik)
- [x] Write CLI helper for recommending Argon2 parameters
      ([kratos#723](https://github.com/ory/kratos/issues/723)) -
      [@Patrik](https://github.com/zepatrik)
- [x] Add possibility to configure the "claims" query parameter in the auth_url
      of OIDC providers to request individial id_token claims
      ([kratos#735](https://github.com/ory/kratos/issues/735))

#### Pull Requests

- [ ] feat: add selinux compatible quickstart config
      ([kratos#889](https://github.com/ory/kratos/pull/889)) -
      [@hackerman](https://github.com/aeneasr)

### [Docs](https://github.com/ory/kratos/labels/docs)

Affects documentation.

#### Issues

- [ ] Document that identity information (traits, etc) are available to token
      holders and backend systems
      ([kratos#43](https://github.com/ory/kratos/issues/43)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] Config JSON Schema needs example values
      ([kratos#179](https://github.com/ory/kratos/issues/179)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Elaborate on security practices against DoS and Brute Force
      ([kratos#134](https://github.com/ory/kratos/issues/134)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Building From Source fails
      ([kratos#711](https://github.com/ory/kratos/issues/711))

### [Rfc](https://github.com/ory/kratos/labels/rfc)

A request for comments to discuss and share ideas.

#### Issues

- [ ] Introduce prevent extension in Identity JSON schema
      ([kratos#47](https://github.com/ory/kratos/issues/47))
