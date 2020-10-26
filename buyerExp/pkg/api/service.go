package api

type SubscribeReq struct{
	NoticeURL 	string				`json:"url"`
	Mail	 	string				`json:"mail"`
}

type SubscribeResp struct {
	Id          int					`json:"user_id"`
}

type ConfirmReq struct{
	Hash 	string
}

type ConfirmResp struct{
	Message string					`json:"message"`
}

