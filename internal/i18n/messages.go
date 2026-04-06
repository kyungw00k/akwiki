package i18n

const (
	// Root
	MsgRootShort Key = "MsgRootShort"
	MsgRootLong  Key = "MsgRootLong"

	// Init
	MsgInitShort  Key = "MsgInitShort"
	MsgInitCreate Key = "MsgInitCreate"
	MsgInitSkip   Key = "MsgInitSkip"
	MsgInitDone   Key = "MsgInitDone"
	MsgInitNext   Key = "MsgInitNext"

	// Build
	MsgBuildShort    Key = "MsgBuildShort"
	MsgBuildBuilding Key = "MsgBuildBuilding"
	MsgBuildDone     Key = "MsgBuildDone"

	// Dev
	MsgDevShort      Key = "MsgDevShort"
	MsgDevServing    Key = "MsgDevServing"
	MsgDevRebuilding Key = "MsgDevRebuilding"
	MsgDevRebuilt    Key = "MsgDevRebuilt"

	// Serve
	MsgServeShort   Key = "MsgServeShort"
	MsgServeServing Key = "MsgServeServing"

	// Note
	MsgNoteShort  Key = "MsgNoteShort"
	MsgNoteAdded  Key = "MsgNoteAdded"
	ErrNoteNoText Key = "ErrNoteNoText"

	// Flags
	FlagPortUsage Key = "FlagPortUsage"

	// Errors
	ErrConfigLoad Key = "ErrConfigLoad"
	ErrBuildFail  Key = "ErrBuildFail"
)

var ko = map[Key]string{
	MsgRootShort:     "개인 위키 정적 사이트 생성기",
	MsgRootLong:      "akwiki는 마크다운 파일로 정적 위키 사이트를 생성합니다.\nakngs의 위키(https://wiki.g15e.com)에서 영감을 받았습니다.",
	MsgInitShort:     "새 위키 생성",
	MsgInitCreate:    "  생성 %s",
	MsgInitSkip:      "  건너뜀 %s (이미 존재)",
	MsgInitDone:      "\n위키가 %s에 초기화되었습니다",
	MsgInitNext:      "다음 단계:\n  cd %s\n  akwiki dev",
	MsgBuildShort:    "정적 사이트 빌드",
	MsgBuildBuilding: "빌드 중...",
	MsgBuildDone:     "%s에 완료 → %s/",
	MsgDevShort:      "라이브 리로드 개발 서버 시작",
	MsgDevServing:    "http://localhost%s 에서 서빙 중",
	MsgDevRebuilding: "리빌드 중...",
	MsgDevRebuilt:    "%s에 리빌드 완료",
	MsgServeShort:    "빌드된 사이트 서빙",
	MsgServeServing:  "%s를 http://localhost%s 에서 서빙 중",
	MsgNoteShort:     "일지 작성",
	MsgNoteAdded:     "작성 완료: %s",
	ErrNoteNoText:    "내용을 입력하세요",
	FlagPortUsage:    "서버 포트",
	ErrConfigLoad:    "설정 로드 실패: %v",
	ErrBuildFail:     "빌드 실패: %v",
}

var en = map[Key]string{
	MsgRootShort:     "Personal wiki static site generator",
	MsgRootLong:      "akwiki generates static wiki sites from markdown files.\nInspired by akngs's wiki (https://wiki.g15e.com).",
	MsgInitShort:     "Create a new wiki",
	MsgInitCreate:    "  create %s",
	MsgInitSkip:      "  skip %s (already exists)",
	MsgInitDone:      "\nWiki initialized in %s",
	MsgInitNext:      "Next steps:\n  cd %s\n  akwiki dev",
	MsgBuildShort:    "Build static site",
	MsgBuildBuilding: "Building...",
	MsgBuildDone:     "Done in %s → %s/",
	MsgDevShort:      "Start development server with live reload",
	MsgDevServing:    "Serving at http://localhost%s",
	MsgDevRebuilding: "Rebuilding...",
	MsgDevRebuilt:    "Rebuilt in %s",
	MsgServeShort:    "Serve the built site",
	MsgServeServing:  "Serving %s at http://localhost%s",
	MsgNoteShort:     "Write a note",
	MsgNoteAdded:     "Note added to %s",
	ErrNoteNoText:    "note text is required",
	FlagPortUsage:    "server port",
	ErrConfigLoad:    "failed to load config: %v",
	ErrBuildFail:     "build failed: %v",
}
