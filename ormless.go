// 1div0
// 2019.10.13

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// ОБЩИЕ ОБЪЯВЛЕНИЯ

package main

// Раздел импорта
import (
    "fmt" // пакет для форматированного ввода вывода
    "os"
    "github.com/spf13/viper" // для чтения файлов конфигурации, подтягивает безумное количество новых модулей (ок. 160)
    "net/http" // пакет для поддержки HTTP протокола
    "golang.org/x/crypto/acme/autocert" // пакет для работы с LetsEncrypt
    "strings" // пакет для работы с  UTF-8 строками
    //"encoding/json"
    _ "github.com/denisenkom/go-mssqldb"
    "database/sql"
    "github.com/a1div0/oauth"
    "github.com/a1div0/oauth_yandex"
    "github.com/a1div0/oauth_google"
    "github.com/a1div0/oauth_github"
    "github.com/a1div0/websrv"
    "github.com/a1div0/dbman"
)
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

type MsSql struct {
    ConnectionString string
}

var g_cfg *viper.Viper
var g_oauth oauth.OAuthCollect
var g_webpages websrv.WebPages
var g_dbman dbman.DataBaseManager
var g_db MsSql

func main() {
    var err error

    err = config_init()
    if err != nil {
        fmt.Println("Load config file: ", err)
        return
    }

    g_db := MsSql {
        ConnectionString: g_cfg.GetString("DbConnectionString"),
    }

    err = g_dbman.Init(&g_db, g_cfg.GetString("CommandParametersFileName"), g_cfg.GetString("CommandPathPrefix"), g_cfg.GetInt("ParametersCountLimit"))
    if err != nil {
        fmt.Println("Database init: ", err)
        return
    }

    err = g_webpages.Init(g_cfg.GetString("WebPagesFolder"))
    if err != nil {
        fmt.Println("Webpages init: ", err)
        return
    }

    err = g_oauth.Init(this_application_url() + g_cfg.GetString("OAuthVerificationCodePath"), this_application_url(), false)
    if err != nil {
        fmt.Println("Load oauth: ", err)
        return
    }

    oauth_ya := oauth_yandex.OAuthYandex {
        ClientId: g_cfg.GetString("OAuthYandex.ClientId"),
        ClientPsw: g_cfg.GetString("OAuthYandex.ClientSecret"),
    }
    g_oauth.AddService(&oauth_ya)

    oauth_go := oauth_google.OAuthGoogle {
        ClientId: g_cfg.GetString("OAuthGoogle.ClientId"),
        ClientSecret: g_cfg.GetString("OAuthGoogle.ClientSecret"),
    }
    g_oauth.AddService(&oauth_go)

    oauth_github := oauth_github.OAuthGitHub {
        ClientId: g_cfg.GetString("OAuthGitHub.ClientId"),
        ClientSecret: g_cfg.GetString("OAuthGitHub.ClientSecret"),
    }
    g_oauth.AddService(&oauth_github)

    err = web_server_go()
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("Web server started successfully")
    }
}

func this_application_url() string {
    var result string

    with_s := g_cfg.GetBool("WebPublic.Https.Enabled")
    s := ""

    if with_s {
        s = "e"
    }

    return "http" + s + "://" + g_cfg.GetString("WebPublic.DomainName")
}

func web_server_go() error {

    var err error

    http_ip := g_cfg.GetString("WebPublic.HttpIp")
    if http_ip == "localhost" || http_ip == "127.0.0.1" {
        http_ip = ""
    }
    http_ip_port := fmt.Sprintf("%s:%d", http_ip, g_cfg.GetInt("WebPublic.HttpPort"))
    fmt.Println("Starting web server. Listen: %s", http_ip_port)

    if (g_cfg.GetBool("WebPublic.Https.Enabled")) {

        https_ip := g_cfg.GetString("WebPublic.Https.Ip")
        if https_ip == "localhost" || https_ip == "127.0.0.1" {
            https_ip = ""
        }
        https_ip_port := fmt.Sprintf("%s:%d", https_ip, g_cfg.GetInt("WebPublic.Https.Port"))
        fmt.Println("Starting SSL. Listen: %s", https_ip_port)
        ssl_provider := g_cfg.GetString("WebPublic.Https.Provider")

        if ssl_provider == "letsencrypt" {

            mux := http.NewServeMux()
        	mux.HandleFunc("/", HomeHandler)

        	certManager := autocert.Manager{
        		Prompt: autocert.AcceptTOS,
        		Cache:  autocert.DirCache(g_cfg.GetString("WebPublic.Https.CertFolder")),
                HostPolicy: autocert.HostWhitelist(g_cfg.GetString("WebPublic.DomainName")),
        	}

        	server := &http.Server{
        		Addr: https_ip_port,
        		Handler: mux,
                TLSConfig: certManager.TLSConfig(),
        	}

        	go http.ListenAndServe(http_ip_port, certManager.HTTPHandler(nil))

            err = server.ListenAndServeTLS("", "")

        } else if ssl_provider == "custom" {

            http.HandleFunc("/", HomeHandler)
            go http.ListenAndServe(http_ip_port, http.HandlerFunc(redirectToHttps))
            err = http.ListenAndServeTLS(https_ip_port, g_cfg.GetString("WebPublic.Https.CertPemFilename"), g_cfg.GetString("WebPublic.Https.KeyPemFilename"), nil)

        } else {
            fmt.Errorf("Unknown SSL provider %s. Must be 'letsencrypt' or 'custom'.", ssl_provider)
        }

    } else {
        http.HandleFunc("/", HomeHandler) // установим роутер
        err = http.ListenAndServe(http_ip_port, nil) // задаем слушать порт
    }

    return err
}

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
    if g_cfg.GetBool("WebPublic.Https.Enabled") {
        http.Redirect(w, r, this_application_url() + r.RequestURI, http.StatusMovedPermanently)
    } else {
        fmt.Fprintf(w, "HTTP not enabled")
    }
}

func config_init() (error) {

    config_file_name := "config-develop"
    config_full_file_name := config_file_name + ".yml"

    if _, err := os.Stat(config_full_file_name); err == nil {
    } else if os.IsNotExist(err) {
        config_file_name = "config"
    } else {
        return err
    }

    fmt.Println("Config data load from: ", config_file_name)

    g_cfg = viper.New()
	g_cfg.SetConfigName(config_file_name)
	g_cfg.AddConfigPath(".")
	g_cfg.AutomaticEnv()
	return g_cfg.ReadInConfig()
}

// При получении запроса от клиента
func HomeHandler(w http.ResponseWriter, r *http.Request) {

    if (r.URL.Path == "/") {
        g_webpages.Page(w, "/index.html")
    } else if (r.URL.Path == g_cfg.GetString("JsSettingsPath")) {
        GetSettings(w)
    } else if (strings.HasPrefix(r.URL.Path, g_cfg.GetString("CommandPathPrefix"))) {
        err := ExecuteCommand(w, r)
        if (err != nil) {
            g_webpages.Response500(w, err)
        }
    } else {
        g_webpages.Page(w, r.URL.Path)
    }
}

// Клиент запрашивает специальный файл с настройками "/js/_settings.js"
func GetSettings(w http.ResponseWriter) {
    oauth_settings := fmt.Sprintf("const THIS_APP_URL = \"%s\";\n", this_application_url())
    oauth_settings += g_oauth.GetSettingsJS()
    fmt.Fprintf(w, oauth_settings)
}

// Клиент вызывает команду
func ExecuteCommand(w http.ResponseWriter, r *http.Request) (error) {

    var err error
    r.ParseForm() //анализ аргументов

    if (r.URL.Path == g_cfg.GetString("OAuthVerificationCodePath")) {
        err = g_oauth.OnRecieveVerificationCode(w, r, &g_dbman);
    } else {

        user_id, err := g_oauth.CheckAuth(w, r)
        if (err != nil) {
            return err
        }
        // if (user_id == 0) {
        //     return nil;
        // }

        fmt.Println("")
        fmt.Println("RequestURI=", r.RequestURI)

        err = g_dbman.ExecuteCommand(w, r, user_id);
        if err != nil {
            return err
        }
    }

    return err
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
                                                        //...sql.NamedArg
func (sdb *MsSql) SqlExecute(cmd_name string, cmd_arg ...interface{}) (*sql.Rows, error) {

    var (
        arg_named sql.NamedArg
        ok bool
    )

    argument_text := ""

    for _, arg_i := range cmd_arg {
fmt.Printf("Raw arg_i: %+v\n", arg_i)
        if (argument_text != "") {
            argument_text = argument_text + ","
        }
        if arg_named, ok = arg_i.(sql.NamedArg); ok {
            argument_text = argument_text + "@" + arg_named.Name
        }else{
            return nil, fmt.Errorf("Arguments must be type is sq.NamedArg!")
        }
    }

    query_text := "EXECUTE " + cmd_name + " " + argument_text
fmt.Println(query_text)
    // всё что ниже - можно вынести в отдельный модуль, или часть - до получения результата, а парсинг результата в JSON - оставить здесь
    db, err := sql.Open("sqlserver", g_cfg.GetString("DbConnectionString"))
    if err != nil {
        return nil, err
    }

    rows, err := db.Query(query_text, cmd_arg...)
    return rows, err
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
