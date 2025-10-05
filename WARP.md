# WARP.md

Este arquivo oferece orientações para o WARP (warp.dev) ao trabalhar com código neste repositório.

## Sobre o Projeto

O `wisp` é uma biblioteca Go que oferece **value objects** robustos e imutáveis para modelagem de domínios complexos. O pacote combate a "Obsessão Primitiva" fornecendo tipos seguros para conceitos de negócios brasileiros como CPF, CNPJ, CEP, moedas e timestamps de auditoria.

### Princípios Fundamentais
- **Segurança de Tipos**: Impede criação de dados inválidos
- **Imutabilidade**: Objetos nunca mudam após criação (operações retornam novas instâncias)
- **Validação Embutida**: Regras de negócio centralizadas nos value objects
- **Extensibilidade**: Sistema de registro para valores customizados (unidades, papéis, etc.)

## Comandos de Desenvolvimento

### Teste
```bash
# Executar todos os testes
go test ./...

# Executar testes com verbose
go test -v ./...

# Executar teste de um arquivo específico
go test -v ./ -run TestCPF

# Executar testes com coverage
go test -cover ./...
```

### Build e Linting
```bash
# Verificar compilação
go build ./...

# Análise estática
go vet ./...

# Formatar código
go fmt ./...

# Verificar dependências não utilizadas
go mod tidy
```

### Comandos Específicos para Value Objects
```bash
# Testar um value object específico (ex: CPF)
go test -v ./cpf_test.go

# Ver cobertura de um tipo específico
go test -cover -run TestCPF
```

## Arquitetura do Código

### Estrutura Geral
O projeto segue um padrão **flat** com todos os value objects na raiz do package:
- Cada tipo tem seu próprio arquivo `.go` e `_test.go`
- Não há subpackages, mantendo a API simples
- Todos os tipos estão no package `wisp`

### Categorias dos Value Objects

#### Identificadores
- `UUID`, `NullableUUID`: Identificadores únicos (padrão v7)
- `CPF`, `CNPJ`: Documentos brasileiros com validação de dígitos verificadores
- `Slug`: Strings seguras para URLs

#### Financeiro
- `Money`: Valores monetários precisos (centavos como int64)
- `Currency`: Códigos de moeda validados
- `Percentage`: Porcentagens para cálculos financeiros
- `Discount`: Descontos fixos ou percentuais

#### Temporal e Geográfico
- `Date`, `DateRange`, `BirthDate`: Datas sem fuso horário
- `TimeOfDay`, `TimeRange`, `BusinessHours`: Horários
- `CreatedAt`, `UpdatedAt`, `NullableTime`: Timestamps
- `Timezone`: Fusos horários IANA
- `Latitude`, `Longitude`: Coordenadas geográficas

#### Contato e Endereçamento
- `Email`, `Phone`: Contatos validados
- `CEP`, `UF`: Endereçamento brasileiro
- `IPAddress`, `PortNumber`: Rede

#### Medidas e Quantidades
- `Weight`, `Length`: Medidas físicas com conversão de unidades
- `Quantity`: Valores numéricos com unidades personalizáveis
- `Unit`: Sistema de registro para unidades de medida

#### Auditoria
- `Audit`: Struct embutível com trilha completa (created/updated/archived/deleted)
- `AuditUser`: Identificador de usuário para auditoria
- `Version`: Controle de versão para travamento otimista

#### Tipos Genéricos e Extensíveis
- `Role`: Sistema de papéis de usuário registráveis
- `Flag[T]`, `Status[T]`: Estados binários/múltiplos personalizáveis
- `NonEmptyString`, `PositiveInt`: Primitivos com validação

### Padrões Arquiteturais

#### Factory Pattern
Todos os value objects seguem o padrão:
```go
func NewValueObject(input string) (ValueObject, error) {
    // validação e parsing
    return ValueObject{}, nil
}
```

#### Imutabilidade
Operações que "modificam" retornam novas instâncias:
```go
newMoney, err := money.Add(otherMoney)
```

#### Sistema de Registro
Tipos extensíveis usam registros globais:
```go
func init() {
    wisp.RegisterUnits("KG", "UN", "L")
    wisp.RegisterRoles("ADMIN", "USER")
}
```

#### Zero Values
Cada tipo define seu valor zero de forma explícita:
```go
var EmptyCPF CPF
var ZeroMoney = Money{}
```

### Dependências Importantes
- `github.com/marcelofabianov/fault`: Sistema de erros estruturados
- `github.com/google/uuid`: UUIDs (v7 por padrão)
- `github.com/stretchr/testify`: Framework de testes
- `golang.org/x/text`: Processamento de texto/localização

## Exemplo Prático

O diretório `example/course/domain/` demonstra como usar os value objects em uma entidade de domínio real, mostrando:
- DTO (`NewCourseInput`) com tipos primitivos
- Entidade (`Course`) com value objects do wisp
- Factory function (`NewCourse`) para validação e criação
- Métodos de comportamento que preservam a imutabilidade

## Dicas para Desenvolvimento

### Ao Adicionar Novos Value Objects
1. Crie os arquivos `nome.go` e `nome_test.go`
2. Implemente: `NewXxx()`, `String()`, `IsZero()`, `MarshalJSON()`, `UnmarshalJSON()`
3. Para tipos de banco: implemente `Value()` e `Scan()`
4. Adicione validações específicas do domínio
5. Mantenha imutabilidade em todas as operações

### Testes
- Use testify/suite para organizar testes relacionados
- Teste casos válidos, inválidos e edge cases
- Valide serialização JSON e banco de dados quando aplicável
- Teste imutabilidade das operações

### Tratamento de Erros
- Use o package `fault` para erros estruturados
- Inclua contexto relevante nos erros
- Categorize erros (Invalid, DomainViolation, etc.)