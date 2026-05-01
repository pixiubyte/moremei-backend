package cos

type Config struct {
    BucketURL string `json:",optional"`
    Bucket    string `json:","`
    Region    string `json:","`
    AppId     string `json:","`
    SecretID  string `json:","`
    SecretKey string `json:","`
}
