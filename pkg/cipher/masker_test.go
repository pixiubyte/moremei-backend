package cipher

import "testing"

func TestMasker_Mask(t *testing.T) {
    type args struct {
        plaintext string
        category  Category
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        {
            name: "test mask name",
            args: args{
                plaintext: "张三",
                category:  CategoryName,
            },
            want: "*三",
        },
        {
            name: "test mask name",
            args: args{
                plaintext: "王小六",
                category:  CategoryName,
            },
            want: "**六",
        },
        {
            name: "test mask name",
            args: args{
                plaintext: "上官婉儿",
                category:  CategoryName,
            },
            want: "***儿",
        },
        {
            name: "test mask id num ",
            args: args{
                plaintext: "310000000000000000",
                category:  CategoryIdNum,
            },
            want: "3****************0",
        },
        {
            name: "test mask id num ",
            args: args{
                plaintext: "3100000000000000",
                category:  CategoryIdNum,
            },
            want: "3**************0",
        },
        {
            name: "test mask id num ",
            args: args{
                plaintext: "310000000000000",
                category:  CategoryIdNum,
            },
            want: "***************",
        },
        {
            name: "test mask mobile",
            args: args{
                plaintext: "13012341234",
                category:  CategoryMobile,
            },
            want: "130******34",
        },
        {
            name: "test mask email",
            args: args{
                plaintext: "12345678@qq.com",
                category:  CategoryEmail,
            },
            want: "123*****@qq.com",
        },
        {
            name: "test mask text",
            args: args{
                plaintext: "abed",
                category:  CategoryText,
            },
            want: "****",
        },
        {
            name: "test mask text",
            args: args{
                plaintext: "hello world",
                category:  CategoryText,
            },
            want: "h***d",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Masker.Mask(tt.args.plaintext, tt.args.category); got != tt.want {
                t.Errorf("Mask() = %v, want %v", got, tt.want)
            }
        })
    }
}