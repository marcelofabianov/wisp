# wisp

[![Go Report Card](https://goreportcard.com/badge/github.com/marcelofabianov/wisp)](https://goreportcard.com/report/github.com/marcelofabianov/wisp)
[![Go Reference](https://pkg.go.dev/badge/github.com/marcelofabianov/wisp.svg)](https://pkg.go.dev/github.com/marcelofabianov/wisp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Um conjunto de *value objects* robustos e imutáveis para Go, projetado para modelagem de domínios complexos com segurança de tipos, com foco em cenários de negócios brasileiros.

## Visão Geral

Este pacote nasceu da necessidade de combater a "Obsessão Primitiva" (*Primitive Obsession*), onde conceitos de negócio ricos como CPF, Dinheiro ou um Período de Datas são representados por tipos primitivos (`string`, `int`, `float64`). Essa prática comum leva a código com baixa expressividade e é uma fonte constante de bugs, pois a validação fica espalhada e duplicada por todo o sistema.

`wisp` oferece tipos de valor que encapsulam essas regras de negócio, garantindo que um objeto, uma vez criado, esteja **sempre em um estado válido**. Isso torna o código mais seguro, mais limpo e muito mais fácil de manter.

### Principais Vantagens

* **Segurança de Tipos**: Impede a criação de dados inválidos em tempo de compilação ou na inicialização. Um `wisp.CPF` nunca conterá um valor inválido.
* **Imutabilidade**: Os objetos `wisp` são imutáveis. Operações que modificam o valor (ex: `Add` em `Money`) retornam uma nova instância, tornando o código seguro para uso em ambientes concorrentes e livre de efeitos colaterais.
* **Validação Embutida**: Regras complexas, como o cálculo de dígitos verificadores de CPF/CNPJ e a validação de formatos, estão centralizadas e testadas.
* **API Expressiva**: Métodos como `Formatted()`, `IsMobile()`, `MultiplyByMoney()` e `Age()` tornam o código mais legível e explícito sobre sua intenção.
* **Extensibilidade**: Tipos como `Unit`, `Role` e `FileExtension` permitem que o consumidor do pacote defina seus próprios conjuntos de valores válidos, adaptando a biblioteca ao seu domínio específico.

## Tipos Disponíveis

| Tipo | Descrição |
| :--- | :--- |
| **Identificadores** | |
| `UUID` | Wrapper para `uuid.UUID` (padrão v7) para identificadores únicos. |
| `NullableUUID` | Um `wisp.UUID` que pode ser nulo, ideal para chaves estrangeiras opcionais. |
| `CPF` | CPF brasileiro com validação de dígitos verificadores e formatação. |
| `CNPJ` | CNPJ brasileiro com validação de dígitos verificadores e formatação. |
| `Slug`| Uma string otimizada e segura para ser usada em URLs. |
| **Financeiro** | |
| `Currency` | Código de moeda (ex: BRL) validado a partir de uma lista registrável. |
| `Money` | Representa um valor monetário com segurança, evitando `float64`. |
| `Percentage` | Tipo de porcentagem preciso para cálculos financeiros seguros. |
| `Discount` | Objeto polimórfico para descontos (fixos ou percentuais). |
| **Medidas Físicas** | |
| `Weight`| Medida de massa com unidades (kg, g, lb) e conversão segura. |
| `Length`| Medida de comprimento com unidades (m, cm, ft) e conversão segura. |
| `Quantity`| Valor numérico com unidade de medida extensível e precisão configurável. |
| `Unit` | Sistema de registro para unidades de medida (`KG`, `UN`, etc.). |
| **Rede & Formatos**| |
| `IPAddress`| Endereço de rede IPv4 ou IPv6 validado. |
| `PortNumber`| Número de porta de rede com validação de intervalo (1-65535). |
| `FileExtension` | Extensão de arquivo (ex: "pdf") validada contra uma lista registrável. |
| `MIMEType` | Tipo de mídia (ex: "image/jpeg") validado contra uma lista registrável. |
| `Color` | Representação e validação de cores no formato hexadecimal. |
| **Contato & Endereçamento**| |
| `Email`| Endereço de e-mail validado. |
| `Phone`| Telefone brasileiro (fixo ou móvel) com validação e formatação. |
| `CEP`| CEP brasileiro com validação de formato e formatação. |
| `UF` | Unidade Federativa brasileira validada a partir de uma lista registrável. |
| **Geolocalização** | |
| `Longitude`| Coordenada geográfica de longitude, com validação de intervalo (-180 a 180). |
| `Latitude`| Coordenada geográfica de latitude, com validação de intervalo (-90 a 90). |
| **Temporal** | |
| `Date`| Representa uma data de calendário (YYYY-MM-DD) sem fuso horário. |
| `DateRange` | Um período entre duas datas, com validação de `start <= end`. |
| `BirthDate`| Uma data de nascimento que não pode ser no futuro, com cálculos de idade. |
| `Day` | Um dia do mês (1-31) para eventos recorrentes. |
| `DayOfWeek` | Um dia da semana (Domingo, Segunda, etc.) de forma segura. |
| `TimeOfDay` | Representa uma hora do dia (HH:MM) sem data. |
| `TimeRange` | Um intervalo de tempo entre duas `TimeOfDay`. |
| `BusinessHours` | Modelo completo de horário comercial para uma semana. |
| `Timezone` | Representa um fuso horário IANA (ex: "America/Sao_Paulo") de uma lista registrável. |
| `CreatedAt`, `UpdatedAt` | Timestamps de criação e modificação. |
| `NullableTime`| Um `time.Time` que pode ser nulo, para campos como `deleted_at`. |
| **Auditoria & Domínio** | |
| `Audit` | Struct embutível com a trilha de auditoria completa. |
| `AuditUser`| Identificador de usuário de auditoria (e-mail ou "system"). |
| `Version` | Versão numérica para travamento otimista. |
| `Role` | Sistema de registro extensível para papéis de usuário (`ADMIN`, etc.). |
| `Preferences` | Objeto seguro para armazenar dados JSON flexíveis (chave-valor). |
| `Flag[T]` | Tipo genérico para representar um estado binário com valores customizados. |
| **Primitivos Seguros** | |
| `NonEmptyString` | Uma `string` que garante não ser vazia após remover espaços. |
| `PositiveInt` | Um `int` que garante ser sempre maior que zero. |

## Instalação

```sh
go get github.com/marcelofabianov/wisp
```

### Configuração do Pacote

Alguns tipos em wisp são extensíveis. É recomendado configurar os valores padrão na inicialização da sua aplicação.

```go
func init() {
    // Define a maioridade padrão para 21 anos (padrão é 18)
    wisp.SetLegalAge(21)

    // Define a precisão padrão para quantidades (padrão é 3)
    wisp.SetDefaultPrecision(4)

    // Registra as unidades de medida que seu domínio utilizará
    wisp.RegisterUnits("KG", "UN", "L", "M2", "H")

    // Registra os papéis de usuário que seu domínio utilizará
    wisp.RegisterRoles("ADMIN", "TEACHER", "STUDENT")

    // Registra as extensões de arquivo permitidas
    wisp.RegisterFileExtensions("pdf", "xml", "jpg", "jpeg", "png")

    // Registra os tipos de mídia (MIME types) permitidos para uploads
    wisp.RegisterMIMETypes("image/jpeg", "image/png", "application/pdf", "text/xml")
}
```

### Exemplos de uso

**Exemplo 1: Validação de CPF**

```go
package main

import (
	"fmt"
	"log"

	"github.com/marcelofabianov/wisp"
)

func main() {
	// Cria um CPF a partir de uma string com máscara
	cpf, err := wisp.NewCPF("123.456.789-00") // Exemplo com um CPF inválido
	if err != nil {
		log.Fatal(err) // Ex: invalid CPF check digit 1
	}

	// O valor é armazenado sem a máscara
	fmt.Println("Valor canônico:", cpf.String())

	// Mas pode ser exibido com a formatação
	fmt.Println("Valor formatado:", cpf.Formatted())
}
```

**Exemplo 2: Cálculos Financeiros Seguros**

```go
package main

import (
	"fmt"
	"log"

	"github.com/marcelofabianov/wisp"
)

// Configuração inicial da aplicação
func init() {
	// Registra as unidades que serão usadas no sistema
	wisp.RegisterUnits("KG", "UN")
	// Define a precisão padrão para quantidades (ex: 3 casas decimais)
	wisp.SetDefaultPrecision(3)
}

func main() {
	// Preço por quilo: R$ 10,31
	precoPorKg, _ := wisp.NewMoney(1031, wisp.BRL)

	// Quantidade comprada: 1.57 KG
	quantidade, err := wisp.NewQuantity(1.57, "KG")
	if err != nil {
		log.Fatal(err)
	}

	// Calcula o total
	total, err := quantidade.MultiplyByMoney(precoPorKg)
	if err != nil {
		log.Fatal(err)
	}

	// O resultado é um objeto Money seguro com o arredondamento correto
	// 1.57 * 10.31 = 16.1867 -> arredonda para R$ 16,19
	fmt.Println("Preço Total:", total.String()) // Saída: BRL 16.19
}
```

**Exemplo 3: Modelagem de uma Entidade**

Este exemplo demonstra como os tipos `wisp` se unem para criar uma entidade `Course` segura e expressiva. O pacote [fault](https://github.com/marcelofabianov/fault) presente no exemplo vem para qualificar nossos erros com contexto.

```go
package main

import (
	"fmt"
	"log"

	"github.com/marcelofabianov/fault"
	"github.com/marcelofabianov/wisp"
)

// --- Definições do seu pacote de domínio ---

// Definimos os estados permitidos para o status do curso.
const (
	StatusPublished = "published"
	StatusDraft     = "draft"
)

// NewCourseInput é o DTO que carrega dados brutos e não validados.
type NewCourseInput struct {
	Name           string
	Description    string
	MaxEnrollments int
	InitialStatus  string // O status inicial vem como uma string simples
	CreatedBy      wisp.AuditUser
}

// Course é a Entidade de Domínio, protegida por tipos wisp.
type Course struct {
	ID             wisp.UUID
	Name           wisp.NonEmptyString
	Description    wisp.NonEmptyString
	MaxEnrollments wisp.PositiveInt
	Status         wisp.Flag[string] // O status agora é um wisp.Flag
	wisp.Audit
}

// NewCourse é a Factory que valida os dados brutos e cria uma entidade segura.
func NewCourse(input NewCourseInput) (*Course, error) {
	name, err := wisp.NewNonEmptyString(input.Name)
	if err != nil {
		return nil, fault.Wrap(err, "invalid name")
	}
	description, err := wisp.NewNonEmptyString(input.Description)
	if err != nil {
		return nil, fault.Wrap(err, "invalid description")
	}
	maxEnrollments, err := wisp.NewPositiveInt(input.MaxEnrollments)
	if err != nil {
		return nil, fault.Wrap(err, "invalid max enrollments")
	}

	// Valida e cria o Flag de status
	status, err := wisp.NewFlag(input.InitialStatus, StatusPublished, StatusDraft)
	if err != nil {
		return nil, fault.Wrap(err, "invalid initial status for course")
	}

	id, err := wisp.NewUUID()
	if err != nil {
		return nil, err
	}
	return &Course{
		ID:             id,
		Name:           name,
		Description:    description,
		MaxEnrollments: maxEnrollments,
		Status:         status,
		Audit:          wisp.NewAudit(input.CreatedBy),
	}, nil
}

// Publish é um método de comportamento que altera o estado do curso para publicado.
func (c *Course) Publish(updatedBy wisp.AuditUser) error {
	if c.Status.Is(StatusPublished) {
		return fault.New("course is already published", fault.WithCode(fault.Conflict))
	}
	newStatus, _ := wisp.NewFlag(StatusPublished, StatusPublished, StatusDraft)
	c.Status = newStatus
	c.Audit.Touch(updatedBy)
	return nil
}

// --- Uso prático na aplicação ---
func main() {
	creator, _ := wisp.NewAuditUser("admin@example.com")

	// 1. Criando um novo curso como rascunho
	fmt.Println("Criando um curso como rascunho...")
	input := NewCourseInput{
		Name:           "   Curso de Go   ",
		Description:    "Um curso focado em boas práticas.",
		MaxEnrollments: 50,
		InitialStatus:  StatusDraft, // Inicia como rascunho
		CreatedBy:      creator,
	}

	course, err := NewCourse(input)
	if err != nil {
		log.Fatalf("Falha inesperada: %v", err)
	}

	fmt.Printf("Curso criado com sucesso!\n")
	fmt.Printf("  ID: %s\n", course.ID)
	fmt.Printf("  Status Inicial: %s\n", course.Status.Get())
	// A verificação de estado é explícita e segura
	fmt.Printf("  O curso está público? %t\n\n", course.Status.Is(StatusPublished))

	// 2. Usando um método de comportamento para publicar o curso
	fmt.Println("Publicando o curso...")
	err = course.Publish(creator)
	if err != nil {
		log.Fatalf("Falha ao publicar: %v", err)
	}

	fmt.Printf("Status atualizado: %s\n", course.Status.Get())
	fmt.Printf("  O curso está público? %t\n", course.Status.Is(StatusPublished))
	fmt.Printf("  Versão atualizada: %d\n", course.Audit.Version.Int())
	fmt.Printf("  Atualizado por: %s\n", course.Audit.UpdatedBy)
}
```

## Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para abrir uma issue para discutir uma nova feature ou enviar um pull request.

## Licença

Este projeto é licenciado sob a Licença MIT.
