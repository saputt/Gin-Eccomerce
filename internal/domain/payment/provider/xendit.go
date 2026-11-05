package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"orders/internal/payment"

	"github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/invoice"
	"gorm.io/datatypes"
)

type XenditPaymentService struct {
	client *xendit.APIClient
	repo   payment.PaymentRepository
}

func NewXenditPaymentService(secret string, repo payment.PaymentRepository) *XenditPaymentService {
	client := xendit.NewClient(secret)
	return &XenditPaymentService{
		client: client,
		repo:   repo,
	}
}

func (s *XenditPaymentService) CreateInvoice(ctx context.Context, dto payment.CreateInvoiceInput) (*payment.CreateInvoiceResponse, error) {
	invoiceReq := *invoice.NewCreateInvoiceRequest(dto.OrderID, dto.Amount)

	// invoiceReq.SetSuccessRedirectUrl("")
	// invoiceReq.SetFailureRedirectUrl("")

	resp, r, err := s.client.InvoiceApi.CreateInvoice(ctx).CreateInvoiceRequest(invoiceReq).Execute()

	if err != nil {
		return nil, fmt.Errorf("xendit create invoice error: %v, response: %v", err, r)
	}
	var rawBytes []byte
	if resp != nil {
		var errMarshal error
		rawBytes, errMarshal = json.Marshal(resp)
		if errMarshal != nil {
			return nil, fmt.Errorf("failed to marshal xendit raw response: %w", errMarshal)
		}
	}

	pymnt := payment.Payments{
		OrderId:              resp.ExternalId,
		XenditInvoiceId:      *resp.Id,
		XenditPaymentMethod:  resp.PaymentMethod.String(),
		XenditPaymentChannel: resp.MerchantName,
		Amount:               float32(resp.Amount),
		PaymentUrl:           resp.InvoiceUrl,
		RawResponse:          datatypes.JSON(rawBytes),
	}

	repoErr := s.repo.CreatePayment(ctx, &pymnt)
	if repoErr != nil {
		return nil, repoErr
	}

	return &payment.CreateInvoiceResponse{
		InvoiceID:  *resp.Id,
		InvoiceURL: resp.InvoiceUrl,
		Status:     string(resp.Status),
	}, nil
}
