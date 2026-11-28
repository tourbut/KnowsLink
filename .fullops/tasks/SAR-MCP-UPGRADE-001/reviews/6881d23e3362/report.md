# SAR-MCP-UPGRADE-001 리뷰 보고

- base: 4d3998e52dcb7e31045b83c3e524cac2a6c21dab
- head: 6881d23e336251a9d30a1d65875dc727d4fd633b
- snapshot: C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-fa723ea9

## 요약

원격 `POST /mcp`의 MCP 2026-07-28 dual-era 지원(remote_era.go 신설, remote_mcp.go 연결, D05·system-design 문서, 단위·integration 테스트)을 주인 지시대로 보안 경계 중심으로 1회 리뷰했다.

보안 초점 점검:

1. **헤더·본문 불일치 거부(-32020)가 인증 전인가**: 그렇다. `remoteMCP`는 본문 Strict 파싱 직후, `remotePrincipal`(인증) 전에 `remoteEra`를 부른다. `/mcp` 경로는 `POST /mcp` 하나뿐이고(GET/DELETE는 ServeMux 405), events/subscribe도 `remoteEra` 뒤에서 분기한다. 우회 경로 없음. `remoteEra`는 `json.Unmarshal`, 실행부는 `Strict`(같은 encoding/json, 중복 키는 jcs가 거부)로 같은 키 해석을 써서 헤더 검사와 실제 실행이 다른 `name`을 보는 경우가 없다. Base64 sentinel이 잘못되면 원문 그대로 비교되어 불일치로 거부된다. 인증 전 응답은 정적 정보(지원 버전 목록, 요청 값 echo)만 담아 사용자 정보 노출이 없다.
2. **Mcp-Method/Mcp-Name 누락 허용**: 안전하다. 라우팅·도구 선택은 모두 본문(`in.Method`, `params.name`)을 따르고 헤더는 비교에만 쓴다. 누락 허용은 `shortcut:` 주석으로 한계와 전환 조건(헤더 라우팅 게이트웨이 도입 시)이 적혀 있다.
3. **discover/list 캐시 `ttlMs:300000`·`cacheScope:"public"`**: `server/discover`, `tools/list`(`remoteTools()` 인자 없음), `events/list`(`remoteEventDefinition()` 정적)만 붙고 모두 설치·사용자와 무관한 정적 값이다. `tools/call` 결과에는 붙지 않는다(integration 테스트가 `ttlMs` 부재 확인). 사용자별 정보(agent·권한·도구 결과) 공유 캐시 노출 없음.
4. **legacy initialize 버전 에코**: `remoteVersions`(2026-07-28, 2025-11-25, 2025-06-18, 2025-03-26)에 있는 값만 그대로 돌려주고 그 밖(예: 2024-11-05)은 이전처럼 2025-11-25다. 미지원 버전을 받아들이지 않는다. 변화는 2025-06-18·2025-03-26 요청에 이전(2025-11-25) 대신 요청 버전을 돌려주는 것뿐이며, 서버 동작은 버전별로 갈리지 않아 Dots·Grok Bot·Claude Code 결과 모양은 같다(legacy 결과에 `resultType`·`ttlMs`·`cacheScope`가 없음을 integration이 확인). 미지원 `MCP-Protocol-Version` 헤더에 400을 주는 새 동작은 2025-06-18 명세의 요구와 일치한다.
5. **OAuth·인증·grant·발송 승인 경계**: 바뀌지 않았다. `remotePrincipal` 두 번 확인, `remoteTool`(grant·발송 승인) 경로, 401 challenge, `remoteSubscription` 내부 인증 모두 그대로다. 새 404 `-32601`은 인증 뒤에서만 나온다. modern 결과 가공(`v.(map[string]any)`)은 모든 성공 분기가 map을 돌려 panic 경로가 없다.

검증: `go test -count=1 -run 'TestRemoteModernHeaderChecks' -v ./internal/relay/` PASS, `go test -count=1 -run 'Remote' ./internal/relay/` ok. integration(`-tags=integration`, 격리 Postgres·Redis)은 구현자 QA 보고(qa/report.md)의 기록을 인용하며 다시 돌리지 않았다. lint 기록(lint.json) passed, 경고 2(SIZE-001, LINT-002).

## 발견 사항

- low — `internal/relay/remote_mcp.go:80`: modern `events/subscribe`·`events/unsubscribe` 결과에 `resultType`·`_meta.serverInfo`가 없다. D05 41행 '모든 modern 결과' 서술과 다르다. 보안 영향 없음.
- low — `internal/relay/remote_era.go:50`: `Mcp-Name`을 항상 `params.name`과 비교해 `resources/read`(uri)는 404 대신 400 `-32020`을 받는다. 미제공 메서드라 오류 코드 차이뿐, 우회 아님.
- low — `internal/relay/remote_mcp.go`: SIZE-001 300줄 상한 초과(309줄).

high·medium 없음.

## 결론

accept. 보안 경계 5개 항목 모두 문제 없음. low 3건은 후속 정리 대상이다.
