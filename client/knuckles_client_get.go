package client


type ClientGet struct {
	targetNodeAddress string
}

type ServerMessages struct {
	Ack string `json:"ack"`
}

func NewClientGet(target string) *ClientGet {
	return &ClientGet{
		targetNodeAddress: target
	}
}

func (c *ClientGet) GetData(key, value []byte) {
	var (
		err error
		serverResponse ServerMessages 
	)

	if okKey := c.IsEmpty(key); okKey {
		return "", errors.New("The Key is Empty!")
	}

	jsonGetValue, marshalError := json.Marshal(map[string]any{
		"type": "get",
		"key": key,
	})

	if marshalError != nil {
		return "", marshalError
	}

	k.conn, err = net.Dial("tcp", c.targetNodeAddress)
	if err != nil {
		return "", err
	}

	k.conn.Write(jsonGetValue)

	reply := make([]byte, 2024)
	n, _ := k.conn.Read(reply)
	json.Unmarshal(reply[:n], &serverResponse)

	k.conn.Close()

	return serverResponse.Ack, nil
}

func (c *ClientGet) IsEmpty(bytesToCheck []byte) bool {
	var resultToReturn bool = true
	for _, b := range bytesToCheck {
		if b != 0 {
			resultToReturn = false
		} else {
			resultToReturn = true
		}
	}
	return resultToReturn
}