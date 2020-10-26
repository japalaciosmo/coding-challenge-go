package product

type product struct {
	ProductID  int            `json:"-"`
	UUID       string         `json:"uuid"`
	Name       string         `json:"name"`
	Brand      string         `json:"brand"`
	Stock      int            `json:"stock"`
	SellerUUID string         `json:"seller_uuid,omitempty"`
	Seller     *productSeller `json:"seller,omitempty""`
}

type productSeller struct {
	UUID string `json:"uuid"`
	Link link   `json:"_links"`
}

type link struct {
	Self    ref `json:"self"`
	Related * ref `json:"related,omitempty"`
}

type ref struct {
	Href string `json:"href"`
}
