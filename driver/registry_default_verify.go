package driver

import (
	"context"

	"kratos/identity"
	"kratos/selfservice/flow/verification"
	"kratos/selfservice/strategy/link"
)

func (m *RegistryDefault) VerificationFlowPersister() verification.FlowPersister {
	return m.persister
}

func (m *RegistryDefault) VerificationFlowErrorHandler() *verification.ErrorHandler {
	if m.selfserviceVerifyErrorHandler == nil {
		m.selfserviceVerifyErrorHandler = verification.NewErrorHandler(m)
	}

	return m.selfserviceVerifyErrorHandler
}

func (m *RegistryDefault) VerificationManager() *identity.Manager {
	if m.selfserviceVerifyManager == nil {
		m.selfserviceVerifyManager = identity.NewManager(m)
	}

	return m.selfserviceVerifyManager
}

func (m *RegistryDefault) VerificationHandler() *verification.Handler {
	if m.selfserviceVerifyHandler == nil {
		m.selfserviceVerifyHandler = verification.NewHandler(m)
	}

	return m.selfserviceVerifyHandler
}

func (m *RegistryDefault) LinkSender() *link.Sender {
	if m.selfserviceLinkSender == nil {
		m.selfserviceLinkSender = link.NewSender(m)
	}

	return m.selfserviceLinkSender
}

func (m *RegistryDefault) VerificationStrategies(ctx context.Context) (verificationStrategies verification.Strategies) {
	for _, strategy := range m.selfServiceStrategies() {
		if s, ok := strategy.(verification.Strategy); ok {
			if m.Config(ctx).SelfServiceStrategy(s.VerificationStrategyID()).Enabled {
				verificationStrategies = append(verificationStrategies, s)
			}
		}
	}
	return
}

func (m *RegistryDefault) AllVerificationStrategies() (recoveryStrategies verification.Strategies) {
	for _, strategy := range m.selfServiceStrategies() {
		if s, ok := strategy.(verification.Strategy); ok {
			recoveryStrategies = append(recoveryStrategies, s)
		}
	}
	return
}
