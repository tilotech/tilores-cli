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
    "searchRuleSets": [
      {
        "name": "default",
        "ruleSet": {
          "ruleSetID": "common",
          "satisfiedRules": [
            "R1EXACT"
          ],
          "unsatisfiedRules": []
        }
      }
    ],
    "mutationRuleSetGroups": [
      {
        "name": "default",
        "linkRuleSet": {
          "ruleSetID": "common",
          "satisfiedRules": [
            "R1EXACT"
          ],
          "unsatisfiedRules": []
        },
        "deduplicateRuleSet": {}
      }
    ]
  }
}
`,
	Run: func(_ *cobra.Command, _ []string) {
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
	ruleConfigFile, err := os.Open("./rule-config.json")
	if err != nil {
		return fmt.Errorf("unable to open rule-config.json: %v", err)
	}
	var ruleConfig map[string]interface{}
	err = json.NewDecoder(ruleConfigFile).Decode(&ruleConfig)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	errCh := make(chan error)
	for _, caseFile := range caseFiles {
		wg.Add(1)
		go runCase(caseFile, ruleConfig, wg, errCh)
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
		SearchRuleSets        []responseSearchRuleSet        `json:"searchRuleSets"`
		MutationRuleSetGroups []responseMutationRuleSetGroup `json:"mutationRuleSetGroups"`
	} `json:"expectation"`
}

func runCase(caseFile string, ruleConfig map[string]interface{}, wg *sync.WaitGroup, errCh chan error) {
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

	errors := compareSearchRuleSets(ruleTestCase.Expectation.SearchRuleSets, actual.TiloRes.SimulateRules.SearchRuleSets)
	mutationRuleErrors := compareMutationRuleSetGroups(ruleTestCase.Expectation.MutationRuleSetGroups, actual.TiloRes.SimulateRules.MutationRuleSetGroups)
	errors = append(errors, mutationRuleErrors...)

	if len(errors) != 0 {
		errCh <- fmt.Errorf("case %v failed, errors:\n%v", caseFile, strings.Join(errors, "\n"))
		return
	}
}

func compareSearchRuleSets(expected, actual []responseSearchRuleSet) []string {
	expectedMap, errors := searchToMap(expected, "invalid expectation")
	actualMap, actualErrors := searchToMap(actual, "invalid actual")
	errors = append(errors, actualErrors...)
	errors = append(errors, compareRuleSets(expectedMap, actualMap)...)
	return errors
}

func compareMutationRuleSetGroups(expected, actual []responseMutationRuleSetGroup) []string {
	expectedMap, errors := mutationToMap(expected, "invalid expectation")
	actualMap, actualErrors := mutationToMap(actual, "invalid actual")
	errors = append(errors, actualErrors...)
	errors = append(errors, compareRuleSets(expectedMap, actualMap)...)
	return errors
}

func compareRuleSets(expectedMap, actualMap map[string]map[string]bool) []string {
	errors := make([]string, 0)
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

func searchToMap(searchRuleSets []responseSearchRuleSet, errorPrefix string) (map[string]map[string]bool, []string) {
	errors := make([]string, 0)
	ruleSetMap := map[string]map[string]bool{}
	for _, ruleSet := range searchRuleSets {
		identifier := fmt.Sprintf("%s-%s", ruleSet.Name, ruleSet.RuleSet.RuleSetID)
		_, ok := ruleSetMap[identifier]
		if ok {
			errors = append(errors, fmt.Sprintf("%v, rule set %v only allowed once", errorPrefix, identifier))
		}
		m, rulesErrors := mapRules(ruleSet.RuleSet, errorPrefix)
		errors = append(errors, rulesErrors...)
		ruleSetMap[identifier] = m
	}
	return ruleSetMap, errors
}

func mutationToMap(mutationGroups []responseMutationRuleSetGroup, errorPrefix string) (map[string]map[string]bool, []string) {
	errors := make([]string, 0)
	ruleSetMap := map[string]map[string]bool{}
	for _, group := range mutationGroups {
		identifier := fmt.Sprintf("%s-link-%s", group.Name, group.LinkRuleSet.RuleSetID)
		_, ok := ruleSetMap[identifier]
		if ok {
			errors = append(errors, fmt.Sprintf("%v, rule set %v only allowed once", errorPrefix, identifier))
		}
		m, rulesErrors := mapRules(group.LinkRuleSet, errorPrefix)
		errors = append(errors, rulesErrors...)
		ruleSetMap[identifier] = m

		identifier = fmt.Sprintf("%s-deduplicate-%s", group.Name, group.DeduplicateRuleSet.RuleSetID)
		_, ok = ruleSetMap[identifier]
		if ok {
			errors = append(errors, fmt.Sprintf("%v, rule set %v only allowed once", errorPrefix, identifier))
		}
		m, rulesErrors = mapRules(group.DeduplicateRuleSet, errorPrefix)
		errors = append(errors, rulesErrors...)
		ruleSetMap[identifier] = m
	}
	return ruleSetMap, errors
}

func mapRules(rs ruleSet, errorPrefix string) (map[string]bool, []string) {
	errors := make([]string, 0)
	m := map[string]bool{}
	for _, satisfiedRule := range rs.SatisfiedRules {
		_, ok := m[satisfiedRule]
		if ok {
			errors = append(errors, fmt.Sprintf("%v, rule %v only allowed once", errorPrefix, satisfiedRule))
		}
		m[satisfiedRule] = true
	}
	for _, unsatisfiedRule := range rs.UnsatisfiedRules {
		_, ok := m[unsatisfiedRule]
		if ok {
			errors = append(errors, fmt.Sprintf("%v, rule %v only allowed once", errorPrefix, unsatisfiedRule))
		}
		m[unsatisfiedRule] = false
	}
	return m, errors
}
