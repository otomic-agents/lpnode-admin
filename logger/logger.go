package logger

import (
	"fmt"
	"os"
	"runtime"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog/zt_formatter"
)

var Log *logrus.Logger
var MsgHandle *logrus.Entry
var TxLogger *logrus.Entry
var ServerLog *logrus.Entry
var PeerManagerLog *logrus.Entry
var LatestHeader *logrus.Entry
var Service *logrus.Entry
var Chain *logrus.Entry
var DataCenter *logrus.Entry
var Market *logrus.Entry
var System *logrus.Entry
var Config *logrus.Entry
var Database *logrus.Entry
var Http *logrus.Entry
var TrustedLogger *logrus.Entry
var BlockInfoCenter *logrus.Entry
var DownLoadCenter *logrus.Entry
var TransactionSend *logrus.Entry
var TransactionCenter *logrus.Entry
var Cluster *logrus.Entry

func init() {
	var formatter = &zt_formatter.ZtFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := f.File
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
		Formatter: nested.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			//HideKeys: true,
			FieldsOrder: []string{"模块", "category"},
		},
	}
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetLevel(logrus.DebugLevel)
	if formatter != nil {
		l.SetFormatter(formatter)
	}
	l.WithField("c", "BscLightNode")
	l.SetOutput(os.Stdout)
	Log = l
	TxLogger = Log.WithField("M", "TxModule")
	TrustedLogger = Log.WithField("M", "Trusted")
	ServerLog = Log.WithField("M", "P2PServer")
	PeerManagerLog = Log.WithField("M", "PeerManager")
	MsgHandle = Log.WithField("M", "MsgHandle")
	LatestHeader = Log.WithField("M", "LatestHeader")
	Chain = Log.WithField("M", "Chain")
	Service = Log.WithField("M", "Service")
	DataCenter = Log.WithField("M", "DataCenter")
	Market = Log.WithField("M", "Market")
	System = Log.WithField("M", "System")
	Config = Log.WithField("M", "Config")
	Database = Log.WithField("M", "Database")
	Http = Log.WithField("M", "Http")
	BlockInfoCenter = Log.WithField("M", "BlockInfoCenter")
	DownLoadCenter = Log.WithField("M", "DownLoadCenter")
	TransactionSend = Log.WithField("M", "TransactionSend")
	TransactionCenter = Log.WithField("M", "TransactionCenter")
	Cluster = Log.WithField("M", "Cluster")
}
