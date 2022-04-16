package scraper

import (
	"reflect"
	"strings"
	"testing"

	"github.com/DuGlaser/atc/internal/core"
)

func TestGetProblemSamples(t *testing.T) {
	tests := []struct {
		input    string
		expected []*core.TestCase
	}{
		{input: `
<div id="task-statement">
  <span class="lang-ja">
    <div class="io-style">
      <div>
        <section>
          <h3>入力</h3><p>入力は以下の形式で標準入力から与えられる。</p>
          <pre></pre>
        </section>
      </div>

      <div class="part">
        <section>
          <h3>出力</h3><p>答えを出力せよ。<br>
          条件を満たす数が存在しない場合は <code>-1</code> を出力せよ。</p>
        </section>
      </div>
    </div>
    <div class="part">
      <section>
        <h3>入力例 1 <span>Copy</span></h3><div><span>Copy</span></div><pre id="pre-sample0">123 456 100
        </pre>
      </section>
    </div>
    <div class="part">
      <section>
        <h3>出力例 1 <span>Copy</span></h3><div><span>Copy</span></div><pre id="pre-sample1">200
        </pre>
        <p><code>300</code> や <code>400</code> も正解です。</p>
      </section>
    </div>
    <div class="part">
      <section>
        <h3>入力例 2 <span>Copy</span></h3><div><span>Copy</span></div><pre id="pre-sample2">630 940 314
        </pre>
      </section>
    </div>
    <div class="part">
      <section>
        <h3>出力例 2 <span>Copy</span></h3><div><span>Copy</span></div><pre id="pre-sample3">-1
      </pre></section>
    </div>
  </span>
</div>
`,
			expected: []*core.TestCase{
				{In: "123 456 100", Expected: "200"},
				{In: "630 940 314", Expected: "-1"},
			},
		},
	}

	for _, test := range tests {
		pp, err := NewTaskPage(strings.NewReader(test.input))
		if err != nil {
			t.Fatal(err)
		}

		tc, err := pp.GetTaskTestCases()
		if err != nil {
			t.Fatal(err)
		}

		if len(tc) != len(test.expected) {
			t.Errorf("tc length is wrong. got=%d, want=%d.", len(tc), len(test.expected))
		}

		for i := range tc {
			if !reflect.DeepEqual(tc[i], test.expected[i]) {
				t.Errorf("tc has wrong value. got=%v, want=%v.", tc[i], test.expected[i])
			}
		}
	}
}
