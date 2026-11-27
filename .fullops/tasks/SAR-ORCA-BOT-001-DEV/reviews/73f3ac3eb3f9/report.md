# SAR-ORCA-BOT-001-DEV 리뷰 보고

- base: c556f2cbbdb5095d6f2069fa5f03d062a6ca291b
- head: 73f3ac3eb3f9092da6d1686a8526e5b7d6be5cc4
- snapshot: C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-d77226f1
- 리뷰어 세션: 72a42ba7-8ef2-4414-a92d-0f0a31d4e751 (구현자 997bb188…와 다름)

## 요약

적용 규칙은 인박스 `to_dev.md`의 완료 기준(S1–S6), `targets.json` 경로별 규칙, `.fullops/guide/rules/common/`이다. 관점은 버그·보안·준수 세 가지다. 대상 33개 중 28개를 검토했다. 하네스·증거 5개는 `targets.json` preset에 따라 건너뛰었다.

검토한 것:

- relay grant: `grants.go` 전부. `public_text.go`의 `kind` 파싱·서명·digest와 `ingestPublicText` 안 grant 위치(서명·authority·멱등 뒤, capacity 전). `authority`의 `self:` 판정, `human_http.go` grant-create(5분 재확인·동의)·grant-revoke, `cleanup_admission.go` 소유 검사, `store.go` 상태 로드(`newState()` 위 Unmarshal이라 `Grants` 키가 없어도 nil 맵이 아님).
- 서명 범위: `kind`가 없으면 서명·digest 바이트가 이전과 같다. 있으면 서명되어 kind 바꿔치기·제거는 `invalid_signature`다. reply_to와 같이 쓸 수 없다. 모르는 kind는 `approval_required`다.
- 재생 금지: grant가 만들 때의 보내는 쪽·받는 쪽 키 목록을 기록한다. 그중 하나라도 철회되거나 agent 철회·`self:` 종료·기한이 오면 sweep이 grant를 지운다. 새 키는 grant 키 목록에 없어 통과하지 못한다.
- adapter: `orca.ts` 전부(스키마·JCS·2 KB, argv 단일 operand 검사, `shell:false`, 도어벨 본문에 text·발신자 없음, 실패 시 marker 해제). `inbox.ts` `Wake` 일반화·부분 재시도, `text.ts` `sendKind`·수신 `kind` 검사, `mcp.ts` 진입점, MCP receive 출력.
- 웹 화면은 기존 컴포넌트(`agents`·`choice`·`consent`·`danger-button`)만 쓴다. 영어 문구도 있다.
- 문서: D03 architecture, D05 interface-design, D10 module-design, adapters README, Grok SKILL, dev.md 교훈이 코드 동작과 맞다.

검증 근거:

- `lint.json`(441d1a41): ERROR 0, WARNING 14. make 2개는 Windows에서 실행할 수 없다. head 73f3ac3e는 441d1a41에 하네스 기록 2개(`to_dev.md`·`qa/report.md`)만 더한다.
- QA 보고: Linux clean clone에서 `make test`·`make lint`가 통과했다(`check_compose.py`만 host docker). 격리 relay에서 S1–S6 왕복도 통과했다.
- 리뷰어가 snapshot에서 `go test ./internal/relay -run 'Grant|Translation' -count=1`을 한 번 확인했다: ok. 파일 변경은 없다.

## 발견 사항

critical 0, high 0, medium 0, low 3.

1. low [준수] `adapters/skills/knowslink/SKILL.md:60`: Grok 지시의 `<plugin.js>`·`<this agent folder>`·`<Orca agent>` 자리 표시를 채우는 방법이 없다. S6는 CLI를 직접 실행해 이 공백을 덮지 않았다. 운영 왕복 전에 보완한다.
2. low [버그] `adapters/src/text.ts:245-261`: 이전 버전 adapter는 `kind` 있는 envelope를 버린다. 받는 쪽 plugin을 갱신하기 전에 grant를 만들면 만료까지 재전달과 실패가 반복된다. 배포 순서를 운영 절차에 적는다.
3. low [버그] `adapters/src/orca.ts:109-117`: answer 도어벨의 `--item`을 claim 뒤 peek로 얻는다. 그 사이 `receive`가 메시지를 가져가면 `--item` 없이 울린다. 본문은 노출되지 않는다.

## 결론

accept. 발송별 승인·서명·키 경계를 넓히지 않았다. 예외는 같은 owner `self:` 두 agent 사이의 고정 3형식 grant 하나뿐이다. grant는 회원 세션에서만 만들고 끈다. agent API 경로와 `confirmed:true` 경로는 없다. S1–S6 증명 방법이 지시서와 맞다. low 3건은 병합을 막지 않으며 `PLANS.md` 후속(운영 연결 과제 전)으로 남긴다.
