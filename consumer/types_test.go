package consumer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/adonese/noebs/ebs_fields"
	"github.com/adonese/noebs/utils"
	"github.com/alicebob/miniredis"
)

var mr, _ = miniredis.Run()
var mockRedis = utils.GetRedisClient(mr.Addr())

func Test_cardsFromZ(t *testing.T) {
	lcards := []ebs_fields.CardsRedis{
		{
			PAN:     "1234",
			Expdate: "2209",
			IsMain:  false,
			ID:      1,
		},
	}
	_, err := json.Marshal(lcards)
	if err != nil {
		t.Fatalf("there is an error in testing: %v\n", err)
	}

	fromRedis := []string{`{"pan": "1234", "exp_date": "2209", "id": 1}`}

	tests := []struct {
		name string
		args []string
		want []ebs_fields.CardsRedis
	}{
		{"Successful Test",
			fromRedis,
			lcards,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cardsFromZ(fromRedis); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cardsFromZ() = %v, want %v", got, tt.want)
			} else {
				fmt.Printf("they are: %v - %v", got[0], tt.want[0])
			}
		})
	}
}

func Test_generateCardsIds(t *testing.T) {
	have1 := ebs_fields.CardsRedis{PAN: "1334", Expdate: "2201", ID: 1}
	have2 := ebs_fields.CardsRedis{PAN: "1234", Expdate: "2202", ID: 2}
	have := &[]ebs_fields.CardsRedis{
		have1, have2,
	}
	want := []ebs_fields.CardsRedis{
		{PAN: "1334", Expdate: "2201", ID: 1},
		{PAN: "1234", Expdate: "2202", ID: 2},
	}
	tests := []struct {
		name string
		have *[]ebs_fields.CardsRedis
		want []ebs_fields.CardsRedis
	}{
		{"testing equality", have, want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generateCardsIds(tt.have)
			for i, c := range *tt.have {
				if !reflect.DeepEqual(c, tt.want[i]) {
					t.Errorf("have: %v, want: %v", c, tt.want[i])
				}
			}
		})
	}
}

func Test_newFromBytes(t *testing.T) {
	type args struct {
		d    []byte
		code int
	}
	tests := []struct {
		name    string
		args    args
		want    response
		wantErr bool
	}{
		{"testing response - 200", args{d: []byte(`{ "ebs_response": { "pubKeyValue": "MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBANx4gKYSMv3CrWWsxdPfxDxFvl+Is/0kc1dvMI1yNWDXI3AgdI4127KMUOv7gmwZ6SnRsHX/KAM0IPRe0+Sa0vMCAwEAAQ==", "UUID": "958c8835-9f89-4f96-96a8-7430039c6323", "responseMessage": "Approved", "responseStatus": "Successful", "responseCode": 0, "tranDateTime": "200222113700" } }`), code: 200}, response{Code: 0, Response: "Approved"}, false},
		{"testing response - 200", args{d: []byte(`{ "message": "EBSError", "code": 613, "status": "EBSError", "details": { "UUID": "6cccfb54-640c-495c-8e0c-434b280937a2", "responseMessage": "DUPLICATE_TRANSACTION", "responseStatus": "Failed", "responseCode": 613, "tranDateTime": "200222113700" } }`), code: 502}, response{Code: 613, Response: "DUPLICATE_TRANSACTION"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newFromBytes(tt.args.d, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("newFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
