// 1div0
// 2019.10.13

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// ОБЩИЕ ОБЪЯВЛЕНИЯ

package main

// Раздел импорта
import (
    "fmt" // пакет для форматированного ввода вывода
    "os"
    //"context"
    "net/http" // пакет для поддержки HTTP протокола
    "golang.org/x/crypto/acme/autocert" // пакет для работы с LetsEncrypt
    "strings" // пакет для работы с  UTF-8 строками
    //"strconv"
    //"reflect"
    "encoding/json"
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

type OAuthConfiguration struct {
    ClientId string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
}

type HttpsParam struct {
    Enabled bool `json:"enabled"`
    Port int `json:"port"`
    Provider string `json:"provider"`
    SubscriberEmail string `json:"subscriber_email"`
    CertFolder string `json:"cert_folder"`
    CertPemFilename string `json:"cert_pem_filename"`
    KeyPemFilename string `json:"key_pem_filename"`
}

type WebPublicParameters struct {
    DomainName string `json:"domain_name"`
    HttpPort int `json:"http_port"`
    Https HttpsParam `json:"https"`
}

type Configuration struct {
    WebPublic WebPublicParameters `json:"web_public"`
    WebPagesFolder string `json:"web_pages_folder"`
    OAuthVerificationCodePath string `json:"oauth_verification_code_Path"`
    JsSettingsPath string `json:"js_settings_Path"`
    OAuthYandex OAuthConfiguration `json:"oauth_yandex"`
    OAuthGoogle OAuthConfiguration `json:"oauth_google"`
    OAuthGitHub OAuthConfiguration `json:"oauth_github"`
    CommandPathPrefix string `json:"command_path_prefix"`
    ParametersCountLimit int `json:"parameters_count_limit"`
    CommandParametersFileName string `json:"command_parameters_file_name"`
    DbConnectionString string `json:"db_connection_string"`
}

func (cfg *Configuration) ThisApplicationUrl() string {
    var result string
    var port int
    var default_port int

    if cfg.WebPublic.Https.Enabled {
        result = "https://"
        port = cfg.WebPublic.Https.Port
        default_port = 443
    } else {
        result = "http://"
        port = cfg.WebPublic.HttpPort
        default_port = 80
    }

    result += cfg.WebPublic.DomainName
    if port != default_port {
        result += fmt.Sprintf(":%d", port)
    }

    return result
}

type MsSql struct {
    ConnectionString string
}

var g_cfg Configuration
var g_oauth oauth.OAuthCollect
var g_webpages websrv.WebPages
var g_dbman dbman.DataBaseManager
var g_db MsSql

func main() {
    var err error

    err = config_init(&g_cfg)
    if err != nil {
        fmt.Println("Load config file: ", err)
        return
    }

    g_db := MsSql {
        ConnectionString: g_cfg.DbConnectionString,
    }

    err = g_dbman.Init(&g_db, g_cfg.CommandParametersFileName, g_cfg.CommandPathPrefix, g_cfg.ParametersCountLimit)
    if err != nil {
        fmt.Println("Database init: ", err)
        return
    }

    err = g_webpages.Init(g_cfg.WebPagesFolder)
    if err != nil {
        fmt.Println("Webpages init: ", err)
        return
    }

    err = g_oauth.Init(g_cfg.ThisApplicationUrl() + g_cfg.OAuthVerificationCodePath, g_cfg.ThisApplicationUrl(), false)
    if err != nil {
        fmt.Println("Load oauth: ", err)
        return
    }

    oauth_ya := oauth_yandex.OAuthYandex {
        ClientId: g_cfg.OAuthYandex.ClientId,
        ClientPsw: g_cfg.OAuthYandex.ClientSecret,
    }
    g_oauth.AddService(&oauth_ya)

    oauth_go := oauth_google.OAuthGoogle {
        ClientId: g_cfg.OAuthGoogle.ClientId,
        ClientSecret: g_cfg.OAuthGoogle.ClientSecret,
    }
    g_oauth.AddService(&oauth_go)

    oauth_github := oauth_github.OAuthGitHub {
        ClientId: g_cfg.OAuthGitHub.ClientId,
        ClientSecret: g_cfg.OAuthGitHub.ClientSecret,
    }
    g_oauth.AddService(&oauth_github)

    err = web_server_go(&g_cfg)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("Web server started successfully")
    }
}

func web_server_go(cfg *Configuration) error {

    var err error

    fmt.Println("Starting web server...")
    http_ip_port := fmt.Sprintf(":%d", g_cfg.WebPublic.HttpPort)

    if (cfg.WebPublic.Https.Enabled) {

        fmt.Println("SSL setup...")
        https_ip_port := fmt.Sprintf(":%d", g_cfg.WebPublic.Https.Port)

        if g_cfg.WebPublic.Https.Provider == "letsencrypt" {

            mux := http.NewServeMux()
        	mux.HandleFunc("/", HomeHandler)

        	certManager := autocert.Manager{
        		Prompt: autocert.AcceptTOS,
        		Cache:  autocert.DirCache(cfg.WebPublic.Https.CertFolder),
                HostPolicy: autocert.HostWhitelist(cfg.WebPublic.DomainName),
        	}

        	server := &http.Server{
        		Addr: https_ip_port,
        		Handler: mux,
                TLSConfig: certManager.TLSConfig(),
        	}

        	go http.ListenAndServe(http_ip_port, certManager.HTTPHandler(nil))

            err = server.ListenAndServeTLS("", "")

        } else if g_cfg.WebPublic.Https.Provider == "custom" {

            http.HandleFunc("/", HomeHandler)
            go http.ListenAndServe(http_ip_port, http.HandlerFunc(redirectToHttps))
            err = http.ListenAndServeTLS(https_ip_port, g_cfg.WebPublic.Https.CertPemFilename, g_cfg.WebPublic.Https.KeyPemFilename, nil)

        } else {
            fmt.Errorf("Unknown SSL provider %s. Must be 'letsencrypt' or 'custom'.", g_cfg.WebPublic.Https.Provider)
        }

    } else {
        http.HandleFunc("/", HomeHandler) // установим роутер
        err = http.ListenAndServe(http_ip_port, nil) // задаем слушать порт
    }

    return err
}

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
    if g_cfg.WebPublic.Https.Enabled {
        http.Redirect(w, r, g_cfg.ThisApplicationUrl() + r.RequestURI, http.StatusMovedPermanently)
    } else {
        fmt.Fprintf(w, "HTTP not enabled")
    }
}

func config_init(cfg *Configuration) (error) {

    config_file_name := "config.development.json"

    if _, err := os.Stat(config_file_name); err == nil {
    } else if os.IsNotExist(err) {
        config_file_name = "config.json"
    } else {
        return err
    }

    file, err := os.Open(config_file_name)
    if (err != nil) {
        return err
    }
    defer file.Close()

    fmt.Println("Config data loaded from: ", config_file_name)

    decoder := json.NewDecoder(file)
    return decoder.Decode(cfg)
}

// При получении запроса от клиента
func HomeHandler(w http.ResponseWriter, r *http.Request) {

    if (r.URL.Path == "/") {
        g_webpages.Page(w, "/index.html")
    } else if (r.URL.Path == g_cfg.JsSettingsPath) {
        GetSettings(w)
    } else if (strings.HasPrefix(r.URL.Path, g_cfg.CommandPathPrefix)) {
        err := ExecuteCommand(w, r)
        if (err != nil) {
            g_webpages.Response500(w, err)
        }
    } else {
        g_webpages.Page(w, r.URL.Path)
    }

fmt.Fprintf(w, "Hello secure world!")
}

// Клиент запрашивает специальный файл с настройками "/js/_settings.js"
func GetSettings(w http.ResponseWriter) {
    oauth_settings := fmt.Sprintf("const THIS_APP_URL = \"%s\";\n", g_cfg.ThisApplicationUrl())
    oauth_settings += g_oauth.GetSettingsJS()
    fmt.Fprintf(w, oauth_settings)
}

// Клиент вызывает команду
func ExecuteCommand(w http.ResponseWriter, r *http.Request) (error) {

    var err error
    r.ParseForm() //анализ аргументов

    if (r.URL.Path == g_cfg.OAuthVerificationCodePath) {
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
    db, err := sql.Open("sqlserver", g_cfg.DbConnectionString)
    if err != nil {
        return nil, err
    }

    rows, err := db.Query(query_text, cmd_arg...)
    return rows, err
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
