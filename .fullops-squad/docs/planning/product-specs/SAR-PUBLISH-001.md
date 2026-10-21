---
id: D02
title: 앱 공개 게시용 개인정보처리방침·서비스 약관 정본
status: review
updated: 2026-10-21
owner: designer
tasks: [SAR-PUBLISH-001]
upstream: [D01]
summary: "Google OAuth 앱 공개 게시(D017)에 필요한 /privacy·/terms의 KO/EN 문안 초안, 화면 기준, 홈페이지 링크 위치와 완료 조건을 정한다. 법률 검토 전 초안이다"
---

# 앱 공개 게시용 개인정보처리방침·서비스 약관 — SAR-PUBLISH-001

## 상태와 적용 범위

> 법률 검토 전 초안이다. 이 문안은 법률 자문이 아니다. 아래 '게시 전 확인 항목'이 닫히기 전에는 법적 효력을 갖는 최종본으로 보지 않는다. front matter `status: review`는 D02 인덱스 상태와 맞춘 값이며 법률 검토 완료를 뜻하지 않는다.

2026-10-21 사용자 결정 D017로 Google OAuth 앱(knowslink-auth)을 공개 게시한다. Google 브랜딩 설정에는 홈페이지·개인정보처리방침·서비스 약관 URL이 필요하다. 세 URL은 승인된 도메인 `knowslog.com` 아래에 있어야 한다. 현재 `link.knowslog.com`에는 `/privacy`·`/terms`가 없다(404). 이 문서가 두 페이지의 제품 정본이다. 같은 날 결정 D018로 서버 캘린더 연결을 제거한다. 그래서 이 문안은 제거 뒤 기준으로 쓴다. 제거 규칙은 [SAR-SCHEDULE-001 '철회 — D018'](SAR-SCHEDULE-001.md#철회--d018-2026-10-21)이다.

| URL | 용도 |
|---|---|
| `https://link.knowslog.com/` | 홈페이지(기존 공개 첫페이지) |
| `https://link.knowslog.com/privacy` | 개인정보처리방침(새로 만든다) |
| `https://link.knowslog.com/terms` | 서비스 약관(새로 만든다) |

## 문안의 사실 근거

문안의 수집 항목과 기간은 기준 커밋 `dcf8825`의 코드와 D10([module-design](../../design-docs/module-design.md))에서 확인했다. DEV는 구현 때 같은 사실을 다시 확인한다. 코드와 문안이 다르면 추측하지 않고 designer에게 묻는다.

| 문안 항목 | 근거(코드) |
|---|---|
| Google 로그인은 `openid email`만 요청한다. 이름·사진은 받지 않는다 | `internal/relay/google.go` `Scopes: {openid, email}`, 인증된 이메일만 사용 |
| 계정: Google 계정 식별자·이메일·발급자, 회원이 정한 이름(1–40자), 친구 코드, 언어 | `identity.go` `Member`, `State.Identities` |
| 이메일 코드 로그인: 이메일과 10분 코드 | `identity.go` `Challenge`, `codeTTL` |
| 세션: 최대 12시간, 60분 미사용 시 종료 | `identity.go` `sessionAbsolute`, `sessionIdle` |
| 친구 관계·요청·차단, 보낸 요청의 입력 코드 끝 4자 | `human.go` `Friend` |
| 에이전트 이름·클라이언트 종류·공개 키·자격 해시, 원격 연결의 서버 보관 서명 키·이벤트 수신 주소 | `store.go` `Agent`·`Key`, `remote_auth.go` `remoteInstall`, `remote_events.go` `remoteSubscription` |
| 메시지·전달 기록·승인 요청은 받은 때부터 24시간 뒤 삭제 | `store.go` sweep `24*time.Hour`(Messages·Gates) |
| 연결 해제한 에이전트는 24시간 뒤 삭제 | `connections.go` `revokedRetention` |
| IP 주소는 요청 횟수 제한에만 쓰고 최대 1시간 보관 | `identity.go` `http:ip:`·`send:ip:` bucket, 1시간 정리 |
| 쿠키: 세션·로그인 진행·언어(365일). 분석·광고 스크립트 없음 | `identity.go` `__Host-kl_session`·`__Host-kl_pending`, `google.go` 쿠키, `__Host-kl_lang` |
| 계정 삭제 화면은 없다 | 회원 삭제 경로 없음. 삭제는 이메일 요청으로 처리한다 |
| 캘린더 데이터는 받지 않는다 | D018 제거 뒤 기준 |

## 화면 기준

- 경로: `GET /privacy`, `GET /terms`. 로그인 없이 열린다. 로그인 상태여도 같은 내용이다.
- 스타일: 공개 첫페이지(`landing_pages.go`)의 `memberStyle`·`body.landing` 토큰과 헤더·하단을 재사용한다. 새 색·글꼴·JS·외부 리소스를 추가하지 않는다. 본문 폭은 읽기 폭(약 72ch)으로 제한한다.
- 헤더: 브랜드 'KnowsLink'(`/`로 이동)와 기존 언어 전환(`lang-nav`)만 둔다. 첫페이지의 앵커 메뉴는 넣지 않는다.
- 본문 순서: h1, 시행일 한 줄(`muted`), 다른 정책 페이지 링크 한 줄, 번호 있는 h2 절, 연락처 절.
- 하단: 첫페이지와 같은 하단에 '개인정보처리방침', '서비스 약관', '홈'을 둔다. 현재 페이지 링크는 `aria-current="page"`다.
- 언어: 기존 i18n 구조와 `/lang/{code}` 전환을 쓴다. 영어판에는 한국어가 남지 않는다. `<html lang>`과 `Content-Language`는 현재 언어다. 단, 영어판 하단의 기존 '설치·사용 문서' 링크는 페이지 밖 한국어 문서임을 기존 문구('(in Korean)')로 표시한다.
- 고유 `<title>`: '개인정보처리방침 · KnowsLink' / 'Privacy Policy · KnowsLink', '서비스 약관 · KnowsLink' / 'Terms of Service · KnowsLink'.
- 피할 것: 크림 배경, 알약 버튼, 모노스페이스 라벨, 장식 그라디언트, 이모지 아이콘, 접는 상자(본문 전체가 바로 보여야 한다).
- 반응형·접근성: 1440·390·320px에서 가로 넘침이 없다. 본문 대비 4.5:1, 보이는 focus, 본문 건너뛰기 링크를 유지한다. 이메일 주소는 `mailto:` 링크다.
- 시행일: 코드의 단일 상수 하나로 두 언어가 같은 날짜를 보인다. 값은 coor가 dev 배정 때 정한다(예정 운영 배포일). 배포가 늦어지면 배포 전 상수를 맞춘다.

## 홈페이지와 회원 화면의 링크 위치

Google은 홈페이지에서 개인정보처리방침으로 가는 링크를 요구한다.

| 위치 | 변경 |
|---|---|
| 첫페이지 하단 | 기존 '설치·사용 문서', '맨 위로' 앞에 '개인정보처리방침', '서비스 약관'을 둔다. 순서: 개인정보처리방침 · 서비스 약관 · 설치·사용 문서 · 맨 위로. |
| 첫페이지 `#start` 로그인 절 | 로그인 버튼 아래 `muted` 한 줄: '로그인하기 전에 [서비스 약관]과 [개인정보처리방침]을 확인하세요.' (EN 'Before you sign in, read the [Terms of Service] and the [Privacy Policy].') 동의 체크박스는 두지 않는다. |
| 회원 설정 | 기존 '서비스 안내' 링크 옆에 '개인정보처리방침'·'서비스 약관' 링크를 둔다. 주 메뉴는 바꾸지 않는다. |

## 개인정보처리방침 — 한국어

**h1** 개인정보처리방침

시행일: {시행일}

[서비스 약관]도 함께 확인하세요.

**1. 이 방침의 범위**

이 방침은 KnowsLink(`link.knowslog.com`, 이하 '서비스')가 회원의 개인정보를 어떻게 수집·이용·보관·삭제하는지 설명합니다. KnowsLink는 서로 친구가 된 사람들의 에이전트를 연결하고, 에이전트와 사람이 주고받는 메시지를 전달하는 서비스입니다.

**2. 수집하는 정보**

- 로그인 정보: Google로 로그인하면 Google 계정 고유 식별자와 인증된 이메일 주소를 받습니다. 이름·프로필 사진·연락처는 Google에서 받지 않습니다. 이메일 코드로 로그인하면 이메일 주소와 일회용 확인 코드를 처리합니다.
- 회원이 정하는 정보: 친구에게 보일 이름, 화면 언어.
- 서비스가 만드는 정보: 친구 코드, 친구 요청·수락·차단 상태, 보낸 친구 요청에 입력한 코드의 끝 4자.
- 에이전트 연결 정보: 에이전트 이름, 클라이언트 종류(예: Grok, Dots), 공개 키, 연결 자격의 해시값. 원격으로 연결한 에이전트는 서비스가 그 에이전트의 서명 키와 이벤트 수신 주소를 보관합니다.
- 메시지: 웹과 에이전트에서 보낸 메시지 본문, 보낸 사람·받는 사람·요청 유형·시각·전달 상태, 사람 승인 요청과 결정.
- 접속·보안 정보: 로그인 세션, 요청 횟수 제한에 쓰는 IP 주소, 서비스 쿠키.

KnowsLink는 회원의 캘린더·일정·주소록·위치·결제 정보를 수집하지 않습니다.

**3. 이용 목적**

- 로그인과 계정 확인, 부정 이용 방지.
- 친구 연결과 에이전트 연결, 메시지 전달과 전달 상태 표시.
- 사람 승인이 필요한 요청을 보여 주고 결정을 기록.
- 과도한 요청 제한과 서비스 보안 유지.

KnowsLink는 개인정보를 판매하지 않습니다. 광고에 쓰지 않습니다. 메시지 본문을 AI 모델 학습에 쓰지 않습니다.

**4. 보관 기간과 삭제**

| 정보 | 보관 기간 |
|---|---|
| 계정 정보(로그인 정보, 이름, 친구 코드, 언어, 친구 관계) | 계정 삭제 요청을 처리할 때까지 |
| 이메일 확인 코드 | 발급 후 10분 |
| 로그인 세션 | 최대 12시간. 60분 동안 쓰지 않으면 종료 |
| 메시지 본문, 전달 기록, 승인 요청 | 서비스가 받은 때부터 24시간 뒤 삭제 |
| 연결 해제한 에이전트 정보 | 해제 후 24시간 뒤 삭제 |
| 요청 횟수 제한용 IP 주소 | 최대 1시간 |

장애 복구용 백업에는 정보가 더 오래 남을 수 있습니다. 백업은 {백업 보존 기간}이 지나면 삭제합니다.

**5. 다른 사람에게 전달되는 정보**

- 친구와 친구의 에이전트: 회원이 보낸 메시지와 회원이 정한 이름은 회원이 고른 받는 사람과 그 사람의 에이전트에 전달됩니다.
- 에이전트 호스트: 받는 사람의 에이전트가 외부 서비스(예: Grok, Dots)에서 실행되면 그 서비스가 메시지를 처리합니다. 그 처리에는 해당 서비스의 정책이 적용됩니다.
- 처리 위탁: 서비스 운영에 다음 업체를 이용합니다.

| 업체 | 처리 내용 | 위치 |
|---|---|---|
| Google LLC | Google 로그인 인증 | 미국 |
| Cloudflare, Inc. | 네트워크 전달과 보안(접속 정보가 거쳐 감) | 미국 및 전 세계 |
| {이메일 발송 업체} | 이메일 확인 코드 발송(이메일 코드 로그인을 제공할 때) | {위치} |

KnowsLink는 법령에 따른 요청이 있는 경우를 빼고 위 범위 밖의 제3자에게 개인정보를 제공하지 않습니다.

**6. Google 사용자 데이터**

KnowsLink가 Google API에서 받은 정보를 사용하고 다른 곳으로 옮기는 방식은 Limited Use 요건을 포함한 [Google API Services User Data Policy](https://developers.google.com/terms/api-services-user-data-policy)를 따릅니다. KnowsLink는 Google에서 받은 계정 식별자와 이메일을 로그인과 계정 확인에만 씁니다.

**7. 쿠키**

서비스는 로그인 유지, 로그인 진행 확인, 화면 언어 기억에 필요한 쿠키만 씁니다. 언어 쿠키는 365일 동안 보관됩니다. 분석·광고 쿠키는 쓰지 않습니다.

**8. 회원의 권리**

회원은 자기 개인정보의 열람·정정·삭제·처리 정지를 요청할 수 있습니다. 이름과 언어는 설정에서 바로 바꿀 수 있습니다. 계정 삭제와 그 밖의 요청은 아래 연락처로 보내 주세요. 본인 확인 뒤 요청을 받은 날부터 10일 안에 처리 결과를 알려 드립니다.

**9. 안전 조치**

서비스는 모든 접속을 HTTPS로 받습니다. 연결 자격과 세션 토큰은 해시값으로 보관합니다. 운영 비밀값은 저장소와 화면에 두지 않습니다. 운영 데이터베이스는 외부에 직접 공개하지 않습니다.

**10. 아동**

서비스는 만 14세 미만을 대상으로 하지 않습니다. 만 14세 미만의 정보를 받은 사실을 알게 되면 삭제합니다.

**11. 방침 변경**

방침을 바꾸면 시행 7일 전에 이 페이지에 알립니다. 회원에게 불리한 변경은 시행 30일 전에 알립니다.

**12. 연락처**

개인정보 보호책임자: {보호책임자 표시}
이메일: [knowslog01@gmail.com](mailto:knowslog01@gmail.com)

## Privacy Policy — English

**h1** Privacy Policy

Effective date: {effective date}

Please also read the [Terms of Service].

**1. Scope**

This policy explains how KnowsLink (`link.knowslog.com`, "the Service") collects, uses, keeps, and deletes members' personal information. KnowsLink connects the agents of people who have become friends, and delivers the messages that agents and people send to each other.

**2. Information we collect**

- Sign-in information: When you sign in with Google, we receive your Google account identifier and your verified email address. We do not receive your name, profile photo, or contacts from Google. When you sign in with an email code, we process your email address and a one-time verification code.
- Information you choose: the name your friends see, and your display language.
- Information the Service creates: your friend code, the state of friend requests, acceptances, and blocks, and the last 4 characters of the code you entered in a friend request you sent.
- Agent connection information: agent name, client type (for example, Grok or Dots), public key, and a hash of the connection credential. For agents connected remotely, the Service stores that agent's signing key and its event delivery address.
- Messages: the text you send from the web or from agents; the sender, recipient, request type, time, and delivery state; and requests for human approval and their decisions.
- Access and security information: sign-in sessions, IP addresses used for rate limiting, and Service cookies.

KnowsLink does not collect your calendar, schedule, address book, location, or payment information.

**3. How we use information**

- To sign you in, confirm your account, and prevent misuse.
- To connect friends and agents, deliver messages, and show delivery state.
- To show requests that need a person's approval and record the decision.
- To limit excessive requests and keep the Service secure.

KnowsLink does not sell personal information. We do not use it for advertising. We do not use message text to train AI models.

**4. Retention and deletion**

| Information | Retention |
|---|---|
| Account information (sign-in information, name, friend code, language, friend relationships) | Until we process your account deletion request |
| Email verification code | 10 minutes after it is issued |
| Sign-in session | Up to 12 hours; ends after 60 minutes without use |
| Message text, delivery records, approval requests | Deleted 24 hours after the Service receives them |
| Information about a disconnected agent | Deleted 24 hours after disconnection |
| IP addresses for rate limiting | Up to 1 hour |

Backups for recovering from failures can keep information longer. We delete backups after {backup retention period}.

**5. Information shared with others**

- Your friends and their agents: The messages you send and the name you chose are delivered to the recipient you select and to that person's agent.
- Agent hosts: When a recipient's agent runs on an outside service (for example, Grok or Dots), that service processes the message under its own policies.
- Service providers: We use the following providers to run the Service.

| Provider | Purpose | Location |
|---|---|---|
| Google LLC | Google sign-in authentication | United States |
| Cloudflare, Inc. | Network delivery and security (access information passes through it) | United States and worldwide |
| {email provider} | Sending email verification codes (when email-code sign-in is offered) | {location} |

KnowsLink does not give personal information to other third parties, except when the law requires it.

**6. Google user data**

KnowsLink's use and transfer of information received from Google APIs adheres to the [Google API Services User Data Policy](https://developers.google.com/terms/api-services-user-data-policy), including the Limited Use requirements. KnowsLink uses the account identifier and email received from Google only to sign you in and confirm your account.

**7. Cookies**

The Service uses only the cookies needed to keep you signed in, to check a sign-in in progress, and to remember your display language. The language cookie is kept for 365 days. We do not use analytics or advertising cookies.

**8. Your rights**

You can ask to access, correct, delete, or stop the processing of your personal information. You can change your name and language in Settings at any time. Send account deletion and other requests to the contact below. After we confirm your identity, we tell you the result within 10 days of receiving your request.

**9. Security**

The Service accepts connections only over HTTPS. We store connection credentials and session tokens as hashes. Operational secrets are not kept in the code repository or shown on screen. The production database is not directly exposed to the internet.

**10. Children**

The Service is not intended for children under 14. If we learn that we received information from a child under 14, we delete it.

**11. Changes to this policy**

We announce changes on this page 7 days before they take effect. We announce changes that are unfavorable to members 30 days before they take effect.

**12. Contact**

Privacy officer: {privacy officer}
Email: [knowslog01@gmail.com](mailto:knowslog01@gmail.com)

## 서비스 약관 — 한국어

**h1** 서비스 약관

시행일: {시행일}

[개인정보처리방침]도 함께 확인하세요.

**1. 서비스**

KnowsLink는 서로 친구가 된 사람들의 에이전트를 연결하고 메시지를 전달합니다. 이 약관은 KnowsLink를 쓰는 모든 회원에게 적용됩니다. 로그인하기 전에 이 약관을 읽어 주세요.

**2. 계정**

- 서비스는 만 14세 이상만 쓸 수 있습니다.
- 회원은 자기 Google 계정 또는 자기 이메일로만 로그인합니다.
- 다른 사람을 사칭하는 이름을 쓰지 않습니다.
- 회원은 자기 계정과 연결한 에이전트의 자격 정보를 안전하게 관리합니다.

**3. 에이전트와 메시지**

- 회원은 자기가 연결한 에이전트가 보내는 메시지에 책임을 집니다.
- 메시지는 서로 수락한 친구와 회원 자신의 에이전트 사이에서만 주고받습니다.
- 사람 승인이 필요한 요청은 회원이 승인할 때만 다음 단계로 넘어갑니다.
- KnowsLink는 외부 에이전트 서비스(예: Grok, Dots)를 운영하지 않습니다. 그 서비스의 이용에는 해당 서비스의 약관이 적용됩니다.
- 메시지는 서비스가 받은 때부터 24시간 뒤 삭제됩니다. 중요한 내용은 회원이 따로 보관하세요.

**4. 금지 행위**

회원은 다음 행위를 하지 않습니다.

- 법을 어기거나 다른 사람의 권리를 침해하는 내용 보내기.
- 스팸, 원하지 않는 대량 메시지, 괴롭힘.
- 다른 사람의 계정이나 에이전트를 허락 없이 쓰기.
- 서비스의 보안 장치를 우회하거나, 서비스에 과도한 부하를 주거나, 서비스를 공격하기.

**5. 회원의 콘텐츠**

회원이 보낸 메시지의 권리는 회원에게 있습니다. 회원은 KnowsLink가 그 메시지를 받는 사람에게 전달하고 보관 기간 동안 저장하도록 허락합니다. KnowsLink는 이 목적 밖으로 메시지를 쓰지 않습니다.

**6. 서비스 변경과 중단**

- KnowsLink는 현재 무료로 제공되며 기능이 바뀔 수 있습니다.
- 메시지 전달과 화면 표시가 항상 성공한다고 보장하지 않습니다.
- 서비스를 크게 바꾸거나 끝낼 때는 이 페이지와 서비스 화면에 미리 알립니다.

**7. 이용 제한과 계정 삭제**

- 회원이 이 약관을 어기면 KnowsLink는 해당 계정이나 에이전트 연결을 제한하거나 중지할 수 있습니다.
- 회원은 언제든지 아래 연락처로 계정 삭제를 요청할 수 있습니다.

**8. 책임의 한계**

KnowsLink는 관련 법이 허용하는 범위에서, 회원이 고른 에이전트·외부 서비스의 동작과 회원끼리 주고받은 내용에 책임을 지지 않습니다. KnowsLink의 고의나 중대한 과실로 생긴 손해에는 이 조항을 적용하지 않습니다.

**9. 약관 변경**

약관을 바꾸면 시행 7일 전에 이 페이지에 알립니다. 회원에게 불리한 변경은 시행 30일 전에 알립니다.

**10. 준거법**

이 약관은 대한민국 법을 따릅니다.

**11. 연락처**

이메일: [knowslog01@gmail.com](mailto:knowslog01@gmail.com)

## Terms of Service — English

**h1** Terms of Service

Effective date: {effective date}

Please also read the [Privacy Policy].

**1. The Service**

KnowsLink connects the agents of people who have become friends and delivers their messages. These terms apply to every member who uses KnowsLink. Please read them before you sign in.

**2. Accounts**

- You must be at least 14 years old to use the Service.
- Sign in only with your own Google account or your own email address.
- Do not use a name that impersonates another person.
- Keep your account and the credentials of the agents you connect safe.

**3. Agents and messages**

- You are responsible for the messages that the agents you connect send.
- Messages are exchanged only between friends who have accepted each other and between your own agents.
- A request that needs a person's approval moves forward only when you approve it.
- KnowsLink does not operate outside agent services (for example, Grok or Dots). Their own terms apply when you use them.
- Messages are deleted 24 hours after the Service receives them. Keep anything important somewhere else.

**4. Prohibited conduct**

Do not do any of the following.

- Send content that breaks the law or violates the rights of others.
- Send spam, unwanted bulk messages, or harassment.
- Use another person's account or agent without permission.
- Bypass the Service's security, put excessive load on the Service, or attack it.

**5. Your content**

You keep the rights to the messages you send. You allow KnowsLink to deliver them to the recipient and to store them during the retention period. KnowsLink does not use your messages for any other purpose.

**6. Changes to and suspension of the Service**

- KnowsLink is currently free, and its features can change.
- We do not guarantee that message delivery or display always succeeds.
- Before we make a major change to the Service or end it, we announce it on this page and in the Service.

**7. Restrictions and account deletion**

- If you break these terms, KnowsLink can restrict or suspend your account or agent connections.
- You can ask to delete your account at any time through the contact below.

**8. Limitation of liability**

To the extent the law allows, KnowsLink is not responsible for the behavior of the agents and outside services you choose, or for the content members exchange. This section does not apply to damage caused by KnowsLink's intent or gross negligence.

**9. Changes to these terms**

We announce changes on this page 7 days before they take effect. We announce changes that are unfavorable to members 30 days before they take effect.

**10. Governing law**

These terms are governed by the laws of the Republic of Korea.

**11. Contact**

Email: [knowslog01@gmail.com](mailto:knowslog01@gmail.com)

## 게시 전 확인 항목

`{ }`로 표시한 값은 게시 전에 채운다. DEV는 값이 정해지기 전에는 같은 자리 표시를 화면에 내보내지 않는다. 값이 없으면 coor를 거쳐 묻는다.

| ID | 항목 | 담당 | 메모 |
|---|---|---|---|
| P-01 | 시행일 | coor(dev 배정 때) | 예정 운영 배포일. 두 언어 공통 상수 |
| P-02 | 개인정보 보호책임자 표시(이름 또는 직책) | user | 사용자가 정할 때까지 'KnowsLink 운영자' / 'KnowsLink operator'를 쓴다 |
| P-03 | 백업 보존 기간 | ops | 현재 배포 전 백업 정리 주기를 확인해 적는다. 확인 전에는 4절 백업 문장을 '백업은 정기적으로 교체하며 교체할 때 함께 삭제합니다.' / 'Backups are replaced regularly, and the information in them is deleted when they are replaced.'로 쓴다 |
| P-04 | 이메일 발송 업체·위치 | ops | 운영에서 이메일 코드 로그인이 꺼져 있으면 표의 그 행을 뺀다 |
| P-05 | 법률 검토 | user | 개인정보 보호법의 국외 이전 고지, 처리 위탁 고지, 약관 책임 조항의 표현을 확인한다. 이 문서는 효력을 단정하지 않는다 |

## 완료 조건

| ID | 관찰 가능한 완료 조건 |
|---|---|
| PC-01 | `/privacy`·`/terms`가 로그인 없이 200으로 열린다. 공개 ingress에서도 열린다(`https://link.knowslog.com/privacy`, `/terms`). |
| PC-02 | KO/EN 본문이 이 문서의 문안과 같다. 영어판에 한국어가 남지 않는다(하단 외부 문서 링크의 기존 표시 제외). `{ }` 자리 표시가 화면에 없다. |
| PC-03 | 첫페이지 하단과 `#start`, 회원 설정에서 두 페이지로 간다. 두 페이지끼리 서로 링크한다. |
| PC-04 | 화면 기준(스타일 재사용, 금지 요소 없음, 1440·390·320px 가로 넘침 없음, 고유 title, `aria-current`)을 지킨다. |
| PC-05 | 정책 문안의 수집·보관 사실이 제거 뒤 코드와 같다. 캘린더 문구가 없다. |
| PC-06 | designer 직접 시각 검수(KO/EN, 3폭)와 tester 독립 확인이 있다. |

## 변경 기록

- 2026-10-21: SAR-PUBLISH-001에서 초안을 작성했다(D017·D018). 법률 검토 전이다.
