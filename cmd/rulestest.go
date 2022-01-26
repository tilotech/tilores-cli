package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

// rulesTestCmd represents rules test command
var rulesTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Tests the rules in rule-config.json and checks if results match the expectations.",
	Long: `Tests the rules in rule-config.json and checks if results match the expectations.
Reads the rule configuration from "./rule-config.json" file.
Test cases are read from "./test/rules/*.json" where each file represents a test case and
includes a pair of records and an expectation.
Each expectation is a list of rule sets with their expected satisfied and unsatisfied rules.

Usage example:
tilores-cli rules test
Where ./test/rules/case1.json contains the following:
{
  "recordA": {
    "myCustomField": "same value"
  },
  "recordB": {
    "myCustomField": "same value"
  },
  "expectation": {
    "ruleSets": [
      {
        "ruleSetID": "index",
        "satisfiedRules": [
          "R1EXACT"
        ],
        "unsatisfiedRules": [
          "R2OTHER"
        ]
      }
    ]
  }
}
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := testRules()
		cobra.CheckErr(err)
	},
}

func init() {
	rulesCmd.AddCommand(rulesTestCmd)
}

func testRules() error {
	caseFiles, err := filepath.Glob("./test/rules/*.json")
	if err != nil {
		return err
	}
	if len(caseFiles) == 0 {
		fmt.Println("no test case files found in \"./test/rules/*.json\"")
		return nil
	}
	ruleConfig, err := os.ReadFile("./rule-config.json")
	if err != nil {
		return err
	}
	wg := &sync.WaitGroup{}
	errCh := make(chan error)
	for _, caseFile := range caseFiles {
		wg.Add(1)
		go runCase(caseFile, string(ruleConfig), wg, errCh)
	}
	go func() {
		wg.Wait()
		close(errCh)
	}()
	for err = range errCh {
		fmt.Println(err)
	}
	if err != nil {
		return fmt.Errorf("%vsome test cases did not pass%v", colorRed, colorReset)
	}
	fmt.Printf("%vall tests passed%v\n", colorGreen, colorReset)
	return nil
}

type ruleTestCase struct {
	*RulesSimulateInput
	Expectation struct {
		RuleSets []ruleSet `json:"ruleSets"`
	} `json:"expectation"`
}

func runCase(caseFile string, ruleConfig string, wg *sync.WaitGroup, errCh chan error) {
	defer wg.Done()
	file, err := os.ReadFile(caseFile) //nolint:gosec
	if err != nil {
		errCh <- err
		return
	}

	ruleTestCase := &ruleTestCase{}
	err = json.Unmarshal(file, ruleTestCase)
	if err != nil {
		errCh <- err
		return
	}
	ruleTestCase.RuleConfig = ruleConfig
	actual, err := callTiloTechAPI(ruleTestCase.RulesSimulateInput)
	if err != nil {
		errCh <- err
		return
	}

	errors := compareRuleSets(ruleTestCase.Expectation.RuleSets, actual.TiloRes.SimulateRules.RuleSets)

	if len(errors) != 0 {
		errCh <- fmt.Errorf("case %v failed, errors:\n%v", caseFile, strings.Join(errors, "\n"))
		return
	}
}

func compareRuleSets(expected, actual []ruleSet) []string {
	expectedMap, errors := toMap(expected, "invalid expectation")
	actualMap, actualErrors := toMap(actual, "invalid actual")
	errors = append(errors, actualErrors...)

	for ruleSetID, expectedRules := range expectedMap {
		actualRules, ok := actualMap[ruleSetID]
		if !ok {
			errors = append(errors, fmt.Sprintf("rule set %v expected but not found", ruleSetID))
			continue
		}
		for rule, expectedIsSatisfied := range expectedRules {
			actualIsSatisfied, ok := actualRules[rule]
			if !ok {
				errors = append(errors, fmt.Sprintf("%v: rule %v expected but not found", ruleSetID, rule))
				continue
			}
			if actualIsSatisfied != expectedIsSatisfied {
				errors = append(errors, fmt.Sprintf("%v: rule %v expected to be %v but was %v", ruleSetID, rule, isSatisfiedString(expectedIsSatisfied), isSatisfiedString(actualIsSatisfied)))
				continue
			}
		}
		for rule := range actualRules {
			if _, ok := expectedRules[rule]; !ok {
				errors = append(errors, fmt.Sprintf("%v: rule %v found but not expected", ruleSetID, rule))
			}
		}
	}
	for ruleSetID := range actualMap {
		if _, ok := expectedMap[ruleSetID]; !ok {
			errors = append(errors, fmt.Sprintf("rule set %v found but not expected", ruleSetID))
		}
	}

	return errors
}

func isSatisfiedString(isSatisfied bool) string {
	if isSatisfied {
		return fmt.Sprintf("%vsatisfied%v", colorGreen, colorReset)
	}
	return fmt.Sprintf("%vunsatisfied%v", colorRed, colorReset)
}

func toMap(ruleSets []ruleSet, errorPrefix string) (map[string]map[string]bool, []string) {
	errors := make([]string, 0)
	ruleSetMap := map[string]map[string]bool{}
	for _, ruleSet := range ruleSets {
		_, ok := ruleSetMap[ruleSet.RuleSetID]
		if ok {
			errors = append(errors, fmt.Sprintf("%v, rule set %v only allowed once", errorPrefix, ruleSet.RuleSetID))
		}
		m := map[string]bool{}
		for _, satisfiedRule := range ruleSet.SatisfiedRules {
			_, ok := m[satisfiedRule]
			if ok {
				errors = append(errors, fmt.Sprintf("%v, rule %v only allowed once", errorPrefix, satisfiedRule))
			}
			m[satisfiedRule] = true
		}
		for _, unsatisfiedRule := range ruleSet.UnsatisfiedRules {
			_, ok := m[unsatisfiedRule]
			if ok {
				errors = append(errors, fmt.Sprintf("%v, rule %v only allowed once", errorPrefix, unsatisfiedRule))
			}
			m[unsatisfiedRule] = false
		}
		ruleSetMap[ruleSet.RuleSetID] = m
	}
	return ruleSetMap, errors
}
