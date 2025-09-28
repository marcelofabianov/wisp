package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
	"github.com/marcelofabianov/wisp/example/course/domain"
)

type CourseSuite struct {
	suite.Suite
	creator wisp.AuditUser
	updater wisp.AuditUser
}

func (s *CourseSuite) SetupSuite() {
	s.creator, _ = wisp.NewAuditUser("creator@example.com")
	s.updater, _ = wisp.NewAuditUser("updater@example.com")
}

func TestCourseSuite(t *testing.T) {
	suite.Run(t, new(CourseSuite))
}

func (s *CourseSuite) TestNewCourse() {
	s.Run("should create a new course successfully with valid input", func() {
		input := domain.NewCourseInput{
			Name:                "Go for Production",
			Description:         "A course about writing production-ready Go.",
			EnrollmentLimit:     100,
			EnrollmentStartDate: "2025-10-01",
			EnrollmentEndDate:   "2025-10-31",
			CreatedBy:           s.creator,
		}

		course, err := domain.NewCourse(input)

		s.Require().NoError(err)
		s.Require().NotNil(course)

		s.False(course.ID.IsNil())
		s.Equal("Go for Production", course.Name.String())
		s.Equal(100, course.EnrollmentLimit.Int())
		s.Equal("2025-10-01", course.EnrollmentPeriod.Start().String())
		s.Equal("2025-10-31", course.EnrollmentPeriod.End().String())
		s.Equal(wisp.Version(1), course.Audit.Version)
		s.Equal(s.creator, course.Audit.CreatedBy)
	})

	s.Run("should fail with invalid name", func() {
		input := domain.NewCourseInput{Name: " "}
		_, err := domain.NewCourse(input)
		s.Require().Error(err)
		s.Contains(err.Error(), "invalid course name")
	})

	s.Run("should fail with invalid enrollment limit", func() {
		input := domain.NewCourseInput{
			Name:            "Valid Name",
			Description:     "Valid Description",
			EnrollmentLimit: 0,
		}
		_, err := domain.NewCourse(input)
		s.Require().Error(err)
		s.Contains(err.Error(), "invalid enrollment limit")
	})

	s.Run("should fail with invalid date range", func() {
		input := domain.NewCourseInput{
			Name:                "Valid Name",
			Description:         "Valid Description",
			EnrollmentLimit:     50,
			EnrollmentStartDate: "2025-11-01",
			EnrollmentEndDate:   "2025-10-31",
		}
		_, err := domain.NewCourse(input)
		s.Require().Error(err)
		s.Contains(err.Error(), "invalid enrollment period")
	})
}

func (s *CourseSuite) TestCourse_ChangeName() {
	input := domain.NewCourseInput{
		Name:                "Original Name",
		Description:         "Desc",
		EnrollmentLimit:     10,
		EnrollmentStartDate: "2025-10-01",
		EnrollmentEndDate:   "2025-10-31",
		CreatedBy:           s.creator,
	}
	course, _ := domain.NewCourse(input)
	originalAudit := course.Audit

	time.Sleep(10 * time.Millisecond)

	newName, _ := wisp.NewNonEmptyString("New Name")
	course.ChangeName(newName, s.updater)

	s.Equal("New Name", course.Name.String())
	s.Equal(wisp.Version(2), course.Audit.Version)
	s.Equal(s.updater, course.Audit.UpdatedBy)
	s.True(course.Audit.UpdatedAt.Time().After(originalAudit.UpdatedAt.Time()))
}

func (s *CourseSuite) TestCourse_UpdateEnrollmentLimit() {
	input := domain.NewCourseInput{
		Name:                "Course Name",
		Description:         "Desc",
		EnrollmentLimit:     10,
		EnrollmentStartDate: "2025-10-01",
		EnrollmentEndDate:   "2025-10-31",
		CreatedBy:           s.creator,
	}
	course, _ := domain.NewCourse(input)
	originalAudit := course.Audit

	time.Sleep(10 * time.Millisecond)

	newLimit, _ := wisp.NewPositiveInt(20)
	course.UpdateEnrollmentLimit(newLimit, s.updater)

	s.Equal(20, course.EnrollmentLimit.Int())
	s.Equal(wisp.Version(2), course.Audit.Version)
	s.Equal(s.updater, course.Audit.UpdatedBy)
	s.True(course.Audit.UpdatedAt.Time().After(originalAudit.UpdatedAt.Time()))
}

func (s *CourseSuite) TestCourse_UpdateEnrollmentPeriod() {
	input := domain.NewCourseInput{
		Name:                "Course Name",
		Description:         "Desc",
		EnrollmentLimit:     10,
		EnrollmentStartDate: "2025-10-01",
		EnrollmentEndDate:   "2025-10-31",
		CreatedBy:           s.creator,
	}
	course, _ := domain.NewCourse(input)
	originalAudit := course.Audit

	time.Sleep(10 * time.Millisecond)

	newStart, _ := wisp.NewDate(2025, time.November, 1)
	newEnd, _ := wisp.NewDate(2025, time.November, 30)
	newPeriod, _ := wisp.NewDateRange(newStart, newEnd)

	course.UpdateEnrollmentPeriod(newPeriod, s.updater)

	s.True(course.EnrollmentPeriod.Equals(newPeriod))
	s.Equal(wisp.Version(2), course.Audit.Version)
	s.Equal(s.updater, course.Audit.UpdatedBy)
	s.True(course.Audit.UpdatedAt.Time().After(originalAudit.UpdatedAt.Time()))
}
