package domain

type Wallet struct {
	Id          int    `field:"id"`
	CustomerId  string `field:"customer_id"`
	IsActive    bool   `field:"is_active"`
	ActiveAt    any    `field:"active_at"`
	DisabledAt  any    `field:"disabled_at"`
	Balance     int    `field:"balance"`
	Token       string `field:"token"`
	TxId        int    //transaction id
	ReferenceId string
	Amount      int
	CreatedAt   any
	CreatedBy   string
}
