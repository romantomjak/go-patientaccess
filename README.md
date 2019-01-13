# go-patientaccess

Barebones client for Patient Acccess API for monitoring my local clinic's GP availability

---

```go
package main

import (
    "fmt"
    "github.com/romantomjak/go-patientaccess"
)

func main() {
    client := patientaccess.NewClient()

    token, err := client.GetToken("username", "password")
    if err != nil {
        fmt.Println("Failed to obtain API token:", err)
        return
    }

    patientId, err := client.GetPatientId(token.Token)
    if err != nil {
        fmt.Println("Failed to obtain patient ID", err)
        return
    }

    slots, err := client.GetAppointmentSlots(token.Token, patientId)
    if err != nil {
        fmt.Println("Failed to list appointment types", err)
        return
    }

    fmt.Println("Available appointment types:")
    for _, slot := range slots {
        fmt.Printf(" - %s\n", slot.SlotType.Name)
    }
}
```

## License

MIT
