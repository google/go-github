package github

import (
	"testing"
	"time"
)

const fakePrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCBEJcaraGrNm68LstWKMbwalt+n0kaz1pYkBiz+vuFU2K2BkO8
qrk2VXsJQKtLRJYBsEb5CzbgEPRQFoHbqo7lGlJlT8iLVpd0G0g5A9AfABtMnSCP
B5onUg3qAy+G16EtsG+5xOddq3CGs3Au17l9xPCinrpjUts+yYDnUdgKhwIDAQAB
AoGALJT8fRypEak1yw8m8dYYEgfLHwwKhpZFkP4FanYx17YcDOBRGaSnYZtZarLJ
+K/yWRb5DSHQjMmOevOrW6Oow/twOxgt/Qrdmtf7zJCBDuKeIdb5n2Eo1iCeD3jD
LgdM376i0zfgGOeUQFSjGARmAdXPYuITorXPur+fwqQQosECQQD1dzzLtW63GQSa
kO4zm8szLQQKtt22PiTMnqIEizb9GV8Ut96yRaqGXOPPW+WVK402OLyGxg8PxGQp
Pi2+aplxAkEAhpqJMsE/Zedy5mTeyrZ4twDqIlRzkKOOEsPjyvR6zC5kEN3xtE32
oxt2dh9oAINUabQRl4HsnYdvLkQZ/lKndwJBAIYktnXA3hNzuZ9aisJrZn2+IRJV
4w4gIe4s1u/SwKdKgTqKTUpxZgJtqxV77Bg8J7y0/tpMTJjaQ76CO29vRBECQBZR
u5w55PxNw4SmhzbPyPZ3ZVtp63u5Uw6TgssdfNssehU96B1ArDvhiuQNUx56qF2a
qSwZNBXu5iWizziXlgMCQFodovtv4Wk5WI7BYADWXiHeTCwxuZ4Mp1Z3jD9E7Le6
vvShKgb038c8wUyjsYYEmSI5SsdLAKrSqXsFKFMTKWg=
-----END RSA PRIVATE KEY-----`

func Test_applicationTokenSource_Token(t *testing.T) {
	type fields struct {
		id         string
		privateKey []byte
		expiration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "should error when applicationID is empty",
			fields: fields{
				id: "",
			},
			wantErr: true,
		},
		{
			name: "should error when private key is invalid",
			fields: fields{
				id:         "github-app-id",
				privateKey: []byte("invalid"),
			},
			wantErr: true,
		},
		{
			name: "should return a token",
			fields: fields{
				id:         "github-app-id",
				privateKey: []byte(fakePrivateKey),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := NewApplicationTokenSource(tt.fields.id, tt.fields.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewApplicationTokenSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tr == nil {
				return
			}

			got, err := tr.Token()
			if (err != nil) != tt.wantErr {
				t.Errorf("applicationTokenSource.Token() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.AccessToken == "" {
				t.Errorf("applicationTokenSource.Token() = %v, want non-empty AccessToken", got)
			}
			if got.TokenType != "Bearer" {
				t.Errorf("applicationTokenSource.Token() = %v, want TokenType to be Bearer", got)
			}
			if got.Expiry.IsZero() {
				t.Errorf("applicationTokenSource.Token() = %v, want non-zero Expiry", got)
			}
		})
	}
}
