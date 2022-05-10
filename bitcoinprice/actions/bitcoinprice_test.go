package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gobuffalo/buffalo"
)

func (as *ActionSuite) Test_BitcoinPriceHandlerWithBuffalo() {

	res := as.JSON("/api/v1/btc").Get()

	as.Equal(http.StatusOK, res.Code)
	btcResponse := BitcoinPriceHandlerResponse{}
	if err := json.Unmarshal(res.Body.Bytes(), &btcResponse); err != nil {
		as.Errorf(err, "Failed to unmarshal response as BitcoinPriceHandlerResponse")
	}

	if currencyLength := len(btcResponse.Currency); currencyLength != 3 {
		msg := fmt.Sprintf("Unexpected currency %q with length: %d should be 3 characters long", btcResponse.Currency, currencyLength)
		as.T().Error(msg)
	}

	if btcResponse.Value <= 0 {
		msg := fmt.Sprintf("Currency value cannot be negative or equal to 0 but was %d", btcResponse.Value)
		as.T().Error(msg)
	}
}

// func Test_BitcoinPriceHandler(t *testing.T) {

// 	req, err := http.NewRequest("GET", "/api/v1/btc", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	req.Header.Add("Content-Type", "application/json")

// 	w := httptest.NewRecorder()

// 	res := w.Result()
// 	if res.StatusCode != http.StatusCreated {
// 		t.Fatalf("status code %d different than %d", res.StatusCode, http.StatusCreated)
// 	}
// }

func TestExample(t *testing.T) {

	type testCase struct {
		name           string
		num1           int
		num2           int
		operation      string
		expectedResult int
		expectError    bool
	}

	testCases := []testCase{
		{
			name:           "Add1and2",
			num1:           1,
			num2:           2,
			operation:      "+",
			expectedResult: 3,
			expectError:    false,
		},
		{
			name:           "Multiply1and2",
			num1:           1,
			num2:           2,
			operation:      "*",
			expectedResult: 2,
			expectError:    false,
		},
		{
			name:           "Substract1and2",
			num1:           1,
			num2:           2,
			operation:      "-",
			expectedResult: -1,
			expectError:    false,
		},
		{
			name:           "Divide1and0",
			num1:           1,
			num2:           0,
			operation:      "/",
			expectedResult: 0,
			expectError:    true,
		},
		{
			name:           "Multiply3and5",
			num1:           3,
			num2:           5,
			operation:      "*",
			expectedResult: 15,
			expectError:    false,
		},
	}

	for i := range testCases {
		currentTest := testCases[i]
		t.Run(currentTest.name, func(t *testing.T) {
			expectedResult := currentTest.expectedResult
			gotError := false
			actualResult := 0
			switch currentTest.operation {
			case "+":
				actualResult = currentTest.num1 + currentTest.num2
			case "-":
				actualResult = currentTest.num1 - currentTest.num2
			case "/":
				if currentTest.num2 == 0 {
					gotError = true
				}
			case "*":
				actualResult = currentTest.num1 * currentTest.num2
			default:
				t.Fatalf("Invalid test case unexpected operation %q", currentTest.operation)
			}

			if gotError != currentTest.expectError {
				t.Errorf("Expected error %t but got %t when doing %d %s %d",
					currentTest.expectError,
					gotError,
					currentTest.num1,
					currentTest.operation,
					currentTest.num2)
				t.Fail()
			} else if actualResult != expectedResult {
				t.Errorf("Expected %d but got %d when doing %d %s %d",
					expectedResult,
					actualResult,
					currentTest.num1,
					currentTest.operation,
					currentTest.num2)
			}
		})
	}

}

func TestBitcoinPriceRequestOutput(t *testing.T) {

	req, err := http.NewRequest("GET", "localhost:8080/api/v1/btc", nil)
	if err != nil {
		t.FailNow()
	}

	ctx := buffalo.DefaultContext{req: req}

	if err := BitcoinPriceHandler(ctx); err != nil {

	}
}
