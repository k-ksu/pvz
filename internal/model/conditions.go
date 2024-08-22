package model

var (
	ConditionAccepted Condition = "accepted"
	ConditionGiven    Condition = "given"
	ConditionRefund   Condition = "refund"
)

type Condition string
