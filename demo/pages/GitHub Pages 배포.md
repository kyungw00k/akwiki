---
title: GitHub Pages Deployment
titleKo: GitHub Pages 배포
type: Article
tags:
  - guide
  - deployment
---

# GitHub Pages 배포

akwiki로 만든 위키를 GitHub Pages로 무료 배포하는 방법입니다.

## 자동 설정

`akwiki init`이 `.github/workflows/deploy.yml`을 자동으로 생성합니다. push만 하면 됩니다.

```
git push origin main
```

## GitHub 설정

Repository Settings → Pages → Source를 **GitHub Actions**로 변경하세요.

## 워크플로우 동작

1. `main` 브랜치에 push
2. GitHub Actions가 akwiki 바이너리를 다운로드
3. `akwiki build` 실행 (git 이력 전체 checkout — 날짜 추출에 필요)
4. `dist/` 디렉토리를 GitHub Pages로 배포

## 커스텀 도메인

`.akwiki/config.yml`에서 URL을 설정합니다:

```yaml
site:
  url: "https://wiki.example.com"
```

GitHub Pages Settings에서도 Custom domain을 동일하게 설정하세요.

## base path

프로젝트 페이지(`username.github.io/repo-name`)를 사용할 경우:

```yaml
site:
  url: "https://username.github.io/repo-name"
```

빌드 시 모든 경로에 base path가 자동 적용됩니다.

[[설치 가이드]]에서 akwiki 설치 방법을 확인하세요.
