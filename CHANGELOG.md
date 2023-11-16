# [1.11.0](https://github.com/getstrm/cli/compare/v1.10.1...v1.11.0) (2023-11-16)


### Features

* **feature/pace-32:** renamed to blueprint, handles bare policy violations ([#14](https://github.com/getstrm/cli/issues/14)) ([8e79aa0](https://github.com/getstrm/cli/commit/8e79aa0d31ce7ef16e6e60c68e81383bbc83f389))

## [1.10.1](https://github.com/getstrm/cli/compare/v1.10.0...v1.10.1) (2023-11-15)


### Bug Fixes

* update pace api definitions ([57f3186](https://github.com/getstrm/cli/commit/57f318602de881f991ebf2733b76c80f3e4a6eff))

# [1.10.0](https://github.com/getstrm/cli/compare/v1.9.0...v1.10.0) (2023-11-15)


### Features

* **feature/pace-31:** ref and type removal ([#15](https://github.com/getstrm/cli/issues/15)) ([c154847](https://github.com/getstrm/cli/commit/c1548478bfe7157b36a709cbf63a23918093aa30))

# [1.9.0](https://github.com/getstrm/cli/compare/v1.8.1...v1.9.0) (2023-11-15)


### Features

* **main:** update to pace api v1.0.0-alpha.19 ([0e05743](https://github.com/getstrm/cli/commit/0e05743f13a7f0abbabc79b9b58ad2a956481fee))

## [1.8.1](https://github.com/getstrm/cli/compare/v1.8.0...v1.8.1) (2023-11-14)


### Bug Fixes

* add a go mod tidy to the pace protos version update make ([6708921](https://github.com/getstrm/cli/commit/6708921b340da35863e6cf563ccadbaa8622dd97))
* refandtype to ref and type ([a36a464](https://github.com/getstrm/cli/commit/a36a46473128fe443529b16935de027527686cde))
* remove duplicate openapi spec ([e527daa](https://github.com/getstrm/cli/commit/e527daad2b04f2963a4c54937fea8304ff472d9e))
* update pace protos version; add update script to makefile ([e72559f](https://github.com/getstrm/cli/commit/e72559fdf45e93e685cc59e1bb3ca271ef4552a6))

# [1.8.0](https://github.com/getstrm/cli/compare/v1.7.0...v1.8.0) (2023-11-10)


### Bug Fixes

* **global-policies:** merged with main ([ab043f5](https://github.com/getstrm/cli/commit/ab043f55ee6b056c28870e79e312856c886b6c28))
* **global-policies:** merged with main ([140f6c7](https://github.com/getstrm/cli/commit/140f6c7143f52d32b128aa0cadcf2cffcd823617))


### Features

* **global-policies:** add yml to autocomplete for upsert commands ([fd5986c](https://github.com/getstrm/cli/commit/fd5986c6b2b417c0c0b6a03f895bd8bc3ef76ae2))
* **global-policies:** added global transform interaction ([38c65a3](https://github.com/getstrm/cli/commit/38c65a3b5fb8cfb8bd78259485c5979dc0fe3e14))
* **global-policies:** added global transform interaction ([ced780a](https://github.com/getstrm/cli/commit/ced780a1e3a83a89d74b3949a367f95112b9aee7))
* **global-policies:** more functionality global-transforms ([18e590d](https://github.com/getstrm/cli/commit/18e590d6f5628d1cec405b481c4bc9200de5d2f3))
* **global-policies:** upgrade to latest buf release ([0c080a3](https://github.com/getstrm/cli/commit/0c080a33a84ecdee4960342060f81b0f65f72f93))
* **main:** trigger ci ([0b468a4](https://github.com/getstrm/cli/commit/0b468a489d526a8fcd64cd54d84ba5801c7ff1c3))

# [1.7.0](https://github.com/getstrm/cli/compare/v1.6.0...v1.7.0) (2023-11-09)


### Features

* **pace-24:** upgrade proto definition dependencies ([f910fe1](https://github.com/getstrm/cli/commit/f910fe1db8b67d51d2215f592b90405ebc84360b))

# [1.6.0](https://github.com/getstrm/cli/compare/v1.5.0...v1.6.0) (2023-11-07)


### Features

* **main:** trigger ci ([12e0973](https://github.com/getstrm/cli/commit/12e09734569ab26890a289c8c30bbf95b697ebf1))

# [1.5.0](https://github.com/getstrm/cli/compare/v1.4.1...v1.5.0) (2023-11-07)


### Features

* postgres processing platform integration ([c74c1f4](https://github.com/getstrm/cli/commit/c74c1f440ac7c2b9a2b8e16386b8a3d7a8d5e9bd))

## [1.4.1](https://github.com/getstrm/cli/compare/v1.4.0...v1.4.1) (2023-11-03)


### Bug Fixes

* brew url template ([046ad61](https://github.com/getstrm/cli/commit/046ad619111fc81f4b4574a59557d36bd3c8383a))

# [1.4.0](https://github.com/getstrm/cli/compare/v1.3.0...v1.4.0) (2023-11-02)


### Features

* **main:** primary key data-policy is platform id + policy id ([9a38472](https://github.com/getstrm/cli/commit/9a384726cf74c1093c03a70a83ba813c4c04ed1a))
* **main:** richer table printing ([817ece7](https://github.com/getstrm/cli/commit/817ece7a9d7892712a917cc2413f585a27b200fb))

# [1.3.0](https://github.com/getstrm/cli/compare/v1.2.2...v1.3.0) (2023-11-02)


### Features

* **main:** data-policy completion ([bf71c9c](https://github.com/getstrm/cli/commit/bf71c9c30f87dd6f734dd7c3e2a3eb3b08f4a5dd))
* **main:** primary key data-policy is platform id + policy id ([3f9731c](https://github.com/getstrm/cli/commit/3f9731ca9a2f0c78a746ffd7d75c5915c4042377))

## [1.2.2](https://github.com/getstrm/cli/compare/v1.2.1...v1.2.2) (2023-11-01)


### Bug Fixes

* **main:** removed call to publish_docs script ([f34c68c](https://github.com/getstrm/cli/commit/f34c68c3c55a2e9a51ebaa045bbc7d19a32bdcf1))
