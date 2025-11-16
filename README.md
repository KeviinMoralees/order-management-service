# Order Management Service

Servicio simple de gestión de pedidos con integración a AWS DynamoDB.

## Endpoints

### GET /hello
Endpoint de prueba que retorna un mensaje de bienvenida.

**Response:**
```
Remember, you are the best v1
```

### POST /orders
Crea un nuevo pedido en DynamoDB.

**Request Body:**
```json
{
  "orderName": "Pedido #123",
  "userName": "Juan Pérez"
}
```

**Response (Success):**
```json
{
  "success": true,
  "message": "Order created successfully",
  "orderId": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response (Error):**
```json
{
  "success": false,
  "message": "Error message"
}
```

## Configuración

### Variables de Entorno

- `DYNAMODB_TABLE_NAME`: Nombre de la tabla de DynamoDB (por defecto: "Orders")

### Tabla DynamoDB

La aplicación espera una tabla de DynamoDB con la siguiente configuración:

- **Nombre de tabla**: Orders (o el que configures en `DYNAMODB_TABLE_NAME`)
- **Partition Key**: `orderId` (String)
- **Atributos**:
  - orderId (String) - ID único del pedido
  - orderName (String) - Nombre del pedido
  - userName (String) - Nombre del usuario
  - createdAt (String) - Timestamp de creación

### Crear la tabla desde AWS CLI

```bash
aws dynamodb create-table \
    --table-name Orders \
    --attribute-definitions AttributeName=orderId,AttributeType=S \
    --key-schema AttributeName=orderId,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region us-east-1
```

### Permisos IAM necesarios

La instancia EC2 donde se despliega necesita permisos de DynamoDB:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "dynamodb:PutItem",
        "dynamodb:GetItem",
        "dynamodb:Query",
        "dynamodb:Scan"
      ],
      "Resource": "arn:aws:dynamodb:*:*:table/Orders"
    }
  ]
}
```

## Desarrollo Local

### Compilar
```bash
go build -o order-management-service-ci .
```

### Ejecutar
```bash
export DYNAMODB_TABLE_NAME=Orders
./order-management-service-ci
```

### Probar endpoints

**Test /hello:**
```bash
curl http://localhost:8080/hello
```

**Test /orders:**
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "orderName": "Pedido de prueba",
    "userName": "Kevin Morales"
  }'
```

## Pipeline CI/CD

El proyecto está configurado para desplegarse automáticamente usando:
- **AWS CodeBuild**: Compila la aplicación
- **AWS CodeDeploy**: Despliega en EC2

Los archivos de configuración son:
- `buildspec.yml` - Configuración de CodeBuild
- `appspec.yml` - Configuración de CodeDeploy
- `scripts/restart.sh` - Script de reinicio del servicio

## Estructura del Proyecto

```
.
├── main.go           # Punto de entrada de la aplicación
├── handlers.go       # Handlers HTTP
├── models.go         # Modelos de datos
├── dynamodb.go       # Lógica de DynamoDB
├── go.mod            # Dependencias
├── buildspec.yml     # Config CodeBuild
├── appspec.yml       # Config CodeDeploy
└── scripts/
    └── restart.sh    # Script de despliegue
```

## Notas

- El endpoint `/hello` se mantiene igual para no romper el funcionamiento actual
- Si DynamoDB no está disponible, la app inicia pero el endpoint `/orders` no funcionará
- Los logs se guardan en `/tmp/order-management-service.log` en producción
