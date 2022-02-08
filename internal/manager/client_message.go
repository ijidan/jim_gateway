package manager

type ClientMessage struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

type TextMessage struct {
	Id           uint64 `json:"id"`
	Content      string `json:"content"`
	ToReceiverId string `json:"to_receiver_id"`
	AtUserId  string `json:"at_user_id"`
}

//type LocationMessage struct {
//	Id uint64 `json:"id"`
//	cover_image string `json:"cover_image"`
//	Float lat
//	double lng = 4;
//	string map_link = 5;
//	string desc = 6;
//}
//
//
////表情消息
//message FaceMessage{
//uint64 id = 1;
//string symbol = 2;
//}
//*/

//
//type C2CMessage struct {
//	Type    string `yaml:"type"`
//	Content string `yaml:"content"`
//}
//
//type C2GMessage struct {
//	Type    string `yaml:"type"`
//	Content string `yaml:"content"`
//}
