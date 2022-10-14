package web

type InitRequest struct {
	CustomerXid string ` validate:"required" json:"customer_xid" form:"customer_xid" `
}
