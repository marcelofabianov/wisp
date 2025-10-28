package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/marcelofabianov/wisp"
	"github.com/stretchr/testify/suite"
)

type TypeSuite struct {
	suite.Suite
}

func TestTypeSuite(t *testing.T) {
	suite.Run(t, new(TypeSuite))
}

func (s *TypeSuite) SetupTest() {
	wisp.ClearRegisteredTypes()
}

func (s *TypeSuite) TestRegisterAndValidateType() {
	const (
		TypeInvoice  wisp.Type = "INVOICE"
		TypeReceipt  wisp.Type = "RECEIPT"
		TypeContract wisp.Type = "CONTRACT"
	)
	wisp.RegisterTypes(TypeInvoice, TypeReceipt, TypeContract)

	s.Run("NewType should create a valid type that is registered", func() {
		docType, err := wisp.NewType("invoice")
		s.Require().NoError(err)
		s.Equal(TypeInvoice, docType)
	})

	s.Run("NewType should return EmptyType for empty input", func() {
		docType, err := wisp.NewType("  ")
		s.Require().NoError(err)
		s.Equal(wisp.EmptyType, docType)
		s.True(docType.IsZero())
	})

	s.Run("NewType should fail for a type that is not registered", func() {
		_, err := wisp.NewType("REPORT")
		s.Require().Error(err)
	})

	s.Run("IsValid should work correctly", func() {
		s.True(wisp.Type("INVOICE").IsValid())
		s.False(wisp.Type("MEMO").IsValid())
	})
}

func (s *TypeSuite) TestType_JSON_SQL() {
	const TypePayment wisp.Type = "PAYMENT"
	wisp.RegisterTypes(TypePayment)
	paymentType, _ := wisp.NewType("payment")

	s.Run("JSON Marshaling/Unmarshaling", func() {
		data, err := json.Marshal(paymentType)
		s.Require().NoError(err)
		s.Equal(`"PAYMENT"`, string(data))

		var unmarshaledType wisp.Type
		err = json.Unmarshal(data, &unmarshaledType)
		s.Require().NoError(err)
		s.Equal(paymentType, unmarshaledType)

		invalidJSON := `"REFUND"`
		err = json.Unmarshal([]byte(invalidJSON), &unmarshaledType)
		s.Require().Error(err)
	})

	s.Run("SQL Value/Scan", func() {
		val, err := paymentType.Value()
		s.Require().NoError(err)
		s.Equal("PAYMENT", val.(string))

		var scannedType wisp.Type
		err = scannedType.Scan("PAYMENT")
		s.Require().NoError(err)
		s.Equal(paymentType, scannedType)

		err = scannedType.Scan("VOID")
		s.Require().Error(err)
	})
}
