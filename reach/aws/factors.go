package aws

import (
	"errors"

	"github.com/luhring/reach/reach"
)

func getInstanceStateFactor(factors []reach.Factor) (*reach.Factor, error) {
	for _, factor := range factors {
		if factor.Kind == FactorKindInstanceState {
			return &factor, nil
		}
	}

	return nil, errors.New("no instance state factor found")
}

func getSecurityGroupRulesFactor(factors []reach.Factor) (*reach.Factor, error) {
	for _, factor := range factors {
		if factor.Kind == FactorKindSecurityGroupRules {
			return &factor, nil
		}
	}

	return nil, errors.New("no security group rules factor found")
}

func getNetworkACLRulesFactor(factors []reach.Factor) (*reach.Factor, error) {
	for _, factor := range factors {
		if factor.Kind == FactorKindNetworkACLRules {
			return &factor, nil
		}
	}

	return nil, errors.New("no network ACL rules factor found")
}
