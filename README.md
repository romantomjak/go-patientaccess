# go-patientaccess

Barebones client for Patient Acccess API for monitoring my local clinics GP availability

---

```go
package main

import (
    "fmt"
    "github.com/romantomjak/go-patientaccess"
)

func main() {
    client := patientacess.NewClient()

    token, err := client.GetToken("username", "password")
    if err != nil {
        fmt.Errorf("Failed to obtain API token: %v", err)
    }

    appointments, err = client.GetAppointments(token)
    if err != nil {
        fmt.Errorf("Failed to list appointment types: %v", err)
    }

    fmt.Printf("Available appointment types: %+v", appointments)
}
```

## License

MIT
