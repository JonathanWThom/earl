package payments

import (
	"fmt"
	"github.com/JonathanWThom/coinbase-commerce-go"
	"os"
)

func CreateCharge(accountToken interface{}) {
	client := coinbase.Client(os.Getenv("COINBASE_API_KEY"))
	metadata := make(map[string]interface{})
	metadata["AccountToken"] = accountToken
	charge, err := client.Charge.Create(coinbase.APIChargeData{
		Name:         "1 month unlimited",
		Pricing_type: "fixed_price",
		Local_price:  coinbase.Money{Amount: 2.00, Currency: "USD"},
		Metadata:     metadata,
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v\n", charge)
}
