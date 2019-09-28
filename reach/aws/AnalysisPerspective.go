package aws

import "github.com/luhring/reach/reach"

type AnalysisPerspective struct {
	self                  reach.NetworkPoint
	other                 reach.NetworkPoint
	getSecurityGroupRules func(sg SecurityGroup) []SecurityGroupRule
}

func NewAnalysisPerspectiveSourceOriented(v reach.NetworkVector) AnalysisPerspective {
	return AnalysisPerspective{
		self:  v.Source,
		other: v.Destination,
		getSecurityGroupRules: func(sg SecurityGroup) []SecurityGroupRule {
			return sg.OutboundRules
		},
	}
}

func NewAnalysisPerspectiveDestinationOriented(v reach.NetworkVector) AnalysisPerspective {
	return AnalysisPerspective{
		self:  v.Destination,
		other: v.Source,
		getSecurityGroupRules: func(sg SecurityGroup) []SecurityGroupRule {
			return sg.InboundRules
		},
	}
}
