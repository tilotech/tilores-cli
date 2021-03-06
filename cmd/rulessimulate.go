package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

const tiloTechAPI = "https://api.tilotech.io"

const colorReset = "\033[0m"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"
const colorYellow = "\033[33m"

var (
	asJson bool
)

// rulesSimulateCmd represents rules simulate command
var rulesSimulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulates the rules in rule-config.json and tries to match the provided records.",
	Long: `Simulates the rules in rule-config.json and tries to match the provided records.
Reads the rule configuration from "./rule-config.json" file.
Reads both records to match as json from standard input.

Usage example:
cat records.json | tilores-cli rules simulate
Where records.json contains the following:
{
  "recordA": {
    "myCustomField": "same value"
  },
  "recordB": {
    "myCustomField": "same value"
  }
}
`,
	Run: func(cmd *cobra.Command, args []string) {
		simulateRulesOutput, err := simulateRules()
		cobra.CheckErr(err)
		if asJson {
			err := json.NewEncoder(os.Stdout).Encode(simulateRulesOutput.TiloRes.SimulateRules)
			cobra.CheckErr(err)
		} else {
			printNicely(simulateRulesOutput.TiloRes.SimulateRules.RuleSets)
		}
	},
}

func printNicely(ruleSets []ruleSet) {
	for _, ruleSet := range ruleSets {
		fmt.Printf("Rule Set: %v\n", ruleSet.RuleSetID)
		for _, satisfiedRule := range ruleSet.SatisfiedRules {
			fmt.Printf("%v: %ssatisfied%s\n", satisfiedRule, colorGreen, colorReset)
		}
		for _, unsatisfiedRule := range ruleSet.UnsatisfiedRules {
			fmt.Printf("%v: %sunsatisfied%s\n", unsatisfiedRule, colorRed, colorReset)
		}
		fmt.Println()
	}
}

func init() {
	rulesCmd.AddCommand(rulesSimulateCmd)

	rulesSimulateCmd.Flags().BoolVarP(&asJson, "json", "j", false, "Shows output as JSON")
}

type RulesSimulateInput struct {
	RecordA    map[string]interface{} `json:"recordA"`
	RecordB    map[string]interface{} `json:"recordB"`
	RuleConfig map[string]interface{} `json:"ruleConfig"`
}

type ruleSet struct {
	RuleSetID        string   `json:"ruleSetID"`
	SatisfiedRules   []string `json:"satisfiedRules"`
	UnsatisfiedRules []string `json:"unsatisfiedRules"`
}

type rulesSimulateOutput struct {
	TiloRes struct {
		SimulateRules struct {
			RuleSets []ruleSet `json:"ruleSets"`
		} `json:"simulateRules"`
	} `json:"tiloRes"`
}

type gqlResult struct {
	Errors              []interface{}       `json:"errors"`
	SimulateRulesOutput rulesSimulateOutput `json:"data"`
}

func simulateRules() (*rulesSimulateOutput, error) {
	simulateRulesInput := &RulesSimulateInput{}
	err := json.NewDecoder(os.Stdin).Decode(simulateRulesInput)
	if err != nil {
		return nil, fmt.Errorf("unable to decode input records from standard input: %v", err)
	}

	ruleConfigFile, err := os.Open("./rule-config.json")
	if err != nil {
		return nil, fmt.Errorf("unable to open rule-config.json: %v", err)
	}
	err = json.NewDecoder(ruleConfigFile).Decode(&simulateRulesInput.RuleConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to decode rule-config.json: %v", err)
	}

	return callTiloTechAPI(simulateRulesInput)
}

func callTiloTechAPI(simulateRulesInput *RulesSimulateInput) (*rulesSimulateOutput, error) {
	inputA, err := json.Marshal(simulateRulesInput.RecordA)
	if err != nil {
		return nil, err
	}
	inputB, err := json.Marshal(simulateRulesInput.RecordB)
	if err != nil {
		return nil, err
	}
	ruleConfig, err := json.Marshal(simulateRulesInput.RuleConfig)
	if err != nil {
		return nil, err
	}

	body := struct {
		Query     string      `json:"query"`
		Variables interface{} `json:"variables"`
	}{
		Query: `query simulate($recordA: AWSJSON!, $recordB: AWSJSON!, $ruleConfig: AWSJSON!) {
	tiloRes {
		simulateRules(simulateRulesInput: {
				inputA: $recordA
				inputB: $recordB
				ruleConfig: $ruleConfig
		}) {
			ruleSets {
				ruleSetID
				satisfiedRules
				unsatisfiedRules
			}
		}
	}
}
`,
		Variables: map[string]string{
			"recordA":    string(inputA),
			"recordB":    string(inputB),
			"ruleConfig": string(ruleConfig),
		},
	}

	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal RulesSimulateInput %s, error was %v\n", requestBody, err)
	}

	gqlRes := gqlResult{}
	res, err := http.Post(tiloTechAPI, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("invalid status code %s\n", res.Status)
	}
	err = unmarshalResponse(res, &gqlRes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response %v, error was %v\n", res, err)
	}
	if len(gqlRes.Errors) != 0 {
		return nil, fmt.Errorf("GraphQL errors occured for request %s, errors were %v\n", requestBody, gqlRes.Errors)
	}

	return &gqlRes.SimulateRulesOutput, nil
}

func unmarshalResponse(res *http.Response, v interface{}) error {
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resBody, v)
	if err != nil {
		return err
	}

	return nil
}
