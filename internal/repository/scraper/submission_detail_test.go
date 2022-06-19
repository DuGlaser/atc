package scraper

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetSubmissionStatusMap(t *testing.T) {
	input := `
<div id="main-container" class="container" style="padding-top: 50px">
  <div class="row">
    <div id="contest-nav-tabs" class="col-sm-12 mb-2 cnvtb-fixed">
      <div>
        <small class="contest-duration"></small>
        <small class="back-to-home pull-right"><a href="/home"></a></small>
      </div>
      <ul class="nav nav-tabs"></ul>
    </div>

    <div class="col-sm-12">
      <p>
        <span class="h2"></span>
      </p>
      <hr />
      <p>
        <span class="h4">ソースコード</span>
        <a
          class="btn-text toggle-btn-text source-code-expand-btn"
          data-on-text="拡げる"
          data-off-text="折りたたむ"
          data-target="submission-code"
          >拡げる</a
        >
      </p>
      <div class="div-btn-copy">
        <span
          class="btn-copy btn-pre"
          tabindex="0"
          data-toggle="tooltip"
          data-trigger="manual"
          title=""
          data-target="submission-code-for-copy"
          data-original-title="Copied!"
        ></span>
      </div>
      <div class="div-btn-copy"></div>
      <pre
        id="submission-code"
        class="prettyprint linenums source-code prettyprinted"
        style=""
      ></pre>
      <pre id="for_copy0" class="source-code-for-copy"></pre>
      <h4>提出情報</h4>
      <div class="panel panel-default">
        <table class="table table-bordered table-striped">
          <tbody>
            <tr>
              <th class="col-sm-4">提出日時</th>
              <td class="text-center">
                <time class="fixtime-second">2022-06-18 22:38:51</time>
              </td>
            </tr>
            <tr>
              <th>問題</th>
              <td class="text-center">
                <a href="/contests/abc256/tasks/abc256_e"
                  >E - Takahashi's Anguish</a
                >
              </td>
            </tr>
            <tr>
              <th>ユーザ</th>
              <td class="text-center"></td>
            </tr>
            <tr>
              <th>言語</th>
              <td class="text-center">C++ (GCC 9.2.1)</td>
            </tr>
            <tr>
              <th>得点</th>
              <td class="text-center">0</td>
            </tr>
            <tr>
              <th>コード長</th>
              <td class="text-center">1132 Byte</td>
            </tr>
            <tr>
              <th>結果</th>
              <td id="judge-status" class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="不正解"
                  >WA</span
                >
              </td>
            </tr>

            <tr>
              <th>実行時間</th>
              <td class="text-center">2206 ms</td>
            </tr>
            <tr>
              <th>メモリ</th>
              <td class="text-center">9632 KB</td>
            </tr>
          </tbody>
        </table>
      </div>

      <h4>ジャッジ結果</h4>
      <div class="panel panel-default">
        <div class="table-responsive">
          <table class="table table-bordered table-striped th-center">
            <tbody>
              <tr>
                <th width="10%">セット名</th>

                <th class="text-center">Sample</th>

                <th class="text-center">All</th>
              </tr>
              <tr>
                <th class="no-break">得点 / 配点</th>

                <td class="text-center">0 / 0</td>

                <td class="text-center">0 / 500</td>
              </tr>
              <tr>
                <th class="no-break">結果</th>

                <td class="text-center">
                  <table style="margin: auto">
                    <tbody>
                      <tr>
                        <td>
                          <span
                            class="label label-success"
                            data-toggle="tooltip"
                            data-placement="top"
                            title=""
                            data-original-title="正解"
                            >AC</span
                          >
                        </td>
                        <td style="padding-left: 5px">× 2</td>
                      </tr>
                    </tbody>
                  </table>
                </td>

                <td class="text-center">
                  <table style="margin: auto">
                    <tbody>
                      <tr>
                        <td>
                          <span
                            class="label label-success"
                            data-toggle="tooltip"
                            data-placement="top"
                            title=""
                            data-original-title="正解"
                            >AC</span
                          >
                        </td>
                        <td style="padding-left: 5px">× 2</td>
                      </tr>

                      <tr>
                        <td>
                          <span
                            class="label label-warning"
                            data-toggle="tooltip"
                            data-placement="top"
                            title=""
                            data-original-title="不正解"
                            >WA</span
                          >
                        </td>
                        <td style="padding-left: 5px">× 3</td>
                      </tr>

                      <tr>
                        <td>
                          <span
                            class="label label-warning"
                            data-toggle="tooltip"
                            data-placement="top"
                            title=""
                            data-original-title="実行時間制限超過"
                            >TLE</span
                          >
                        </td>
                        <td style="padding-left: 5px">× 12</td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="panel panel-default">
        <table class="table table-bordered table-striped th-center">
          <thead>
            <tr>
              <th class="no-break">セット名</th>
              <th>テストケース</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td class="text-center">Sample</td>
              <td>00_sample_00.txt, 00_sample_01.txt</td>
            </tr>

            <tr>
              <td class="text-center">All</td>
              <td>
                00_sample_00.txt, 00_sample_01.txt, 01_small_00.txt,
                01_small_01.txt, 01_small_02.txt, 02_random_00.txt,
                02_random_01.txt, 02_random_02.txt, 02_random_03.txt,
                02_random_04.txt, 03_many_cycles_00.txt, 03_many_cycles_01.txt,
                04_large_cycle_00.txt, 04_large_cycle_01.txt,
                05_largest_cycle_00.txt, 06_max_cycles_00.txt, 07_max_rho_00.txt
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="panel panel-default">
        <table class="table table-bordered table-striped th-center">
          <thead>
            <tr>
              <th>ケース名</th>
              <th>結果</th>
              <th>実行時間</th>
              <th>メモリ</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td class="text-center">00_sample_00.txt</td>

              <td class="text-center">
                <span
                  class="label label-success"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="正解"
                  >AC</span
                >
              </td>
              <td class="text-right">6 ms</td>
              <td class="text-right">3476 KB</td>
            </tr>

            <tr>
              <td class="text-center">00_sample_01.txt</td>

              <td class="text-center">
                <span
                  class="label label-success"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="正解"
                  >AC</span
                >
              </td>
              <td class="text-right">5 ms</td>
              <td class="text-right">3668 KB</td>
            </tr>

            <tr>
              <td class="text-center">01_small_00.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="不正解"
                  >WA</span
                >
              </td>
              <td class="text-right">3 ms</td>
              <td class="text-right">3596 KB</td>
            </tr>

            <tr>
              <td class="text-center">01_small_01.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="不正解"
                  >WA</span
                >
              </td>
              <td class="text-right">3 ms</td>
              <td class="text-right">3488 KB</td>
            </tr>

            <tr>
              <td class="text-center">01_small_02.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="不正解"
                  >WA</span
                >
              </td>
              <td class="text-right">3 ms</td>
              <td class="text-right">3536 KB</td>
            </tr>

            <tr>
              <td class="text-center">02_random_00.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="実行時間制限超過"
                  >TLE</span
                >
              </td>
              <td class="text-right">2206 ms</td>
              <td class="text-right">9528 KB</td>
            </tr>

            <tr>
              <td class="text-center">02_random_01.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="実行時間制限超過"
                  >TLE</span
                >
              </td>
              <td class="text-right">2206 ms</td>
              <td class="text-right">9592 KB</td>
            </tr>

            <tr>
              <td class="text-center">02_random_02.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="実行時間制限超過"
                  >TLE</span
                >
              </td>
              <td class="text-right">2206 ms</td>
              <td class="text-right">9572 KB</td>
            </tr>

            <tr>
              <td class="text-center">02_random_03.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="実行時間制限超過"
                  >TLE</span
                >
              </td>
              <td class="text-right">2206 ms</td>
              <td class="text-right">9532 KB</td>
            </tr>

            <tr>
              <td class="text-center">02_random_04.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="実行時間制限超過"
                  >TLE</span
                >
              </td>
              <td class="text-right">2206 ms</td>
              <td class="text-right">9536 KB</td>
            </tr>

            <tr>
              <td class="text-center">03_many_cycles_00.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="実行時間制限超過"
                  >TLE</span
                >
              </td>
              <td class="text-right">2206 ms</td>
              <td class="text-right">9556 KB</td>
            </tr>

            <tr>
              <td class="text-center">03_many_cycles_01.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="実行時間制限超過"
                  >TLE</span
                >
              </td>
              <td class="text-right">2206 ms</td>
              <td class="text-right">9604 KB</td>
            </tr>

            <tr>
              <td class="text-center">04_large_cycle_00.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="実行時間制限超過"
                  >TLE</span
                >
              </td>
              <td class="text-right">2206 ms</td>
              <td class="text-right">9516 KB</td>
            </tr>

            <tr>
              <td class="text-center">04_large_cycle_01.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="実行時間制限超過"
                  >TLE</span
                >
              </td>
              <td class="text-right">2206 ms</td>
              <td class="text-right">9512 KB</td>
            </tr>

            <tr>
              <td class="text-center">05_largest_cycle_00.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="実行時間制限超過"
                  >TLE</span
                >
              </td>
              <td class="text-right">2206 ms</td>
              <td class="text-right">9628 KB</td>
            </tr>

            <tr>
              <td class="text-center">06_max_cycles_00.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="実行時間制限超過"
                  >TLE</span
                >
              </td>
              <td class="text-right">2206 ms</td>
              <td class="text-right">9632 KB</td>
            </tr>

            <tr>
              <td class="text-center">07_max_rho_00.txt</td>

              <td class="text-center">
                <span
                  class="label label-warning"
                  data-toggle="tooltip"
                  data-placement="top"
                  title=""
                  data-original-title="実行時間制限超過"
                  >TLE</span
                >
              </td>
              <td class="text-right">2206 ms</td>
              <td class="text-right">9452 KB</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
  <hr />
</div>
  `

	expected := map[StatusCode]int{
		AC:  2,
		WA:  3,
		TLE: 12,
	}
	sd, err := NewSubmissionDetailPage(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Faild to create SubmissionDetailPage. got=%s", err.Error())
	}

	sm, err := sd.GetSubmissionStatusMap()
	if err != nil {
		t.Fatalf("Faild to get status map. got=%s", err.Error())
	}

	if !reflect.DeepEqual(expected, sm) {
		t.Errorf("Status map has wrong value. got=%v, want=%v", sm, expected)
	}
}
