package domain

import (
	"github.com/marcelofabianov/fault"

	"github.com/marcelofabianov/wisp"
)

// NewCourseInput é o Data Transfer Object (DTO) para carregar dados externos
// para a criação de um novo curso. Ele usa tipos primitivos para a validação inicial.
type NewCourseInput struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	EnrollmentLimit     int    `json:"enrollment_limit"`
	EnrollmentStartDate string `json:"enrollment_start_date"`
	EnrollmentEndDate   string `json:"enrollment_end_date"`
	CreatedBy           wisp.AuditUser
}

// Course é a entidade de domínio que representa um curso.
// Seus campos são protegidos por value objects do pacote wisp,
// garantindo que uma instância de Course esteja sempre em um estado válido.
type Course struct {
	ID               wisp.UUID
	Name             wisp.NonEmptyString
	Description      wisp.NonEmptyString
	EnrollmentLimit  wisp.PositiveInt
	EnrollmentPeriod wisp.DateRange
	wisp.Audit
}

// NewCourse é a factory que cria uma entidade Course válida a partir de um DTO.
// É o único ponto de entrada para a criação de um novo curso, garantindo
// que todas as regras de negócio sejam aplicadas.
func NewCourse(input NewCourseInput) (*Course, error) {
	// Validação e criação dos Value Objects
	name, err := wisp.NewNonEmptyString(input.Name)
	if err != nil {
		return nil, fault.Wrap(err, "invalid course name", fault.WithCode(fault.Invalid))
	}

	description, err := wisp.NewNonEmptyString(input.Description)
	if err != nil {
		return nil, fault.Wrap(err, "invalid course description", fault.WithCode(fault.Invalid))
	}

	enrollmentLimit, err := wisp.NewPositiveInt(input.EnrollmentLimit)
	if err != nil {
		return nil, fault.Wrap(err, "invalid enrollment limit", fault.WithCode(fault.Invalid))
	}

	startDate, err := wisp.ParseDate(input.EnrollmentStartDate)
	if err != nil {
		return nil, fault.Wrap(err, "invalid enrollment start date", fault.WithCode(fault.Invalid))
	}

	endDate, err := wisp.ParseDate(input.EnrollmentEndDate)
	if err != nil {
		return nil, fault.Wrap(err, "invalid enrollment end date", fault.WithCode(fault.Invalid))
	}

	enrollmentPeriod, err := wisp.NewDateRange(startDate, endDate)
	if err != nil {
		return nil, fault.Wrap(err, "invalid enrollment period", fault.WithCode(fault.Invalid))
	}

	id, err := wisp.NewUUID()
	if err != nil {
		return nil, err
	}

	// Montagem da Entidade
	course := &Course{
		ID:               id,
		Name:             name,
		Description:      description,
		EnrollmentLimit:  enrollmentLimit,
		EnrollmentPeriod: enrollmentPeriod,
		Audit:            wisp.NewAudit(input.CreatedBy),
	}

	return course, nil
}

// ChangeName é um método de comportamento da entidade Course.
func (c *Course) ChangeName(newName wisp.NonEmptyString, updatedBy wisp.AuditUser) {
	c.Name = newName
	c.Audit.Touch(updatedBy)
}

// UpdateEnrollmentLimit é outro exemplo de método de comportamento.
func (c *Course) UpdateEnrollmentLimit(newLimit wisp.PositiveInt, updatedBy wisp.AuditUser) {
	c.EnrollmentLimit = newLimit
	c.Audit.Touch(updatedBy)
}

// UpdateEnrollmentPeriod atualiza o período de matrículas do curso.
func (c *Course) UpdateEnrollmentPeriod(newPeriod wisp.DateRange, updatedBy wisp.AuditUser) {
	c.EnrollmentPeriod = newPeriod
	c.Audit.Touch(updatedBy)
}
