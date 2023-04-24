# bc-wallet-common-lib-postgres

## Description 

Library for manage postgresql config and connection.

Library contains:
* common postgresql config struct
* connection manager
* small function-helpers for work with transaction statement

## Usage example

Examples of create connection and write database communication code

### Config and connection
```go
package main

import (
	commonEnvConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/envconfig"
    commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"
)

type VaultWrappedConfig struct {
	*commonVault.BaseConfig
	*commonVaultTokenClient.AuthConfig
}

func main() {
	// vault prepare 
	vaultSrv, err := commonVault.NewService(ctx, vaultCfg, vaultClientSrv)
	if err != nil {
		panic(err)
	}

	_, err = vaultSrv.Login(ctx)
	if err != nil {
		panic(err)
	}

	pgConfig := commonPostgres.PostgresConfig{}
	pgCfgPreparerSrv := commonEnvConfig.NewConfigManager()
	err = pgCfgPreparerSrv.PrepareTo(pgConfig).With(vaultSrv).Do(ctx)
	if err != nil {
		panic(err)
	}

	pgConn := commonPostgres.NewConnection(context.Background(), pgConfig, loggerSvc)
	_, err = pgConn.Connect()
	if err != nil {
		panic(err)
	}
}


```

## Contributors

* Author and maintainer - [@gudron (Alex V Kotelnikov)](https://github.com/gudron)

## Licence

Proprietary license

Switcher to proprietary license from MIT - [CHANGELOG.MD - v0.0.7](./CHANGELOG.md)