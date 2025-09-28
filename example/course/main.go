package main

import (
	"fmt"
	"log"
	"time"

	// Importa o pacote wisp para configurar e usar seus tipos
	// Importa o pacote local do nosso exemplo
	// Ajuste este import para o caminho do seu projeto

	"github.com/marcelofabianov/wisp"
	"github.com/marcelofabianov/wisp/example/course/domain"
)

// init é executado uma vez no início para configurar o pacote wisp.
func init() {
	// Para este exemplo, vamos registrar a role 'ADMIN'.
	wisp.RegisterRoles("ADMIN", "SYSTEM")
}

func main() {
	creator, _ := wisp.NewAuditUser("admin@example.com")
	updater, _ := wisp.NewAuditUser("system")

	// --- 1. Criação do Curso ---
	input := domain.NewCourseInput{
		Name:                "Go Básico",
		Description:         "Fundamentos da linguagem Go.",
		EnrollmentLimit:     50,
		EnrollmentStartDate: "2025-10-01",
		EnrollmentEndDate:   "2025-10-31",
		CreatedBy:           creator,
	}

	fmt.Println("--- CRIANDO NOVO CURSO ---")
	c, err := domain.NewCourse(input)
	if err != nil {
		log.Fatalf("Falha inesperada ao criar curso: %v", err)
	}
	printCourseState("Estado Inicial", c)

	// --- 2. Alteração do Nome ---
	time.Sleep(10 * time.Millisecond) // Pequeno delay para ver a mudança no UpdatedAt
	fmt.Println("\n--- ALTERANDO NOME ---")
	newName, _ := wisp.NewNonEmptyString("Go Básico: Primeiros Passos")
	c.ChangeName(newName, updater)
	printCourseState("Após Alterar Nome", c)

	// --- 3. Alteração do Limite de Matrículas ---
	time.Sleep(10 * time.Millisecond)
	fmt.Println("\n--- ALTERANDO LIMITE DE VAGAS ---")
	newLimit, _ := wisp.NewPositiveInt(75)
	c.UpdateEnrollmentLimit(newLimit, updater)
	printCourseState("Após Alterar Limite", c)

	// --- 4. Alteração do Período de Matrículas ---
	time.Sleep(10 * time.Millisecond)
	fmt.Println("\n--- ALTERANDO PERÍODO DE MATRÍCULAS ---")
	newStart, _ := wisp.NewDate(2025, time.November, 1)
	newEnd, _ := wisp.NewDate(2025, time.November, 30)
	newPeriod, _ := wisp.NewDateRange(newStart, newEnd)
	c.UpdateEnrollmentPeriod(newPeriod, updater)
	printCourseState("Estado Final", c)
}

// printCourseState é uma função auxiliar para exibir o estado atual do curso.
func printCourseState(stage string, c *domain.Course) {
	fmt.Printf("[%s]\n", stage)
	fmt.Printf("  ID: %s\n", c.ID)
	fmt.Printf("  Nome: %s\n", c.Name)
	fmt.Printf("  Limite de Vagas: %d\n", c.EnrollmentLimit.Int())
	fmt.Printf("  Período: de %s a %s\n", c.EnrollmentPeriod.Start(), c.EnrollmentPeriod.End())
	fmt.Printf("  Auditoria:\n")
	fmt.Printf("    Versão: %d\n", c.Audit.Version.Int())
	fmt.Printf("    Criado por: %s em %s\n", c.Audit.CreatedBy, c.Audit.CreatedAt.Time().Format(time.RFC3339))
	fmt.Printf("    Atualizado por: %s em %s\n", c.Audit.UpdatedBy, c.Audit.UpdatedAt.Time().Format(time.RFC3339))
}
