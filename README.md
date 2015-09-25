# generator-aurelia-auth-go
A generator for backend of an aurelia app with JWT Authentication (login, register, logout) pages as well as a simple user profile page

## Frontend/ui
An [Aurelia](http://aurelia.io/) frontend/ui that works with this is https://github.com/francoishill/generator-aurelia-auth-ui


## Current features

### Basic workflow
- Authentication with JWT (login, register, logout)
- User profile details and changing of Full Name
- Setting with auto-reload (using github.com/francoishill/golang-common-ddd)
- Mocking out stuff: look how it's done in `Context/RouterContext/RouterContext.go` method `getUserRepository()`

### DDD and Repository pattern
Using Repository pattern with Domain Driven Design (golang package https://github.com/francoishill/golang-common-ddd/).


### Interface centric
Revolves mainly around golang `interface{}` to keep things loosely coupled


### MVC (or at least the M & C)

#### Routers & Controller
A router defines what URL pattern maps to a Controllers.

A controller maps to a single url pattern and can handle multiple method types. Look at the login controller for an example:
```
type controller struct{}

func (c *controller) Post(w http.ResponseWriter, r *http.Request, ctx *RouterContext) {
    ctx.AuthenticationService.BaseLoginHandler(w, r)
}

func NewAuthLoginController() *controller {
    return &controller{}
}
```

Note the signature `Post(..., ctx *RouterContext)`. It could also be `Get(..., ctx *RouterContext)`. For supported methods, have a look at the `Routers/Setup/Controller.go` file.

This `RouterContext` is what holds all our Logging, Services, Repositories, Helpers, etc.


#### Model
This will be the Repositories and Entities


### Db Migration example
Example included in file `Db/MysqlMigrations/2015-09-21--10-00_initial_db_setup.sql`.
