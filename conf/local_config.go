package conf

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/samber/lo"
	"go.yaml.in/yaml/v3"
)

// **
// ** 加载本地配置文件
// **

var defaultPath = "config/config.yaml"

func optionPaths(ops ...ConfigOption) []string {
	paths := make([]string, len(ops))
	if len(ops) == 0 {
		return []string{defaultPath}
	}
	for i, op := range ops {
		paths[i] = op()
	}
	return paths
}

type Config[T any] struct {
	filePaths []string
	data      T
	mtx       sync.RWMutex
}

func NewConfig[T any](ops ...ConfigOption) *Config[T] {
	return &Config[T]{
		filePaths: optionPaths(ops...),
	}
}

func (dr *Config[T]) Load() error {
	dr.mtx.Lock()
	defer dr.mtx.Unlock()
	// 读取所有的配置文件生成map
	resultMap := make(map[string]any)
	for _, path := range dr.filePaths {
		yamlFile, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Error reading file: %v", err)
			continue
		}

		var fileData map[string]any
		// 解析 YAML 文件
		if err = yaml.Unmarshal(yamlFile, &fileData); err != nil {
			log.Printf("Error unmarshal file: %v", err)
			continue
		}
		resultMap = lo.Assign(resultMap, fileData)
	}
	// 转化成json，然后映射成结构体
	jsonData, err := yaml.Marshal(resultMap)
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
		return err
	}
	if err := yaml.Unmarshal(jsonData, &dr.data); err != nil {
		log.Printf("Error unmarshaling data: %v", err)
		return err
	}

	log.Printf("Data reloaded success")
	return nil
}

func (dr *Config[T]) Get() T {
	dr.mtx.RLock()
	defer dr.mtx.RUnlock()
	return dr.data
}

func (dr *Config[T]) Watch() {
	for _, path := range dr.filePaths {
		initWG := sync.WaitGroup{}
		initWG.Add(1)
		go func() {
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				log.Printf("Error creating watcher: %s", err)
				os.Exit(0)
			}
			defer watcher.Close()

			eventsWG := sync.WaitGroup{}
			eventsWG.Add(1)
			go func() {
				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							eventsWG.Done()
							return
						}
						// 只处理写入事件
						if event.Op&fsnotify.Write == fsnotify.Write {
							log.Printf("File modified: %s", event.Name)
							// 小延迟确保文件完全写入
							time.Sleep(100 * time.Millisecond)
							dr.Load()
						}
					case err, ok := <-watcher.Errors:
						if !ok {
							eventsWG.Done()
							return
						}
						log.Printf("Watcher error: %v", err)
					}
				}
			}()

			err = watcher.Add(path)
			if err != nil {
				log.Printf("Error adding file to watcher: %s", err)
				initWG.Done()
				return
			}
			initWG.Done()
			eventsWG.Wait()
		}()
		initWG.Wait()
	}
}
