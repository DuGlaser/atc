package scraper

import (
	"strings"
	"testing"
)

func TestGetProblemIds(t *testing.T) {
	test := `
<div id="contest-statement">
  <span>
    <span>
      <h3>配点</h3>
      <section>
        <div>
          <table>
            <thead>
              <tr>
                <th>問題</th>
                <th>点数</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>A</td>
                <td>100</td>
              </tr>
              <tr>
                <td>B</td>
                <td>200</td>
              </tr>
              <tr>
                <td>C</td>
                <td>300</td>
              </tr>
              <tr>
                <td>D</td>
                <td>400</td>
              </tr>
              <tr>
                <td>E</td>
                <td>500</td>
              </tr>
              <tr>
                <td>F</td>
                <td>500</td>
              </tr>
              <tr>
                <td>G</td>
                <td>600</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <h3>賞金</h3>
      <section>
        <div>
          <table>
            <thead>
              <tr>
                <th>順位</th>
                <th>金額</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>1位</td>
                <td></td>
              </tr>
              <tr>
                <td>2位</td>
                <td></td>
              </tr>
              <tr>
                <td>3位</td>
                <td></td>
              </tr>
              <tr>
                <td>4位</td>
                <td></td>
              </tr>
              <tr>
                <td>5位</td>
                <td></td>
              </tr>
              <tr>
                <td>6位</td>
                <td></td>
              </tr>
              <tr>
                <td>7位</td>
                <td></td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </span>
  </span>
</div>
  `

	expect := []string{"a", "b", "c", "d", "e", "f", "g"}

	cp, err := NewContestPage(strings.NewReader(test))
	if err != nil {
		t.Fatal("Failed to create ContestPage.")
	}

	ids := cp.GetProblemIds()

	if len(ids) != len(expect) {
		t.Fatalf("ids has wrong value. got=%v.", ids)
	}

	for i, _ := range ids {
		id := ids[i].DisplayID

		if id != expect[i] {
			t.Errorf("ids has wrong value. got=%v, want=%v.", ids, expect)
		}
	}
}
