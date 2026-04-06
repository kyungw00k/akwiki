package cli

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kyungw00k/akwiki/internal/builder"
	"github.com/kyungw00k/akwiki/internal/config"
	"github.com/kyungw00k/akwiki/internal/i18n"
	"github.com/spf13/cobra"
)

var devPort string

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: i18n.T(i18n.MsgDevShort),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootDir := "."
		cfg, err := config.Load(rootDir)
		if err != nil {
			return fmt.Errorf(i18n.T(i18n.ErrConfigLoad), err)
		}
		outDir := filepath.Join(rootDir, cfg.Build.OutDir)

		fmt.Println(i18n.T(i18n.MsgBuildBuilding))
		if err := builder.Build(rootDir, outDir); err != nil {
			return fmt.Errorf(i18n.T(i18n.ErrBuildFail), err)
		}

		go watchAndRebuild(rootDir, outDir)

		addr := ":" + devPort
		fmt.Println(i18n.Tf(i18n.MsgDevServing, addr))
		return http.ListenAndServe(addr, http.FileServer(http.Dir(outDir)))
	},
}

func init() {
	devCmd.Flags().StringVarP(&devPort, "port", "p", "3000", i18n.T(i18n.FlagPortUsage))
	rootCmd.AddCommand(devCmd)
}

func watchAndRebuild(rootDir, outDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("watch error: %v", err)
		return
	}
	defer watcher.Close()

	// Watch pages and config directories
	for _, dir := range []string{
		filepath.Join(rootDir, "pages"),
		filepath.Join(rootDir, ".akwiki"),
	} {
		watcher.Add(dir)
	}

	var debounce <-chan time.Time
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
				debounce = time.After(300 * time.Millisecond)
			}
		case <-debounce:
			fmt.Println(i18n.T(i18n.MsgDevRebuilding))
			start := time.Now()
			if err := builder.Build(rootDir, outDir); err != nil {
				log.Printf("build error: %v", err)
			} else {
				fmt.Println(i18n.Tf(i18n.MsgDevRebuilt, time.Since(start).Round(time.Millisecond)))
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("watch error: %v", err)
		}
	}
}
