package system

import "context"

type AuditService struct{}

func (s *AuditService) VerifyAuditChain(ctx context.Context) (ChainVerifyResult, error) {
	return VerifyAuditChain(ctx)
}
