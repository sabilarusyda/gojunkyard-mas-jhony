# Form
Package form implements request struct binding, filter, and validator by using struct tag. This package is backed by these repositories.
1. `form`: github.com/go-playground/form
2. `validator`: github.com/go-playground/validator
3. `mold`: github.com/go-playground/mold

### Dictionary
1. `Binding`: Converting user request parameter into a struct. `(form)`
2. `Filter`: Normalizing user parameters such as space triming, capitalize, etc. `(mold)`
3. `Validation`: Check is user request has a valid format or not. `validator`
4. `Tag`: Last column in struct definition. Started by  \`...\`
5. `Tag-form`: mapping of user request parameters `form:"id"`. [#docs](github.com/go-playground/form)
6. `Tag-mod`: filter request parameters `mod:"id"`. [#docs](github.com/go-playground/mold)
7. `Tag-validate`: filter request parameters `validate:"required,eq=m|eq=f"`. [#docs](github.com/go-playground/validator)

## How to use?
There are 3 ways to use this package.
1. Binder
2. Validator

### 1. Binder
Binder is used to perform form binding, filtering, and validation.

**Bind**
You can use bind feature to perform form binding, filtering, and validation
> Example:
```
type  User  struct {
	ID int64  `form:"id"`
	Name string  `form:"name" mod:"trim"`
	Email string  `form:"email" mod:"trim"`
	Gender string  `form:"gender" validate:"required,eq=m|eq=f"`
	BirthDate time.Time `form:"bdate"`
}

func handleGetUser(w http.ResponseWriter, r *http.Request)) {
	var got User
	err  :=  form.Bind(&got, r)
	if err !=  nil {
		// do something error
	}
	// do something ok
}
```

**BindFlag**
If you only use several feature, such as `filter` only. BindFlag is the best practice use because it will skip the other feature.
> Example:
```
type  User  struct {
	ID int64  `form:"id"`
	Name string  `form:"name" mod:"trim"`
	Email string  `form:"email" mod:"trim"`
	Gender string  `form:"gender"`
	BirthDate time.Time `form:"bdate"`
}

func handleGetUser(w http.ResponseWriter, r *http.Request)) {
	var got User
	err  :=  form.BindFlag(form.Bfilter, &got, r)
	if err !=  nil {
		// do something error
	}
	// do something ok
}
```
> Notes:
> `Bnone` only does a form binding
> `Bfilter` performs form binding and filter
> `Bvalidate` performs form binding and validation
> `Bstd` performs form binding, filter, and validation

### 2. Validator
Validator is used to performs filtering and validation.

**Validate**
You can use validate feature to perform filtering and validation
> Example:
```
type  User  struct {
	ID int64  `form:"id"`
	Name string  `form:"name" mod:"trim"`
	Email string  `form:"email" mod:"trim"`
	Gender string  `form:"gender" validate:"required,eq=m|eq=f"`
	BirthDate time.Time `form:"bdate"`
}

func handleGetUser(w http.ResponseWriter, r *http.Request)) {
	user := User {
		ID: 1,
		Name: "My Name",
		Email: "my@email.com",
		Gender: "m"
	}
	err  :=  form.Validate(&user, r)
	if err !=  nil {
		// do something error
	}
	// do something ok
}
```

**ValidateFlag**
You can use validate feature to perform filtering and validation
> Example:
```
type  User  struct {
	ID int64  `form:"id"`
	Name string  `form:"name" mod:"trim"`
	Email string  `form:"email" mod:"trim"`
	Gender string  `form:"gender" validate:"required,eq=m|eq=f"`
	BirthDate time.Time `form:"bdate"`
}

func handleGetUser(w http.ResponseWriter, r *http.Request)) {
	user := User {
		ID: 1,
		Name: "My Name",
		Email: "my@email.com",
		Gender: "m"
	}
	err  :=  form.ValidateFlag(form.Vfilter, &user, r)
	if err !=  nil {
		// do something error
	}
	// do something ok
}
```
> Notes:
> `Vnone` only does a validation
> `Vfilter` performs filtering and validation
