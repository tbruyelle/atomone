package cli

import (
	"fmt"
	"reflect" // #nosec
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/cosmos/cosmos-sdk/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/atomone-hub/atomone/x/gov/types"
)

const (
	draftProposalFileName = "draft_proposal.json"
)

// Prompt prompts the user for all values of the given type.
// data is the struct to be filled
// namePrefix is the name to be displayed as "Enter <namePrefix> <field>"
func Prompt[T any](data T, namePrefix string) (T, error) {
	v := reflect.ValueOf(&data).Elem()
	if v.Kind() == reflect.Interface {
		v = reflect.ValueOf(data)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}

	for i := 0; i < v.NumField(); i++ {
		// if the field is a struct skip or not slice of string or int then skip
		switch v.Field(i).Kind() {
		case reflect.Struct:
			// TODO(@julienrbrt) in the future we can add a recursive call to Prompt
			continue
		case reflect.Slice:
			if v.Field(i).Type().Elem().Kind() != reflect.String && v.Field(i).Type().Elem().Kind() != reflect.Int {
				continue
			}
		}

		// create prompts
		prompt := promptui.Prompt{
			Label:    fmt.Sprintf("Enter %s %s", namePrefix, strings.ToLower(client.CamelCaseToString(v.Type().Field(i).Name))),
			Validate: client.ValidatePromptNotEmpty,
		}

		fieldName := strings.ToLower(v.Type().Field(i).Name)

		if strings.EqualFold(fieldName, "authority") {
			// pre-fill with gov address
			prompt.Default = authtypes.NewModuleAddress(types.ModuleName).String()
			prompt.Validate = client.ValidatePromptAddress
		}

		// TODO(@julienrbrt) use scalar annotation instead of dumb string name matching
		if strings.Contains(fieldName, "addr") ||
			strings.Contains(fieldName, "sender") ||
			strings.Contains(fieldName, "voter") ||
			strings.Contains(fieldName, "depositor") ||
			strings.Contains(fieldName, "granter") ||
			strings.Contains(fieldName, "grantee") ||
			strings.Contains(fieldName, "recipient") {
			prompt.Validate = client.ValidatePromptAddress
		}

		result, err := prompt.Run()
		if err != nil {
			return data, fmt.Errorf("failed to prompt for %s: %w", fieldName, err)
		}

		switch v.Field(i).Kind() {
		case reflect.String:
			v.Field(i).SetString(result)
		case reflect.Int:
			resultInt, err := strconv.ParseInt(result, 10, 0)
			if err != nil {
				return data, fmt.Errorf("invalid value for int: %w", err)
			}
			// If a value was successfully parsed the ranges of:
			//      [minInt,     maxInt]
			// are within the ranges of:
			//      [minInt64, maxInt64]
			// of which on 64-bit machines, which are most common,
			// int==int64
			v.Field(i).SetInt(resultInt)
		case reflect.Slice:
			switch v.Field(i).Type().Elem().Kind() {
			case reflect.String:
				v.Field(i).Set(reflect.ValueOf([]string{result}))
			case reflect.Int:
				resultInt, err := strconv.ParseInt(result, 10, 0)
				if err != nil {
					return data, fmt.Errorf("invalid value for int: %w", err)
				}

				v.Field(i).Set(reflect.ValueOf([]int{int(resultInt)}))
			}
		default:
			// skip any other types
			continue
		}
	}

	return data, nil
}
