package client

import (
	"net"
	"time"
)

type ClientSet struct {
	targetNodeAddress string 
	sleepInterval time.Duration
	conn net.Conn
}

type ServerMessages struct {
	Ack string `json:"ack"`
}

func NewClientSet(address string, sleepTime time.Duration) *ClientSet {
	return &ClientSet{
		targetNodeAddress: address,
		sleepInterval: sleepTime,
	}
}

func (c *ClientSet) Set(key, value []byte) error {
	var response ServerMessages

	fmt.Printf("yooo")

	if okKey := k.IsEmpty(key); okKey {
		return errors.New("The Key is Empty!")
	}

	if okValue := k.IsEmpty(value); okValue {
		return errors.New("The Value is Empty!")
	}

	jsonValue, _ := json.Marshal(map[string]any{
		"type": "set",
		"key": key,
		"value": value,
	})

	fmt.Println(string(jsonValue))

	var err error
	for {
		time.Sleep(c.sleepInterval)

		fmt.Printf("foo")

		c.conn, err = net.Dial("tcp", c.targetNodeAddress)
		if err != nil {
			return err
		}
		c.conn.Write([]byte(jsonValue))

		reply := make([]byte, 2024)
		n, _ := k.conn.Read(reply)
		json.Unmarshal(reply[:n], &response)

		c.conn.Close()
	}
	return nil
}