package main

import "os"

const initConfigFile = `
{
  "db": {
    "user": "root",
    "pass": "1234",
    "host": "localhost",
    "port": 3306,
    "name": "demo",
    "query_values": {
      "charset": ["utf8mb4"],
      "parseTime": ["true"],
      "loc": ["UTC"]
    }
  },
  "serve_addr": ":8000",
  "use_case_timeout": "3s"
}
`

func main() {
	localConfigFile, err := os.Create("config.local.json")
	if err != nil {
		panic(err)
	}

	localConfigFile.WriteString(initConfigFile)
}
