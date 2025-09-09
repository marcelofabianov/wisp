# wisp

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Um conjunto de *value objects* robustos e imutáveis para Go, projetado para modelagem de domínios complexos com segurança de tipos, com foco em cenários de negócios brasileiros.

## Visão Geral

Este pacote nasceu da necessidade de combater a "Obsessão Primitiva" (*Primitive Obsession*), onde conceitos de negócio ricos como CPF, Dinheiro ou um Período de Datas são representados por tipos primitivos (`string`, `int`, `float64`). Essa prática comum leva a código com baixa expressividade e é uma fonte constante de bugs, pois a validação fica espalhada e duplicada por todo o sistema.

`wisp` oferece tipos de valor que encapsulam essas regras de negócio, garantindo que um objeto, uma vez criado, esteja **sempre em um estado válido**. Isso torna o código mais seguro, mais limpo e muito mais fácil de manter.

### Principais Vantagens

* **Segurança de Tipos**: Impede a criação de dados inválidos em tempo de compilação ou na inicialização. Um `wisp.CPF` nunca conterá um valor inválido.
* **Imutabilidade**: Os objetos `wisp` são imutáveis. Operações que modificam o valor (ex: `Add` em `Money`) retornam uma nova instância, tornando o código seguro para uso em ambientes concorrentes e livre de efeitos colaterais.
* **Validação Embutida**: Regras complexas, como o cálculo de dígitos verificadores de CPF/CNPJ e a validação de formatos, estão centralizadas e testadas.
* **API Expressiva**: Métodos como `Formatted()`, `IsMobile()`, `ApplyTo(money)` e `Age()` tornam o código mais legível e explícito sobre sua intenção.
* **Extensibilidade**: Tipos como `Quantity` e `Unit` permitem que o consumidor do pacote defina suas próprias unidades de medida, adaptando a biblioteca ao seu domínio específico.

## Tipos Disponíveis

| Tipo | Descrição |
| :--- | :--- |
| **Identificadores** | |
| `UUID` | Wrapper para `uuid.UUID` (padrão v7) para identificadores únicos. |
| `CPF` | CPF brasileiro com validação de dígitos verificadores e formatação. |
| `CNPJ` | CNPJ brasileiro com validação de dígitos verificadores e formatação. |
| **Financeiro** | |
| `Currency` | Código de moeda (BRL, USD, EUR) validado a partir de uma lista. |
| `Money` | Representa um valor monetário com segurança, evitando `float64`. |
| `Percentage` | Tipo de porcentagem preciso, baseado em `int64`, para cálculos financeiros. |
| **Quantidades** | |
| `Unit` | Sistema de registro extensível para unidades de medida (`KG`, `UN`, etc.). |
| `Quantity`| Valor numérico com unidade de medida, com precisão configurável. |
| **Contato & Endereçamento**| |
| `Email`| Endereço de e-mail validado via parser da biblioteca padrão. |
| `Phone`| Telefone brasileiro (fixo ou móvel) com validação e formatação. |
| `CEP`| CEP brasileiro com validação de formato e formatação. |
| `UF` | Unidade Federativa brasileira validada a partir de uma lista. |
| **Temporal** | |
| `Date`| Representa uma data de calendário (YYYY-MM-DD) sem fuso horário. |
| `DateRange` | Um período entre duas datas, com validação de `start <= end`. |
| `BirthDate`| Uma data de nascimento que não pode ser no futuro, com cálculos de idade. |
| `Day` | Um dia do mês (1-31) para eventos recorrentes. |
| **Auditoria & Técnicos**| |
| `Version` | Versão numérica para travamento otimista. |
| `AuditUser`| Identificador de usuário de auditoria (e-mail ou "system"). |
| `CreatedAt` | Timestamp de criação. |
| `UpdatedAt` | Timestamp de modificação com método `Touch()`. |
| `NullableTime`| Um `time.Time` que pode ser nulo, para campos como `deleted_at`. |

## Instalação

```sh
go get [github.com/marcelofabianov/wisp](https://github.com/marcelofabianov/wisp)
```

### Configuração do Pacote

Alguns tipos em wisp são extensíveis. É recomendado configurar os valores padrão na inicialização da sua aplicação (em uma função init() ou no início da main()).

```go
func init() {
    // Define a maioridade padrão para 21 anos
    wisp.SetLegalAge(21)

    // Define a precisão padrão para quantidades
    wisp.SetDefaultPrecision(4)

    // Registra as unidades de medida que seu domínio utilizará
    wisp.RegisterUnits("KG", "UN", "L", "M2", "H")
}
```

### Exemplos de uso

**Exemplo 1: Validação de CPF**

```go
package main

import (
	"fmt"
	"log"

	"[github.com/marcelofabianov/wisp](https://github.com/marcelofabianov/wisp)"
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

	"[github.com/marcelofabianov/wisp](https://github.com/marcelofabianov/wisp)"
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

## Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para abrir uma issue para discutir uma nova feature ou enviar um pull request.

## Licença

Este projeto é licenciado sob a Licença MIT.
