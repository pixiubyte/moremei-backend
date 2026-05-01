package cipher

type (
    Config struct {
        Key string `json:"Key"`
        Iv  string `json:"Iv"`
    }

    ManagerConfig map[string]*Config
)
