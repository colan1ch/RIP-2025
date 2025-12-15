package apitypes

type QueryResult struct {
    QueryID      int                      `json:"query_id"`
    ExecutionTime float64                 `json:"execution_time"`
    IndexResults []IndexQueryResult `json:"index_results"`
    Timestamp    float64                  `json:"timestamp"`
}

type IndexQueryResult struct {
    IndexID      int    `json:"index_id"`
    ReceivedRows int    `json:"received_rows"`
    Success      bool   `json:"success"`
    Error        string `json:"error,omitempty"`
}