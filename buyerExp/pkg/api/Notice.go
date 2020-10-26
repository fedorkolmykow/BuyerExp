package api

type Notice struct{
	URL 		string				`json:"url" db:"url"`
	Price	 	float64				`json:"price" db:"price"`
}

func (n *Notice) Validate() (err error){
	return
}
