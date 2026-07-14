# Changelog

## [2.0.1](https://github.com/dokod-fr/quadboard/compare/v2.0.0...v2.0.1) (2026-07-14)


### Bug Fixes

* ci process to create release, docker and metadata. Dockerfile using taskfile ([1aedb05](https://github.com/dokod-fr/quadboard/commit/1aedb05fb67c7a44373a23975b7574fa314e70a2))

## [2.0.0](https://github.com/dokod-fr/quadboard/compare/v1.0.0...v2.0.0) (2026-07-14)


### ⚠ BREAKING CHANGES

* Issue #14 Use base URL instead of Redirect URL

### Features

* **auth:** Add retries on connection with OIDC ([c89dffb](https://github.com/dokod-fr/quadboard/commit/c89dffb4e728dd1f65ad5bc914c0b749e879f9be))
* Issue [#14](https://github.com/dokod-fr/quadboard/issues/14) Use base URL instead of Redirect URL ([bfdb76f](https://github.com/dokod-fr/quadboard/commit/bfdb76f41e49266d84c97e9af412e48c221d3f33))


### Bug Fixes

* [#1](https://github.com/dokod-fr/quadboard/issues/1) default look-for directory for quadlet ([4e60bd7](https://github.com/dokod-fr/quadboard/commit/4e60bd778afa95d9feef7911922afff8ab11d482))
* issue [#13](https://github.com/dokod-fr/quadboard/issues/13) default quadlet dirs to look-for ([6eead5a](https://github.com/dokod-fr/quadboard/commit/6eead5abba517ea7c14233a16c53d0c90783b72d))

## 1.0.0 (2026-07-13)


### Features

* Add label management ([8fbef56](https://github.com/dokod-fr/quadboard/commit/8fbef568df543e311b007f3ccc59879384f29fd9))
* bootstrap application ([e636575](https://github.com/dokod-fr/quadboard/commit/e636575343494fe56a3e859dbd1ca3aa2bf2c2d6))
* **cli:** initial Cobra setup with logging and build metadata ([538f34c](https://github.com/dokod-fr/quadboard/commit/538f34cfdf4a7f63f1327d19dfb79f042b38d675))
* **config:** init config though yaml or environment vars ([5a105c9](https://github.com/dokod-fr/quadboard/commit/5a105c972c6733aeff066675bf2d941d989053dc))
* **config:** introduce configuration system with defaults and validation ([548ccea](https://github.com/dokod-fr/quadboard/commit/548ccea446e718a47d83d932c1b6bc9754dac122))
* **config:** Optional secure config for auth args ([0640d98](https://github.com/dokod-fr/quadboard/commit/0640d9894c0bc68dba555c714e93bbcee7736343))
* **config:** Update config and logging ([e737bc8](https://github.com/dokod-fr/quadboard/commit/e737bc88a18eddf1703daa95810c9e81769d2131))
* **domain:** introduce resource model and dashboard resource card ([50de201](https://github.com/dokod-fr/quadboard/commit/50de201ebd18917c8e499add9068173d228e8e26))
* **http:** Add auth oidc ([70a4066](https://github.com/dokod-fr/quadboard/commit/70a40667ec1b1fd2b3300a1bdf70cd5dc837c3ff))
* **http:** introduce chi server with health and version endpoints ([bf03680](https://github.com/dokod-fr/quadboard/commit/bf03680ee2c1b81dd88c3c3021d153f15852825f))
* **http:** Manage assets ([1f07a02](https://github.com/dokod-fr/quadboard/commit/1f07a0207605105accd1e74f671e975df4e226af))
* Initialize v0 ([59513c5](https://github.com/dokod-fr/quadboard/commit/59513c5933dc26619e7b874705869b789c99c489))
* Initialize v0 ([074d0c5](https://github.com/dokod-fr/quadboard/commit/074d0c561e2df4929ccbc751c399c20e39781f5a))
* **logic:** Add first version of quadlet  providers ([9ae4684](https://github.com/dokod-fr/quadboard/commit/9ae46841e1a223a6f63c7198465efbebde585def))
* **ui:** add templ-based dashboard with layout and home page ([2fe4c9c](https://github.com/dokod-fr/quadboard/commit/2fe4c9c50c80681f7764410ed7bfb66453ebe643))
* **ui:** First step ui design ([653fbf1](https://github.com/dokod-fr/quadboard/commit/653fbf1d2fd79c0ea0549eca405cdc7e8ec0923d))
* **ui:** idea from AI for design ([9aa05d7](https://github.com/dokod-fr/quadboard/commit/9aa05d7acc697efca8e23284b4963d4dfd181511))
* **ui:** update [#1](https://github.com/dokod-fr/quadboard/issues/1) to be closer to the branding ([3bb2384](https://github.com/dokod-fr/quadboard/commit/3bb23846bc98967c0c5b9c9bbfe2292922aed5c4))


### Bug Fixes

* **authen:** Fix PKCE exchange and add mock with Keycloak ([93d4e39](https://github.com/dokod-fr/quadboard/commit/93d4e39396c877bf9f4b62c6b6d92dead2075ff9))
