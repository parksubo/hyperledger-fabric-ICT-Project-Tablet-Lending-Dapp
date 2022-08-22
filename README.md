# Tablet 실행방법 및 순서

## 1. Network 실행

network 경로에서 (다음 쉘스크립트 순차 실행):

- ./startnetwork.sh
- ./createchannel.sh
- ./setAnchorPeerUpdate.sh
- ./deployCC.sh

## 2. CCP 생성

application/ccp 경로에서 (다음 쉘스크립트 실행):

- ./ccp-generate.sh

## 3. 인증서 가져오기 및 지갑 생성

application 경로에서 (다음 쉘스크립트 순차 실행):

- ./getCert.sh
- ./addToWallet.js

## 4. node.js 모듈 설치

application 경로에서

- npm install

## 5. 서버 실행

application 경로에서

- node addWallet.js : 지갑 발행 (테스트 위함)
- node server.js : tablet 어플리케이션 웹서버 실행
