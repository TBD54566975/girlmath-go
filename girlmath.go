package girlmath

import (
	"errors"
	"math/big"
	"strconv"
	"strings"
)

const (
	BTCSubunitsPerUnit  = 100000000
	BTCSubunitsSigDigs  = 8
	USDSubunitsPerUnit  = 100
	USDSubunitsSigDigs  = 2
	KESSubunitsPerUnit  = 100
	KESSubunitsSigDigs  = 2
	USDCSubunitsPerUnit = 100
	USDCSubunitsSigDigs = 6
	MXNSubunitsPerUnit  = 100
	MXNSubunitsSigDigs  = 2
)

type CurrencyConstants struct {
	SubunitsPerUnit int
	SigDigs         int
}

func GetConversionConstants(currencyCode string) (CurrencyConstants, error) {
	switch currencyCode {
	case "USDC":
		return CurrencyConstants{SubunitsPerUnit: USDCSubunitsPerUnit, SigDigs: USDCSubunitsSigDigs}, nil
	case "USD":
		return CurrencyConstants{SubunitsPerUnit: USDSubunitsPerUnit, SigDigs: USDSubunitsSigDigs}, nil
	case "BTC":
		return CurrencyConstants{SubunitsPerUnit: BTCSubunitsPerUnit, SigDigs: BTCSubunitsSigDigs}, nil
	case "KES":
		return CurrencyConstants{SubunitsPerUnit: KESSubunitsPerUnit, SigDigs: KESSubunitsSigDigs}, nil
	case "MXN":
		return CurrencyConstants{SubunitsPerUnit: MXNSubunitsPerUnit, SigDigs: MXNSubunitsSigDigs}, nil
	default:
		return CurrencyConstants{}, errors.New("unexpected currency")
	}
}

/**
 * @param payinCurrency starting currency code that customer wants to convert to @param payoutCurrency
 * @param payoutCurrency currency code that customer wants to end up with
 * @param payinAmountSubunits amount subunits of @param payinCurrency being used to purchase @param payoutCurrency
 * @param payinCurrencyUnitsPerPayoutCurrencyUnit price of 1 whole @param payoutCurrency unit, in terms of @param payinCurency units.
 * Say payin currency is USD, payout currency is BTC, the spot price would be written as 30,741.70 USD/BTC, so 30,741.70 would be this arg.
 * If payin currency is BTC, and payout currency is USD, the spot price would be written as 0.0000325291 BTC/USD, so 0.0000325291 would be this arg.
 * @returns Number of @param payoutCurrency subunits that can be bought with @param payinAmountSubunits amount of @param payinCurency.
 */
func CalculatePayoutAmountSubunitsWithPayinCurrencySpotPrice(payinCurrency, payoutCurrency string, payinAmountSubunits int64, payinCurrencyUnitsPerPayoutCurrencyUnit string) (*big.Int, error) {
	payinConstants, err := GetConversionConstants(payinCurrency)
	if err != nil {
		return nil, err
	}
	payoutConstants, err := GetConversionConstants(payoutCurrency)
	if err != nil {
		return nil, err
	}

	payinCurrencyUnitsPerPayoutCurrencyUnit = strings.ReplaceAll(payinCurrencyUnitsPerPayoutCurrencyUnit, ",", "")
	payinUnitsPerPayoutUnitStripped, err := strconv.ParseFloat(payinCurrencyUnitsPerPayoutCurrencyUnit, 64)
	if err != nil {
		return nil, err
	}

	payinSubunitsPerPayoutUnit := payinUnitsPerPayoutUnitStripped * float64(payinConstants.SubunitsPerUnit)
	payoutAmountUnits := float64(payinAmountSubunits) / payinSubunitsPerPayoutUnit
	payoutAmountSubunits := big.NewInt(int64(payoutAmountUnits * float64(payoutConstants.SubunitsPerUnit)))

	return payoutAmountSubunits, nil
}

/**
 * @param payinCurrency starting currency code that customer wants to convert to @param payoutCurrency
 * @param payoutCurrency currency code that customer wants to end up with
 * @param payinAmountSubunits amount subunits of @param payinCurrency being used to purchase @param payoutCurrency
 * @param payoutCurrencyUnitsPerPayinCurrencyUnit price of 1 whole @param payinCurency unit, in terms of @param payoutCurrency units.
 * Say payin currency is USD, payout currency is BTC, the spot price would be written as 0.0000325291 BTC/USD, so 0.0000325291 would be this arg.
 * If payin currency is BTC, and payout currency is USD, the spot price would be written as 30,741.70 USD/BTC, so 30,741.70 would be this arg.
 * @returns Number of @param payoutCurrency subunits that can be bought with @param payinAmountSubunits amount of @param payinCurency.
 */
func CalculatePayoutAmountSubunitsWithPayoutCurrencySpotPrice(payinCurrency, payoutCurrency string, payinAmountSubunits int64, payoutCurrencyUnitsPerPayinCurrencyUnit string) (*big.Int, error) {
	payinConstants, err := GetConversionConstants(payinCurrency)
	if err != nil {
		return nil, err
	}
	payoutConstants, err := GetConversionConstants(payoutCurrency)
	if err != nil {
		return nil, err
	}

	payoutCurrencyUnitsPerPayinCurrencyUnit = strings.ReplaceAll(payoutCurrencyUnitsPerPayinCurrencyUnit, ",", "")
	payoutUnitPerPayinUnitsStripped, err := strconv.ParseFloat(payoutCurrencyUnitsPerPayinCurrencyUnit, 64)
	if err != nil {
		return nil, err
	}

	payoutUnitsPerPayinUnit := float64(payinConstants.SubunitsPerUnit) / payoutUnitPerPayinUnitsStripped
	payoutAmountUnits := float64(payinAmountSubunits) / payoutUnitsPerPayinUnit
	payoutAmountSubunits := big.NewInt(int64(payoutAmountUnits * float64(payoutConstants.SubunitsPerUnit)))

	return payoutAmountSubunits, nil
}

/**
 * Converts a number @param amountSubunits into a unit amount string, with a decimal point for overflow subunits.
 * @param amountSubunits starting amount subunits of @param currencyCode
 * @param currencyCode referring to @param amountSubunits
 * @returns a whole unit amount of @param currencyCode with extra subunits after a decimal point
 */
func ConvertSubunitsToUnits(amountSubunits int64, currencyCode string) (string, error) {
	constants, err := GetConversionConstants(currencyCode)
	if err != nil {
		return "", err
	}

	amountUnits := amountSubunits / int64(constants.SubunitsPerUnit)
	remainingSubunits := amountSubunits % int64(constants.SubunitsPerUnit)
	subunitsString := strconv.FormatInt(remainingSubunits, 10)
	subunitsLength := len(subunitsString)

	if remainingSubunits == 0 {
		subunitsString = strings.Repeat("0", constants.SigDigs)
	} else if subunitsLength < constants.SigDigs {
		subunitsString = strings.Repeat("0", constants.SigDigs-subunitsLength) + subunitsString
	} else if subunitsLength > constants.SigDigs {
		subunitsString = subunitsString[:constants.SigDigs]
	}

	return strconv.FormatInt(amountUnits, 10) + "." + subunitsString, nil
}

/**
 * Converts a number @param amountUnits into a subunit amount string
 * @param amountUnits starting amount subunits of @param currencyCode
 * @param currencyCode referring to @param amountUnits
 * @returns subunits string of @param currencyCode that @param amountUnits contains
 */
func ConvertUnitsToSubunits(amountUnits string, currencyCode string) (string, error) {
	constants, err := GetConversionConstants(currencyCode)
	if err != nil {
		return "", err
	}

	amountUnits = strings.ReplaceAll(amountUnits, ",", "")
	parts := strings.Split(amountUnits, ".")
	majorSegment := parts[0]
	var minorSegment string
	if len(parts) > 1 {
		minorSegment = parts[1]
	}

	if len(minorSegment) < constants.SigDigs {
		minorSegment += strings.Repeat("0", constants.SigDigs-len(minorSegment))
	} else if len(minorSegment) > constants.SigDigs {
		minorSegment = minorSegment[:constants.SigDigs]
	}

	if majorSegment == "" {
		majorSegment = "0"
	}

	subunits, err := strconv.ParseInt(majorSegment+minorSegment, 10, 64)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(subunits, 10), nil
}
