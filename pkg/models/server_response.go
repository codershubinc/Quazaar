package models

type ServerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type UniServerResponse struct {
	Success     string `json:"success"`
	Distributor string `json:"distributor"`
	Type        string `json:"type"`
	MsgOf       string `json:"msg_of"`
	Action      string `json:"action"`
	Data        any    `json:"data,omitempty"`
}

////////////////////////////////////////////////////////////////////////
///	 distributor will like (ssr , broadcast , event atc)
///  type is , check internal (system , player)
///   others are just sub commands , refer those in specific dir in internal
/// *i know this docs are not good bt keep it for now*
////////////////////////////////////////////////////////////////////////
// ~~ distributor
// ssr server events
// broadcast ~ events which go to the all the clients (connected)
// responce ~ will go for client which request the action
