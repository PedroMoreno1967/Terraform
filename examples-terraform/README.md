# Ejemplos Terraform Corregidos

Estos ejemplos están listos para desplegar y equivalen a los ejemplos convertidos a Bicep en el repo `BICEP/examples-bicep`.

## Carpetas

- `network/vnet`
- `storage/storage-account`
- `app-service/linux-webapp`
- `kubernetes/aks-basic`
- `sql/sql-database`
- `security/key-vault`

## Uso básico

1. Crear Resource Group:

```bash
az group create -n <rg-name> -l <location>
```

2. Ejecutar Terraform en la carpeta del ejemplo:

```bash
terraform init
terraform plan -var="resource_group_name=<rg-name>" -var="location=<location>"
terraform apply -var="resource_group_name=<rg-name>" -var="location=<location>"
```

Variables especiales:
- SQL: `administrator_login_password`
- Key Vault: `admin_object_id`
