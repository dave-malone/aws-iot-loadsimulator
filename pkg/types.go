package awsiotloadsimulator

type SimulationRequest struct {
	StartClientNumber int `json:"start_client_num"`
	ClientCount       int `json:"client_count"`
	MessagesPerClient int `json:"messages_per_client"`
}
