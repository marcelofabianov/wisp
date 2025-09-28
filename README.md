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

**Exemplo 3: Modelando a Entidade `Course`**

Este é um guia que irá lhe mostrar como usar os value objects do `wisp` pode construir uma entidade de domínio ` Course` que é segura, expressiva e robusta desde a sua criação.

 O pacote [fault](https://github.com/marcelofabianov/fault) presente no exemplo vem para qualificar nossos erros com contexto.

 O objetivo é simples: parar de usar tipos primitivos como `string` e `int` para representar conceitos de negócio complexos e, com isso, eliminar uma série de bugs comuns.

 **Passo 1: A Porta de Entrada - O DTO (`NewCourseInput`)**

 Primeiro, precisamos de uma forma de receber dados do "mundo exterior" (uma API, um formulário, etc.). Para isso, usamos uma `struct` simples, conhecida como DTO (Data Transfer Object). Note que ela usa tipos primitivos (`string`, `int`), pois neste ponto, os dados ainda são brutos e não confiáveis.

 ```go
 type NewCourseInput struct {
	Name           string
	Description    string
	MinEnrollments int
	MaxEnrollments int
	CreatedBy      wisp.AuditUser
}
```

**Passo 2: O Coração do Domínio - A Entidade (Course)**

Agora, a parte mais importante: a entidade `Course`. Esta `struct` representa um curso dentro do nosso sistema. Diferente do DTO, aqui nós usamos os tipos `wisp` para garantir que um `Course`, uma vez criado, esteja sempre em um estado válido.

```go
// Course é a Entidade de Domínio, protegida por tipos wisp.
type Course struct {
	ID          wisp.UUID
	Name        wisp.NonEmptyString
	Description wisp.NonEmptyString
	Enrollments wisp.RangedValue // Este tipo é a chave da nossa nova lógica
	wisp.Audit
}
```

Vamos analisar cada campo `wisp`:
- ID `wisp.UUID`: Garante que todo curso terá um identificador único e em um formato válido.
- Name `wisp.NonEmptyString`: Garante que o nome do curso nunca será uma string vazia ou contendo apenas espaços. A validação e a limpeza (trim) são feitas na criação.
- Enrollments `wisp.RangedValue`: Este é o nosso value object mais poderoso aqui. Ele encapsula a regra de negócio de que as matrículas têm um valor mínimo para o curso acontecer e um máximo. O RangedValue gerencia o número current (atual), min (mínimo) e max (máximo) de matrículas, garantindo que o valor atual nunca saia desses limites.
- `wisp.Audit`: Com uma única linha, embutimos todos os campos de auditoria (CreatedAt, CreatedBy, Version, etc.) na nossa entidade.

**Passo 3: O Guardião - A Factory (`NewCourse`)**

Como transformamos os dados brutos do `NewCourseInput` na nossa entidade `Course` segura? Através de uma "Factory Function". Esta função é o **guardião do nosso domínio**. É o único ponto de entrada permitido para criar um `Course`, e é aqui que toda a validação acontece.

```go
// NewCourse é a Factory que valida os dados brutos e cria uma entidade segura.
func NewCourse(input NewCourseInput) (*Course, error) {
	// Usamos os construtores do wisp para validar cada campo do DTO
	name, err := wisp.NewNonEmptyString(input.Name)
	if err != nil {
		return nil, fault.Wrap(err, "invalid name")
	}
	description, err := wisp.NewNonEmptyString(input.Description)
	if err != nil {
		return nil, fault.Wrap(err, "invalid description")
	}

	// Aqui, criamos o RangedValue. O valor inicial (current) de matrículas é 0.
	enrollments, err := wisp.NewRangedValue(0, int64(input.MinEnrollments), int64(input.MaxEnrollments))
	if err != nil {
		return nil, fault.Wrap(err, "invalid enrollments range")
	}

	id, err := wisp.NewUUID()
	if err != nil {
		return nil, err
	}

	// Se todas as validações passaram, montamos a entidade.
	return &Course{
		ID:          id,
		Name:        name,
		Description: description,
		Enrollments: enrollments,
		Audit:       wisp.NewAudit(input.CreatedBy),
	}, nil
}
```

Se a função `NewCourse` retornar sem erro, você tem a garantia matemática de que o objeto `Course` é válido.

**Passo 4: Dando Vida à Entidade (Comportamento)**

Entidades não são apenas dados, elas têm comportamento. Veja como os métodos do `Course` ficam simples e legíveis, pois a lógica complexa já está nos value objects.

```go
// EnrollStudent é um método de comportamento que adiciona uma matrícula.
func (c *Course) EnrollStudent(updatedBy wisp.AuditUser) error {
	// A validação de limite máximo é delegada para o RangedValue!
	newEnrollments, err := c.Enrollments.Add(1)
	if err != nil {
		return err // Retornará wisp.ErrValueExceedsMax se o curso estiver cheio
	}
	c.Enrollments = newEnrollments
	c.Audit.Touch(updatedBy)
	return nil
}

// CanStart verifica se o curso atingiu o número mínimo de matrículas.
func (c *Course) CanStart() bool {
	return c.Enrollments.Current() >= c.Enrollments.Min()
}
```

Note como `EnrollStudent` é simples. Ele apenas chama `c.Enrollments.Add(1)`. Toda a complexidade de "verificar se a soma excede o máximo" está escondida dentro do `RangedValue`, onde deve estar.

**Colocando Tudo Junto: Exemplo Prático**

Agora, vamos ver o código completo em ação.

```go
package main

import (
	"errors"
	"fmt"
	"log"
    // ... imports de Course e wisp
)

func main() {
	creator, _ := wisp.NewAuditUser("admin@example.com")

	// 1. Definindo os parâmetros para um novo curso
	input := NewCourseInput{
		Name:           "Curso de Testes em Go",
		Description:    "Aprendendo a testar com testify.",
		MinEnrollments: 5,  // Precisa de 5 alunos para começar
		MaxEnrollments: 7,  // Limite de 7 alunos
		CreatedBy:      creator,
	}

	fmt.Println("Criando um curso...")
	course, err := NewCourse(input)
	if err != nil {
		log.Fatalf("Falha inesperada: %v", err)
	}

	fmt.Printf("Curso criado!\n")
	fmt.Printf("  Matrículas: %d (Mín: %d, Máx: %d)\n",
		course.Enrollments.Current(), course.Enrollments.Min(), course.Enrollments.Max())
	fmt.Printf("  O curso pode começar? %t\n\n", course.CanStart())

	// 2. Simulando matrículas
	fmt.Println("Matriculando 5 alunos para atingir o mínimo...")
	for i := 0; i < 5; i++ {
		course.EnrollStudent(creator)
	}
	fmt.Printf("  Matrículas atuais: %d\n", course.Enrollments.Current())
	fmt.Printf("  O curso pode começar? %t\n\n", course.CanStart())

	// 3. Enchendo a turma
	fmt.Println("Matriculando mais 2 alunos para encher a turma...")
	course.EnrollStudent(creator)
	course.EnrollStudent(creator)
	fmt.Printf("  A turma está cheia? %t\n\n", course.Enrollments.IsAtMax())

	// 4. Tentando matricular um aluno extra (onde o wisp brilha!)
	fmt.Println("Tentando matricular o 8º aluno...")
	err = course.EnrollStudent(creator)
	if err != nil {
		// Verificamos o erro específico retornado pelo RangedValue
		if errors.Is(err, wisp.ErrValueExceedsMax) {
			fmt.Println("SUCESSO: Erro esperado recebido -> Turma cheia!")
		}
	}
}
```

## Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para abrir uma issue para discutir uma nova feature ou enviar um pull request.

## Licença

Este projeto é licenciado sob a Licença MIT.
