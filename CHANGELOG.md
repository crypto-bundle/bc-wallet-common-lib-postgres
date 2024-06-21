# Change Log

## [v0.0.9] - 21.06.2024
### Changed
* Removed zap.Logger dependency. Replaced by std logger
* Fixed bug in postgres DSN format in connection flow
* Added License banner to all *.go files
* Changed MIT License to MIT NON-AI License

## [v0.0.8] - 16.04.2024
### Changed
* Bump golang version 1.19 -> 1.22

## [v0.0.7] - 15.04.2024
### Changed
* Changed env variables name for supporting kubernetes naming standard:
  * DB_HOST -> POSTGRESQL_SERVICE_HOST
  * DB_PORT -> POSTGRESQL_SERVICE_PORT

## [v0.0.6] - 14.04.2024
### Added
* Added support of healthcheck flow, which required by [lib-healthcheck](https://github.com/crypto-bundle/bc-wallet-common-lib-healthcheck)

## [v0.0.5] - 09.02.2024
### Changed
* Added info about helm-chart to [CHANGELOG.md](./CHANGELOG.md) file
* Added deployment commands of helm-chart to Makefile
* Tx-statement changes
  * Added _BeginReadCommittedTxRollbackOnError_ helper function
  * Added _BeginReadUncommittedTxRollbackOnError_ helper function
* Changed content of LICENSE file
* Fixed mistake in code examples
* Changed go-namespace

## [v0.0.4] - 08.05.2023
### Added
* Added PostgreSQL helm-chart for local development. Chart cloned from [official Bitnami repository](https://github.com/bitnami/charts/tree/main/bitnami/postgresql)

## [v0.0.3] - 22.04.2023 19:48 MSK
### Changed
* Modified PostgresSQL connection config - added _secret_ golang tag.
* Some changes in [README.me](./README.md) file
  * Added contributors sections
  * Added code-examples
* Tx-statement changes
  * Helper function _BeginTxWithUnlessCommittedRollback_ renamed to _BeginTxWithRollbackOnError_

## [v0.0.2] - 06.02.2023 22:04 MSK
### Changed
* Lib-postgresl moved to another repository - https://github.com/crypto-bundle/bc-wallet-common-lib-postgres
* Added MIT license
* Transactional-statement helper-functions changes
  * Added BeginTxWithUnlessCommittedRollback helper
### Fixed
* Bug in tx-statement management - missing err return

## [v0.0.1] - 08.03.2022 00:48 MSK
### Added
* Added Postgresql connection wrapper 
* Added transactional-statement helper functions