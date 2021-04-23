package models

import (
	"testing"
)

func TestValidate(t *testing.T) {

	type subsOut struct {
		ha1    string
		ha1b   string
		passed bool
	}

	type predefinedSubs struct {
		subsIn Subscribers
		want   subsOut
	}

	checkValidate := func(t *testing.T, p predefinedSubs) {
		t.Helper()
		t.Logf("name is %s", p.subsIn.Username)
		err := p.subsIn.Validate()
		if err != nil && p.want.passed {
			t.Errorf("validation failed: %v", err)
		}
		if p.subsIn.ha1 != p.want.ha1 {
			t.Errorf("ha1 got %s want %s", p.subsIn.ha1, p.want.ha1)
		}
		if p.subsIn.ha1b != p.want.ha1b {
			t.Errorf("ha1b got %s want %s", p.subsIn.ha1b, p.want.ha1b)
		}
	}

	t.Run("PredefinedSubs", func(t *testing.T) {
		// t.Parallel() we are fast enough..

		subsTest := []predefinedSubs{
			{
				Subscribers{1, "userA", "arsperger.com", "superpass", "", ""},
				subsOut{"9772dbdffc612464a6c6e109af68c4ac", "", true},
			},
			{
				Subscribers{2, "userA@arsperger.com", "arsperger.com", "superpass", "", ""},
				subsOut{"", "590f240561a6b8f415f0cf080c29ba3e", true},
			},
			{
				Subscribers{3, "arsen@arsper@ger.com", "arsperger.com", "superpass", "", ""},
				subsOut{"", "", false},
			},
			{
				Subscribers{4, "aaa#aa", "arsperger.com", "superpass", "", ""},
				subsOut{"", "", false},
			},
			{
				Subscribers{5, "", "arsperger.com", "superpass", "", ""},
				subsOut{"", "", false},
			},
		}

		for i := range subsTest {
			checkValidate(t, subsTest[i])
		}
	})

}
