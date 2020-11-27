package bitcoin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScriptFromAddressUncompressed(t *testing.T) {
	testKey := "0499f8239bfe10eb0f5e53d543635a423c96529dd85fa4bad42049a0b435ebdd"

	// Change this to false to see the error case
	address, err := GetAddressFromPrivateKeyString(testKey, true)
	if err != nil {
		t.Error()
	}
	lockingScript, err := ScriptFromAddress(address)
	if err != nil {
		t.Error()
	}

	var utxo = &Utxo{
		TxID:         "b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
		Vout:         0,
		ScriptPubKey: lockingScript,
		Satoshis:     1000}

	wifString, err := PrivateKeyToWifString(testKey)
	if err != nil {
		t.Error()
	}

	payTo := &PayToAddress{
		Address:  address, // back to self
		Satoshis: 500,
	}

	rawTx, err := CreateTxUsingWif(
		[]*Utxo{utxo},
		[]*PayToAddress{payTo},
		[]OpReturnData{{[]byte("data")}},
		wifString,
	)

	if err != nil {
		t.Errorf("error occurred: %s", err.Error())
	}

	assert.NotPanics(t, func() { rawTx.ToString() }, nil)
}

// TestScriptFromAddress will test the method ScriptFromAddress()
func TestScriptFromAddress(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		inputAddress   string
		expectedScript string
		expectedError  bool
	}{
		{"", "", true},
		{"0", "", true},
		{"1234567", "", true},
		{"1HRVqUGDzpZSMVuNSZxJVaB9xjneEShfA7", "76a914b424110292f4ea2ac92beb9e83cf5e6f0fa2996388ac", false},
		{"13Rj7G3pn2GgG8KE6SFXLc7dCJdLNnNK7M", "76a9141a9d62736746f85ca872dc555ff51b1fed2471e288ac", false},
	}

	// Run tests
	for _, test := range tests {
		if script, err := ScriptFromAddress(test.inputAddress); err != nil && !test.expectedError {
			t.Fatalf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.inputAddress, err.Error())
		} else if err == nil && test.expectedError {
			t.Fatalf("%s Failed: [%v] inputted and error was expected", t.Name(), test.inputAddress)
		} else if script != test.expectedScript {
			t.Fatalf("%s Failed: [%v] inputted [%s] expected but failed comparison of scripts, got: %s", t.Name(), test.inputAddress, test.expectedScript, script)
		}
	}
}

// ExampleScriptFromAddress example using ScriptFromAddress()
func ExampleScriptFromAddress() {
	script, err := ScriptFromAddress("1HRVqUGDzpZSMVuNSZxJVaB9xjneEShfA7")
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("script generated: %s", script)
	// Output:script generated: 76a914b424110292f4ea2ac92beb9e83cf5e6f0fa2996388ac
}

// BenchmarkScriptFromAddress benchmarks the method ScriptFromAddress()
func BenchmarkScriptFromAddress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ScriptFromAddress("1HRVqUGDzpZSMVuNSZxJVaB9xjneEShfA7")
	}
}
