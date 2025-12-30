/* Este se enncargara de tener las configuraciones que nosotros vamos a necesitar cuando generemos nuestro docker */

SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN= 0;

SET @@GLOBAL.GTID_PURGED= '';

CREATE DATABASE IF NOT EXISTS `go_backend_user`;