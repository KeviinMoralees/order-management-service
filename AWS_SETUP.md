# Configuración de AWS para Order Management Service

## Paso 1: Crear la Tabla DynamoDB

### Opción A: Desde la Consola de AWS

1. Ve a **DynamoDB** en la consola de AWS
2. Click en **"Create table"**
3. Configuración:
   - **Table name**: `Orders`
   - **Partition key**: `orderId` (String)
   - **Table settings**: Default settings (On-demand)
4. Click en **"Create table"**

### Opción B: Desde AWS CLI

```bash
aws dynamodb create-table \
    --table-name Orders \
    --attribute-definitions AttributeName=orderId,AttributeType=S \
    --key-schema AttributeName=orderId,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region us-east-1
```

Verifica que se creó:
```bash
aws dynamodb describe-table --table-name Orders --region us-east-1
```

## Paso 2: Configurar Permisos IAM

Tu instancia EC2 necesita permisos para acceder a DynamoDB.

### Opción A: Agregar política al rol IAM existente de tu EC2

1. Ve a **IAM** → **Roles**
2. Busca el rol que usa tu instancia EC2 (probablemente algo como `EC2-CodeDeploy-Role`)
3. Click en **"Add permissions"** → **"Create inline policy"**
4. En JSON, pega:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "DynamoDBAccess",
      "Effect": "Allow",
      "Action": [
        "dynamodb:PutItem",
        "dynamodb:GetItem",
        "dynamodb:Query",
        "dynamodb:Scan",
        "dynamodb:UpdateItem"
      ],
      "Resource": [
        "arn:aws:dynamodb:us-east-1:*:table/Orders"
      ]
    }
  ]
}
```

5. Nombra la política como `DynamoDB-Orders-Access`
6. Click en **"Create policy"**

### Opción B: Crear un nuevo rol (si no tienes uno)

```bash
# Crear la política
aws iam create-policy \
    --policy-name DynamoDB-Orders-Access \
    --policy-document '{
      "Version": "2012-10-17",
      "Statement": [{
        "Effect": "Allow",
        "Action": [
          "dynamodb:PutItem",
          "dynamodb:GetItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ],
        "Resource": "arn:aws:dynamodb:*:*:table/Orders"
      }]
    }'

# Adjuntar la política al rol de tu EC2
aws iam attach-role-policy \
    --role-name TU-ROL-EC2 \
    --policy-arn arn:aws:iam::TU-ACCOUNT-ID:policy/DynamoDB-Orders-Access
```

## Paso 3: Configurar Variable de Entorno (Opcional)

Si quieres usar un nombre de tabla diferente a "Orders":

### En EC2 (durante el despliegue)

Edita el script `scripts/restart.sh` y agrega antes de ejecutar la app:

```bash
export DYNAMODB_TABLE_NAME=MiTablaPedidos
```

### O configura en tu pipeline

Agrega en tu configuración de CodeDeploy o en el user data de EC2:

```bash
echo 'export DYNAMODB_TABLE_NAME=Orders' >> /etc/environment
```

## Paso 4: Desplegar la Nueva Versión

Tu pipeline de CI/CD se encargará de:
1. Compilar el código con las nuevas dependencias
2. Crear el artefacto con `go.mod` y `go.sum` actualizados
3. Desplegar en tu EC2

**No necesitas modificar `buildspec.yml` ni `appspec.yml`** - ya funcionan correctamente.

## Paso 5: Verificar el Despliegue

Después del despliegue, verifica:

```bash
# SSH a tu instancia EC2
ssh ec2-user@TU-EC2-IP

# Ver los logs
tail -f /tmp/order-management-service.log

# Probar el endpoint
curl http://localhost:8080/hello
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"orderName": "Test Order", "userName": "Test User"}'
```

## Troubleshooting

### Error: "Could not initialize DynamoDB"

**Causa**: La aplicación no puede conectarse a DynamoDB.

**Soluciones**:
1. Verifica que la tabla existe:
   ```bash
   aws dynamodb describe-table --table-name Orders
   ```

2. Verifica que el rol IAM tiene permisos:
   ```bash
   aws iam list-attached-role-policies --role-name TU-ROL-EC2
   ```

3. Verifica la región en la configuración de AWS

### Error: "AccessDeniedException"

**Causa**: El rol IAM no tiene permisos suficientes.

**Solución**: Revisa el Paso 2 y asegúrate de que la política está correctamente adjunta al rol.

### La aplicación inicia pero /orders no funciona

**Causa**: DynamoDB no está disponible pero la app permite iniciar de todas formas.

**Solución**: Revisa los logs en `/tmp/order-management-service.log` para ver el error específico.

## Monitoreo

### Ver items en DynamoDB

```bash
aws dynamodb scan --table-name Orders
```

### Ver logs en CloudWatch (si está configurado)

```bash
aws logs tail /aws/ec2/order-management-service --follow
```

## Costos Estimados

- **DynamoDB (On-demand)**: ~$0.25 por millón de escrituras
- Sin cambios en los costos de EC2, CodeBuild o CodeDeploy existentes

## Próximos Pasos (Opcional)

- Agregar endpoints para listar pedidos (GET /orders)
- Agregar endpoint para obtener un pedido específico (GET /orders/:id)
- Implementar paginación
- Agregar índices secundarios en DynamoDB para búsquedas más eficientes

