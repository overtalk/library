package mysql

import (
	"fmt"

	"github.com/gohouse/converter"
)

const (
	defaultSavePath = "./model"
	defaultHost     = "127.0.0.1"
	defaultPort     = 3306
)

var (
	t2t *converter.Table2Struct
)

func init() {
	// 初始化
	t2t = converter.NewTable2Struct()
	// 个性化配置
	t2t.Config(&converter.T2tConfig{
		// 如果字段首字母本来就是大写, 就不添加tag, 默认false添加, true不添加
		RmTagIfUcFirsted: false,
		// tag的字段名字是否转换为小写, 如果本身有大写字母的话, 默认false不转
		TagToLower: false,
		// 字段首字母大写的同时, 是否要把其他字母转换为小写,默认false不转换
		UcFirstOnly: false,
		//// 每个struct放入单独的文件,默认false,放入同一个文件(暂未提供)
		//SeperatFile: false,
	})
}

type ModelBuilder struct {
	SavePath string // default: `./model`, 一般采用相对路径
	Host     string // default: `127.0.0.1`
	Port     int    // default: `3306`
	User     string
	Pwd      string
	DB       string
	Tables   []string // 生成的model会存储在SavePath之下，文件名就是表名
}

func (m *ModelBuilder) Build() {
	// 检查
	m.check()

	for _, table := range m.Tables {
		fmt.Printf("convert table[%s]...\n", table)
		filePath := fmt.Sprintf("%s/%s.go", m.SavePath, table)
		if err := t2t.
			// 指定某个表,如果不指定,则默认全部表都迁移
			Table(table).
			// 表前缀
			//Prefix("prefix_").
			// 是否添加json tag
			//EnableJsonTag(true).
			// 生成struct的包名(默认为空的话, 则取名为: package model)
			PackageName("model").
			// tag字段的key值,默认是orm
			TagKey("json").
			// 是否添加结构体方法获取表名
			RealNameMethod("TableName").
			// 生成的结构体保存路径
			SavePath(filePath).
			// 数据库dsn,这里可以使用 t2t.DB() 代替,参数为 *sql.DB 对象
			Dsn(m.dns()).
			// 执行
			Run(); err != nil {
			fmt.Printf("failed to convert table[%s], error = %v\n", table, err)
			break
		}
	}
}

// check : 设置默认值
func (m *ModelBuilder) check() {
	if len(m.SavePath) == 0 {
		m.SavePath = defaultSavePath
	}

	if len(m.Host) == 0 {
		m.Host = defaultHost
	}

	if m.Port == 0 {
		m.Port = defaultPort
	}
}

func (m *ModelBuilder) dns() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", m.User, m.Pwd, m.Host, m.Port, m.DB)
}
