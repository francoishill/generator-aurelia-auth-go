var generators = require('yeoman-generator');
var fs = require('fs');

function trimEndingSlash(url) {
    while (url.length > 0 && url[url.length - 1] === "/") {
        url = url.substr(0, url.length - 1);
    }
    return url;
}

function isInt(x) {
   var y = parseInt(x, 10);
   return !isNaN(y) && x == y && x.toString() == y.toString();
}

function validateIntAnswer(ans) {
    if (isInt(ans))
        return true;
    return "Invalid integer value";
}

function sanitizePath(p) {
    var s = p.replace(/\\/g, "/");
    while (s.length > 0 && s[s.length-1] == '/') {
        s = s.substr(0, s.length-1);
    }
    return s;
}

function trimStartingSlash(s) {
    while (s.length > 0 && s[0] == '/')
        s = s.substr(1);
    return s;
}

function getRelativePathFromGoPath() {
    var GOPATH = process.env.GOPATH;
    if (!GOPATH || GOPATH.trim() == '') {
        throw new Error("GOPATH is not defined.");
    }
    GOPATH = sanitizePath(GOPATH);
    
    var currentWorkinkDir = sanitizePath(process.cwd());
    if (currentWorkinkDir.length <= GOPATH.length
        || currentWorkinkDir.substr(0, GOPATH.length) != GOPATH) {
        throw new Error("Current Working Directory ('" + currentWorkinkDir + "') must be inside your GOPATH ('" + GOPATH + "').");
    }

    var relToGoPath = trimStartingSlash(currentWorkinkDir.substr(GOPATH.length));
    if (relToGoPath.length < 4 || relToGoPath.substr(0, 4) != 'src/') {
        throw new Error("Current Working Directory ('" + currentWorkinkDir + "') must be inside GOPATH/src, GOPATH ('" + GOPATH + "')")
    }

    var relToSrc = relToGoPath.substr(4);
    if (relToSrc.trim().length = 0) {
        throw new Error("Invalid Current Working Directory ('" + currentWorkinkDir + "'), must be a subfolder inside GOPATH/src/")
    }

    return relToSrc;
}

var currentDirRelativeToGopathSrc = getRelativePathFromGoPath();
console.log("Found your current working dir will be '" + currentDirRelativeToGopathSrc + "' relative to GOPATH/src.")

module.exports = generators.Base.extend({
    prompting: function () {
        var done = this.async();

        var allVariables = {};

        var nonTemplateFilesRelativePathsToCopy = [
        ];

        var templateFileRelativePathsToCopy = [
            'run_gin.bat',
            'server.go',

            'config/server.gcfg',

            'Authentication/DefaultAuthenticationService/service.go',
            'Authentication/DefaultAuthUserHelperService/service.go',
            'Context/RouterContext/RouterContext.go',
            'Controllers/Auth/Login/controller.go',
            'Controllers/Auth/Logout/controller.go',
            'Controllers/Auth/Register/controller.go',
            'Controllers/UserDetails/controller.go',
            'Db/MysqlMigrations/2015-09-21--10-00_initial_db_setup.sql',
            'Entities/User/User.go',
            'Entities/User/UserRepository.go',
            'Interface/Authentication/AuthenticationService.go',
            'Repositories/User/DbUserRepository.go',
            'Repositories/User/MockUserRepository.go',
            'Routers/Routers.go',
            'Routers/Setup/Controller.go',
            'Routers/Setup/ControllerMethod.go',
            'Routers/Setup/negroniMiddleware.go',
            'Routers/Setup/RegisterRouters.go',
            'Routers/Setup/Router.go',
            'Routers/Setup/RouterBuilder.go',
            'Settings/DefaultSettings/config.go',
            'Settings/DefaultSettings/settings.go',
            'Settings/DefaultSettings/settingValues.go',
            'Settings/Settings.go'
        ];

        var questions = [
            {
                type    : 'input',
                name    : 'frontendurl',
                message : 'Url of your UI server (used to enable CORS)',
                default : 'http://localhost:12401'
            },
            {
                type    : 'input',
                name    : 'backendUrl',
                message : 'Backend url of your server',
                default : 'http://localhost:12301'
            },
            {
                type    : 'input',
                name    : 'jwtExpirationHours',
                message : 'Number of hours for JWT token expiry',
                default : '8760',
                validate: validateIntAnswer
            }
        ];

        this.prompt( questions, function( answers ) {
            allVariables.OWN_GO_IMPORT_PATH = currentDirRelativeToGopathSrc;
            allVariables.FRONTEND_URL = trimEndingSlash(answers.frontendurl);
            allVariables.BACKEND_URL = trimEndingSlash(answers.backendUrl);
            allVariables.JWT_EXPIRATION_HOURS = answers.jwtExpirationHours;

            for (var i = 0; i < templateFileRelativePathsToCopy.length; i++) {
                var relPath = templateFileRelativePathsToCopy[i];
                this.fs.copyTpl(this.templatePath(relPath), this.destinationPath(relPath), allVariables);
            }

            for (var i = 0; i < nonTemplateFilesRelativePathsToCopy.length; i++) {
                var relPath = nonTemplateFilesRelativePathsToCopy[i];
                this.fs.copy(this.templatePath(relPath), this.destinationPath(relPath));
            }

            fs.writeFileSync(this.destinationPath('.gitignore'), [
                'config/server.gcfg',
                'config/jwt.rsa',
                'config/jwt.rsa.pub'
            ].join("\n"));

            done();
        }.bind(this));
    },
    install: function () {
        var that = this;

        var goGetProcess = this.spawnCommand('go', ['get', './...']);
        goGetProcess.on('close', function(code1) {
            if (code1 !== 0) {
                that.log.error('GO GET process exited with code ' + code1);
                return;
            }

            /*var npmProcess = that.spawnCommand('npm', ['install', '-y']);
            npmProcess.on('close', function(code2) {
                if (code2 !== 0) {
                    that.log.error('npm process exited with code ' + code2);
                    return;
                }
            })*/
        });
    }
});
