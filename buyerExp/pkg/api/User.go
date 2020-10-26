package api

type User struct{
	Mail	 	string				`json:"mail"`
	Hash		[]byte				`json:"hash,omitempty"`
}
