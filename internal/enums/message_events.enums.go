package enums

type MessageInboundEvents string

const (
	INBOUND_MSG  MessageInboundEvents = "MSG"
	INBOUND_JOIN MessageInboundEvents = "JOIN"
	INBOUND_LEFT MessageInboundEvents = "LEFT"
)

type MessageOutboundEvents string

const (
	OUTBOUND_MSG   MessageOutboundEvents = "MSG"
	OUTBOUND_READ  MessageOutboundEvents = "READ"
	OUTBOUND_CHATS MessageOutboundEvents = "CHATS"
	OUTBOUND_CONN  MessageOutboundEvents = "CONN"
)
