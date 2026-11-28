# SAR-MSG-TTL-1H-001-DEV 리뷰 보고

- base: a6ee3d7e07fbae272e6a9aa3fb0fe650d4bcb7f5
- head: 4ef924cdbd8bb258d95502c6e65620742153e454
- snapshot: C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-95788b93

## 요약

재작업 1차(D076)를 리뷰 초점 ①–⑥으로 봤다. high 0, medium 1, low 3이다.

- ① 1시간은 `selfAgents`(같은 owner, 양쪽 활성·미철회, `authority` scope `self:`)일 때만 준다. relay가 from/to로 판정하고 송신 입력(exp·kind)으로 넓힐 수 없다. 철회 agent는 `authority`가 거부하고, 다른 owner·친구·legacy pair는 text 180초·`/v1/send` 300초다. grant 판정도 같은 함수로 묶였고 기존 조건과 같다.
- ② `envelopeTTL`은 `relay.approval.request`면 300초로 고정한다. 사람 승인 대기는 바뀌지 않았다.
- ③ text: exp>now+1시간만 `ttl_too_long`이고, 1시간 안은 쌍 상한으로 줄여 receipt Exp에 실제 기한을 넣는다. 170초 클라이언트는 그대로, 3590초 클라이언트는 친구 사이 180초로 줄어든다. reply는 adapter·원격 MCP 모두 부모 receipt Exp로 자르고 relay도 부모 Exp 초과를 거부한다. `/v1/send`는 clamp 없이 거부만 한다(기존 300초 클라이언트 동작 유지).
- ④ Receipt 구조에 본문이 없고 로그 변경도 없다. Redis 값 TTL(`ephemeral.go:102`)과 sweep 삭제(`store.go:335`)가 같은 `messageRetention`을 쓴다.
- ⑤ QA 보고 R3 서술(전역 100건 큐, Redis 128MiB noeviction, 운영 사용량 미확인)이 코드와 맞다. 위험 자체는 medium으로 남긴다(아래 1).
- ⑥ 방침 ko·en 문구가 동작과 같고 2026-12-28 전체 1시간 전환 고지가 두 언어·랜딩·대화 도움말에 있다. `policyEffective`=2026-11-28.
- `remote_mcp.go`→`remote_tools.go` 분리는 이동이며 차이는 TTL 한 줄·설명 한 줄·`remoteReply` 도입뿐이다. `Mcp-Name` 비교를 `params.name` method로 한정한 변경도 맞다.

검증(lite): `go test ./internal/relay/ -count=1` 통과, 표적 `-run 'TestDeliveryCeilingBySelfAgents|PublicText|Grant|Policy|RemoteEra|Landing|HumanPages'` 14건 PASS, `go vet` 통과. integration 태그·adapter 테스트는 다시 돌리지 않았다.

## 발견 사항

1. medium `internal/relay/public_text.go:25` — 1시간 보관으로 한 member가 꺼 둔 자기 agent에 100건을 보내 전역 미배달 큐(capacity.go:96)를 최대 1시간 막을 수 있다. R3에 적힌 알려진 위험이고 정책 변경은 지시상 후속이다. 2026-12-28 전체 전환 전에 member별 큐 상한이 필요하다.
2. low `internal/relay/store.go:335` — Exp=수락+1시간인 웹·원격 MCP 발송은 만료 순간 receipt까지 지워진다. 송신자가 `failed:expired`를 보지 못하고 없는 ID로 본다.
3. low `internal/relay/landing_pages.go:50` — 공개 문구 '내 에이전트끼리 1시간'은 SelfConnect 미동의 legacy member의 자기 agent(180초)를 빼고 말한다. 보관이 더 짧은 쪽이라 개인정보 위험은 없다.
4. low `adapters/src/text.ts:67` — 새 adapter(3590초)는 이 커밋 이전 relay에서 `ttl_too_long`이다. 배포 때 relay 먼저, adapter 갱신 뒤를 ops에 남긴다.

## 결론

accept. 요청 범위(자기 agent 사이만 1시간, 그 밖 유지, relay clamp, 승인 300초 유지, 방침 고지)를 정확히 구현했다. medium 1건은 지시가 후속으로 둔 용량 정책이며 이번 병합을 막지 않는다.
