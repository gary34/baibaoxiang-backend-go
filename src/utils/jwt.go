package utils

import (
	jwt "github.com/dgrijalva/jwt-go"
)

func GenJwt(claims jwt.Claims, jwtSignKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtSignKey))
	return ss, err
}

func ParseJwt(tokenString, jwtSignKey string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSignKey), nil
	})
}

// tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"

// type MyCustomClaims struct {
//     Foo string `json:"foo"`
//     jwt.StandardClaims
// }

// // sample token is expired.  override time so it parses as valid
// at(time.Unix(0, 0), func() {
//     token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//         return []byte("AllYourBase"), nil
//     })

//     if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
//         fmt.Printf("%v %v", claims.Foo, claims.StandardClaims.ExpiresAt)
//     } else {
//         fmt.Println(err)
//     }
// })
