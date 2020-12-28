/*

Package server provides a Glint server implementation.

Example:

        s := &server.Server{
                DataDir:          "/var/lib/glint/",
                Host:             "glintcore.net",
                Port:             "443",
                TLSCertFile:      "/etc/letsencrypt/live/glintcore.net/cert.pem",
                TLSKeyFile:       "/etc/letsencrypt/live/glintcore.net/privkey.pem",
                PostgresHost:     "localhost",
                PostgresPort:     "5432",
                PostgresUser:     "glint",
                PostgresPassword: "password_goes_here",
                PostgresDBName:   "glint",
        }
        err := s.ListenAndServe()

*/
package server
